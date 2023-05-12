package main

import (
	"errors"
	"net/http"
	"time"

	"greenlight.wook.net/internal/data"
	"greenlight.wook.net/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 데이터베이스에 사용자 레코드가 생성된 후 사용자에 대한 새 활성화 토큰을 생성합니다.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		// 이제 이메일 템플릿에 전달할 데이터가 여러 개 있으므로 데이터의 '보유 구조'
		//역할을 할 맵을 만듭니다. 여기에는 사용자의 ID와 함께 활성화 토큰의 일반 텍스트 버전이 포함됩니다.
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		// 위의 map 을 동적 데이터로 전달하여 환영 이메일을 보냅니다.
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// 요청 본문에서 일반 텍스트 활성화 토큰을 파싱합니다.
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// 클라이언트가 제공한 일반 텍스트 토큰의 유효성을 검사합니다.
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// GetForToken() 메서드를 사용하여 토큰과 연결된 사용자의 세부 정보를 검색합니다
	// (잠시 후에 생성할 것입니다). 일치하는 레코드가 발견되지 않으면 클라이언트가
	// 제공한 토큰이 유효하지 않음을 알립니다.
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// 사용자의 활성화 상태를 업데이트합니다.
	user.Activated = true
	// 업데이트된 사용자 기록을 데이터베이스에 저장하고 movie 레코드와 동일한
	// 방식으로 편집 충돌이 있는지 확인합니다.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// 모든 것이 성공적으로 진행되면 사용자의 모든 활성화 토큰을 삭제합니다.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// 업데이트된 사용자 세부 정보를 JSON 응답으로 클라이언트에 보냅니다.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

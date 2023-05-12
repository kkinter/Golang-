package main

import (
	"errors"
	"net/http"

	"greenlight.wook.net/internal/data"
	"greenlight.wook.net/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// 요청 본문에서 예상되는 데이터를 저장할 익명 구조체를 생성합니다.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 요청 본문을 익명 구조체로 파싱합니다.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// 요청 본문에서 새 User 구조체에 데이터를 복사합니다. Activated 필드는
	// 기본적으로 0 값인 false를 갖기 때문에 꼭 필요한 것은 아니지만 Activated 필드를 false로 설정했습니다.
	// 하지만 이렇게 명시적으로 설정하면 코드를 읽는 모든 사람에게 의도를 명확하게 전달하는 데 도움이 됩니다.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Password.Set() 메서드를 사용하여 해시 및 일반 텍스트 비밀번호를 생성하고 저장합니다.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	// 사용자 구조체의 유효성을 검사하고 검사 중 하나라도 실패하면 오류 메시지를 클라이언트에 반환합니다.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 사용자 데이터를 데이터베이스에 삽입합니다.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// ErrDuplicateEmail 오류가 발생하면 v.AddError() 메서드를 사용하여
		// 유효성 검사기 인스턴스에 메시지를 수동으로 추가한 다음 failedValidationResponse()
		// 헬퍼를 호출합니다.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Mailer에서 Send() 메서드를 호출하여 사용자의 이메일 주소, 템플릿 파일 이름,
	// 새 사용자의 데이터가 포함된 User 구조체를 전달합니다.
	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

package main

import (
	"errors"
	"fmt"
	"net/http"

	"greenlight.wook.net/internal/data"
	"greenlight.wook.net/internal/validator"
)

func (app *application) creteMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResopnse(w, r, err)
		return
	}

	// movie 변수에는 Movie 구조체에 대한 *포인터*가 포함되어 있습니다.
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// 새 유효성 검사기를 초기화합니다.
	v := validator.New()

	// 검사 중 하나라도 실패하면 ValidateMovie() 함수를 호출하고
	// 오류가 포함된 응답을 반환합니다.
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 유효성이 검사된 movies 구조체에 대한 포인터를 전달하여 movies 모델에서
	// Insert() 메서드를 호출합니다. 그러면 데이터베이스에 레코드가 생성되고
	// 시스템에서 생성된 정보로 movie 구조체가 업데이트됩니다.
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// HTTP 응답을 보낼 때 클라이언트가 새로 생성된 리소스를 찾을 수 있는 URL을 알 수 있도록,
	// Location 헤더를 포함하려고 합니다. 빈 http.Header 맵을 만든 다음
	// Set() 메서드를 사용하여 새 Location 헤더를 추가하고 URL에 새 movie 에
	// 대해 시스템에서 생성된 ID를 보간합니다.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// 201 Created 상태 코드, 응답 본문의 movie 데이터, Location 헤더가
	// 포함된 JSON 응답을 작성합니다.
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Get() 메서드를 호출하여 특정 동영상에 대한 데이터를 가져옵니다.
	// 또한 errors.Is() 함수를 사용하여 data.ErrRecordNotFound 오류를 반환하는지
	// 확인해야 하며, 이 경우 클라이언트에 404 찾을 수 없음 응답을 전송합니다.
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResopnse(w, r, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

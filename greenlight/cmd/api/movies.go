package main

import (
	"fmt"
	"net/http"
	"time"

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

	// input 구조체의 값을 새 Movie 구조체로 복사합니다.
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

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		// 새로운 notFoundResponse() 헬퍼를 사용합니다.
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Avengers",
		Runtime:   120,
		Genres:    []string{"action", "war", "hero"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// 새로운 serverErrorResponse() 헬퍼를 사용합니다.
		app.serverErrorResponse(w, r, err)
	}

}

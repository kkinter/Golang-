package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.wook.net/internal/data"
)

func (app *application) creteMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
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

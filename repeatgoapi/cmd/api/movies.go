package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"wook.net/internal/data"
)

func (app *application) creteMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := chi.URLParam(r, "id")

	movie := data.Movie{
		ID:        pathParams,
		CreatedAt: time.Now(),
		Title:     "Avengers",
		Runtime:   120,
		Genres:    []string{"action", "war", "hero"},
		Version:   1,
	}

	err := app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "서버에 문제가 발생하여 요청을 처리할 수 없습니다.", http.StatusInternalServerError)
	}
}

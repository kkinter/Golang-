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
		http.NotFound(w, r)
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

	// 일반 동영상 구조체를 전달하는 대신 envelope{"movie": movie} 인스턴스를 생성하고
	// 이를 writeJSON()에 전달합니다.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "서버에 문제가 발생하여 요청을 처리할 수 없습니다.", http.StatusInternalServerError)
	}

}

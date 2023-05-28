package main

import (
	"fmt"
	"net/http"
	"pracweb/internal/data"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (app *application) healthzHandler(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status      string `json:"status"`
		Environment string `json:"environment"`
		Version     string `json:"version"`
	}{
		Status:      "available",
		Environment: app.config.env,
		Version:     version,
	}

	err := app.writeJSON(w, http.StatusOK, payload, nil)
	if err != nil {
		app.logger.Error("write json error", zap.Error(err))
		http.Error(w, "서버에 이상이 생겨 요청을 처리할 수 없습니다", http.StatusInternalServerError)
	}
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))

	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        int64(userID),
		CreatedAt: time.Now(),
		Title:     "anything",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Error("write json error", zap.Error(err))
		http.Error(w, "서버에 이상이 생겨 요청을 처리할 수 없습니다", http.StatusInternalServerError)
	}
}

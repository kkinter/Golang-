package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Get("/v1/healthcheck", app.healthcheckHandler)
	mux.Get("/v1/movies", app.creteMovieHandler)
	mux.Get("/v1/movies/{id}", app.showMovieHandler)

	return mux
}

package main

import "net/http"

type Env struct {
	Error string `json:"error"`
}

// print log error
func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	var env Env

	err := app.writeJSON(w, status, env.Error, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}

}

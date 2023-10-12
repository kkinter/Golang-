package main

import (
	"net/http"
)

// health check
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "서버에 문제가 발생하여 요청을 처리할 수 없습니다.", http.StatusInternalServerError)
	}
}

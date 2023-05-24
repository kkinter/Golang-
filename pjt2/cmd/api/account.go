package main

import (
	"context"
	"log"
	"net/http"
	db "simple_bank/db/sql"
)

type createAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}

func (app *application) createAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	log.Println("fine")

	account, err := app.store.Queries.CreateAccount(context.Background(), arg)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, account)
}

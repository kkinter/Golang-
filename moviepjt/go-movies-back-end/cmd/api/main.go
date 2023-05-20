package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 8000

type application struct {
	Domain string
}

// entry point for app
func main() {
	// app config
	var app application

	// read from command line e.g) flag

	// connect to the db

	app.Domain = "example.com"

	// start a web server
	log.Println("Starting app on port:", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

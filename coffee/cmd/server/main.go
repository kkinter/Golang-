package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var port = os.Getenv("PORT")

type Config struct {
	Port string
}

type Application struct {
	Config Config
}

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}

	return srv.ListenAndServe()
}

func main() {
	var cfg Config
	cfg.Port = port

	db.openDB()

	app := &Application{
		Config: cfg,
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

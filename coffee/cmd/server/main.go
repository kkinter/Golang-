package main

import (
	"coffee-api/db"
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

	dbConn, err := db.ConnectPostgres(os.Getenv("DSN"))
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	defer dbConn.DB.Close()

	app := &Application{
		Config: cfg,
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

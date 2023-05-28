package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *zap.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "API PORT")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	// zap logger init
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync()

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Sugar().Infof("Starting %s server on %d", cfg.env, cfg.port)
	// logger.Info(fmt.Sprintf("Starting %s server on %d", cfg.env, cfg.port))
	err := srv.ListenAndServe()

	logger.Fatal("Something went terribly wrong",
		zap.String("context", "main"),
		zap.Error(err))

}

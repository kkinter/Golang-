package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Router  chi.Router
	Service CommentService
	Server  *http.Server
}

func NewHandler(service CommentService) *Handler {
	h := &Handler{
		Service: service,
	}

	h.Router = chi.NewRouter()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})

	h.Router.Post("/api/v1/comment", h.PostComment)
	h.Router.Get("/api/v1/comment/{id}", h.GetComment)
	h.Router.Put("/api/v1/comment/{id}", h.UpdateComment)
	h.Router.Delete("/api/v1/comment/{id}", h.DeleteComment)
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
			//return err
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("shut down gracefully")

	return nil
}

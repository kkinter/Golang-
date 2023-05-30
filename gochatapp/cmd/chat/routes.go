package main

import (
	"fmt"
	"gochatapp/internel/redisrepo"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func routes() http.Handler {
	redisClient := redisrepo.InitialiseRedis()
	defer redisClient.Close()

	redisrepo.CreateFetchChatBetweenIndex()

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "healthz: OK")
	})

	mux.Post("/register", registerHandler)
	mux.Post("/login", loginHandler)
	mux.Post("/verfiy-contact", verifyContactHandler)
	mux.Get("/chat-history", chatHistoryHandler)
	mux.Get("/contact-list", contactListHandler)

	return mux
}

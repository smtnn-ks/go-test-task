package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type secretsType struct {
	jwtAccessSecret  string
	jwtRefreshSecret string
}

var secrets secretsType

func InitializeServer(port, jwtAccessSecret, jwtRefreshSecret string) {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/{userId}", signupHandler)
	router.Get("/", signinHandler)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	secrets = secretsType{jwtAccessSecret, jwtRefreshSecret}

	log.Printf("Running on port %v...", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error on running a server:", err)
	}
}

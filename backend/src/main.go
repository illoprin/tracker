package main

import (
	"log"
	"net/http"
	"os"
	"tracker-backend/src/app"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("host is undefined")
	}

	// TODO: init logger (slog)

	// TODO: init mongo connection
	// TODO: init redis connection

	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("{\"ping\":\"pong\"}"))
	})

	// create app instance
	app := app.NewApp(router, port)

	app.Run()
}

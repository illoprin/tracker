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

	host := os.Getenv("HTTP_SERVER")
	if host == "" {
		log.Fatalf("host is undefined")
	}

	// TODO: Init logger (slog)
	// TODO: Init storage (go-mongo-model)

	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("{\"ping\":\"pong\"}"))
	})

	// create app instance
	app := app.NewApp(router, host)

	app.Run()
}

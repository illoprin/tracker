package rest

import (
	"tracker-backend/internal/interfaces/rest/handlers"

	"github.com/go-chi/chi/v5"
)

func NewAppRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/ping", handlers.Ping)

	return r
}

package routers

import (
	"tracker-backend/internal/interfaces/rest/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterGenreRoutes() chi.Router {
	r := chi.NewRouter()
	h := handlers.NewGenreHandler()

	r.Get("/", h.GetAll)
	return r
}

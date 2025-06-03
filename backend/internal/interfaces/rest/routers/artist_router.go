package routers

import (
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/handlers"
	"tracker-backend/internal/interfaces/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterArtistRoutes(s *services.ArtistService, mw middleware.MiddlewareFunc) chi.Router {
	h := handlers.NewArtistHandler(s)
	r := chi.NewRouter()

	r.Use(mw)
	r.Post("/", h.Create)

	r.Get("/my", h.GetMy)
	r.Get("/{id}", h.GetByID)
	r.Get("/{id}/stats", h.GetStats)
	r.Get("/{id}/albums", h.GetAlbums)
	r.Get("/{id}/popular", h.GetPopularTracks)

	r.Delete("/{id}", h.DeleteByID)

	return r
}

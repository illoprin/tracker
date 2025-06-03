package routers

import (
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/handlers"
	"tracker-backend/internal/interfaces/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterAlbumRoutes(s *services.AlbumService, mw middleware.MiddlewareFunc) chi.Router {
	h := handlers.NewAlbumHandler(s)
	r := chi.NewRouter()

	r.Use(mw)
	r.Post("/", h.Create)
	r.Post("/{id}/tracks", h.CreateTrack)
	r.Get("/{id}/tracks", h.GetTracksByAlbumID)
	r.Get("/{id}", h.GetByID)
	r.Patch("/{id}", h.UpdateByID)
	r.Delete("/{id}", h.DeleteByID)
	return r
}

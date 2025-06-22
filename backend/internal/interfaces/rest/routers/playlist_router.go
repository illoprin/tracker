package routers

import (
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/handlers"
	"tracker-backend/internal/interfaces/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterPlaylistRoutes(
	s *services.PlaylistService,
	auth middleware.MiddlewareFunc,
) chi.Router {
	h := handlers.NewPlaylistHandler(s)

	r := chi.NewRouter()

	r.Use(auth)
	r.Post("/", h.Create)
	r.Get("/my", h.My)
	r.Get("/{id}", h.Get)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Get("/{id}/tracks", h.GetTracks)
	r.Patch("/{id}/tracks/{trackId}", h.AddTrack)
	r.Delete("/{id}/tracks/{trackId}", h.DeleteTrack)

	return r
}

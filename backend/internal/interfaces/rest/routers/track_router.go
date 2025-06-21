package routers

import (
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/handlers"
	"tracker-backend/internal/interfaces/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterTrackRoutes(s *services.TrackService, mw middleware.MiddlewareFunc) chi.Router {
	h := handlers.NewTrackHandler(s)

	r := chi.NewRouter()
	r.Use(mw)
	r.Get("/{id}/stream", h.GetFile)
	r.Get("/{id}", h.GetMetadata)
	return r
}

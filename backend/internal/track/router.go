package track

import (
	"tracker-backend/internal/auth"

	"github.com/go-chi/chi/v5"
)

func RegisterTrackRoutes(r chi.Router, s *TrackService, authMiddleware auth.MiddlewareFunc) {
	h := NewTrackHandler(s)

	r.Route("/track", func(r chi.Router) {
		r.Group(func(rm chi.Router) {
			rm.Use(authMiddleware)
			rm.Post("/", h.Create)
		})
		r.Get("/{id}/stream", h.StreamTrack)
		r.Get("/{id}", h.GetByID)
	})
}

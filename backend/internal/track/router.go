package track

import "github.com/go-chi/chi/v5"

func RegisterTrackRoutes(r chi.Router, s *TrackService) {
	h := NewTrackHandler(s)

	r.Route("/track", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}/stream", h.StreamTrack)
		r.Get("/{id}", h.GetByID)
	})
}

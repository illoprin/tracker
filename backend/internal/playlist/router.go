package playlist

import (
	"tracker-backend/internal/auth"

	"github.com/go-chi/chi/v5"
)

func RegisterPlaylistRoutes(router chi.Router, service *PlaylistService, authMiddleware auth.MiddlewareFunc) {
	h := NewPlaylistHandler(service)
	router.Route("/playlist", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", h.MyPlaylists)       // get my playlists
		r.Get("/{id}/tracks", h.Tracks) // get tracks in playlist
		r.Put("/{id}/tracks/{trackID}", h.AddTrack)
		r.Delete("/{id}/tracks/{trackID}", h.RemoveTrack)
	})
}

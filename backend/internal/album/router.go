package album

import (
	albumTracks "tracker-backend/internal/album/tracks"
	"tracker-backend/internal/auth"

	"github.com/go-chi/chi/v5"
)

func RegisterAlbumRoutes(
	router chi.Router,
	albumSvc *AlbumService,
	albumTracksSvc *albumTracks.AlbumTracksService,
	authMiddleware auth.MiddlewareFunc,
) {
	h := NewAlbumHandler(albumSvc)
	ht := albumTracks.NewAlbumTracksHandler(albumTracksSvc)

	router.Route("/album", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/{id}/tracks", ht.GetAlbumTracks)

		r.Group(func(mr chi.Router) {
			mr.Use(authMiddleware)
			mr.Post("/", h.Create)
			mr.Put("/{id}", h.Update)
			mr.Delete("/{id}", h.Delete)
		})
		r.Get("/{id}", h.GetByID)
	})
}

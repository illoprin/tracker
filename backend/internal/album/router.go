package album

import (
	"tracker-backend/internal/auth"

	"github.com/go-chi/chi/v5"
)

func RegisterAlbumRoutes(
	router chi.Router, s *AlbumService,
	authMiddleware auth.MiddlewareFunc,
) {
	h := NewAlbumHandler(s)

	router.Route("/album", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Group(func(mr chi.Router) {
			mr.Use(authMiddleware)
			mr.Post("/", h.Create)
			mr.Put("/{id}", h.Update)
			mr.Delete("/{id}", h.Delete)
		})
		r.Get("/{id}", h.GetByID)
	})
	router.Get("/artist/{id}/albums", h.GetByArtistID)
}

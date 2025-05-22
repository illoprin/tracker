package artist

import (
	"tracker-backend/internal/auth/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterArtistRoutes(r chi.Router, service *ArtistService, authMiddleware middleware.MiddlewareFunc) {
	h := NewArtistHandler(service)

	r.Route("/artist", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Post("/", h.Create)
			r.Get("/my", h.GetByUserID) // GET /artists/my
		})

		r.Get("/{id}", h.GetByID) // public access

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})

	})
}

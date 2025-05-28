package artist

import (
	artistAlbums "tracker-backend/internal/artist/albums"
	"tracker-backend/internal/auth"

	"github.com/go-chi/chi/v5"
)

func RegisterArtistRoutes(
	r chi.Router,
	service *ArtistService, artistAlbumsService *artistAlbums.ArtistAlbumsService,
	authMiddleware auth.MiddlewareFunc,
) {
	h := NewArtistHandler(service)
	ha := artistAlbums.NewArtistAlbumsHandler(artistAlbumsService)

	r.Route("/artist", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Post("/", h.Create)
			r.Get("/my", h.GetByUserID) // GET /artists/my
		})

		r.Get("/{id}", h.GetByID) // public access

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/{id}/albums", ha.GetAlbums)
			r.Put("/{id}", h.Update)
			r.Put("/{id}/avatar", h.UpdateAvatar)
			r.Delete("/{id}", h.Delete)
		})

	})
}

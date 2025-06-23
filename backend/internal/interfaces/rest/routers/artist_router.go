package routers

import (
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/handlers"
	"tracker-backend/internal/interfaces/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterArtistRoutes(
	artistSvc *services.ArtistService,
	albumSvc *services.AlbumService,
	recommendationsSvc *services.RecommendationsService,
	mw middleware.MiddlewareFunc,
) chi.Router {
	h := handlers.NewArtistHandler(artistSvc, albumSvc, recommendationsSvc)
	r := chi.NewRouter()

	r.Use(mw)
	r.Post("/", h.Create)

	r.Get("/my", h.GetMy)
	r.Get("/liked", h.GetLiked)
	r.Get("/{id}", h.GetByID)
	r.Get("/{id}/wave", h.GetWave)
	r.Post("/{id}/like", h.Like)
	r.Get("/{id}/stats", h.GetStats)
	r.Get("/{id}/albums", h.GetAlbums)
	r.Post("/{id}/album", h.PushAlbum)
	r.Get("/{id}/popular", h.GetPopularTracks)

	r.Delete("/{id}", h.DeleteByID)

	return r
}

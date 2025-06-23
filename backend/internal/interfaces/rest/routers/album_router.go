package routers

import (
	"tracker-backend/internal/domain/repository/schemas"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/handlers"
	"tracker-backend/internal/interfaces/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterAlbumRoutes(
	s *services.AlbumService,
	t *services.TrackService,
	rcm *services.RecommendationsService,
	mw middleware.MiddlewareFunc,
) chi.Router {
	h := handlers.NewAlbumHandler(s, t, rcm)
	r := chi.NewRouter()

	r.Use(mw)

	r.Group(func(rm chi.Router) {
		rm.Use(middleware.Role(schemas.RoleModerator))
		rm.Get("/moderation", h.GetUnapproved)
		r.Post("/{id}/moderate", h.Moderate)
	})

	r.Get("/liked", h.GetLiked)
	r.Get("/{id}", h.GetByID)
	r.Patch("/{id}", h.UpdateByID)
	r.Delete("/{id}", h.DeleteByID)
	r.Get("/{id}/wave", h.GetWave)
	r.Post("/{id}/like", h.Like)
	r.Post("/{id}/publish", h.Publish)
	r.Post("/{id}/tracks", h.CreateTrack)
	r.Get("/{id}/tracks", h.GetTracksByAlbumID)
	return r
}

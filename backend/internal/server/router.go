package server

import (
	"tracker-backend/internal/album"
	"tracker-backend/internal/app/dependencies"
	"tracker-backend/internal/artist"
	"tracker-backend/internal/auth/middleware"
	"tracker-backend/internal/user"

	"github.com/go-chi/chi/v5"
)

func NewAppRouter(deps *dependencies.Dependencies) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/ping", HandlePing)
	authMiddleware := middleware.Authorization(deps.UserService)

	user.RegisterUserRoutes(router, deps.UserService, authMiddleware)
	artist.RegisterArtistRoutes(router, deps.ArtistService, authMiddleware)
	album.RegisterAlbumRoutes(router, deps.AlbumService, authMiddleware)

	return router
}

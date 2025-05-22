package server

import (
	"tracker-backend/internal/app/setup"
	"tracker-backend/internal/artist"
	"tracker-backend/internal/auth/middleware"
	"tracker-backend/internal/user"

	"github.com/go-chi/chi/v5"
)

func NewAppRouter(deps *setup.Dependencies) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/ping", HandlePing)
	authMiddleware := middleware.Authorization(deps.UserService, deps.UserService.JwtSecret)

	user.RegisterUserRoutes(router, deps.UserService, authMiddleware)
	artist.RegisterArtistRoutes(router, deps.ArtistService, authMiddleware)

	return router
}

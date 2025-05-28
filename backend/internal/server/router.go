package server

import (
	"tracker-backend/internal/album"
	"tracker-backend/internal/app/dependencies"
	"tracker-backend/internal/artist"
	"tracker-backend/internal/auth/middleware"
	"tracker-backend/internal/genre"
	"tracker-backend/internal/track"
	"tracker-backend/internal/user"

	"github.com/go-chi/chi/v5"
)

func NewAppRouter(deps *dependencies.Dependencies) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/ping", HandlePing)
	authMiddleware := middleware.Authorization(deps.UserService)

	router.Get("/genres", genre.GetAllGenres)
	user.RegisterUserRoutes(router, deps.UserService, authMiddleware)
	artist.RegisterArtistRoutes(router, deps.ArtistService, deps.ArtistAlbumsService, authMiddleware)
	album.RegisterAlbumRoutes(router, deps.AlbumService, deps.AlbumTracksService, authMiddleware)
	track.RegisterTrackRoutes(router, deps.TrackService, authMiddleware)

	return router
}

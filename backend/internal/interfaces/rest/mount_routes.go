package rest

import (
	"tracker-backend/internal/infrastructure/dependencies"
	"tracker-backend/internal/interfaces/rest/handlers"
	"tracker-backend/internal/interfaces/rest/middleware"
	"tracker-backend/internal/interfaces/rest/routers"

	"github.com/go-chi/chi/v5"
)

func MountAppRoutes(r chi.Router, deps *dependencies.Dependencies) chi.Router {
	authMiddleware := middleware.Authorization(deps.AuthSvc)
	searchHandler := handlers.NewSearchHandler(deps.SearchSvc)
	r.Route("/api/v1", func(ar chi.Router) {
		ar.Get("/ping", handlers.Ping)
		ar.Get("/search", searchHandler.GlobalSearch)

		r.Mount("/genres", routers.RegisterGenreRoutes())

		// PERF
		ar.Mount("/auth", routers.RegisterAuthRoutes(deps.AuthSvc))

		// PERF
		ar.Mount("/user", routers.RegisterUserRoutes(
			deps.UserSvc,
			deps.AuthSvc,
			authMiddleware,
		))

		// PERF
		ar.Mount("/artists", routers.RegisterArtistRoutes(
			deps.ArtistSvc,
			deps.AlbumSvc,
			authMiddleware,
		))

		// PERF
		ar.Mount("/albums", routers.RegisterAlbumRoutes(
			deps.AlbumSvc,
			deps.TrackSvc,
			authMiddleware,
		))

		// PERF
		ar.Mount("/tracks", routers.RegisterTrackRoutes(
			deps.TrackSvc,
			authMiddleware,
		))

		ar.Mount("/playlists", routers.RegisterPlaylistRoutes(
			deps.PlaylistSvc,
			authMiddleware,
		))
	})
	return r
}

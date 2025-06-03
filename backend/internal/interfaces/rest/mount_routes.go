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
	r.Route("/api", func(ar chi.Router) {
		ar.Get("/ping", handlers.Ping)
		// PERF
		ar.Mount("/auth", routers.RegisterAuthRoutes(deps.AuthSvc))
		// PERF
		ar.Mount("/user", routers.RegisterUserRoutes(deps.UserSvc, deps.AuthSvc, authMiddleware))
	})
	return r
}

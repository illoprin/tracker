package routers

import (
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(s *services.AuthorizationService) chi.Router {
	h := handlers.NewAuthHandler(s)

	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/verify", h.Verify)
	r.Post("/logout", h.Logout)
	return r
}

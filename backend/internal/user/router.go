package user

import (
	"tracker-backend/internal/auth"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, service *UserService, authMiddleware auth.MiddlewareFunc) {
	h := NewUserHandler(service)

	r.Route("/user", func(r chi.Router) {
		r.Post("/", h.Register)
		r.Post("/login", h.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/me", h.Me)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})
}

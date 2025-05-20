package userRouter

import (
	"tracker-backend/internal/middleware/authorization"
	"tracker-backend/internal/user/handler"
	"tracker-backend/internal/user/service"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, service *userService.UserService) {
	h := userHandler.NewUserHandler(service)

	r.Route("/user", func(r chi.Router) {
		r.Post("/", h.Register)
		r.Post("/login", h.Login)

		r.Group(func(r chi.Router) {
			r.Use(authorization.Middleware(service))

			r.Get("/", h.Me)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})
}

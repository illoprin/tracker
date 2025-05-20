package router

import (
	"tracker-backend/internal/handler/ping"
	"tracker-backend/internal/service"
	userRouter "tracker-backend/internal/user/router"

	"github.com/go-chi/chi/v5"
)

func NewAppRouter(deps *service.Dependencies) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/ping", ping.HandlePing)

	userRouter.RegisterUserRoutes(router, deps.UserService)

	return router
}

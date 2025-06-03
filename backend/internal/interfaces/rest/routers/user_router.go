package routers

import (
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/handlers"
	"tracker-backend/internal/interfaces/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(
	userSvc *services.UserService,
	authSvc *services.AuthorizationService,
	mw middleware.MiddlewareFunc,
) chi.Router {
	h := handlers.NewUserHandler(userSvc, authSvc)

	r := chi.NewRouter()
	r.Use(middleware.Authorization(authSvc))
	r.Get("/me", h.Me)
	r.Patch("/me", h.UpdateMe)
	r.Delete("/me", h.DeleteMe)
	return r
}

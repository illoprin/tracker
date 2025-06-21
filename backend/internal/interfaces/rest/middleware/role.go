package middleware

import (
	"net/http"
	"tracker-backend/internal/interfaces/rest/utils/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/render"
)

func Role(minRole int) MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// get context keys
			ctx := r.Context()
			userRole := ctx.Value(UserRoleKey).(int)

			// check role
			if userRole < minRole {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error(service.ErrForbidden.Error()))
				return
			}

			// next
			h.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"tracker-backend/internal/domain/repository/schemas"
	"tracker-backend/internal/interfaces/rest/utils/response"
	"tracker-backend/internal/pkg/service"
	authToken "tracker-backend/internal/pkg/token"

	"github.com/go-chi/render"
)

type AuthorizationProvider interface {
	Verify(ctx context.Context, token string) (*schemas.User, *authToken.JWTClaims, bool, error)
}

func Authorization(p AuthorizationProvider) MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// configure logger
			_logger := slog.With(slog.String("func", "middleware.Authorization"))

			// get auth header
			header := r.Header.Get("Authorization")
			if header == "" {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("authorization required"))
				return
			}

			// decode token
			token := header[len("Bearer "):]
			if token == "" {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("authorization required"))
				return
			}
			_logger.Debug("authorization", slog.String("token", token))

			// validate user session
			_, claims, valid, err := p.Verify(r.Context(), token)
			if err != nil {
				if errors.Is(err, service.ErrInternal) {
					render.Status(r, http.StatusInternalServerError)
					render.JSON(w, r, response.Error(err.Error()))
					return
				}
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("authorization required"))
				return
			}

			// validate token
			if !valid {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("you need to refresh token"))
				return
			}

			// set context keys
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			req := r.WithContext(ctx)

			// next
			h.ServeHTTP(w, req)
		})
	}
}

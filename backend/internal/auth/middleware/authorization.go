package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/config"
	"tracker-backend/internal/pkg/response"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

func Authorization(userProvider auth.UserProvider) auth.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// configure logger
			logger := slog.With(slog.String("function", "middleware.Authorization"))

			// get authorization header
			tokenHeader := r.Header.Get("Authorization")
			if tokenHeader == "" {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("access denied"))
				return
			}

			tokenStr := tokenHeader[len("Bearer "):]
			claims := jwt.MapClaims{}

			// parse token
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv(config.JWTSecretEnvName)), nil
			})

			// check token validness
			if err != nil || !token.Valid {
				logger.Info("access denied", slog.String("error", err.Error()))
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("invalid token"))
				return
			}

			logger.Info("jwt payload decoded", "claims", claims)

			// get user by id
			user, err := userProvider.GetAuthDTOByID(r.Context(), claims["id"].(string), int(claims["role"].(float64)))
			if err != nil {
				logger.Info("failed to get user dto", slog.String("error", err.Error()))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("user not found"))
				return
			}

			// set user id and role to context
			ctx := context.WithValue(
				r.Context(), auth.UserIDKey, user.ID,
			)
			ctx = context.WithValue(
				ctx, auth.UserRoleKey, user.Role,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

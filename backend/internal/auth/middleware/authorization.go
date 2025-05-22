package middleware

import (
	"context"
	"net/http"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/pkg/response"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

func Authorization(userProvider auth.UserProvider, jwtSecret string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				return []byte(jwtSecret), nil
			})

			// check token validness
			if err != nil || !token.Valid {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("invalid token"))
				return
			}

			// get user by id
			user, err := userProvider.GetAuthDTOByID(r.Context(), claims["id"].(string), claims["role"].(string))
			if err != nil {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("user not found"))
				return
			}

			// set user id and role to context
			ctx := context.WithValue(
				r.Context(), auth.UserIDKey, user.ID,
			)
			ctx = context.WithValue(
				ctx, auth.UserRoleKey, claims["role"],
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

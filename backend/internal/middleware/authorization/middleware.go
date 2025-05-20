package authorization

import (
	"context"
	"net/http"
	"strings"
	"tracker-backend/internal/lib/handler/response"
	"tracker-backend/internal/middleware"
	userService "tracker-backend/internal/user/service"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

func Middleware(userService *userService.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get authorization header
			tokenHeader := r.Header.Get("Authorization")
			if tokenHeader == "" {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("access denied"))
				return
			}

			// split header by 'Bearer' and '$token'
			parts := strings.Split(tokenHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("access denied"))
				return
			}

			// get token
			tokenStr := parts[1]
			claims := jwt.MapClaims{}

			// parse token
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(userService.JwtSecret), nil
			})

			// check token validness
			if err != nil || !token.Valid {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, response.Error("access denied"))
				return
			}

			user, err := userService.GetByID(r.Context(), claims["id"].(string))
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(
				r.Context(), middleware.UserIDKey, user.ID,
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

package middleware

import "net/http"

type AuthorizationProvider interface {
	IsValidUser(id string, email string, role int) (bool, error)
}

func Authorization(p AuthorizationProvider) MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO
			// get auth header
			// decode token
			// validate user session
			// set context keys
			h.ServeHTTP(w, r)
		})
	}
}

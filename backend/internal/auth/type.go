package auth

import (
	"errors"
	"net/http"
)

var (
	ErrAccessDenied = errors.New("user does not own this resource")
)

type MiddlewareFunc func(http.Handler) http.Handler

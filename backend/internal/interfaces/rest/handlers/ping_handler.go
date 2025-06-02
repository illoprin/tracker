package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"ping": "pong",
	})
}

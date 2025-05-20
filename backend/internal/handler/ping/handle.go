package ping

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Ping string `json:"ping"`
}

func HandlePing(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Ping: "pong",
	})
}

package handlers

import (
	"net/http"
	"tracker-backend/internal/domain/dtos"

	"github.com/go-chi/render"
)

type GenreHandler struct {
}

func NewGenreHandler() *GenreHandler {
	return &GenreHandler{}
}

func (h *GenreHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string][]string{"genres": dtos.AllowedGenres})
}

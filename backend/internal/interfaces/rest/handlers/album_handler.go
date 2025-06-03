package handlers

import (
	"net/http"
	"tracker-backend/internal/domain/services"

	"github.com/go-chi/render"
)

type AlbumHandler struct {
	aSvc *services.AlbumService
}

func NewAlbumHandler(s *services.AlbumService) *AlbumHandler {
	return &AlbumHandler{
		aSvc: s,
	}
}

func (h *AlbumHandler) Create(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}
func (h *AlbumHandler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}
func (h *AlbumHandler) GetTracksByAlbumID(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}
func (h *AlbumHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}
func (h *AlbumHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}
func (h *AlbumHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}

package playlist

import (
	"errors"
	"net/http"
	"tracker-backend/internal/pkg/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type PlaylistHandler struct {
	Service   *PlaylistService
	validator *validator.Validate
}

func NewPlaylistHandler(service *PlaylistService) *PlaylistHandler {
	v := validator.New()
	return &PlaylistHandler{
		Service:   service,
		validator: v,
	}
}

func (h *PlaylistHandler) AddTrack(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "id")
	trackID := chi.URLParam(r, "trackID")

	// validate param
	if err := h.validator.Var(trackID, "required,uuid4"); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid track ID"))
		return
	}

	// execute service function
	updated, err := h.Service.PushTrackLink(r.Context(), playlistID, trackID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("playlist not found"))
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to add track"))
		return
	}

	render.JSON(w, r, updated)
}

func (h *PlaylistHandler) RemoveTrack(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "id")
	trackID := chi.URLParam(r, "trackID")

	// validate param
	if err := h.validator.Var(trackID, "required,uuid4"); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid track ID"))
		return
	}

	// execute service function
	updated, err := h.Service.RemoveTrackLink(r.Context(), playlistID, trackID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("playlist not found"))
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to remove track"))
		return
	}

	render.JSON(w, r, updated)
}

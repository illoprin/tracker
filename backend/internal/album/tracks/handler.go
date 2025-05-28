package albumTracks

import (
	"errors"
	"net/http"
	"tracker-backend/internal/album"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/pkg/response"
	"tracker-backend/internal/track"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type AlbumTracksHandler struct {
	albumTracksService *AlbumTracksService
}

func NewAlbumTracksHandler(ats *AlbumTracksService) *AlbumTracksHandler {
	return &AlbumTracksHandler{
		albumTracksService: ats,
	}
}

func (h *AlbumTracksHandler) GetAlbumTracks(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)
	userRole := ctx.Value(auth.UserRoleKey).(int)

	// url param
	albumID := chi.URLParam(r, "id")

	if albumID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to get url param"))
	}

	// execute service function
	tracks, err := h.albumTracksService.GetTracksByID(ctx, albumID, userID, userRole)
	if err != nil {
		if errors.Is(err, album.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else {
			render.Status(r, http.StatusBadRequest)
		}
		render.JSON(w, r, response.Error("failed to get url param"))
		return
	}

	tracksResponse := make([]track.TrackResponse, len(tracks))
	for i, t := range tracks {
		tracksResponse[i] = t.ToResponse()
	}

	// return tracks
	render.JSON(w, r, tracks)
}

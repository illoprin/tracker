package artistAlbums

import (
	"net/http"
	albumType "tracker-backend/internal/album/type"
	"tracker-backend/internal/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ArtistAlbumsHandler struct {
	Service *ArtistAlbumsService
}

func NewArtistAlbumsHandler(
	artistAlbumsSvc *ArtistAlbumsService,
) *ArtistAlbumsHandler {
	return &ArtistAlbumsHandler{
		Service: artistAlbumsSvc,
	}
}

func (h *ArtistAlbumsHandler) GetAlbums(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	artistID := chi.URLParam(r, "artistID")

	// execute service function
	albums, err := h.Service.GetByArtistID(ctx, artistID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get albums"))
		return
	}

	// execution result to response
	responseAlbums := make([]albumType.AlbumResponse, len(albums))
	for i, a := range albums {
		responseAlbums[i] = a.ToResponse()
	}

	// send response
	render.JSON(w, r, responseAlbums)
}

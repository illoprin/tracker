package artistAlbums

import (
	"log/slog"
	"net/http"
	albumType "tracker-backend/internal/album/type"
	"tracker-backend/internal/auth"
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
	artistID := chi.URLParam(r, "id")
	userID := ctx.Value(auth.UserIDKey).(string)
	userRole := ctx.Value(auth.UserRoleKey).(int)

	slog.Info("artist albums requested", slog.Group("info",
		slog.String("artistID", artistID),
		slog.String("userID", userID),
		slog.Int("userRole", userRole),
	))

	// execute service function
	albums, err := h.Service.GetByArtistID(ctx, artistID, userID, userRole)

	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
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

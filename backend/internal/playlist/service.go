package playlist

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PlaylistService struct {
	Col *mongo.Collection
}

var (
	ErrNotFound = errors.New("playlist not found")
)

func NewPlaylistService(
	playlistCol *mongo.Collection,
) *PlaylistService {
	return &PlaylistService{
		Col: playlistCol,
	}
}

func (s *PlaylistService) Create(
	ctx context.Context, req PlaylistCreateRequest,
) (*Playlist, error) {
	// configure logger
	logger := slog.With(slog.String("function", "playlist.PlaylistService.Create"))

	playlist := &Playlist{
		ID:        uuid.NewString(),
		Name:      req.Name,
		UserID:    req.UserID,
		IsDefault: req.IsDefault,
		IsPublic:  req.IsPublic,
	}

	_, err := s.Col.InsertOne(ctx, playlist)

	if err != nil {
		return nil, errors.New("failed to insert")
	}

	logger.Info("playlist created",
		slog.Group("info",
			slog.String("userID", req.UserID),
			slog.Bool("isDefault", req.IsDefault),
			slog.String("id", playlist.ID),
		),
	)

	return playlist, nil
}

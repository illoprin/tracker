package albumTracks

import (
	"context"
	"errors"
	"log/slog"
	albumType "tracker-backend/internal/album/type"
	"tracker-backend/internal/auth/ownership"
	"tracker-backend/internal/pkg/service"
	"tracker-backend/internal/track"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AlbumTracksService struct {
	TracksCol        *mongo.Collection
	AlbumsCol        *mongo.Collection
	ownershipService *ownership.OwnershipService
}

func NewAlbumTracksService(
	tracksCol,
	albumsCol *mongo.Collection,
	ownershipService *ownership.OwnershipService,
) *AlbumTracksService {

	return &AlbumTracksService{
		TracksCol:        tracksCol,
		AlbumsCol:        albumsCol,
		ownershipService: ownershipService,
	}
}

func (s *AlbumTracksService) GetTracksByID(
	ctx context.Context, albumID, userID string, userRole int,
) ([]track.Track, error) {
	// configure logger
	logger := slog.With(slog.String("function", "albumTracks.AlbumTracksService.GetTracksByID"))

	// check album existence
	res := s.AlbumsCol.FindOne(ctx, bson.M{"id": albumID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		return nil, errors.New("failed to check album existence")
	}

	// if user is not owner of album and its status 'OnModeration' or it is hidden
	if isOwn, err := s.ownershipService.IsAlbumOwner(ctx, userID, albumID); !isOwn {
		if err != nil {
			return nil, errors.New("failed to check owner")
		}
		var albumDecoded albumType.Album
		res.Decode(&albumDecoded)
		if albumDecoded.Status == albumType.StatusOnModeration || albumDecoded.IsHidden {
			return nil, service.ErrAccessDenied
		}
	}

	// find tracks
	cur, err := s.TracksCol.Find(ctx, bson.M{"album": albumID})
	if err != nil {
		logger.Warn("failed to find tracks", slog.String("error", err.Error()))
		return nil, errors.New("failed to get tracks")
	}
	defer cur.Close(ctx)

	// decode cursor
	tracks := make([]track.Track, cur.RemainingBatchLength())
	if err = cur.All(ctx, &tracks); err != nil {
		logger.Warn("failed to decode cursor by found tracks", slog.String("error", err.Error()))
		return nil, errors.New("failed to decode tracks cursor")
	}

	// return tracks
	return tracks, nil
}

func (s *AlbumTracksService) IsAnyTracksInAlbum(
	ctx context.Context, albumID string,
) (bool, error) {
	// configure logger
	logger := slog.With(slog.String("function", "albumTracks.AlbumTracksService.IsAnyTracksInAlbum"))

	// count tracks by album id
	count, err := s.TracksCol.CountDocuments(ctx, bson.M{"album": albumID})
	if err != nil {
		logger.Warn("failed to execute find function", slog.String("error", err.Error()))
		return false, errors.New("failed to count tracks in album")
	}

	// return result
	return count > 0, nil
}

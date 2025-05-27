package albumTracks

import (
	"context"
	"errors"
	"log/slog"
	"tracker-backend/internal/album"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/auth/ownership"
	"tracker-backend/internal/track"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AlbumTracksService struct {
	TracksCol        *mongo.Collection
	AlbumsCol        *mongo.Collection
	ownershipService *ownership.OwnershipService
}

func NewAlbumTracksService(tracksCol, albumsCol *mongo.Collection, ownershipService *ownership.OwnershipService) *AlbumTracksService {

	return &AlbumTracksService{
		TracksCol:        tracksCol,
		AlbumsCol:        albumsCol,
		ownershipService: ownershipService,
	}
}

func (s *AlbumTracksService) GetTracksByID(
	ctx context.Context, albumID, userID string,
) ([]track.Track, error) {
	// configure logger
	logger := slog.With(slog.String("function", "albumTracks.AlbumTracksService.GetTracksByID"))

	// check album existence
	res := s.AlbumsCol.FindOne(ctx, bson.M{"id": albumID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, album.ErrNotFound
		}
		return nil, errors.New("failed to check album existence")
	}

	// decode albumEnt data
	var albumEnt album.Album // album entry
	if err := res.Decode(&albumEnt); err != nil {
		return nil, errors.New("failed to decode album")
	}

	// check ownership
	isOwn, err := s.ownershipService.IsAlbumOwner(ctx, userID, albumID)
	if err != nil {
		return nil, errors.New("failed to check album ownership")
	}

	// check album status
	// if user who made the request is not owner -> return error
	// if user is owner -> return tracks
	if !isOwn && albumEnt.Status != album.StatusPublic {
		return nil, auth.ErrAccessDenied
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

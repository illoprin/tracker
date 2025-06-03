package services

import (
	"context"
	"errors"
	"log/slog"
	"tracker-backend/internal/pkg/logger"
	"tracker-backend/internal/pkg/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OwnershipChecker interface {
	IsAlbumOwner(ctx context.Context, userId string, albumId string) (bool, error)
	IsTrackOwner(ctx context.Context, userId string, trackId string) (bool, error)
	IsArtistOwner(ctx context.Context, userId string, artistId string) (bool, error)
}

type OwnershipService struct {
	albumsCol *mongo.Collection
	artistCol *mongo.Collection
	tracksCol *mongo.Collection
}

func NewOwnershipService(
	albumsCol *mongo.Collection,
	artistCol *mongo.Collection,
	tracksCol *mongo.Collection,
) *OwnershipService {
	return &OwnershipService{
		albumsCol: albumsCol,
		artistCol: artistCol,
		tracksCol: tracksCol,
	}
}

func (s *OwnershipService) isOwner(
	ctx context.Context,
	col *mongo.Collection,
	userId string,
	entityId string,
	_logger *slog.Logger,
) (bool, error) {
	count, err := col.CountDocuments(ctx, bson.M{
		"ownerId": userId,
		"id":      entityId,
	})
	if err != nil {
		_logger.Error("failed to count", logger.ErrorAttr(err))
		return false, err
	}
	return count == 1, nil
}

func (s *OwnershipService) IsTrackOwner(
	ctx context.Context,
	userId string,
	trackId string,
) (bool, error) {
	// configure logger
	_logger := slog.With(
		"func", "services.OwnershipService.IsTrackOwner",
		"trackId", trackId,
		"userId", userId,
	)

	// find track and get album id
	var track struct {
		AlbumId string `bson:"albumId"`
	}
	err := s.tracksCol.FindOne(ctx, bson.M{"id": trackId}).Decode(&track)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, service.ErrNotFound
		}
		_logger.Error("failed to find", logger.ErrorAttr(err))
		return false, err
	}
	// count albums with id = albumId ownerId = userId
	return s.isOwner(ctx, s.albumsCol, userId, track.AlbumId, _logger)
}

func (s *OwnershipService) IsArtistOwner(
	ctx context.Context,
	userId string,
	artistId string,
) (bool, error) {
	// configure logger
	_logger := slog.With(
		"func", "services.OwnershipService.IsArtistOwner",
		"artistId", artistId,
		"userId", userId,
	)
	return s.isOwner(ctx, s.artistCol, userId, artistId, _logger)
}

func (s *OwnershipService) IsAlbumOwner(
	ctx context.Context,
	userId string,
	albumId string,
) (bool, error) {
	// configure logger
	_logger := slog.With(
		"func", "services.OwnershipService.IsArtistOwner",
		"albumId", albumId,
		"userId", userId,
	)
	return s.isOwner(ctx, s.albumsCol, userId, albumId, _logger)
}

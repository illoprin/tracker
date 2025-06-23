package services

import (
	"context"
	"errors"
	"log/slog"
	"tracker-backend/internal/domain/dtos"
	"tracker-backend/internal/domain/repository/schemas"
	"tracker-backend/internal/pkg/logger"
	"tracker-backend/internal/pkg/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RecommendationsService struct {
	genreSvc  *GenreService
	tracksCol *mongo.Collection
}

func NewRecommendationsService(
	genreSvc *GenreService,
	tracksCol *mongo.Collection,
) *RecommendationsService {
	return &RecommendationsService{
		genreSvc:  genreSvc,
		tracksCol: tracksCol,
	}
}

// GetForResource returns tracks recommendations based on requested resource
func (svc *RecommendationsService) GetForResource(
	ctx context.Context,
	resourceType string,
	resourceID string,
	limit int,
	page int,
	excludeTrackIDs []string,
) ([]dtos.TrackResponse, error) {
	// configure logger
	_logger := slog.With(
		slog.String("func", "services.RecommendationsService.GetForResource"),
		slog.String("resourceType", resourceType),
		slog.String("resourceID", resourceID),
		slog.Int("limit", limit),
		slog.Int("page", page),
	)

	// TODO: exclude tracks from resource

	// get genre model
	genreModel, err := svc.genreSvc.GetPopularGenres(ctx, resourceType, resourceID)
	if err != nil {
		_logger.Error("Failed to get genre model", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// get recommendations
	return svc.GetFromModel(ctx, genreModel, limit, page, excludeTrackIDs)
}

func (svc *RecommendationsService) GetForTrack(
	ctx context.Context,
	trackID string,
	limit int,
	page int,
) ([]dtos.TrackResponse, error) {

	// configure logger
	_logger := slog.With(
		slog.String("func", "services.RecommendationsService.GetForResource"),
	)

	// get track genres
	var track struct {
		Genres []string `bson:"genres"`
	}
	err := svc.tracksCol.FindOne(ctx, bson.M{"id": trackID}).Decode(&track)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		_logger.Error("failed to update album status", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return svc.GetFromModel(ctx, track.Genres, limit, page, []string{trackID})
}

// GetFromModel returns tracks based on genres
func (svc *RecommendationsService) GetFromModel(
	ctx context.Context,
	genreModel []string,
	limit int,
	page int,
	excludeTrackIDs []string,
) ([]dtos.TrackResponse, error) {
	_logger := slog.With(
		slog.String("func", "services.RecommendationsService.GetFromModel"),
		slog.Int("limit", limit),
		slog.Int("page", page),
	)

	if len(genreModel) == 0 {
		_logger.Warn("Empty genre model provided")
		return []dtos.TrackResponse{}, nil
	}

	// define pagination
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64((page - 1) * limit))

	// create array of ids
	filter := bson.M{
		"genres": bson.M{"$in": genreModel},
		"id":     bson.M{"$nin": excludeTrackIDs},
	}

	cursor, err := svc.tracksCol.Find(ctx, filter, opts)
	if err != nil {
		_logger.Error("failed to find recommendations", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	defer cursor.Close(ctx)

	var tracks []struct {
		ID string `bson:"id"`
	}
	if err := cursor.All(ctx, &tracks); err != nil {
		_logger.Error("failed to decode recommendations", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	_logger.Debug("successfully generated recommendations",
		slog.Int("foundTracks", len(tracks)),
		slog.Any("usedGenres", genreModel),
	)

	// make array of ids
	trackIds := make([]string, 0, len(tracks))
	for _, t := range tracks {
		trackIds = append(trackIds, t.ID)
	}

	// create aggregation pipeline
	match := bson.M{"id": bson.M{"$in": trackIds}}
	matchAlbum := bson.M{"album.isPublic": true, "album.isApproved": true}
	pipeline := schemas.GetTracksDetailsPipeline(match, matchAlbum)

	// execute aggregation
	cursor, err = svc.tracksCol.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("failed to aggregate", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// decode cursor
	var tracksWithStats []dtos.TrackResponse
	if err := cursor.All(ctx, &tracksWithStats); err != nil {
		_logger.Error("failed to aggregate", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return tracksWithStats, nil
}

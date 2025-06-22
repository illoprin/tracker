package services

import (
	"context"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"time"
	"tracker-backend/internal/config"
	"tracker-backend/internal/domain/dtos"
	"tracker-backend/internal/domain/repository/schemas"
	"tracker-backend/internal/infrastructure/storage"
	"tracker-backend/internal/pkg/logger"
	"tracker-backend/internal/pkg/service"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PlaylistService struct {
	playlistsCol *mongo.Collection
	tracksCol    *mongo.Collection
	usersCol     *mongo.Collection
}

func NewPlaylistService(
	playlistCol *mongo.Collection,
	tracksCol *mongo.Collection,
	usersCol *mongo.Collection,
) *PlaylistService {
	return &PlaylistService{
		playlistsCol: playlistCol,
		tracksCol:    tracksCol,
		usersCol:     usersCol,
	}
}

func (svc *PlaylistService) Create(
	ctx context.Context,
	userId string,
	req dtos.PlaylistCreateRequest,
	cover *multipart.File,
	coverHeader *multipart.FileHeader,
	hasCover bool,
) (*dtos.PlaylistResponse, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.Create"))

	// check playlists with similar name for this user
	count, err := svc.playlistsCol.CountDocuments(
		ctx, bson.M{"name": req.Name, "ownerId": userId},
	)
	if err != nil {
		_logger.Error("failed to count similar documents", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	if count > 0 {
		return nil, service.ErrExists
	}

	// upload cover
	coverPath := ""
	if hasCover {
		uploadDir := path.Join(os.Getenv(config.StaticDirEnvName), config.PlaylistsCoversDir)
		coverPath, err = storage.UploadFile(coverHeader, cover, uploadDir)
		if err != nil {
			return nil, err
		}
	}

	// define document
	p := schemas.Playlist{
		ID:          uuid.NewString(),
		OwnerID:     userId,
		Name:        req.Name,
		Description: req.Description,
		Cover:       coverPath,
		IsDefault:   req.IsDefault,
		IsPublic:    req.IsPublic,
		TrackIDs:    []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// insert document
	_, err = svc.playlistsCol.InsertOne(ctx, p)
	if err != nil {
		_logger.Error("failed to insert document", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// return dto
	res, err := svc.GetMetadata(ctx, userId, p.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *PlaylistService) GetMetadata(
	ctx context.Context,
	userId string,
	playlistId string,
) (*dtos.PlaylistResponse, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.GetMetadata"))

	// prepare aggregation pipeline
	match := bson.M{"id": playlistId}
	pipeline := schemas.GetPlaylistsAggregationPipeline(match)
	cursor, err := svc.playlistsCol.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("failed to aggregate", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	defer cursor.Close(ctx)

	// decode cursor
	var p dtos.PlaylistResponse
	if cursor.Next(ctx) {
		if err := cursor.Decode(&p); err != nil {
			_logger.Error("failed to decode playlist", logger.ErrorAttr(err))
			return nil, service.ErrInternal
		}
	} else {
		return nil, service.ErrNotFound
	}

	// check ownership and isPublic flag
	if p.IsPublic == false && userId != p.OwnerInfo.ID {
		return nil, service.ErrForbidden
	}

	// TODO: set cover by last track album cover

	return &p, nil
}

func (svc *PlaylistService) Update(
	ctx context.Context,
	userId string,
	req dtos.PlaylistUpdateRequest,
) (*dtos.PlaylistResponse, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.Update"))
	_ = _logger

	return nil, nil
}

func (svc *PlaylistService) Delete(
	ctx context.Context,
	userId string,
	playlistId string,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.Delete"))
	_ = _logger
	return nil

}

func (svc *PlaylistService) GetTracks(
	ctx context.Context,
	userId string,
	playlistId string,
) ([]dtos.TrackResponse, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.GetTracks"))
	_ = _logger
	return nil, nil

}

func (svc *PlaylistService) AddTrack(
	ctx context.Context,
	userId string,
	playlistId string,
	trackId string,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.AddTrack"))
	_ = _logger
	return nil

}

func (svc *PlaylistService) RemoveTrack(
	ctx context.Context,
	userId string,
	playlistId string,
	trackId string,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.RemoveTrack"))
	_ = _logger
	return nil

}

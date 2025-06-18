package services

import (
	"context"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
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

type TrackService struct {
	tc *mongo.Collection
	ac *mongo.Collection
	oc OwnershipChecker
}

func NewTrackService(
	tracksCol *mongo.Collection,
	albumsCol *mongo.Collection,
	ownershipChecker OwnershipChecker,
) *TrackService {
	return &TrackService{
		tc: tracksCol,
		ac: albumsCol,
		oc: ownershipChecker,
	}
}

func (svc *TrackService) PushTrack(
	ctx context.Context,
	userId string,
	req dtos.TrackCreateRequest,
	file *multipart.File,
	fileHeader *multipart.FileHeader,
) (*schemas.Track, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.TrackService.PushTrack"))

	// check album owner
	if isOwn, err := svc.oc.IsAlbumOwner(ctx, userId, req.AlbumID); !isOwn {
		if err != nil {
			_logger.Error("failed to check ownership", logger.ErrorAttr(err))
			return nil, service.ErrInternal
		}
		return nil, service.ErrForbidden
	}

	// check track file
	if err := storage.ValidateFile(fileHeader, storage.AllowedAudioExtensions); err != nil {
		return nil, err
	}

	// write track file and save it's name
	filePath := path.Join(os.Getenv(config.StaticDirEnvName), config.AudioDir)
	fullPath, err := storage.UploadFile(fileHeader, file, filePath)
	if err != nil {
		return nil, err
	}
	fileName := filepath.Base(fullPath)

	// check tracks with similar names in album
	count, err := svc.tc.CountDocuments(ctx, bson.M{"name": req.Name, "albumId": req.AlbumID})
	if err != nil {
		_logger.Error("failed to count documents", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	// track with similar name exists
	if count > 0 {
		return nil, service.ErrExists
	}

	// create document
	track := schemas.Track{
		ID:        uuid.NewString(),
		Name:      req.Name,
		OwnerID:   userId,
		AlbumID:   req.AlbumID,
		Genres:    req.Genres,
		Duration:  req.Duration,
		AudioFile: fileName,
		CreatedAt: time.Now(),
	}

	// push document
	_, err = svc.tc.InsertOne(ctx, track)
	if err != nil {
		_logger.Error("failed to insert", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// return document
	return &track, nil
}

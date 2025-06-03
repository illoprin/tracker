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

type AlbumFlusher interface {
	FlushAlbumData(ctx context.Context, id string) error
}

type AlbumService struct {
	albumsCol *mongo.Collection
	tracksCol *mongo.Collection
	af        AlbumFlusher
	oc        OwnershipChecker
}

func NewAlbumService(
	albumsCol *mongo.Collection,
	tracksCol *mongo.Collection,
	albumFlusher AlbumFlusher,
	ownershipChecker OwnershipChecker,
) *AlbumService {
	return &AlbumService{
		albumsCol: albumsCol,
		tracksCol: tracksCol,
		af:        albumFlusher,
		oc:        ownershipChecker,
	}
}

func (svc *AlbumService) Create(
	ctx context.Context,
	userId string,
	req dtos.AlbumCreateRequest,
	cover multipart.File,
	fileHeader *multipart.FileHeader,
) (*schemas.Album, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.Create"))

	// check artist ownership
	isOwn, err := svc.oc.IsArtistOwner(ctx, userId, req.ArtistID)
	if err != nil {
		return nil, service.ErrInternal
	}
	if !isOwn {
		return nil, service.ErrForbidden
	}

	// find similar documents
	filter := bson.M{"$and": bson.A{
		bson.M{"ownerId": userId},
		bson.M{"name": req.Name},
	}}
	count, err := svc.albumsCol.CountDocuments(ctx, filter)
	if err != nil {
		_logger.Error("failed to count documents", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	if count > 0 {
		return nil, service.ErrExists
	}

	// upload file
	uploadDir := path.Join(os.Getenv(config.StaticDirEnvName), config.CoversDir)
	coverPath, _ := storage.UploadFile(fileHeader, &cover, uploadDir)

	// create schema
	a := schemas.Album{
		ID:         uuid.NewString(),
		ArtistID:   req.ArtistID,
		OwnerID:    userId,
		Name:       req.Name,
		Year:       req.Year,
		Cover:      coverPath,
		Type:       req.Type,
		IsPublic:   false,
		IsApproved: true,
		Moderation: schemas.AlbumModeration{
			Status:  schemas.AlbumStatusApproved,
			Comment: "auto approve",
		},
		CreatedAt: time.Now(),
	}

	_, err = svc.albumsCol.InsertOne(ctx, a)
	if err != nil {
		_logger.Error("failed to insert", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return &a, err
}

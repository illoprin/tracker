package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	coverHeader *multipart.FileHeader,
	hasCover bool,
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
	coverPath := ""
	if !hasCover {
		randomCover := rand.Intn(2)
		randomCoverFile := fmt.Sprintf("cover_default_%d.jpg", randomCover)
		coverPath = filepath.Join(os.Getenv(config.StaticDirEnvName), config.CoversDir, randomCoverFile)
	} else {
		uploadDir := path.Join(os.Getenv(config.StaticDirEnvName), config.CoversDir)
		coverPath, _ = storage.UploadFile(coverHeader, &cover, uploadDir)
	}

	// create schema
	// WARN: shitcode
	yearCreated, _ := strconv.Atoi(req.Year)
	a := schemas.Album{
		ID:         uuid.NewString(),
		ArtistID:   req.ArtistID,
		OwnerID:    userId,
		Name:       req.Name,
		Year:       yearCreated,
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

func (svc *AlbumService) DeleteByID(ctx context.Context, userId, albumId string) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.DeleteByID"))

	// check album owner
	if isOwn, err := svc.oc.IsAlbumOwner(ctx, userId, albumId); !isOwn {
		if err != nil {
			_logger.Error("failed to check ownership", logger.ErrorAttr(err))
			return service.ErrInternal
		}
		return service.ErrForbidden
	}

	// flush data
	err := svc.af.FlushAlbumData(ctx, albumId)
	if err != nil {
		if !errors.Is(err, service.ErrNotFound) {
			return service.ErrInternal
		}
		return err
	}

	// delete document
	_, err = svc.albumsCol.DeleteOne(ctx, bson.M{"id": albumId})
	if err != nil {
		if !errors.Is(err, service.ErrNotFound) {
			return service.ErrInternal
		}
		return err
	}

	return nil
}

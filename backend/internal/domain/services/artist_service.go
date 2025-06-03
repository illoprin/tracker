package services

import (
	"context"
	"errors"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"time"
	"tracker-backend/internal/config"
	"tracker-backend/internal/domain/repository/schemas"
	"tracker-backend/internal/infrastructure/storage"
	"tracker-backend/internal/pkg/logger"
	"tracker-backend/internal/pkg/service"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ArtistFlusher interface {
	FlushArtistData(ctx context.Context, id string) error
}

type ArtistService struct {
	artistCol *mongo.Collection
	albumsCol *mongo.Collection
	af        ArtistFlusher
	oc        OwnershipChecker
}

func NewArtistService(
	artistCol *mongo.Collection,
	albumsCol *mongo.Collection,
	artistFlusher ArtistFlusher,
	ownershipChecker OwnershipChecker,
) *ArtistService {
	return &ArtistService{
		artistCol: artistCol,
		albumsCol: albumsCol,
		af:        artistFlusher,
		oc:        ownershipChecker,
	}
}

func (svc *ArtistService) Create(
	ctx context.Context,
	userId string,
	name string,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (*schemas.Artist, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.ArtistService.Create"))

	// count similar documents
	filter := bson.M{"$and": bson.A{bson.M{"name": name}, bson.M{"ownerId": userId}}}
	count, err := svc.artistCol.CountDocuments(ctx, filter)
	if err != nil {
		_logger.Error("failed to count documents", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	if count > 0 {
		return nil, service.ErrExists
	}

	// upload file
	uploadDir := path.Join(os.Getenv(config.StaticDirEnvName), config.AvatarsDir)
	avatarPath, _ := storage.UploadFile(fileHeader, &file, uploadDir)

	// create schema
	artist := schemas.Artist{
		ID:        uuid.NewString(),
		OwnerID:   userId,
		Name:      name,
		Avatar:    avatarPath,
		CreatedAt: time.Now(),
	}

	// insert into collection
	_, err = svc.artistCol.InsertOne(ctx, artist)
	if err != nil {
		_logger.Error("failed to insert", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return &artist, err
}

func (svc *ArtistService) GetByUserID(
	ctx context.Context,
	userId string,
) ([]schemas.Artist, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.ArtistService.GetByUserID"))

	// find all
	cur, err := svc.artistCol.Find(ctx, bson.M{"ownerId": userId})
	if err != nil {
		_logger.Error("failed to find", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	if cur.RemainingBatchLength() < 1 {
		return []schemas.Artist{}, nil
	}

	// decode result
	var all []schemas.Artist = []schemas.Artist{}
	if err := cur.All(ctx, &all); err != nil {
		_logger.Error("failed to decode cursor", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return all, nil
}

func (svc *ArtistService) GetByID(
	ctx context.Context,
	id string,
) (*schemas.Artist, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.ArtistService.GetByID"))

	var a schemas.Artist
	err := svc.artistCol.FindOne(ctx, bson.M{"id": id}).Decode(&a)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		_logger.Error("failed to get artist", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return &a, nil
}

func (svc *ArtistService) GetAlbums(
	ctx context.Context,
	userId string,
	id string,
) ([]schemas.Album, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.ArtistService.GetAlbums"))

	// check ownership
	isOwn, err := svc.oc.IsArtistOwner(ctx, userId, id)
	if err != nil {
		return nil, service.ErrInternal
	}

	// if user is owner -> return all albums
	// else return public albums only
	filter := bson.M{}
	if !isOwn {
		filter = bson.M{"artistId": id, "isApproved": true, "isPublic": true}
	} else {
		filter = bson.M{"artistId": id}
	}

	// find albums
	cur, err := svc.albumsCol.Find(ctx, filter)
	if err != nil {
		_logger.Error("failed to find albums", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	if cur.RemainingBatchLength() < 1 {
		return []schemas.Album{}, nil
	}

	// decode cursor
	var all []schemas.Album
	if err := cur.Decode(&all); err != nil {
		_logger.Error("failed to decode cursor", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return all, nil
}

func (svc *ArtistService) DeleteByID(
	ctx context.Context,
	userId string,
	id string,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.ArtistService.GetAlbums"))

	// flush related data
	if err := svc.af.FlushArtistData(ctx, id); err != nil {
		return service.ErrInternal
	}

	// delete artist document
	res, err := svc.artistCol.DeleteOne(ctx, bson.M{"ownerId": userId, "id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return service.ErrNotFound
		}
		_logger.Error("failed to update user", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	if res.DeletedCount == 0 {
		return service.ErrNotFound
	}

	return nil
}

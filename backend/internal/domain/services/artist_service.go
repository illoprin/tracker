package services

import (
	"context"
	"errors"
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
	avatarFile multipart.File,
	avatarFileHeader *multipart.FileHeader,
	hasAvatar bool,
	bannerFile multipart.File,
	bannerFileHeader *multipart.FileHeader,
	hasBanner bool,
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

	// upload files
	avatarUploadDir := path.Join(os.Getenv(config.StaticDirEnvName), config.AvatarsDir)
	bannerUploadDir := path.Join(os.Getenv(config.StaticDirEnvName), config.BannersDir)
	avatarPath := filepath.Join(avatarUploadDir, "avatar_default.jpg")
	if hasAvatar {
		avatarPath, _ = storage.UploadFile(avatarFileHeader, &avatarFile, avatarUploadDir)
	}
	bannerPath := filepath.Join(bannerUploadDir, "banner_default.jpg")
	if hasBanner {
		bannerPath, _ = storage.UploadFile(bannerFileHeader, &bannerFile, bannerUploadDir)
	}

	// create schema
	artist := schemas.Artist{
		ID:        uuid.NewString(),
		OwnerID:   userId,
		Name:      name,
		Avatar:    avatarPath,
		Banner:    bannerPath,
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
) ([]dtos.AlbumWithStats, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.ArtistService.GetAlbums"))

	// check ownership
	isOwn, err := svc.oc.IsArtistOwner(ctx, userId, id)
	if err != nil {
		return nil, service.ErrInternal
	}

	// if user is owner -> return all albums
	// else return public albums only
	var match bson.M = make(bson.M, 3)
	if !isOwn {
		match["artistId"] = id
		match["isApproved"] = true
		match["isPublic"] = true
	} else {
		match["artistId"] = id
	}

	// prepare aggregation pipeline
	pipeline := schemas.GetAlbumsAggregationPipeline(match)

	// sort albums by creating year
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"year": -1}})

	// execute aggregation
	cursor, err := svc.albumsCol.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("failed to aggregate album stats", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	defer cursor.Close(ctx)

	// get next cursor
	if !cursor.Next(ctx) {
		return []dtos.AlbumWithStats{}, nil
	}

	// decode result
	var result []struct {
		schemas.Album `bson:",inline"`
		TrackCount    int    `bson:"trackCount" json:"trackCount"`
		TotalDuration int    `bson:"totalDuration" json:"totalDuration"`
		ArtistName    string `bson:"artistName" json:"artistName"`
		ArtistAvatar  string `bson:"artistAvatar" json:"artistAvatar"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		_logger.Error("failed to decode album stats", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	albumsWithStats := make([]dtos.AlbumWithStats, 0, len(result))
	for _, a := range result {
		// create response struct
		aws := dtos.AlbumWithStats{
			Album:         a.Album,
			TrackCount:    a.TrackCount,
			TotalDuration: a.TotalDuration,
			ArtistName:    a.ArtistName,
			ArtistAvatar:  a.ArtistAvatar,
		}
		albumsWithStats = append(albumsWithStats, aws)
	}

	return albumsWithStats, nil
}

func (svc *ArtistService) DeleteByID(
	ctx context.Context,
	userId string,
	artistId string,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.ArtistService.GetAlbums"))

	// check artist owner
	if isOwn, err := svc.oc.IsArtistOwner(ctx, userId, artistId); !isOwn {
		if err != nil {
			_logger.Error("failed to check ownership", logger.ErrorAttr(err))
			return service.ErrInternal
		}
		return service.ErrForbidden
	}

	// flush data
	err := svc.af.FlushArtistData(ctx, artistId)
	if err != nil {
		if !errors.Is(err, service.ErrNotFound) {
			return service.ErrInternal
		}
		return err
	}

	// delete document
	_, err = svc.artistCol.DeleteOne(ctx, bson.M{"id": artistId})
	if err != nil {
		_logger.Error("failed to delete document", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	return nil
}

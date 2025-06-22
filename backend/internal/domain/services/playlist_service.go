package services

import (
	"context"
	"errors"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func (svc *PlaylistService) CreateDefault(
	ctx context.Context,
	userId string,
) (string, error) {

	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.Create"))

	coverPath := filepath.Join(
		os.Getenv(config.StaticDirEnvName),
		config.PlaylistsCoversDir,
		"liked_playlist_default.jpg",
	)

	// define document
	p := schemas.Playlist{
		ID:          uuid.NewString(),
		OwnerID:     userId,
		Name:        "Мой выбор",
		Description: "Здесь ваши любимые песни",
		Cover:       coverPath,
		IsDefault:   true,
		IsPublic:    false,
		TrackIDs:    []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// insert document
	_, err := svc.playlistsCol.InsertOne(ctx, p)
	if err != nil {
		_logger.Error("failed to insert document", logger.ErrorAttr(err))
		return "", service.ErrInternal
	}

	return p.ID, nil
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
	if !p.IsPublic && userId != p.OwnerInfo.ID {
		return nil, service.ErrNotFound
	}

	// TODO: set cover by last track album cover

	return &p, nil
}

func (svc *PlaylistService) Update(
	ctx context.Context,
	userId string,
	playlistId string,
	req dtos.PlaylistUpdateRequest,
	cover *multipart.File,
	coverHeader *multipart.FileHeader,
	hasCover bool,
) (*dtos.PlaylistResponse, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.Update"))

	// get playlist
	filter := bson.M{"id": playlistId, "ownerId": userId, "isDefault": false}
	var p struct {
		Cover string `bson:"cover"`
	}
	err := svc.playlistsCol.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		_logger.Error("failed to get playlist", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// define updates
	updates := bson.M{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsPublic != "" {
		updates["isPublic"] = req.IsPublic == "1"
	}

	// update cover if needed
	if hasCover {
		// remove old file
		if p.Cover != "" && !strings.Contains(p.Cover, "default") {
			oldCoverFile := filepath.Join(
				os.Getenv(config.StaticDirEnvName),
				config.PlaylistsCoversDir,
				filepath.Base(p.Cover),
			)
			if err := os.Remove(oldCoverFile); err != nil {
				_logger.Error("failed to delete old cover", logger.ErrorAttr(err))
				return nil, service.ErrInternal
			}
		}

		// upload new one
		uploadDir := path.Join(os.Getenv(config.StaticDirEnvName), config.PlaylistsCoversDir)
		filePath, err := storage.UploadFile(coverHeader, cover, uploadDir)
		if err != nil {
			_logger.Error("failed to upload new cover", logger.ErrorAttr(err))
			return nil, service.ErrInternal
		}
		updates["cover"] = filePath
	}
	updates["updatedAt"] = time.Now()

	// apply updates
	_, err = svc.playlistsCol.UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		_logger.Error("failed to update document", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// return dto
	res, err := svc.GetMetadata(ctx, userId, playlistId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *PlaylistService) Delete(
	ctx context.Context,
	userId string,
	playlistId string,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.Delete"))

	// check cover
	filter := bson.M{"id": playlistId, "ownerId": userId, "isDefault": false}
	var p struct {
		Cover string `bson:"cover"`
	}
	err := svc.playlistsCol.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return service.ErrNotFound
		}
		_logger.Error("failed to get playlist", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	// delete cover if exists and not default
	if p.Cover != "" && !strings.Contains(p.Cover, "default") {
		oldCoverFile := filepath.Join(
			os.Getenv(config.StaticDirEnvName),
			config.PlaylistsCoversDir,
			filepath.Base(p.Cover),
		)
		if err := os.Remove(oldCoverFile); err != nil {
			_logger.Error("failed to delete old cover", logger.ErrorAttr(err))
			return service.ErrInternal
		}
	}

	// delete document
	_, err = svc.playlistsCol.DeleteOne(ctx, filter)
	if err != nil {
		_logger.Error("failed to delete document", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	return nil
}

func (svc *PlaylistService) GetTracks(
	ctx context.Context,
	userId string,
	playlistId string,
	limit int,
) ([]dtos.TrackResponse, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.GetTracks"))

	// get track ids
	var p struct {
		TrackIds []string `bson:"trackIds"`
		IsPublic bool     `bson:"isPublic"`
		OwnerID  string   `bson:"ownerId"`
	}
	err := svc.playlistsCol.FindOne(ctx, bson.M{"id": playlistId}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		_logger.Error("failed to get playlist", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// check access to album
	if userId != p.OwnerID && !p.IsPublic {
		return nil, service.ErrNotFound
	}

	// prepare aggregation pipeline
	match := bson.M{"id": bson.M{"$in": p.TrackIds}}
	pipeline := schemas.GetTracksDetailsPipeline(match, bson.M{})
	// use limitation
	pipeline = append(pipeline, bson.M{"$limit": limit})

	// execute aggregation
	cursor, err := svc.tracksCol.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("failed to aggregate tracks with details", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		return []dtos.TrackResponse{}, nil
	}

	// decode result
	var tracks []dtos.TrackResponse
	if err := cursor.All(ctx, &tracks); err != nil {
		_logger.Error("failed to decode tracks with details", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// return result
	return tracks, nil
}

func (svc *PlaylistService) AddTrack(
	ctx context.Context,
	userId string,
	playlistId string,
	trackId string,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.PlaylistService.AddTrack"))

	// check track existing
	count, err := svc.tracksCol.CountDocuments(ctx, bson.M{"id": trackId})
	if err != nil {
		_logger.Error("failed to count tracks", logger.ErrorAttr(err))
		return service.ErrInternal
	}
	if count != 1 {
		return service.ErrForbidden
	}

	// prepare updates
	update := bson.M{
		"$addToSet": bson.M{"trackIds": trackId},
		"$set":      bson.M{"updatedAt": time.Now()},
	}
	res, err := svc.playlistsCol.UpdateOne(
		ctx,
		bson.M{"id": playlistId, "ownerId": userId},
		update,
	)
	if err != nil {
		_logger.Error("failed to update", logger.ErrorAttr(err))
		return service.ErrInternal
	}
	if res.MatchedCount == 0 {
		return service.ErrNotFound
	}

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

	// prepare updates
	update := bson.M{
		"$pull": bson.M{"trackIds": trackId},
		"$set":  bson.M{"updatedAt": time.Now()},
	}

	// execute update
	res, err := svc.playlistsCol.UpdateOne(
		ctx,
		bson.M{"id": playlistId, "ownerId": userId},
		update,
	)
	if err != nil {
		_logger.Error("failed to update", logger.ErrorAttr(err))
		return service.ErrInternal
	}
	if res.MatchedCount == 0 {
		return service.ErrNotFound
	}

	return nil
}

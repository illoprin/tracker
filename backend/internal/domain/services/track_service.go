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

	// update album moderation status
	_, err = svc.ac.UpdateOne(
		ctx,
		bson.M{"id": req.AlbumID},
		bson.M{"$set": bson.M{
			"isPublic":           false,
			"isApproved":         false,
			"moderation.status":  schemas.AlbumStatusPending,
			"moderation.comment": "",
		}},
	)
	if err != nil {
		_logger.Error("failed to update album status", logger.ErrorAttr(err))
	}

	// return document
	return &track, nil
}

func (svc *TrackService) GetMetadataByID(
	ctx context.Context,
	userId string,
	userRole int,
	trackId string,
) (*dtos.TrackResponse, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.TrackService.GetMetadataByID"))

	// find track
	var t schemas.Track
	err := svc.tc.FindOne(ctx, bson.M{"id": trackId}).Decode(&t)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		_logger.Error("failed to find track", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// get related album&artist
	var a struct {
		OwnerID    string `bson:"ownerId"`
		AlbumName  string `bson:"name"`
		IsApproved bool   `bson:"isApproved"`
		IsPublic   bool   `bson:"isPublic"`
		Cover      string `bson:"cover"`
		Artist     struct {
			ID   string `bson:"id"`
			Name string `bson:"name"`
		} `bson:"artist"`
	}

	// create pipeline for aggregation request
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"id": t.AlbumID,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"localField":   "artistId",
				"foreignField": "id",
				"from":         "artists",
				"as":           "artist",
			},
		},
		bson.M{
			"$unwind": "$artist",
		},
	}

	// execute aggregation request
	cur, err := svc.ac.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("failed to find album", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	defer cur.Close(ctx)

	// decode cursor
	if cur.Next(ctx) {
		if err := cur.Decode(&a); err != nil {
			_logger.Error("failed to decode cursor", logger.ErrorAttr(err))
			return nil, service.ErrInternal
		}
	} else {
		return nil, service.ErrNotFound
	}

	// check status of album
	if !IsAlbumAccessAllowed(userId, a.OwnerID, a.IsPublic, a.IsApproved, userRole) {
		return nil, service.ErrForbidden
	}

	return &dtos.TrackResponse{
		ID:         trackId,
		Name:       t.Name,
		AlbumID:    t.AlbumID,
		AlbumName:  a.AlbumName,
		ArtistID:   a.Artist.ID,
		ArtistName: a.Artist.Name,
		Cover:      a.Cover,
		Duration:   t.Duration,
		Genres:     t.Genres,
	}, nil
}

func (svc *TrackService) GetAudioPathByID(
	ctx context.Context,
	userId string,
	userRole int,
	trackId string,
) (string, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.TrackService.GetAudioPathByID"))

	// find track
	var t schemas.Track
	err := svc.tc.FindOne(ctx, bson.M{"id": trackId}).Decode(&t)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", service.ErrNotFound
		}
		_logger.Error("failed to find track", logger.ErrorAttr(err))
		return "", service.ErrInternal
	}

	// checkup related album
	var a struct {
		OwnerID    string `bson:"ownerId"`
		IsPublic   bool   `bson:"isPublic"`
		IsApproved bool   `bson:"isApproved"`
	}
	err = svc.ac.FindOne(ctx, bson.M{
		"id": t.AlbumID,
	}).Decode(&a)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", service.ErrNotFound
		}
		_logger.Error("failed to find album", logger.ErrorAttr(err))
		return "", service.ErrInternal
	}

	// check status of album
	if !IsAlbumAccessAllowed(userId, a.OwnerID, a.IsPublic, a.IsApproved, userRole) {
		return "", service.ErrForbidden
	}

	return t.AudioFile, nil
}

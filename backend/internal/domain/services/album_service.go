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
	usersCol  *mongo.Collection
	af        AlbumFlusher
	oc        OwnershipChecker
}

func NewAlbumService(
	albumsCol, tracksCol, usersCol *mongo.Collection,
	albumFlusher AlbumFlusher,
	ownershipChecker OwnershipChecker,
) *AlbumService {
	return &AlbumService{
		albumsCol: albumsCol,
		tracksCol: tracksCol,
		usersCol:  usersCol,
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
	yearInt, _ := strconv.Atoi(req.Year)
	a := schemas.Album{
		ID:         uuid.NewString(),
		ArtistID:   req.ArtistID,
		OwnerID:    userId,
		Name:       req.Name,
		Year:       yearInt,
		Cover:      coverPath,
		Type:       req.Type,
		IsPublic:   false,
		IsApproved: true,
		Moderation: schemas.AlbumModeration{
			Status:  schemas.AlbumStatusApproved,
			Comment: "auto approve",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func (svc *AlbumService) GetByID(ctx context.Context, userId string, userRole int, id string) (*dtos.AlbumWithStats, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.GetByID"))

	// NOTE: a regular user cannot see albums that have not passed moderation

	// prepare aggregate pipeline
	match := bson.M{"id": id}
	pipeline := schemas.GetAlbumsAggregationPipeline(match)

	// execute aggregation
	cursor, err := svc.albumsCol.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("failed to aggregate album stats", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	defer cursor.Close(ctx)

	// get next cursor
	if !cursor.Next(ctx) {
		return nil, service.ErrNotFound
	}

	// decode result
	var result struct {
		schemas.Album `bson:",inline"`
		TrackCount    int    `bson:"trackCount" json:"trackCount"`
		TotalDuration int    `bson:"totalDuration" json:"totalDuration"`
		ArtistName    string `bson:"artistName" json:"artistName"`
		ArtistAvatar  string `bson:"artistAvatar" json:"artistAvatar"`
	}

	if err := cursor.Decode(&result); err != nil {
		_logger.Error("failed to decode album stats", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// check access rights
	if !IsAlbumAccessAllowed(
		userId,
		result.OwnerID,
		result.IsPublic,
		result.IsApproved,
		userRole,
	) {
		return nil, service.ErrForbidden
	}

	// create response struct
	albumWithStats := &dtos.AlbumWithStats{
		Album:         result.Album,
		TrackCount:    result.TrackCount,
		TotalDuration: result.TotalDuration,
		ArtistName:    result.ArtistName,
		ArtistAvatar:  result.ArtistAvatar,
	}

	return albumWithStats, nil
}

func (svc *AlbumService) GetUnapproved(
	ctx context.Context,
	limit int,
	artistId string,
	searchQuery string,
) ([]dtos.AlbumWithStats, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.GetUnapproved"))

	// prepare aggregate pipeline
	match := bson.M{"isApproved": false, "isPublic": true}

	// use filter and search query
	if artistId != "" {
		match["artistId"] = artistId
	}
	if searchQuery != "" {
		match["name"] = bson.M{"$regex": searchQuery, "$options": "i"}
	}
	pipeline := schemas.GetAlbumsAggregationPipeline(match)

	// use limit option
	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

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

func (svc *AlbumService) Moderate(
	ctx context.Context,
	albumId string,
	req dtos.AlbumModerationRequest,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.Moderate"))

	// define updates
	isApproved := req.Status == schemas.AlbumStatusApproved
	updates := bson.M{
		"isApproved":         isApproved,
		"moderation.status":  req.Status,
		"moderation.comment": req.Comment,
	}

	filter := bson.M{"id": albumId, "isPublic": true}

	// apply updates
	_, err := svc.albumsCol.UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return service.ErrNotFound
		}
		_logger.Error("failed to update album status", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	return nil
}

func (svc *AlbumService) GetTracks(
	ctx context.Context,
	userId string,
	userRole int,
	id string,
) ([]dtos.TrackResponse, error) {

	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.GetTracks"))

	// check access to album
	if _, err := svc.GetByID(ctx, userId, userRole, id); err != nil {
		return nil, err
	}

	// prepare aggregation pipeline
	match := bson.M{
		"albumId": id,
	}
	pipeline := schemas.GetTracksDetailsPipeline(match, bson.M{})

	// execute aggregation
	cursor, err := svc.tracksCol.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("failed to aggregate tracks with details", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	defer cursor.Close(ctx)

	// decode response
	var tracks []dtos.TrackResponse
	if err := cursor.All(ctx, &tracks); err != nil {
		_logger.Error("failed to decode tracks with details", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return tracks, nil
}

func (svc *AlbumService) Publish(
	ctx context.Context,
	userId string,
	albumId string,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.Publish"))

	// check album owner
	if isOwn, err := svc.oc.IsAlbumOwner(ctx, userId, albumId); !isOwn {
		if err != nil {
			_logger.Error("failed to check ownership", logger.ErrorAttr(err))
			return service.ErrInternal
		}
		return service.ErrForbidden
	}

	// check count of tracks in album
	count, err := svc.tracksCol.CountDocuments(ctx, bson.M{"albumId": albumId})
	if err != nil {
		_logger.Error("failed count tracks in album", logger.ErrorAttr(err))
		return service.ErrInternal
	}
	if count <= 0 {
		return service.ErrForbidden
	}

	// update isPublic flag
	updates := bson.M{"$set": bson.M{"isPublic": true}}
	_, err = svc.albumsCol.UpdateOne(ctx, bson.M{"id": albumId}, updates)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return service.ErrNotFound
		}
		_logger.Error("failed to update album status", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	return nil
}

func (svc *AlbumService) GetLiked(
	ctx context.Context,
	userId string,
) ([]dtos.AlbumWithStats, error) {

	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.GetLiked"))

	// get liked albums ids
	var u struct {
		LikedAlbums []string `bson:"likedAlbums"`
	}
	err := svc.usersCol.FindOne(ctx, bson.M{"id": userId}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		_logger.Error("failed to find user", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// prepare aggregation pipeline
	filter := bson.M{"id": bson.M{"$in": u.LikedAlbums}}
	pipeline := schemas.GetAlbumsAggregationPipeline(filter)

	// execute aggregation
	cursor, err := svc.albumsCol.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("failed to aggregate", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}
	defer cursor.Close(ctx)

	// decode result
	var result []struct {
		schemas.Album `bson:",inline"`
		TrackCount    int    `bson:"trackCount" json:"trackCount"`
		TotalDuration int    `bson:"totalDuration" json:"totalDuration"`
		ArtistName    string `bson:"artistName" json:"artistName"`
		ArtistAvatar  string `bson:"artistAvatar" json:"artistAvatar"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		_logger.Error("failed to decode albums", logger.ErrorAttr(err))
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

func (svc *AlbumService) Like(
	ctx context.Context,
	userId string, userRole int, albumId string,
) error {

	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.Like"))

	// check access to album
	_, err := svc.GetByID(ctx, userId, userRole, albumId)
	if err != nil {
		return err
	}

	// push album to set
	updates := bson.M{
		"$addToSet": bson.M{"likedAlbums": albumId},
	}
	_, err = svc.usersCol.UpdateOne(
		ctx,
		bson.M{"id": userId},
		updates,
	)
	if err != nil {
		_logger.Error("failed to update album status", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	return nil
}

func (svc *AlbumService) updateAlbum(ctx context.Context, albumId string, set bson.M) error {

	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.updateAlbum"))

	filter := bson.M{"id": albumId}

	// update moderation status
	set["updatedAt"] = time.Now()
	set["isPublic"] = false
	set["isApproved"] = false
	set["moderation.status"] = schemas.AlbumStatusPending
	set["moderation.comment"] = ""

	_, err := svc.albumsCol.UpdateOne(ctx, filter, bson.M{"$set": set})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return service.ErrNotFound
		}
		_logger.Error("failed to update album status", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	return nil
}

func (svc *AlbumService) UpdateByID(
	ctx context.Context,
	userId string,
	albumId string,
	req dtos.AlbumUpdateRequest,
	cover *multipart.File,
	coverHeader *multipart.FileHeader,
	hasCover bool,
) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AlbumService.UpdateByID"))

	// get album
	var a schemas.Album
	err := svc.albumsCol.FindOne(ctx, bson.M{"id": albumId}).Decode(&a)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return service.ErrNotFound
		}
		_logger.Error("failed to find album", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	if a.OwnerID != userId {
		return service.ErrForbidden
	}

	// check updates
	updates := bson.M{}

	// name
	if req.Name != "" {
		updates["name"] = req.Name
	}

	// year
	if req.Year != "" {
		year, _ := strconv.Atoi(req.Year)
		updates["year"] = year
	}

	// type
	if req.Type != "" {
		updates["type"] = req.Type
	}

	// update cover if needed
	if hasCover {
		// delete old cover file if needed
		oldFileName := filepath.Base(a.Cover)
		oldFilePath := filepath.Join(
			os.Getenv(config.StaticDirEnvName),
			config.CoversDir,
			oldFileName,
		)
		if err := os.Remove(oldFilePath); err != nil {
			_logger.Error("failed to delete old cover", logger.ErrorAttr(err))
		}

		// upload new cover
		uploadDir := path.Join(
			os.Getenv(config.StaticDirEnvName),
			config.CoversDir,
		)
		if filePath, err := storage.UploadFile(coverHeader, cover, uploadDir); err != nil {
			return err
		} else {
			// define new file url in updates
			updates["cover"] = filePath
		}
	}

	if len(updates) <= 0 {
		return nil
	}

	if updates["name"] != "" {
		// check albums with same name of this user
		count, err := svc.albumsCol.CountDocuments(ctx, bson.M{"name": req.Name, "ownerId": userId})
		if err != nil {
			_logger.Error("failed count documents", logger.ErrorAttr(err))
			return service.ErrInternal
		}
		if count > 0 {
			return service.ErrExists
		}
	}

	// apply updates
	err = svc.updateAlbum(ctx, albumId, updates)
	if err != nil {
		return err
	}

	return nil
}

func IsAlbumAccessAllowed(userId, ownerId string, isPublic, isApproved bool, userRole int) bool {
	isOwn := userId == ownerId
	isPermittedRegular := isPublic && isApproved
	isModerator := userRole >= schemas.RoleModerator
	if !isOwn && !isPermittedRegular && !isModerator {
		return false
	}
	return true
}

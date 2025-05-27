package track

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/auth/ownership"
	"tracker-backend/internal/config"
	uploadfile "tracker-backend/internal/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrNotFound         = errors.New("track not found")
	ErrInvalidFileType  = errors.New("invalid file type")
	ErrFileUploadFailed = errors.New("file upload failed")
	ErrUndefinedAlbum   = errors.New("related album is undefined")
)

type TrackService struct {
	Collection       *mongo.Collection
	ownershipService *ownership.OwnershipService
	AlbumCollection  *mongo.Collection
}

// NewService creates new service for tracks
func NewTrackService(ctx context.Context, db *mongo.Database) *TrackService {
	// create collections
	col := db.Collection("tracks")
	albumCol := db.Collection("albums")

	// create indices
	err := EnsureIndexes(ctx, col)
	if err != nil {
		panic(err.Error())
	}

	// create ownership service
	ownershipService := ownership.NewOwnershipService(db.Collection("artists"), albumCol)

	return &TrackService{
		Collection:       col,
		AlbumCollection:  albumCol,
		ownershipService: ownershipService,
	}
}

// Create creates new track and creates file
func (s *TrackService) Create(
	ctx context.Context,
	req *CreateTrackRequest,
	audioFile *multipart.File,
	fileHeader *multipart.FileHeader,
) (*Track, error) {
	// configure logger
	logger := slog.With(slog.String("function", "track.TrackService.Create"))

	// check album existence
	err := s.AlbumCollection.FindOne(ctx, bson.M{"id": req.AlbumID}).Err()
	if err != nil {
		return nil, ErrUndefinedAlbum
	}

	// check file type
	if err := uploadfile.ValidateFile(fileHeader, uploadfile.AllowedAudioExtensions); err != nil {
		return nil, ErrInvalidFileType
	}

	// create file path
	filePath := path.Join(os.Getenv(config.PublicDirPathEnvName), config.AudioDir)
	// upload file
	audioFilePath, err := uploadfile.UploadFile(
		fileHeader,
		audioFile,
		filePath,
		uploadfile.AllowedAudioExtensions,
	)
	if err != nil {
		logger.Error("failed to upload file", slog.String("error", err.Error()))
		return nil, err
	}

	// create new document
	track := &Track{
		ID:        uuid.NewString(),
		Title:     req.Title,
		Genre:     req.Genre,
		Duration:  req.Duration,
		AudioFile: audioFilePath,
		AlbumID:   req.AlbumID,
		CreatedAt: time.Now(),
	}

	// insert to collection
	_, err = s.Collection.InsertOne(ctx, track)
	if err != nil {
		logger.Error("failed to insert", slog.String("error", err.Error()))
		// delete related file if error occurred
		os.Remove(filePath)
		return nil, ErrFileUploadFailed
	}

	logger.Info("track uploaded",
		slog.Group("info",
			slog.String("albumID", req.AlbumID),
			slog.String("id", track.ID),
			slog.String("audioPath", audioFilePath),
		),
	)

	return track, nil
}

// GetByID get track document by id
func (s *TrackService) GetByID(
	ctx context.Context, id string,
) (*Track, error) {
	// configure logger
	logger := slog.With(slog.String("function", "track.TrackService.GetByID"))
	_ = logger

	var track Track
	err := s.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&track)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return &track, nil
}

// GetFilePathByID returns full file path to track
func (s *TrackService) GetFilePathByID(ctx context.Context, id string) (string, error) {
	// configure logger
	logger := slog.With(slog.String("function", "track.TrackService.GetFilePathByID"))
	_ = logger

	// get track
	track, err := s.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	// define filepath
	filePath := filepath.Join(
		os.Getenv(config.PublicDirPathEnvName),
		config.AudioDir,
		track.AudioFile,
	)

	// check file existence
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("audio file not found")
	}

	return filePath, nil
}

// Delete removes track by id
func (s *TrackService) Delete(
	ctx context.Context, id string, userID string,
) error {
	// configure logger
	logger := slog.With(slog.String("function", "track.TrackService.Delete"))
	_ = logger

	// check ownership
	isOwner, err := s.ownershipService.IsTrackOwner(ctx, id, userID)
	if err != nil {
		return errors.New("failed to check ownership")
	}
	if !isOwner {
		return auth.ErrAccessDenied
	}

	// find document
	var track Track
	if err := s.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&track); err != nil {
		return errors.New("failed to find")
	}

	// remove track
	filePath := filepath.Join(os.Getenv(config.PublicDirPathEnvName), config.AudioDir, track.AudioFile)
	if err := os.Remove(filePath); err != nil {
		return errors.New("failed to delete audio")
	}

	// remove document
	res, err := s.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return errors.New("failed to remove")
	}
	if res.DeletedCount < 1 {
		return ErrNotFound
	}

	return nil
}

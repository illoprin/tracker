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
	"tracker-backend/internal/auth/ownership"
	"tracker-backend/internal/config"
	uploadfile "tracker-backend/internal/pkg/file"
	"tracker-backend/internal/pkg/service"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TrackService struct {
	Col              *mongo.Collection
	ownershipService *ownership.OwnershipService
	AlbumChecker
}

type AlbumChecker interface {
	CheckExistence(ctx context.Context, albumID string) (bool, error)
}

// NewService creates new service for tracks
func NewTrackService(
	tracksCollection *mongo.Collection,
	ownershipService *ownership.OwnershipService,
	albumChecker AlbumChecker,
) *TrackService {
	return &TrackService{
		Col:              tracksCollection,
		AlbumChecker:     albumChecker,
		ownershipService: ownershipService,
	}
}

// Create creates new track and creates file
func (s *TrackService) Create(
	ctx context.Context,
	userID string,
	req *CreateTrackRequest,
	audioFile *multipart.File,
	fileHeader *multipart.FileHeader,
) (*Track, error) {
	// configure logger
	logger := slog.With(slog.String("function", "track.TrackService.Create"))

	// check album existence
	if exists, err := s.AlbumChecker.CheckExistence(ctx, req.AlbumID); !exists {
		if err != nil {
			logger.Warn("failed to check album existence", slog.String("error", err.Error()))
			return nil, errors.New("failed to check album existence")
		}
		return nil, service.ErrNotFound
	}

	// check album owner
	if isOwn, err := s.ownershipService.IsAlbumOwner(ctx, userID, req.AlbumID); !isOwn {
		if err != nil {
			return nil, errors.New("failed to check album owner")
		}
		return nil, service.ErrAccessDenied
	}

	// check file type
	if err := uploadfile.ValidateFile(fileHeader, uploadfile.AllowedAudioExtensions); err != nil {
		return nil, err
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
		AudioFile: filepath.Base(audioFilePath),
		AlbumID:   req.AlbumID,
		CreatedAt: time.Now(),
	}

	// insert to collection
	_, err = s.Col.InsertOne(ctx, track)
	if err != nil {
		logger.Error("failed to insert", slog.String("error", err.Error()))
		// delete related file if error occurred
		os.Remove(filePath)
		return nil, service.ErrUploadFailed
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

	var track Track
	err := s.Col.FindOne(ctx, bson.M{"id": id}).Decode(&track)
	if err != nil {
		logger.Warn("failed to find track metadata", slog.String("error", err.Error()))
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return &track, nil
}

// GetFilePathByID returns full file path to track
func (s *TrackService) GetFilePathByID(ctx context.Context, id string) (string, error) {
	// configure logger
	logger := slog.With(slog.String("function", "track.TrackService.GetFilePathByID"))

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
		logger.Warn("failed to get track file", slog.String("error", err.Error()))
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

	// check ownership
	isOwner, err := s.ownershipService.IsTrackOwner(ctx, id, userID)
	if err != nil {
		return errors.New("failed to check ownership")
	}
	if !isOwner {
		return service.ErrAccessDenied
	}

	// find document
	var track Track
	if err := s.Col.FindOne(ctx, bson.M{"id": id}).Decode(&track); err != nil {
		return service.ErrNotFound
	}

	// remove track
	filePath := filepath.Join(os.Getenv(config.PublicDirPathEnvName), config.AudioDir, track.AudioFile)
	if err := os.Remove(filePath); err != nil {
		logger.Warn("failed to remove file", slog.String("error", err.Error()))
		return errors.New("failed to delete audio")
	}

	// remove document
	res, err := s.Col.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		logger.Warn("failed to delete document", slog.String("error", err.Error()))
		return errors.New("failed to remove")
	}
	if res.DeletedCount < 1 {
		return service.ErrNotFound
	}

	return nil
}

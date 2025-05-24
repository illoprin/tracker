package track

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
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

var (
	allowedExtensions = map[string]bool{
		".mp3": true,
		".wav": true,
		".m4a": true,
	}
)

type TrackService struct {
	Collection      *mongo.Collection
	AlbumCollection *mongo.Collection
	uploadPath      string
}

// NewService создает новый сервис для работы с треками
func NewTrackService(ctx context.Context, db *mongo.Database, uploadPath string) *TrackService {
	// create collections
	col := db.Collection("tracks")
	albumCol := db.Collection("album")

	// create indices
	err := EnsureIndexes(ctx, col)
	if err != nil {
		panic(err.Error())
	}

	return &TrackService{
		Collection:      col,
		AlbumCollection: albumCol,
		uploadPath:      uploadPath,
	}
}

// CreateTrack creates new track and creates file
func (s *TrackService) CreateTrack(
	ctx context.Context,
	req *CreateTrackRequest,
	audioFile *multipart.File,
	fileHeader *multipart.FileHeader,
) (*Track, error) {

	// check album existence
	err := s.AlbumCollection.FindOne(ctx, bson.M{"id": req.AlbumID}).Err()
	if err != nil {
		return nil, ErrUndefinedAlbum
	}

	// check file type
	if err := uploadfile.ValidateFile(fileHeader, allowedExtensions); err != nil {
		return nil, ErrInvalidFileType
	}

	// TODO: define duration and audio path
	filePath := uploadfile.GetFullPath(s.uploadPath)
	audioFilePath, err := uploadfile.UploadFile(
		fileHeader,
		audioFile,
		filePath,
		allowedExtensions,
	)
	duration := req.Duration
	if err != nil {
		return nil, err
	}

	// create new document
	track := &Track{
		ID:        uuid.NewString(),
		Title:     req.Title,
		Genre:     req.Genre,
		Duration:  duration,
		AudioFile: audioFilePath,
		AlbumID:   req.AlbumID,
		CreatedAt: time.Now(),
	}

	// insert to collection
	_, err = s.Collection.InsertOne(ctx, track)
	if err != nil {
		// delete related file if error occurred
		os.Remove(filePath)
		return nil, ErrFileUploadFailed
	}

	return track, nil
}

// GetTrackByID get track document by id
func (s *TrackService) GetTrackByID(
	ctx context.Context, id string,
) (*Track, error) {
	var track Track
	err := s.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&track)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return &track, nil
}

// GetTrackFilePath returns full file path to track
func (s *TrackService) GetTrackFilePath(ctx context.Context, id string) (string, error) {
	track, err := s.GetTrackByID(ctx, id)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(
		uploadfile.GetFullPath(s.uploadPath),
		track.AudioFile,
	)

	// check file existence
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("audio file not found")
	}

	return filePath, nil
}

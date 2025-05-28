package artist

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
	artistType "tracker-backend/internal/artist/type"
	"tracker-backend/internal/config"
	uploadfile "tracker-backend/internal/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotFound  = errors.New("artist not found or not owned by the user")
	ErrNameTaken = errors.New("name is already in use")
)

type ArtistService struct {
	Col *mongo.Collection
}

// NewArtistService new artist service instance
func NewArtistService(artistCol *mongo.Collection) *ArtistService {
	return &ArtistService{Col: artistCol}
}

// Create new artist from CreateRequest
func (s *ArtistService) Create(
	ctx context.Context, userID string, req artistType.CreateRequest,
) (*artistType.Artist, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	artist := &artistType.Artist{
		ID:         uuid.New().String(),
		Name:       req.Name,
		UserID:     userID,
		AvatarPath: "public/avatars/avatar_default.jpg",
		CreatedAt:  time.Now(),
	}

	_, err := s.Col.InsertOne(ctx, artist)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrNameTaken
		}
		return nil, err
	}

	return artist, nil
}

// Delete deletes artist
func (s *ArtistService) Delete(
	ctx context.Context, artistID, userID string,
) error {
	filter := bson.M{"id": artistID, "userID": userID}
	res, err := s.Col.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// TODO: delete related tracks and albums

	if res.DeletedCount == 0 {
		return errors.New("artist not found or not owned by the user")
	}
	return nil
}

// Update updates artist name
func (s *ArtistService) Update(
	ctx context.Context,
	artistID string, userID string,
	req artistType.UpdateRequest,
) (*artistType.Artist, error) {

	update := bson.M{}
	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.AvatarPath != nil {
		update["avatarPath"] = *req.AvatarPath
	}

	// check name uniqueness
	if name, ok := update["name"]; ok {
		count, err := s.Col.CountDocuments(ctx, bson.M{"name": name, "id": bson.M{"$ne": artistID}})

		if err != nil {
			return nil, fmt.Errorf("error while checking name uniqueness")
		}

		if count > 0 {
			return nil, ErrNameTaken
		}
	}

	filter := bson.M{"id": artistID, "userID": userID}

	var artist artistType.Artist
	err := s.Col.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$set": update}, // FIX
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&artist)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &artist, nil
}

// UpdateAvatar updates artist avatar from form file
func (s *ArtistService) UpdateAvatar(
	ctx context.Context,
	userID string, artistID string,
	file *multipart.File, fileHeader *multipart.FileHeader,
) (*artistType.Artist, error) {
	// validate file
	if err := uploadfile.ValidateFile(fileHeader, uploadfile.AllowedImageExtensions); err != nil {
		return nil, err
	}
	filter := bson.M{"userID": userID, "id": artistID}

	// find artist and decode
	var artist artistType.Artist
	err := s.Col.FindOne(ctx, filter).Decode(&artist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// check avatar file path
	if !strings.Contains(artist.AvatarPath, "avatar_default") {
		// delete old one
		if _, err := os.Stat(artist.AvatarPath); err == nil {
			if err := os.Remove(artist.AvatarPath); err != nil {
				return nil, fmt.Errorf("failed to delete old file")
			}
		}
	}

	// upload new file
	fileDir := path.Join(os.Getenv(config.PublicDirPathEnvName), config.AvatarsDir)
	newFilePath, err := uploadfile.UploadFile(fileHeader, file, fileDir, uploadfile.AllowedImageExtensions)
	if err != nil {
		return nil, err
	}

	// create new update data
	update := bson.M{"avatarPath": newFilePath}

	var updatedArtist artistType.Artist
	err = s.Col.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedArtist)

	if err != nil {
		return nil, fmt.Errorf("failed to update")
	}

	return &updatedArtist, nil
}

// GetByID returns artist by id
func (s *ArtistService) GetByID(
	ctx context.Context, artistID string,
) (*artistType.Artist, error) {
	var artist artistType.Artist
	err := s.Col.FindOne(ctx, bson.M{"id": artistID}).Decode(&artist)
	if err != nil {
		return nil, ErrNotFound
	}
	return &artist, nil
}

// GetByUserID returns all artists by userID
func (s *ArtistService) GetByUserID(
	ctx context.Context, userID string,
) ([]artistType.Artist, error) {
	cursor, err := s.Col.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		return nil, errors.New("failed to find matches")
	}
	defer cursor.Close(ctx)

	var artists []artistType.Artist
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, errors.New("failed to parse result")
	}
	return artists, nil
}

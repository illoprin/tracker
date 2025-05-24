package artist

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func NewArtistService(ctx context.Context, db *mongo.Database) *ArtistService {
	// create collection
	col := db.Collection("artists")

	// create indices
	err := EnsureIndexes(ctx, col)
	if err != nil {
		panic(err.Error())
	}

	return &ArtistService{Col: col}
}

// creates new artist
func (s *ArtistService) Create(
	ctx context.Context, userID string, req CreateRequest,
) (*Artist, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	artist := &Artist{
		ID:         uuid.New().String(),
		Name:       req.Name,
		UserID:     userID,
		AvatarPath: "/upload/avatar_default.jpg",
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

// deletes artist
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

// updates artists name and avatarPath
func (s *ArtistService) Update(
	ctx context.Context,
	artistID string, userID string,
	req UpdateRequest,
) (*Artist, error) {

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

	var artist Artist
	err := s.Col.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$set": update}, // FIX
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&artist)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &artist, nil
}

// returns artist by id
func (s *ArtistService) GetByID(
	ctx context.Context, artistID string,
) (*Artist, error) {
	var artist Artist
	err := s.Col.FindOne(ctx, bson.M{"id": artistID}).Decode(&artist)
	if err != nil {
		return nil, ErrNotFound
	}
	return &artist, nil
}

// returns all artists by userID
func (s *ArtistService) GetByUserID(
	ctx context.Context, userID string,
) ([]Artist, error) {
	cursor, err := s.Col.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var artists []Artist
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

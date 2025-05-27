package ownership

import (
	"context"
	"errors"
	"log/slog"
	"tracker-backend/internal/artist"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OwnershipService struct {
	artistsCol *mongo.Collection
	albumsCol  *mongo.Collection
}

func NewOwnershipService(artistsCol, albumsCol *mongo.Collection) *OwnershipService {
	return &OwnershipService{
		albumsCol:  albumsCol,
		artistsCol: artistsCol,
	}
}

// IsAlbumOwner checks album ownership
func (s *OwnershipService) IsAlbumOwner(ctx context.Context, userID, albumID string) (bool, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"id":       albumID,
				"artistID": bson.M{"$exists": true},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "artists",
				"localField":   "artistID",
				"foreignField": "id",
				"as":           "artist",
			},
		},
		{
			"$unwind": "$artist",
		},
		{
			"$match": bson.M{
				"artist.userID": userID,
			},
		},
		{
			"$count": "matches",
		},
	}

	logger := slog.With(slog.String("function", "ownership.OwnershipService.IsArtistOwner"))

	// process aggregate by pipeline
	cursor, err := s.albumsCol.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Warn("failed to aggregate", slog.String("error", err.Error()))
		return false, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Matches int `bson:"matches"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			logger.Warn("failed to decode", slog.String("error", err.Error()))
			return false, err
		}
		// check matches and return result
		return result.Matches > 0, nil
	}

	return false, nil
}

// IsArtistOwner returns true if user is owns artist
func (s *OwnershipService) IsArtistOwner(
	ctx context.Context, userID, artistID string,
) (bool, error) {
	// configure logger
	logger := slog.With(slog.String("function", "ownership.OwnershipService.IsArtistOwner"))

	// get artist
	var artistEntry artist.Artist
	err := s.artistsCol.FindOne(
		ctx, bson.M{"userID": userID, "id": artistID},
	).Decode(&artistEntry)

	// process errors
	if err != nil {
		logger.Warn("failed to findone", slog.String("error", err.Error()))
		if err == mongo.ErrNoDocuments {
			return false, artist.ErrNotFound
		}
		return false, errors.New("failed to find artist")
	}

	// check ownership
	return artistEntry.UserID == userID, nil
}

func (s *OwnershipService) IsTrackOwner(
	ctx context.Context, trackID, userID string,
) (bool, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"userID": userID,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "albums",
				"localField":   "id",
				"foreignField": "artistID",
				"as":           "albums",
			},
		},
		{
			"$unwind": "$albums",
		},
		{
			"$lookup": bson.M{
				"from":         "tracks",
				"localField":   "albums.id",
				"foreignField": "album",
				"as":           "tracks",
			},
		},
		{
			"$unwind": "$tracks",
		},
		{
			"$match": bson.M{
				"tracks.id": trackID,
			},
		},
		{
			"$project": bson.M{
				"id": 0,
				"isOwn": bson.M{
					// здесь должен быть массив
					"$eq": bson.A{"$userID", userID},
				},
			},
		},
		{
			"$limit": 1,
		},
	}

	logger := slog.With(slog.String("function", "ownership.OwnershipService.IsTrackOwner"))

	cur, err := s.artistsCol.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Warn("failed to aggregate", slog.String("error", err.Error()))
		return false, err
	}
	defer cur.Close(ctx)

	var res struct {
		isOwn bool
	}

	if cur.Next(ctx) {
		if err := cur.Decode(&res); err != nil {
			logger.Warn("failed to decode", slog.String("error", err.Error()))
			return false, err
		}
		return res.isOwn, nil
	}

	return false, nil
}

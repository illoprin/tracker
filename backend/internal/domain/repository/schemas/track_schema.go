package schemas

import (
	"context"
	"time"
	"tracker-backend/internal/domain/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Track struct {
	ID        string    `bson:"id" json:"id"`
	OwnerID   string    `bson:"ownerId" json:"ownerId"`
	AlbumID   string    `bson:"albumId" json:"albumId"`
	Name      string    `bson:"name" json:"name"`
	Genres    []string  `bson:"genres" json:"genres"`
	Duration  int       `bson:"duration" json:"duration"`
	AudioFile string    `bson:"audioFile" json:"audioFile"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func EnsureTrackIndices(ctx context.Context, col *mongo.Collection) error {
	// unique index by name and album id
	nameAlbumIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "albumId", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("name_albumId_unique"),
	}
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{nameAlbumIndex, utils.UniqueIDIndex()})
	return err
}

func GetTracksDetailsPipeline(match bson.M) bson.A {
	return bson.A{
		bson.M{
			"$match": match,
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "albums",
				"localField":   "albumId",
				"foreignField": "id",
				"as":           "album",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$album",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "artists",
				"localField":   "album.artistId",
				"foreignField": "id",
				"as":           "artist",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$artist",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$project": bson.M{
				"id":         1,
				"name":       1,
				"albumId":    1,
				"duration":   1,
				"genres":     1,
				"artistId":   "$album.artistId",
				"albumName":  "$album.name",
				"artistName": "$artist.name",
				"cover":      "$album.cover",
			},
		},
	}
}

package schemas

import (
	"context"
	"time"
	"tracker-backend/internal/domain/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Playlist struct {
	ID          string    `bson:"id" json:"id"`
	OwnerID     string    `bson:"ownerId" json:"ownerId"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Cover       string    `bson:"cover" json:"cover"` // cover path
	IsDefault   bool      `bson:"isDefault" json:"isDefault"`
	IsPublic    bool      `bson:"isPublic" json:"isPublic"`
	TrackIDs    []string  `bson:"trackIds" json:"trackIds"` // tracks links
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

func GetPlaylistsAggregationPipeline(match bson.M) bson.A {
	pipeline := bson.A{
		bson.M{
			"$match": match,
		},
		bson.M{
			"$lookup": bson.M{
				"localField":   "ownerId",
				"foreignField": "id",
				"from":         "users",
				"as":           "owner",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$owner",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$project": bson.M{
				"ownerInfo": bson.M{
					"id":   "$owner.id",
					"name": "$owner.login",
				},
				"id":          1,
				"name":        1,
				"description": 1,
				"cover":       1,
				"isPublic":    1,
				"createdAt":   1,
				"updatedAt":   1,
			},
		},
	}
	return pipeline
}

func EnsurePlaylistIndices(ctx context.Context, col *mongo.Collection) error {
	// unique index userID and name
	nameUserIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "ownerId", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("name_user_index"),
	}

	isPublicTracksIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "trackIds", Value: 1},
			{Key: "isPublic", Value: 1},
		},
		Options: options.Index().SetName("tracks_isPublic_index"),
	}

	_, err := col.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			nameUserIndex,
			isPublicTracksIndex,
			utils.UniqueIDIndex(),
		},
	)
	return err
}

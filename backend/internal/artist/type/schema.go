package artistType

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Artist struct {
	ID         string    `bson:"id" json:"id"`
	Name       string    `bson:"name" json:"name"`
	UserID     string    `bson:"userID" json:"userID"`
	AvatarPath string    `bson:"avatarPath" json:"avatarPath"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}

// create indices
func EnsureIndexes(ctx context.Context, col *mongo.Collection) error {
	// unique index by name within one user
	artistNameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "userID", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("name_userID_unique"),
	}

	// index for fast search by name
	userIDIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "userID", Value: 1}},
		Options: options.Index().SetName("userID_index"),
	}

	// unique index by id string
	idIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("id_unique"),
	}

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		artistNameIndex,
		userIDIndex,
		idIndex,
	})
	return err
}

package schemas

import (
	"context"
	"time"
	"tracker-backend/internal/domain/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Artist struct {
	ID        string    `bson:"id" json:"id"`
	OwnerID   string    `bson:"ownerId" json:"ownerId"`
	Name      string    `bson:"name" json:"name"`
	Avatar    string    `bson:"avatar" json:"avatar"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func EnsureArtistIndices(ctx context.Context, col *mongo.Collection) error {
	// unique index by name within one user
	artistNameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "ownerId", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("name_userID_unique"),
	}

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		artistNameIndex, utils.UniqueIDIndex(),
	})
	return err
}

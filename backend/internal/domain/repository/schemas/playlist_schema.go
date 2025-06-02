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
	ID        string    `bson:"id" json:"id"`
	UserID    string    `bson:"userId" json:"userId"`
	Name      string    `bson:"name" json:"name"`
	IsPublic  bool      `bson:"isPublic" json:"isPublic"`
	IsDefault bool      `bson:"isDefault" json:"isDefault"`
	Cover     string    `bson:"cover" json:"cover"`
	Tracks    []string  `bson:"trackIds" json:"trackIds"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func EnsurePlaylistIndices(ctx context.Context, col *mongo.Collection) error {
	// unique index userID and name
	nameUserIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "userId", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("name_user_index"),
	}

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{nameUserIndex, utils.UniqueIDIndex()})
	return err
}

package playlist

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Playlist struct {
	ID        string   `bson:"id"`
	Name      string   `bson:"name"`
	UserID    string   `bson:"userID"`
	IsDefault bool     `bson:"isDefault"`
	IsPublic  bool     `bson:"isPublic"`
	TrackIDs  []string `bson:"trackIDs"`
}

func EnsureIndexes(ctx context.Context, col *mongo.Collection) error {
	// unique index userID and name
	nameUserIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "userID", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("name_user_index"),
	}

	// unique index by id string
	idIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("id_unique"),
	}

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{nameUserIndex, idIndex})
	return err
}

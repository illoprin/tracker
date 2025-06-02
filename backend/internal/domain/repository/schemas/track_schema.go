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
	AlbumID   string    `bson:"albumId" json:"albumId"`
	OwnerID   string    `bson:"ownerId" json:"ownerId"`
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

package track

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Track represents a music track in the system
type Track struct {
	ID        string    `bson:"id"`
	Title     string    `bson:"title"`
	Duration  int       `bson:"duration"` // duration in seconds
	Genre     []string  `bson:"genre"`
	AudioFile string    `bson:"audioFile"`
	AlbumID   string    `bson:"album"`
	CreatedAt time.Time `bson:"createdAt"`
}

func EnsureIndexes(ctx context.Context, col *mongo.Collection) error {
	// unique index by name and album id
	nameAlbumIndex := mongo.IndexModel{
		Keys:    bson.M{"name": 1, "album": 1},
		Options: options.Index().SetUnique(true).SetName("name_album_unique"),
	}

	_, err := col.Indexes().CreateOne(ctx, nameAlbumIndex)
	return err
}

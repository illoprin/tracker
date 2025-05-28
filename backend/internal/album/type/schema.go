package albumType

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Album struct {
	ID        string    `bson:"id"`
	Title     string    `bson:"title"`
	ArtistID  string    `bson:"artistID"`
	Year      int       `bson:"year"`
	CoverPath string    `bson:"coverPath"`
	Genres    []string  `bson:"genres"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"createdAt"`
}

const (
	StatusPublic       = "Public"
	StatusHidden       = "Hidden"
	StatusOnModeration = "OnModeration"
)

func EnsureIndexes(
	ctx context.Context, col *mongo.Collection,
) error {
	artistIDTitleIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "artistID", Value: 1},
			{Key: "title", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	// unique index by id string
	idIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("id_unique"),
	}

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{artistIDTitleIndex, idIndex})
	return err
}

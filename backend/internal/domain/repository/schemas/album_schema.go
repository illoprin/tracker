package schemas

import (
	"context"
	"time"
	"tracker-backend/internal/domain/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	AlbumTypeAlbum  = "album"
	AlbumTypeSingle = "single"
)

type AlbumModeration struct {
	Status  string `bson:"status"`
	Comment string `bson:"comment"`
}

type Album struct {
	ID         string          `bson:"id" json:"id"`
	ArtistID   string          `bson:"artistId" json:"artistId"`
	OwnerID    string          `bson:"ownerId" json:"ownerId"`
	Name       string          `bson:"name" json:"name"`
	Year       int             `bson:"year" json:"year"`
	Cover      string          `bson:"cover" json:"cover"` // cover path
	Type       string          `bson:"type" json:"type"`
	IsPublic   bool            `bson:"isPublic" json:"isPublic"`
	IsApproved bool            `bson:"isApproved" json:"isApproved"`
	Moderation AlbumModeration `bson:"moderation" json:"moderation"`
	CreatedAt  time.Time       `bson:"createdAt" json:"createdAt"`
}

func EnsureAlbumIndices(ctx context.Context, col *mongo.Collection) error {
	artistIDTitleIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "artistId", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{artistIDTitleIndex, utils.UniqueIDIndex()})
	return err
}

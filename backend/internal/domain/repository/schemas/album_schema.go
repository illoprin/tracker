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

var (
	AlbumStatusApproved = "approved"
	AlbumStatusPending  = "pending"
	AlbumStatusRejected = "rejected"
)

type AlbumModeration struct {
	Status  string `bson:"status" json:"status"`
	Comment string `bson:"comment" json:"comment"`
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

func GetAlbumsAggregationPipeline(match bson.M) bson.A {
	return bson.A{
		bson.M{
			"$match": match,
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "tracks",
				"localField":   "id",
				"foreignField": "albumId",
				"as":           "tracks",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "artists",
				"localField":   "artistId",
				"foreignField": "id",
				"as":           "artist",
			},
		},
		bson.M{
			"$unwind": "$artist",
		},
		bson.M{
			"$addFields": bson.M{
				"trackCount":    bson.M{"$size": "$tracks"},
				"totalDuration": bson.M{"$sum": "$tracks.duration"},
				"artistName":    "$artist.name",
				"artistAvatar":  "$artist.avatar",
			},
		},
		bson.M{
			"$project": bson.M{
				"tracks": 0, // exclude array of tracks from result
				"artist": 0, // exclude artist metadata from result
			},
		},
	}
}

package schemas

import (
	"context"
	"time"
	"tracker-backend/internal/domain/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID              string    `bson:"id" json:"id"`
	Email           string    `bson:"email" json:"email"`
	Login           string    `bson:"login" json:"login"`
	LikedArtists    []string  `bson:"likedArtists" json:"likedArtists"`
	LikedAlbums     []string  `bson:"likedAlbums" json:"likedAlbums"`
	LikedPlaylistId string    `bson:"likedPlaylistId" json:"likedPlaylistId"`
	PasswordHash    string    `bson:"passwordHash" json:"passwordHash"`
	Role            int       `bson:"role" json:"role"`
	CreatedAt       time.Time `bson:"createdAt" json:"createdAt"`
}

func EnsureUserIndices(ctx context.Context, col *mongo.Collection) error {

	// unique index by user login
	loginIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "login", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("login_unique"),
	}

	// unique index by user email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("email_unique"),
	}

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{loginIndex, emailIndex, utils.UniqueIDIndex()})
	return err
}

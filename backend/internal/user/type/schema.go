package userType

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID               string    `bson:"id" json:"id"`
	Login            string    `bson:"login" json:"login"`
	Email            string    `bson:"email" json:"email"`
	PasswordHash     string    `bson:"passwordHash" json:"passwordHash"`
	MyChoicePlaylist string    `bson:"myChoicePlaylist" json:"myChoicePlaylist"`
	CreatedAt        time.Time `bson:"createdAt" json:"createdAt"`
	Role             int       `bson:"role" json:"role"`
}

func EnsureIndexes(ctx context.Context, col *mongo.Collection) error {
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
	// unique index by id string
	idIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("id_unique"),
	}

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{loginIndex, emailIndex, idIndex})
	return err
}

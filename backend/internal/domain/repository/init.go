package repository

import (
	"context"
	"tracker-backend/internal/domain/repository/schemas"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitRepository(ctx context.Context, db *mongo.Database) (*Repository, error) {
	userCol := db.Collection("users")
	if err := schemas.EnsureUserIndices(ctx, userCol); err != nil {
		return nil, err
	}

	return &Repository{
		UsersCol: userCol,
	}, nil
}

package service

import (
	"context"
	"tracker-backend/internal/env"
	userService "tracker-backend/internal/user/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Dependencies struct {
	UserService *userService.UserService
}

func InitDependencies(
	ctx context.Context, db *mongo.Database, cfg *env.Config,
) *Dependencies {
	userService := userService.NewUserService(ctx, db, cfg.JWTSecret)

	return &Dependencies{
		UserService: userService,
	}
}

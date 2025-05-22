package setup

import (
	"context"
	"tracker-backend/internal/artist"
	"tracker-backend/internal/config"
	"tracker-backend/internal/user"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Dependencies struct {
	UserService   *user.UserService
	ArtistService *artist.ArtistService
}

func InitDependencies(
	ctx context.Context, db *mongo.Database, cfg *config.Config,
) *Dependencies {
	userService := user.NewUserService(ctx, db, cfg.JWTSecret)
	artistService := artist.NewArtistService(ctx, db)

	return &Dependencies{
		UserService:   userService,
		ArtistService: artistService,
	}
}

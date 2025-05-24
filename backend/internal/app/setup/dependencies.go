package setup

import (
	"context"
	"os"
	"tracker-backend/internal/artist"
	"tracker-backend/internal/config"
	"tracker-backend/internal/track"
	"tracker-backend/internal/user"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Dependencies struct {
	UserService   *user.UserService
	ArtistService *artist.ArtistService
	TrackService  *track.TrackService
}

func InitDependencies(
	ctx context.Context, db *mongo.Database,
) *Dependencies {
	userService := user.NewUserService(
		ctx, db,
		os.Getenv(config.JWTSecretEnvName),
	)
	artistService := artist.NewArtistService(ctx, db)
	trackService := track.NewTrackService(ctx, db, "audio")

	return &Dependencies{
		UserService:   userService,
		ArtistService: artistService,
		TrackService:  trackService,
	}
}

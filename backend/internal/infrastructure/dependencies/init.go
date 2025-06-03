package dependencies

import (
	"tracker-backend/internal/domain/repository"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/infrastructure/redis"
)

func InitDependencies(
	repo *repository.Repository,
	redisClient *redis.RedisClient,
) *Dependencies {
	flusher := services.NewFlushService(
		repo.TracksCol,
		repo.ArtistsCol,
		repo.AlbumsCol,
		repo.PlaylistsCol,
	)
	ownershipService := services.NewOwnershipService(
		repo.AlbumsCol,
		repo.ArtistsCol,
		repo.TracksCol,
	)
	authSvc := services.NewAuthorizationService(repo.UsersCol, redisClient)
	userSvc := services.NewUserService(repo.UsersCol, flusher)
	artistSvc := services.NewArtistService(
		repo.ArtistsCol,
		repo.AlbumsCol,
		flusher,
		ownershipService,
	)
	albumSvc := services.NewAlbumService(
		repo.AlbumsCol,
		repo.TracksCol,
		flusher,
		ownershipService,
	)

	return &Dependencies{
		AuthSvc:   authSvc,
		UserSvc:   userSvc,
		ArtistSvc: artistSvc,
		AlbumSvc:  albumSvc,
	}
}

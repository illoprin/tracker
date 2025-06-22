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
	playlistSvc := services.NewPlaylistService(
		repo.PlaylistsCol,
		repo.TracksCol,
		repo.UsersCol,
	)
	authSvc := services.NewAuthorizationService(
		repo.UsersCol,
		redisClient,
		playlistSvc,
	)
	userSvc := services.NewUserService(
		repo.UsersCol,
		flusher,
	)
	artistSvc := services.NewArtistService(
		repo.ArtistsCol,
		repo.AlbumsCol,
		repo.UsersCol,
		flusher,
		ownershipService,
	)
	albumSvc := services.NewAlbumService(
		repo.AlbumsCol,
		repo.TracksCol,
		repo.UsersCol,
		flusher,
		ownershipService,
	)
	trackSvc := services.NewTrackService(
		repo.TracksCol,
		repo.AlbumsCol,
		ownershipService,
	)
	searchSvc := services.NewSearchService(
		repo.TracksCol,
		repo.AlbumsCol,
		repo.ArtistsCol,
	)

	return &Dependencies{
		AuthSvc:     authSvc,
		UserSvc:     userSvc,
		ArtistSvc:   artistSvc,
		AlbumSvc:    albumSvc,
		TrackSvc:    trackSvc,
		SearchSvc:   searchSvc,
		PlaylistSvc: playlistSvc,
	}
}

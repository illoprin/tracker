package dependencies

import (
	"context"
	"os"
	"tracker-backend/internal/album"
	albumTracks "tracker-backend/internal/album/tracks"
	"tracker-backend/internal/app/repository"
	"tracker-backend/internal/artist"
	"tracker-backend/internal/auth/ownership"
	"tracker-backend/internal/config"
	"tracker-backend/internal/playlist"
	"tracker-backend/internal/track"
	"tracker-backend/internal/user"
)

type Dependencies struct {
	UserService        *user.UserService
	ArtistService      *artist.ArtistService
	TrackService       *track.TrackService
	AlbumService       *album.AlbumService
	AlbumTracksService *albumTracks.AlbumTracksService
	PlaylistService    *playlist.PlaylistService
}

func InitDependencies(
	ctx context.Context, repo *repository.Repository,
) *Dependencies {
	// WARN: init repository and create dependencies
	ownershipService := ownership.NewOwnershipService(
		repo.ArtistsCollection, repo.ArtistsCollection,
	)

	playlistService := playlist.NewPlaylistService(repo.PlaylistsCollection)
	userService := user.NewUserService(
		ctx, repo.UsersCollection,
		os.Getenv(config.JWTSecretEnvName),
		playlistService,
	)
	artistService := artist.NewArtistService(repo.ArtistsCollection)
	albumTracksService := albumTracks.NewAlbumTracksService(repo.TracksCollection, repo.AlbumsCollection, ownershipService)
	albumService := album.NewAlbumService(repo.AlbumsCollection, ownershipService)
	trackService := track.NewTrackService(repo.TracksCollection, ownershipService, albumService)

	return &Dependencies{
		PlaylistService:    playlistService,
		UserService:        userService,
		ArtistService:      artistService,
		AlbumTracksService: albumTracksService,
		TrackService:       trackService,
		AlbumService:       albumService,
	}
}

package dependencies

import (
	"context"
	"tracker-backend/internal/album"
	albumTracks "tracker-backend/internal/album/tracks"
	"tracker-backend/internal/app/repository"
	"tracker-backend/internal/artist"
	artistAlbums "tracker-backend/internal/artist/albums"
	"tracker-backend/internal/auth/ownership"
	"tracker-backend/internal/playlist"
	"tracker-backend/internal/track"
	"tracker-backend/internal/user"
)

type Dependencies struct {
	*user.UserService
	*artist.ArtistService
	*artistAlbums.ArtistAlbumsService
	*track.TrackService
	*album.AlbumService
	*albumTracks.AlbumTracksService
	*playlist.PlaylistService
}

func InitDependencies(
	ctx context.Context, repo *repository.Repository,
) *Dependencies {
	ownershipService := ownership.NewOwnershipService(
		repo.AlbumsCollection, repo.ArtistsCollection,
	)

	playlistService := playlist.NewPlaylistService(repo.PlaylistsCollection)
	userService := user.NewUserService(
		ctx, repo.UsersCollection,
		playlistService,
	)
	artistAlbumsService := artistAlbums.NewArtistAlbumsService(repo.AlbumsCollection, ownershipService)
	artistService := artist.NewArtistService(repo.ArtistsCollection)
	albumTracksService := albumTracks.NewAlbumTracksService(repo.TracksCollection, repo.AlbumsCollection, ownershipService)
	albumService := album.NewAlbumService(repo.AlbumsCollection, albumTracksService, ownershipService)
	trackService := track.NewTrackService(repo.TracksCollection, ownershipService, albumService)

	return &Dependencies{
		PlaylistService:     playlistService,
		UserService:         userService,
		ArtistAlbumsService: artistAlbumsService,
		ArtistService:       artistService,
		AlbumTracksService:  albumTracksService,
		TrackService:        trackService,
		AlbumService:        albumService,
	}
}

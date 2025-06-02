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

	artistCol := db.Collection("artists")
	if err := schemas.EnsureArtistIndices(ctx, artistCol); err != nil {
		return nil, err
	}

	albumCol := db.Collection("albums")
	if err := schemas.EnsureAlbumIndices(ctx, albumCol); err != nil {
		return nil, err
	}

	trackCol := db.Collection("tracks")
	if err := schemas.EnsureTrackIndices(ctx, trackCol); err != nil {
		return nil, err
	}

	playlistCol := db.Collection("playlists")
	if err := schemas.EnsurePlaylistIndices(ctx, playlistCol); err != nil {
		return nil, err
	}

	listenCol := db.Collection("listens")

	return &Repository{
		UsersCol:     userCol,
		ArtistsCol:   artistCol,
		AlbumsCol:    albumCol,
		TracksCol:    trackCol,
		PlaylistsCol: playlistCol,
		ListensCol:   listenCol,
	}, nil
}

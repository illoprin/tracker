package repository

import (
	"context"
	"tracker-backend/internal/album"
	"tracker-backend/internal/artist"
	"tracker-backend/internal/playlist"
	"tracker-backend/internal/track"
	"tracker-backend/internal/user"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// TODO: setup collection and ensure indices

type Repository struct {
	PlaylistsCollection *mongo.Collection
	UsersCollection     *mongo.Collection
	ArtistsCollection   *mongo.Collection
	AlbumsCollection    *mongo.Collection
	TracksCollection    *mongo.Collection
}

func MustInitRepository(ctx context.Context, db *mongo.Database) *Repository {

	playlistsCollection := db.Collection("playlists")
	usersCollection := db.Collection("users")
	artistsCollection := db.Collection("artists")
	albumsCollection := db.Collection("albums")
	tracksCollection := db.Collection("tracks")

	// ensure albums indices
	if err := album.EnsureIndexes(ctx, albumsCollection); err != nil {
		panic(err.Error())
	}
	// ensure tracks indices
	if err := track.EnsureIndexes(ctx, tracksCollection); err != nil {
		panic(err.Error())
	}
	// ensure tracks indices
	if err := artist.EnsureIndexes(ctx, artistsCollection); err != nil {
		panic(err.Error())
	}
	// ensure users indices
	if err := user.EnsureIndexes(ctx, usersCollection); err != nil {
		panic(err.Error())
	}
	// ensure playlists indices
	if err := playlist.EnsureIndexes(ctx, playlistsCollection); err != nil {
		panic(err.Error())
	}

	return &Repository{
		PlaylistsCollection: playlistsCollection,
		UsersCollection:     usersCollection,
		ArtistsCollection:   artistsCollection,
		AlbumsCollection:    albumsCollection,
		TracksCollection:    tracksCollection,
	}
}

package repository

import (
	"context"
	albumType "tracker-backend/internal/album/type"
	artistType "tracker-backend/internal/artist/type"
	playlistType "tracker-backend/internal/playlist/type"
	"tracker-backend/internal/track"
	userType "tracker-backend/internal/user/type"

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
	if err := albumType.EnsureIndexes(ctx, albumsCollection); err != nil {
		panic(err.Error())
	}
	// ensure tracks indices
	if err := track.EnsureIndexes(ctx, tracksCollection); err != nil {
		panic(err.Error())
	}
	// ensure tracks indices
	if err := artistType.EnsureIndexes(ctx, artistsCollection); err != nil {
		panic(err.Error())
	}
	// ensure users indices
	if err := userType.EnsureIndexes(ctx, usersCollection); err != nil {
		panic(err.Error())
	}
	// ensure playlists indices
	if err := playlistType.EnsureIndexes(ctx, playlistsCollection); err != nil {
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

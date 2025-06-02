package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repository struct {
	UsersCol     *mongo.Collection
	ArtistsCol   *mongo.Collection
	AlbumsCol    *mongo.Collection
	TracksCol    *mongo.Collection
	PlaylistsCol *mongo.Collection
	ListensCol   *mongo.Collection
}

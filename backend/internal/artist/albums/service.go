package artistAlbums

import (
	"context"
	"errors"
	albumType "tracker-backend/internal/album/type"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ArtistAlbumsService struct {
	albumsCol *mongo.Collection
}

func NewArtistAlbumsService(
	albumsCol *mongo.Collection,
) *ArtistAlbumsService {
	return &ArtistAlbumsService{
		albumsCol: albumsCol,
	}
}

func (s *ArtistAlbumsService) GetByArtistID(
	ctx context.Context, artistID string,
) ([]albumType.Album, error) {

	findOptions := options.Find().SetSort(bson.D{
		{Key: "year", Value: -1},
		{Key: "title", Value: -1},
	})

	cursor, err := s.albumsCol.Find(ctx, bson.M{"artistID": artistID}, findOptions)
	if err != nil {
		return nil, errors.New("failed to find albums")
	}
	defer cursor.Close(ctx)

	var albums []albumType.Album
	if err := cursor.All(ctx, &albums); err != nil {
		return nil, errors.New("failed to decode cursor")
	}
	return albums, nil
}

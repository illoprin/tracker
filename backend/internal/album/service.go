package album

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AlbumService struct {
	Col *mongo.Collection
}

func NewAlbumService(ctx context.Context, db *mongo.Database) *AlbumService {
	col := db.Collection("albums")

	return &AlbumService{
		Col: col,
	}
}

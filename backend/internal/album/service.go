package album

import (
	"context"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AlbumService struct {
	Col *mongo.Collection
}

func NewAlbumService(ctx context.Context, db *mongo.Database) *AlbumService {
	col := db.Collection("albums")

	err := EnsureIndexes(ctx, col)
	if err != nil {
		panic(err.Error())
	}

	return &AlbumService{
		Col: col,
	}
}

func Create(
	ctx context.Context,
	userID string,
	req *AlbumCreateRequest,
) (*Album, error) {
	return nil, nil
}

func UpdateCover(
	ctx context.Context,
	userID string,
	albumID string,
	file *multipart.File,
	fileHeader *multipart.FileHeader,
) (*Album, error) {
	return nil, nil
}

func Update(
	ctx context.Context,
	userID string,
	albumID string,
	req *AlbumUpdateRequest,
) (*Album, error) {
	return nil, nil
}

func Delete(ctx context.Context, userID string, albumID string) error {
	return nil
}

func GetByArtistID(ctx context.Context, albumID string) error {
	return nil
}

func GetByID(ctx context.Context, albumID string) (*Album, error) {
	return nil, nil
}

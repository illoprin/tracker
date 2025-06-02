package artistAlbums

import (
	"context"
	"errors"
	albumType "tracker-backend/internal/album/type"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/auth/ownership"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ArtistAlbumsService struct {
	albumsCol        *mongo.Collection
	ownershipService *ownership.OwnershipService
}

func NewArtistAlbumsService(
	albumsCol *mongo.Collection,
	ownershipSrv *ownership.OwnershipService,
) *ArtistAlbumsService {
	return &ArtistAlbumsService{
		albumsCol:        albumsCol,
		ownershipService: ownershipSrv,
	}
}

func (s *ArtistAlbumsService) GetByArtistID(
	ctx context.Context, artistID string, userID string, userRole int,
) ([]albumType.Album, error) {

	findOptions := options.Find().SetSort(bson.D{
		{Key: "year", Value: -1},
		{Key: "title", Value: -1},
	})

	filter := bson.M{
		"artistID": artistID,
	}
	isOwn, _ := s.ownershipService.IsArtistOwner(ctx, userID, artistID)

	// if user is not moderator -> show only moderated and public albums
	if userRole <= auth.RoleCustomer && !isOwn {
		filter["isHidden"] = false
		filter["status"] = albumType.StatusModerated
	}

	cursor, err := s.albumsCol.Find(ctx, filter, findOptions)
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

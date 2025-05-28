package album

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"
	"time"
	albumType "tracker-backend/internal/album/type"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/auth/ownership"
	uploadfile "tracker-backend/internal/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AlbumService struct {
	Col              *mongo.Collection
	trackChecker     TrackChecker
	ownershipService *ownership.OwnershipService
}

type TrackChecker interface {
	IsAnyTracksInAlbum(ctx context.Context, albumID string) (bool, error)
}

var (
	ErrTitleTaken = errors.New("album with this title already exists")
	ErrNotFound   = errors.New("album not found")
)

func NewAlbumService(
	albumsCol *mongo.Collection,
	trackChecker TrackChecker,
	ownershipService *ownership.OwnershipService,
) *AlbumService {

	return &AlbumService{
		Col:              albumsCol,
		trackChecker:     trackChecker,
		ownershipService: ownershipService,
	}
}

func (s *AlbumService) Create(
	ctx context.Context,
	userID string,
	req *albumType.AlbumCreateRequest,
) (*albumType.Album, error) {

	// validate artist ownership first (you need to implement this check)
	isOwn, err := s.ownershipService.IsArtistOwner(ctx, userID, req.ArtistID)
	if err != nil {
		return nil, errors.New("failed to check ownership")
	}
	if !isOwn {
		return nil, auth.ErrAccessDenied
	}

	// album is hidden by default
	album := &albumType.Album{
		ID:        uuid.NewString(),
		Title:     req.Title,
		ArtistID:  req.ArtistID,
		Year:      req.Year,
		Genres:    req.Genres,
		Status:    albumType.StatusHidden,
		CreatedAt: time.Now(),
	}

	_, err = s.Col.InsertOne(ctx, album)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrTitleTaken
		}
		return nil, err
	}

	return album, nil
}

func (s *AlbumService) UpdateCover(
	ctx context.Context,
	userID string,
	albumID string,
	file *multipart.File,
	fileHeader *multipart.FileHeader,
) (*albumType.Album, error) {

	isOwner, err := s.ownershipService.IsAlbumOwner(ctx, userID, albumID)
	if err != nil {
		return nil, errors.New("failed check album ownership")
	}
	if !isOwner {
		return nil, auth.ErrAccessDenied
	}

	filter := bson.M{"id": albumID}

	if err := uploadfile.ValidateFile(fileHeader, uploadfile.AllowedImageExtensions); err != nil {
		return nil, err
	}

	coverPath, err := uploadfile.UploadFile(
		fileHeader,
		file,
		"uploads/covers",
		uploadfile.AllowedImageExtensions,
	)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": bson.M{"coverPath": coverPath}}
	var updatedAlbum albumType.Album
	err = s.Col.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedAlbum)

	if err != nil {
		return nil, errors.New("album not found or not owned by user")
	}

	return &updatedAlbum, nil
}

func (s *AlbumService) Update(
	ctx context.Context,
	userID string,
	albumID string,
	req *albumType.AlbumUpdateRequest,
) (*albumType.Album, error) {
	isOwner, err := s.ownershipService.IsAlbumOwner(ctx, userID, albumID)
	if err != nil {
		return nil, errors.New("failed check album ownership")
	}
	if !isOwner {
		return nil, auth.ErrAccessDenied
	}

	filter := bson.M{"id": albumID}
	updates := bson.M{}

	if req.Year != nil {
		updates["year"] = *req.Year
	}
	if len(req.Genres) > 0 {
		updates["genres"] = req.Genres
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) <= 0 {
		return nil, err
	}

	// TODO: check the status being set
	// if album has no tracks - it cannot be published
	// if album is on moderation - the status cannot be changed
	// moderation is required for any changes in album

	// album change pipeline:
	// change/creation -> moderation -> publication

	var album albumType.Album
	err = s.Col.FindOneAndUpdate(ctx,
		filter, bson.M{"$set": updates},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&album)

	if err != nil {
		// check not found error
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		// check constraint error
		if writeErr, ok := err.(mongo.WriteException); ok {
			for _, we := range writeErr.WriteErrors {
				// login constraint error
				if strings.Contains(we.Message, "title") {
					return nil, ErrTitleTaken
				}
			}
		}
		return nil, errors.New("failed to update")
	}

	return &album, nil
}

func (s *AlbumService) GetByID(
	ctx context.Context, albumID string,
) (*albumType.Album, error) {
	var album albumType.Album
	err := s.Col.FindOne(ctx, bson.M{"id": albumID}).Decode(&album)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &album, nil
}

func (s *AlbumService) CheckExistence(ctx context.Context, albumID string) (bool, error) {
	if a, err := s.GetByID(ctx, albumID); a == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (s *AlbumService) Delete(
	ctx context.Context, userID string, albumID string,
) error {
	isOwn, err := s.ownershipService.IsAlbumOwner(ctx, userID, albumID)
	if err != nil {
		return errors.New("failed check album owner")
	}
	if !isOwn {
		return auth.ErrAccessDenied
	}
	res, err := s.Col.DeleteOne(ctx, bson.M{"id": albumID})
	if err != nil {
		return errors.New("failed to delete")
	}
	if res.DeletedCount < 1 {
		return ErrNotFound
	}
	return nil
}

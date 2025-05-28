package playlist

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"tracker-backend/internal/pkg/service"
	playlistType "tracker-backend/internal/playlist/type"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PlaylistService struct {
	Col *mongo.Collection
}

func NewPlaylistService(
	playlistCol *mongo.Collection,
) *PlaylistService {
	return &PlaylistService{
		Col: playlistCol,
	}
}

func (s *PlaylistService) Create(
	ctx context.Context, req playlistType.PlaylistCreateRequest,
) (*playlistType.Playlist, error) {
	// configure logger
	logger := slog.With(slog.String("function", "playlist.PlaylistService.Create"))

	playlist := &playlistType.Playlist{
		ID:        uuid.NewString(),
		Name:      req.Name,
		UserID:    req.UserID,
		IsDefault: req.IsDefault,
		IsPublic:  req.IsPublic,
		UpdatedAt: time.Now(),
	}

	_, err := s.Col.InsertOne(ctx, playlist)

	if err != nil {
		return nil, errors.New("failed to insert")
	}

	logger.Info("playlist created",
		slog.Group("info",
			slog.String("userID", req.UserID),
			slog.Bool("isDefault", req.IsDefault),
			slog.String("id", playlist.ID),
		),
	)

	return playlist, nil
}

func (s *PlaylistService) PushTrackLink(
	ctx context.Context, playlistID string, trackID string,
) (*playlistType.Playlist, error) {
	// configure logger
	logger := slog.With(slog.String("function", "playlist.PlaylistService.PushTrack"))

	update := bson.M{
		"$addToSet": bson.M{"trackIDs": trackID}, // Используем addToSet для избежания дубликатов
		"$set":      bson.M{"updatedAt": time.Now()},
	}

	// find and update playlist
	res := s.Col.FindOneAndUpdate(
		ctx,
		bson.M{"id": playlistID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		logger.Warn("failed to update the playlist", slog.String("error", res.Err().Error()))
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		return nil, errors.New("failed to update playlist")
	}

	// decode updated playlist
	var updatedPlaylist playlistType.Playlist
	err := res.Decode(&updatedPlaylist)
	if err != nil {
		logger.Warn("failed to decode updated playlist", slog.String("error", err.Error()))
		return nil, errors.New("failed to decode updated playlist")
	}

	// return updated playlist
	return &updatedPlaylist, nil

}

func (s *PlaylistService) RemoveTrackLink(
	ctx context.Context, playlistID string, trackID string,
) (*playlistType.Playlist, error) {
	// configure logger
	logger := slog.With(slog.String("function", "playlist.PlaylistService.PushTrack"))

	// define updates
	update := bson.M{
		"$pull": bson.M{"trackIDs": trackID},
		"$set":  bson.M{"updatedAt": time.Now()},
	}

	// find and update playlist
	res := s.Col.FindOneAndUpdate(
		ctx,
		bson.M{"id": playlistID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		logger.Warn("failed to update the playlist", slog.String("error", res.Err().Error()))
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		return nil, errors.New("failed to update playlist")
	}

	// decode updated playlist
	var updatedPlaylist playlistType.Playlist
	err := res.Decode(&updatedPlaylist)
	if err != nil {
		logger.Warn("failed to decode updated playlist", slog.String("error", err.Error()))
		return nil, errors.New("failed to decode updated playlist")
	}

	// return updated playlist
	return &updatedPlaylist, nil
}

// TODO: remove invalid links
/*
func (s *PlaylistService) RemoveInvalidLinks(ctx context.Context) error {
	// 1. Получаем все существующие ID треков
	trackCol := s.Col.Database().Collection("tracks")
	existingTrackIDs, err := getExistingTrackIDs(ctx, trackCol)
	if err != nil {
		return err
	}

	// 2. Находим плейлисты с несуществующими треками
	pipeline := []bson.M{
		{
			"$project": bson.M{
				"trackIDs": 1,
				"invalidTracks": bson.M{
					"$setDifference": bson.A{
						"$trackIDs",
						existingTrackIDs,
					},
				},
			},
		},
		{
			"$match": bson.M{
				"invalidTracks": bson.M{"$ne": bson.A{}},
			},
		},
	}

	cursor, err := s.Col.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// 3. Обновляем плейлисты, удаляя невалидные ссылки
	for cursor.Next(ctx) {
		var result struct {
			ID           string   `bson:"_id"`
			InvalidTracks []string `bson:"invalidTracks"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}

		_, err := s.Col.UpdateByID(
			ctx,
			result.ID,
			bson.M{"$pull": bson.M{"trackIDs": bson.M{"$in": result.InvalidTracks}}},
		)
		if err != nil {
			// Логируем ошибку, но продолжаем обработку
			continue
		}
	}

	return nil
}

func getExistingTrackIDs(ctx context.Context, col *mongo.Collection) (bson.A, error) {
	cursor, err := col.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids bson.A
	for cursor.Next(ctx) {
		var doc struct {
			ID string `bson:"id"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		ids = append(ids, doc.ID)
	}
	return ids, nil
}
*/

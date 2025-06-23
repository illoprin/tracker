package services

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"
	"tracker-backend/internal/infrastructure/redis"
	"tracker-backend/internal/pkg/logger"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GenreService struct {
	tracksCol   *mongo.Collection
	redisClient *redis.RedisClient
}

func NewGenreService(
	tracksCol *mongo.Collection,
	redisClient *redis.RedisClient,
) *GenreService {
	return &GenreService{
		tracksCol:   tracksCol,
		redisClient: redisClient,
	}
}

func getCacheKey(resourceType, resourceId string) string {
	return fmt.Sprintf("genre_model:%s:%s", resourceType, resourceId)
}

func (svc *GenreService) GetPopularGenres(
	ctx context.Context,
	resourceType, resourceId string,
) ([]string, error) {
	_logger := slog.With(
		slog.String("func", "services.GenreService.GetPopularGenres"),
		slog.String("resourceType", resourceType),
		slog.String("resourceID", resourceId),
	)

	// define cache key
	cacheKey := getCacheKey(resourceType, resourceId)

	// try to get genre model from redis
	var genreModel []string
	err := svc.redisClient.GetJSON(ctx, cacheKey, &genreModel)
	if err == nil {
		_logger.Debug("successfully get from redis",
			slog.Any("genres", genreModel),
		)
		return genreModel, nil
	}

	var pipeline []bson.M

	switch resourceType {
	case "album":
		pipeline = svc.getAlbumGenresPipeline(resourceId)
	case "artist":
		pipeline = svc.getArtistGenresPipeline(resourceId)
	case "playlist":
		pipeline = svc.getPlaylistGenresPipeline(resourceId)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	cursor, err := svc.tracksCol.Aggregate(ctx, pipeline)
	if err != nil {
		_logger.Error("Failed to aggregate genres", logger.ErrorAttr(err))
		return nil, fmt.Errorf("failed to aggregate genres: %w", err)
	}
	defer cursor.Close(ctx)

	var genres []struct {
		Genre string `bson:"_id"`
		Count int    `bson:"count"`
	}

	if err := cursor.All(ctx, &genres); err != nil {
		_logger.Error("Failed to decode genres", logger.ErrorAttr(err))
		return nil, fmt.Errorf("failed to decode genres: %w", err)
	}

	// Сортируем по убыванию популярности
	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Count > genres[j].Count
	})

	result := make([]string, len(genres))
	for i, g := range genres {
		result[i] = g.Genre
	}

	_logger.Debug("Successfully got popular genres",
		slog.Int("count", len(result)),
	)

	// set genres to redis
	err = svc.redisClient.SetJSON(ctx, cacheKey, result, time.Hour*5)
	if err != nil {
		_logger.Warn("failed to push redis", logger.ErrorAttr(err))
	}

	return result, nil
}

// Пайплайн для получения жанров альбома
func (svc *GenreService) getAlbumGenresPipeline(albumID string) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"albumId": albumID,
			},
		},
		{
			"$unwind": "$genres",
		},
		{
			"$group": bson.M{
				"_id":   "$genres",
				"count": bson.M{"$sum": 1},
			},
		},
	}
}

// Пайплайн для получения жанров артиста (через его альбомы)
func (svc *GenreService) getArtistGenresPipeline(artistID string) []bson.M {
	return []bson.M{
		{
			"$lookup": bson.M{
				"from": "albums",
				"let":  bson.M{"artistId": artistID},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []string{"$artistId", "$$artistId"},
							},
						},
					},
				},
				"as": "albums",
			},
		},
		{
			"$unwind": "$albums",
		},
		{
			"$match": bson.M{
				"albumId": "$albums.id",
			},
		},
		{
			"$unwind": "$genres",
		},
		{
			"$group": bson.M{
				"_id":   "$genres",
				"count": bson.M{"$sum": 1},
			},
		},
	}
}

// Пайплайн для получения жанров плейлиста
func (svc *GenreService) getPlaylistGenresPipeline(playlistID string) []bson.M {
	return []bson.M{
		{
			"$lookup": bson.M{
				"from": "playlists",
				"let":  bson.M{"playlistId": playlistID},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []string{"$id", "$$playlistId"},
							},
						},
					},
				},
				"as": "playlist",
			},
		},
		{
			"$unwind": "$playlist",
		},
		{
			"$match": bson.M{
				"id": bson.M{"$in": "$playlist.trackIds"},
			},
		},
		{
			"$unwind": "$genres",
		},
		{
			"$group": bson.M{
				"_id":   "$genres",
				"count": bson.M{"$sum": 1},
			},
		},
	}
}

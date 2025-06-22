package services

import (
	"context"
	"log/slog"
	"sync"
	"time"
	"tracker-backend/internal/domain/dtos"
	"tracker-backend/internal/domain/repository/schemas"
	"tracker-backend/internal/pkg/logger"
	"tracker-backend/internal/pkg/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SearchService struct {
	tracksCol  *mongo.Collection
	albumsCol  *mongo.Collection
	artistsCol *mongo.Collection
}

func NewSearchService(
	tracksCol, albumsCol, artistsCol *mongo.Collection,
) *SearchService {
	return &SearchService{
		tracksCol:  tracksCol,
		albumsCol:  albumsCol,
		artistsCol: artistsCol,
	}
}

func (svc *SearchService) GlobalSearch(
	ctx context.Context,
	query string,
	limitTracks int,
	generalLimit int,
) (*dtos.SearchResult, error) {
	_logger := slog.With(
		slog.String("func", "services.SearchService.GlobalSearch"),
		slog.String("query", query),
	)

	// Создаем контекст с таймаутом для всех запросов
	searchCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Запускаем все запросы параллельно
	var wg sync.WaitGroup
	wg.Add(3)

	var tracks []dtos.TrackResponse
	var albums []dtos.AlbumWithStats
	var artists []schemas.Artist

	var trackErr, albumErr, artistErr error

	// Поиск треков
	go func() {
		defer wg.Done()
		tracks, trackErr = svc.searchTracks(searchCtx, query, limitTracks)
	}()

	// Поиск альбомов
	go func() {
		defer wg.Done()
		albums, albumErr = svc.searchAlbums(searchCtx, query, generalLimit)
	}()

	// Поиск исполнителей
	go func() {
		defer wg.Done()
		artists, artistErr = svc.searchArtists(searchCtx, query, generalLimit)
	}()

	wg.Wait()

	// Обработка ошибок
	if trackErr != nil {
		_logger.Error("track search failed", logger.ErrorAttr(trackErr))
	}
	if albumErr != nil {
		_logger.Error("album search failed", logger.ErrorAttr(albumErr))
	}
	if artistErr != nil {
		_logger.Error("artist search failed", logger.ErrorAttr(artistErr))
	}

	// Если все запросы завершились с ошибкой
	if trackErr != nil && albumErr != nil && artistErr != nil {
		return nil, service.ErrInternal
	}

	return &dtos.SearchResult{
		Tracks:  tracks,
		Albums:  albums,
		Artists: artists,
	}, nil
}

func (svc *SearchService) searchTracks(
	ctx context.Context,
	query string,
	limit int,
) ([]dtos.TrackResponse, error) {

	matchTrack := bson.M{"name": bson.M{"$regex": query, "$options": "i"}}
	matchAlbum := bson.M{"album.isPublic": true, "album.isApproved": true}
	limitNode := bson.M{"$limit": limit}
	pipeline := schemas.GetTracksDetailsPipeline(matchTrack, matchAlbum)
	pipeline = append(pipeline, limitNode)

	slog.Debug("tracks aggregation pipeline", "pipeline", pipeline)

	cursor, err := svc.tracksCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		slog.Debug("small batch")
		return []dtos.TrackResponse{}, nil
	}

	var tracks []dtos.TrackResponse
	if err := cursor.All(ctx, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (svc *SearchService) searchAlbums(
	ctx context.Context,
	query string,
	limit int,
) ([]dtos.AlbumWithStats, error) {
	matchAlbum := bson.M{
		"isPublic":   true,
		"isApproved": true,
		"name":       bson.M{"$regex": query, "$options": "i"},
	}
	limitNode := bson.M{"$limit": limit}
	pipeline := schemas.GetAlbumsAggregationPipeline(matchAlbum)
	pipeline = append(pipeline, limitNode)

	cursor, err := svc.albumsCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		return []dtos.AlbumWithStats{}, nil
	}

	// decode result
	var result []struct {
		schemas.Album `bson:",inline"`
		TrackCount    int    `bson:"trackCount" json:"trackCount"`
		TotalDuration int    `bson:"totalDuration" json:"totalDuration"`
		ArtistName    string `bson:"artistName" json:"artistName"`
		ArtistAvatar  string `bson:"artistAvatar" json:"artistAvatar"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	albumsWithStats := make([]dtos.AlbumWithStats, 0, len(result))
	for _, a := range result {
		// create response struct
		aws := dtos.AlbumWithStats{
			Album:         a.Album,
			TrackCount:    a.TrackCount,
			TotalDuration: a.TotalDuration,
			ArtistName:    a.ArtistName,
			ArtistAvatar:  a.ArtistAvatar,
		}
		albumsWithStats = append(albumsWithStats, aws)
	}

	return albumsWithStats, nil
}

func (svc *SearchService) searchArtists(
	ctx context.Context,
	query string,
	limit int,
) ([]schemas.Artist, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"name": bson.M{"$regex": query, "$options": "i"},
			},
		},
		{
			"$limit": limit,
		},
	}

	cursor, err := svc.artistsCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		return []schemas.Artist{}, nil
	}

	var artists []schemas.Artist
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, err
	}

	return artists, nil
}

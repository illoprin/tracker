package services

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FlushService struct {
	tc  *mongo.Collection // tracks collection
	arc *mongo.Collection // artist collection
	alc *mongo.Collection // album collection
	pc  *mongo.Collection // playlist collection
}

func NewFlushService(
	tracksCol *mongo.Collection,
	artistCol *mongo.Collection,
	albumsCol *mongo.Collection,
	playlistCol *mongo.Collection,
) *FlushService {
	return &FlushService{
		tc:  tracksCol,
		arc: artistCol,
		alc: albumsCol,
		pc:  playlistCol,
	}
}

// FlushAlbumData deletes track (document and file)
func (svc *FlushService) FlushTrack(ctx context.Context, id string) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.FlushService.FlushTrack"))
	_ = _logger
	// TODO
	return nil
}

// FlushAlbumData deletes album related data (tracks, files)
func (svc *FlushService) FlushAlbumData(ctx context.Context, id string) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.FlushService.FlushAlbumData"))
	_ = _logger
	// TODO
	return nil
}

// FlushArtistData deletes artist related data (albums, tracks, files)
func (svc *FlushService) FlushArtistData(ctx context.Context, id string) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.FlushService.FlushArtistData"))
	_ = _logger
	// TODO
	return nil
}

// FlushArtistData deletes user related data (artists, albums, tracks, files)
func (svc *FlushService) FlushUserData(ctx context.Context, id string) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.FlushService.FlushUserData"))
	_ = _logger
	// TODO
	return nil
}

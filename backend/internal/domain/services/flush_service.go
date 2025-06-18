package services

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"tracker-backend/internal/config"
	"tracker-backend/internal/pkg/logger"
	"tracker-backend/internal/pkg/service"

	"go.mongodb.org/mongo-driver/v2/bson"
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

// FlushAlbumData deletes track file
func (svc *FlushService) FlushTrackData(ctx context.Context, audioFile string) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.FlushService.FlushTrack"))

	// delete file
	filePath := filepath.Join(os.Getenv(config.StaticDirEnvName), config.AudioDir, audioFile)
	if err := os.Remove(filePath); err != nil {
		_logger.Error("failed to delete track file by path",
			slog.String("filepath", filePath),
			logger.ErrorAttr(err),
		)
		return err
	}

	return nil
}

// FlushAlbumData deletes album related data (tracks, files)
func (svc *FlushService) FlushAlbumData(ctx context.Context, id string) error {
	// configure logger
	_logger := slog.With(
		slog.String("func", "services.FlushService.FlushAlbumData"),
		slog.String("albumId", id),
	)

	// find album
	var album struct {
		CoverPath string `bson:"coverPath"`
	}
	err := svc.alc.FindOne(ctx, bson.M{"id": id}).Decode(&album)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_logger.Error("album not found", logger.ErrorAttr(err))
			return service.ErrNotFound
		}
		_logger.Error("failed to find album", logger.ErrorAttr(err))
		return err
	}

	// find related tracks
	var tracks []struct {
		ID        string `bson:"id"`
		AudioFile string `bson:"audioFile"`
	}
	cur, err := svc.tc.Find(ctx, bson.M{"albumId": id})
	if err != nil {
		_logger.Error("failed to find tracks", logger.ErrorAttr(err))
		return err
	}
	defer cur.Close(ctx)

	// decode tracks to array
	if err := cur.All(ctx, &tracks); err != nil {
		_logger.Error("failed to decode tracks", logger.ErrorAttr(err))
		return err
	}

	if len(tracks) == 0 {
		_logger.Debug("nothing to delete")
		return nil
	}

	// delete tracks
	for _, track := range tracks {
		// delete file
		err := svc.FlushTrackData(ctx, track.AudioFile)
		if err != nil {
			return err
		}
		// delete document
		if _, err = svc.tc.DeleteOne(ctx, bson.M{"id": track.ID}); err != nil {
			_logger.Error("failed to delete track",
				slog.String("id", track.ID),
				logger.ErrorAttr(err),
			)
			return err
		}
	}

	// find artist to get files path
	if err := svc.alc.FindOne(ctx, bson.M{"id": id}).Decode(&album); err != nil {
		_logger.Error("failed to find album", logger.ErrorAttr(err))
		return err
	}

	// delete cover
	fileName := filepath.Base(album.CoverPath)
	filePath := filepath.Join(os.Getenv(config.StaticDirEnvName), config.CoversDir, fileName)
	if err := os.Remove(filePath); err != nil {
		_logger.Error("failed to delete album cover", logger.ErrorAttr(err))
		return err
	}

	return nil
}

// FlushArtistData deletes artist related data (albums, tracks, files)
func (svc *FlushService) FlushArtistData(ctx context.Context, id string) error {
	// configure logger
	_logger := slog.With(
		slog.String("func", "services.FlushService.FlushArtistData"),
		slog.String("artistId", id),
	)

	// find albums
	cur, err := svc.alc.Find(ctx, bson.M{"artistId": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_logger.Error("albums not found")
			return err
		}
		_logger.Error("failed to find albums",
			logger.ErrorAttr(err),
		)
		return err
	}
	defer cur.Close(ctx)

	// decode albums data
	var albums []struct {
		ID string `bson:"id"`
	}
	if err := cur.All(ctx, &albums); err != nil {
		_logger.Error("failed to decode albums", logger.ErrorAttr(err))
		return err
	}

	if len(albums) == 0 {
		_logger.Debug("nothing to delete")
		return nil
	}

	// flush albums data
	for _, album := range albums {
		err := svc.FlushAlbumData(ctx, album.ID)
		if err != nil {
			_logger.Error("failed to flush album data", logger.ErrorAttr(err))
			return err
		}
	}

	// flush albums
	albumsIds := make([]string, 0, len(albums))
	for _, album := range albums {
		albumsIds = append(albumsIds, album.ID)
	}
	_, err = svc.alc.DeleteMany(ctx, bson.M{"id": bson.M{"$in": albumsIds}})
	if err != nil {
		_logger.Error("failed to delete album documents", logger.ErrorAttr(err))
		return err
	}

	// find artist to get files path
	var artist struct {
		AvatarPath string `bson:"avatarPath"`
		BannerPath string `bson:"bannerPath"`
	}
	if err := svc.arc.FindOne(ctx, bson.M{"id": id}).Decode(&artist); err != nil {
		_logger.Error("failed to find artist", logger.ErrorAttr(err))
		return err
	}

	// delete avatar
	fileName := filepath.Base(artist.AvatarPath)
	filePath := filepath.Join(os.Getenv(config.StaticDirEnvName), config.AvatarsDir, fileName)
	if err := os.Remove(filePath); err != nil {
		_logger.Error("failed to delete avatar", logger.ErrorAttr(err))
		return err
	}

	// delete banner
	fileName = filepath.Base(artist.BannerPath)
	filePath = filepath.Join(os.Getenv(config.StaticDirEnvName), config.BannersDir, fileName)
	if err := os.Remove(filePath); err != nil {
		_logger.Error("failed to delete avatar", logger.ErrorAttr(err))
		return err
	}

	return nil
}

func (svc *FlushService) FlushPlaylistData(ctx context.Context, coverPath string) error {
	_logger := slog.With(
		slog.String("func", "services.FlushService.FlushPlaylistData"),
	)
	_ = _logger
	return nil
}

func (svc *FlushService) flushUserArtists(ctx context.Context, id string) error {
	// configure logger
	_logger := slog.With(
		slog.String("func", "services.FlushService.flushUserArtists"),
		slog.String("userId", id),
	)

	// find all artists
	cur, err := svc.arc.Find(ctx, bson.M{"ownerId": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_logger.Error("artists not found")
			return err
		}
		_logger.Error("failed to find artists",
			logger.ErrorAttr(err),
		)
		return err
	}
	defer cur.Close(ctx)

	// decode artists
	var artists []struct {
		ID string `bson:"id"`
	}
	if err := cur.All(ctx, &artists); err != nil {
		_logger.Error("failed to decode artists", logger.ErrorAttr(err))
		return err
	}

	if len(artists) == 0 {
		_logger.Debug("user has no artists")
		return nil
	}

	// delete artists data
	for _, artist := range artists {
		err = svc.FlushArtistData(ctx, artist.ID)
		if err != nil {
			return err
		}
	}

	// delete artist documents
	artistIds := make([]string, 0, len(artists))
	for _, artist := range artists {
		artistIds = append(artistIds, artist.ID)
	}
	_, err = svc.arc.DeleteMany(ctx, bson.M{"id": bson.M{"$in": artistIds}})
	if err != nil {
		_logger.Error("failed to delete artits documents", logger.ErrorAttr(err))
		return err
	}

	return nil
}

func (svc *FlushService) flushUserPlaylists(ctx context.Context, id string) error {
	// configure logger
	_logger := slog.With(
		slog.String("func", "services.FlushService.flushUserPlaylists"),
		slog.String("userId", id),
	)

	// find playlists
	cur, err := svc.pc.Find(ctx, bson.M{"ownerId": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_logger.Error("artists not found")
			return err
		}
		_logger.Error("failed to find artists",
			logger.ErrorAttr(err),
		)
		return err
	}
	defer cur.Close(ctx)

	// decode playlists
	var playlists []struct {
		ID        string `bson:"id"`
		CoverPath string `bson:"coverPath"`
	}
	if err = cur.All(ctx, &playlists); err != nil {
		_logger.Error("failed to decode playlists cursor", logger.ErrorAttr(err))
		return err
	}

	if len(playlists) == 0 {
		_logger.Debug("nothing to delete")
		return nil
	}

	// delete playlits data
	for _, playlist := range playlists {
		// playlist may not have cover
		if playlist.CoverPath != "" {
			filename := filepath.Base(playlist.CoverPath)
			filePath := filepath.Join(os.Getenv(config.StaticDirEnvName), config.CoversDir, filename)
			if err := os.Remove(filePath); err != nil {
				_logger.Error("failed to delete playlist cover", logger.ErrorAttr(err))
			}
		}
	}

	// delete playlists
	playlistsIDs := make([]string, 0, len(playlists))
	for _, playlist := range playlists {
		playlistsIDs = append(playlistsIDs, playlist.ID)
	}
	_, err = svc.pc.DeleteMany(ctx, bson.M{"id": bson.M{"$in": playlistsIDs}})
	if err != nil {
		_logger.Error("failed to delete playlists documents", logger.ErrorAttr(err))
		return err
	}

	return nil
}

// FlushArtistData deletes user related data (artists, albums, tracks, files)
func (svc *FlushService) FlushUserData(ctx context.Context, id string) error {

	// flush artists
	err := svc.flushUserArtists(ctx, id)
	if err != nil {
		return err
	}

	// flush playlists
	err = svc.flushUserPlaylists(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"tracker-backend/internal/config"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/infrastructure/storage"
	"tracker-backend/internal/interfaces/rest/middleware"
	"tracker-backend/internal/interfaces/rest/utils/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type TrackHandler struct {
	tSvc *services.TrackService
}

func NewTrackHandler(t *services.TrackService) *TrackHandler {
	return &TrackHandler{
		tSvc: t,
	}
}

func (h *TrackHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	userRole := ctx.Value(middleware.UserRoleKey).(int)
	trackId := chi.URLParam(r, "id")

	// execute service function
	audioFile, err := h.tSvc.GetAudioPathByID(ctx, userId, userRole, trackId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else if errors.Is(err, service.ErrForbidden) {
			render.Status(r, http.StatusForbidden)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	filePath := filepath.Join(os.Getenv(config.StaticDirEnvName), config.AudioDir, audioFile)

	// check file existence
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get file stats"))
		return
	}

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to open file"))
		return
	}
	defer file.Close()

	// set headers
	contentType := storage.GetAudioContentTypeByExtension(filepath.Ext(audioFile))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// send whole file
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		// log the error, but do not send the response
		// because the headers have already been sent.
		slog.Error("error streaming file", slog.String("error", err.Error()))
	}
}

func (h *TrackHandler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	userRole := ctx.Value(middleware.UserRoleKey).(int)
	trackId := chi.URLParam(r, "id")

	// execute service function
	a, err := h.tSvc.GetMetadataByID(ctx, userId, userRole, trackId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else if errors.Is(err, service.ErrForbidden) {
			render.Status(r, http.StatusForbidden)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return content
	render.JSON(w, r, a)
}

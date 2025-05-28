package track

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"tracker-backend/internal/auth"
	uploadfile "tracker-backend/internal/pkg/file"
	"tracker-backend/internal/pkg/response"
	"tracker-backend/internal/pkg/service"
)

type TrackHandler struct {
	service   *TrackService
	validator *validator.Validate
}

// NewTrackHandler создает новый обработчик для треков
func NewTrackHandler(service *TrackService) *TrackHandler {
	return &TrackHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateTrack обрабатывает загрузку нового трека
func (h *TrackHandler) Create(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)

	// parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB максимум
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse multipart form"))
		return
	}

	// create request from form data
	req := &CreateTrackRequest{
		Title:   r.FormValue("title"),
		Genre:   strings.Split(r.FormValue("title"), ","),
		AlbumID: r.FormValue("albumId"),
	}

	// validate request
	if err := h.validator.Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r,
			response.ValidationErrorsResp(err.(validator.ValidationErrors)),
		)
		return
	}

	// get audio file
	audioFile, fileHeader, err := r.FormFile("audio")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("audio file required"))
		return
	}
	defer audioFile.Close()

	// create track document and save file
	track, err := h.service.Create(ctx, userID, req, &audioFile, fileHeader)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, track.ToResponse())
}

// GetTrack получает информацию о треке по ID
func (h *TrackHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	trackID := chi.URLParam(r, "id")
	if trackID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("missing track id"))
		return
	}

	track, err := h.service.GetByID(ctx, trackID)
	if err != nil {
		if err == service.ErrNotFound {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("track not found"))
		} else {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to get track"))
		}
		return
	}

	render.JSON(w, r, track.ToResponse())
}

// StreamTrack handle track stream
func (h *TrackHandler) StreamTrack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	trackID := chi.URLParam(r, "id")
	if trackID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("missing track id"))
		return
	}

	// get file path
	filePath, err := h.service.GetFilePathByID(ctx, trackID)
	if err != nil {
		if err == service.ErrNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("track not found"))
		} else {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to get track file"))
		}
		return
	}

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to open track file"))
		return
	}
	defer file.Close()

	// getting info about file
	fileInfo, err := file.Stat()
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to get file info"))
		return
	}

	// determine MIME type by extension
	ext := strings.ToLower(filepath.Ext(filePath))
	contentType := uploadfile.GetAudioContentTypeByExtension(ext)

	// set necessary headers
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	w.Header().Set("Accept-Ranges", "bytes")

	// handling Range Requests to Support Partial Loading
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		h.handleRangeRequest(w, r, file, fileInfo.Size(), rangeHeader)
		return
	}

	// send whole file if range is not defined
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		// log the error, but do not send the response
		// because the headers have already been sent.
		slog.Error("error streaming file", slog.String("error", err.Error()))
	}
}

func (h *TrackHandler) Update(w http.ResponseWriter, r *http.Request) {
	// update name, genre or file
	// if file updated - delete old one
	render.Status(r, http.StatusNotImplemented)
}

func (h *TrackHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)
	trackID := chi.URLParam(r, "id")

	err := h.service.Delete(ctx, trackID, userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
	}

	render.Status(r, http.StatusNoContent)
}

// handleRangeRequest обрабатывает запросы с заголовком Range
func (h *TrackHandler) handleRangeRequest(w http.ResponseWriter, r *http.Request, file *os.File, fileSize int64, rangeHeader string) {
	// parse Range header (example: "bytes=0-1023")
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	ranges := strings.TrimPrefix(rangeHeader, "bytes=")
	rangeParts := strings.Split(ranges, "-")

	if len(rangeParts) != 2 {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	var start, end int64
	var err error

	// parse start byte position
	if rangeParts[0] != "" {
		start, err = strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
	}

	// parse end byte position
	if rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
	} else {
		end = fileSize - 1
	}

	// check correction of Range
	if start < 0 || end >= fileSize || start > end {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// set read position on file
	_, err = file.Seek(start, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set headers
	contentLength := end - start + 1
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.WriteHeader(http.StatusPartialContent)

	// send the requested part of the file
	_, err = io.CopyN(w, file, contentLength)
	if err != nil {
		fmt.Printf("Error streaming partial file: %v\n", err)
	}
}

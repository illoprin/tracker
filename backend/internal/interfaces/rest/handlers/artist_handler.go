package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/infrastructure/storage"
	"tracker-backend/internal/interfaces/rest/middleware"
	"tracker-backend/internal/interfaces/rest/utils/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ArtistHandler struct {
	aSvc *services.ArtistService
}

func NewArtistHandler(s *services.ArtistService) *ArtistHandler {
	return &ArtistHandler{
		aSvc: s,
	}
}

func (h *ArtistHandler) Create(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)

	err := r.ParseMultipartForm(storage.MaxFormSize << 20) // 32 MB
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse multipart form"))
		return
	}

	// parse form
	name := r.FormValue("name")
	if len(name) <= 3 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("artist name required"))
		return
	}
	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("avatar image required"))
		return
	}

	// validate file
	err = storage.ValidateFile(fileHeader, storage.AllowedImageExtensions)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// execute service function
	a, err := h.aSvc.Create(ctx, userId, name, file, fileHeader)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			render.Status(r, http.StatusInternalServerError)
		} else if errors.Is(err, service.ErrExists) {
			render.Status(r, http.StatusConflict)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return response
	render.JSON(w, r, a)
}

func (h *ArtistHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)

	// execute service function
	all, err := h.aSvc.GetByUserID(ctx, userId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return response
	render.JSON(w, r, all)
}

func (h *ArtistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("url param required"))
		return
	}

	// execute service function
	a, err := h.aSvc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, a)
}

func (h *ArtistHandler) GetAlbums(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)

	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("url param required"))
		return
	}

	a, err := h.aSvc.GetAlbums(ctx, userId, id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, a)
}

func (h *ArtistHandler) GetPopularTracks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	_ = limit

	// TODO
	render.Status(r, http.StatusNotImplemented)
}

func (h *ArtistHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}

func (h *ArtistHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("url param required"))
		return
	}

	// execute service function
	err := h.aSvc.DeleteByID(ctx, userId, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return status
	render.Status(r, http.StatusNoContent)
}

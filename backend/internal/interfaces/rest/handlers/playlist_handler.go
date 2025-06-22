package handlers

import (
	"errors"
	"net/http"
	"tracker-backend/internal/domain/dtos"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/infrastructure/storage"
	"tracker-backend/internal/interfaces/rest/middleware"
	"tracker-backend/internal/interfaces/rest/utils/request"
	"tracker-backend/internal/interfaces/rest/utils/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type PlaylistHandler struct {
	pSvc *services.PlaylistService
	v    *validator.Validate
}

func NewPlaylistHandler(svc *services.PlaylistService) *PlaylistHandler {
	v := validator.New()
	return &PlaylistHandler{
		pSvc: svc,
		v:    v,
	}
}

func (h *PlaylistHandler) Create(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)

	// parse form
	err := r.ParseMultipartForm(storage.MaxFormSize)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse multipart form"))
		return
	}

	// parse form data
	req := dtos.PlaylistCreateRequest{
		Name:        r.FormValue("name"),
		IsDefault:   false,
		Description: r.FormValue("description"),
		IsPublic:    r.FormValue("isPublic") == "1",
	}

	// parse form file
	hasCover := true
	file, fileHeader, err := r.FormFile("cover")
	if err != nil {
		hasCover = false
	}
	// validate file if it has
	if hasCover {
		if err := storage.ValidateFile(fileHeader, storage.AllowedImageExtensions); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
	}

	// validate body
	if !request.ValidateBody(w, r, h.v, req) {
		return
	}

	a, err := h.pSvc.Create(ctx, userId, req, &file, fileHeader, hasCover)
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

func (h *PlaylistHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	a, err := h.pSvc.GetMetadata(ctx, userId, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return response
	render.JSON(w, r, a)

}

func (h *PlaylistHandler) GetTracks(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")
	_ = userId
	_ = id

}

func (h *PlaylistHandler) Update(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	_ = userId
	_ = id

}

func (h *PlaylistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	_ = userId
	_ = id

}

func (h *PlaylistHandler) AddTrack(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")
	trackId := chi.URLParam(r, "trackId")

	_ = userId
	_ = id
	_ = trackId

}

func (h *PlaylistHandler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")
	trackId := chi.URLParam(r, "trackId")

	_ = userId
	_ = id
	_ = trackId

}

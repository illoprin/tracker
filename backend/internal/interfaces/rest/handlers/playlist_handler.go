package handlers

import (
	"errors"
	"net/http"
	"strconv"
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
	rSvc *services.RecommendationsService
	v    *validator.Validate
}

func NewPlaylistHandler(
	svc *services.PlaylistService,
	rcmSvc *services.RecommendationsService,
) *PlaylistHandler {
	v := validator.New()
	return &PlaylistHandler{
		pSvc: svc,
		rSvc: rcmSvc,
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

	if hasCover {
		// validate form file
		if err := storage.ValidateFile(fileHeader, storage.AllowedImageExtensions); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
		}
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

func (h *PlaylistHandler) My(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)

	// execute service function
	p, err := h.pSvc.GetMy(ctx, userId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return response
	render.JSON(w, r, p)
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
	p, err := h.pSvc.GetMetadata(ctx, userId, id)
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
	render.JSON(w, r, p)
}

func (h *PlaylistHandler) GetTracks(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")
	limitStr := chi.URLParam(r, "limit")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// execute service function
	tracks, err := h.pSvc.GetTracks(ctx, userId, id, limit)
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
	render.JSON(w, r, tracks)
}

func (h *PlaylistHandler) GetWave(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	playlistId := chi.URLParam(r, "id")

	// get limit query param
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	// get page query param
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	tracks, err := h.rSvc.GetForResource(ctx, "playlist", playlistId, limit, page, []string{})
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, tracks)
}

func (h *PlaylistHandler) Update(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	// parse form
	err := r.ParseMultipartForm(storage.MaxFormSize)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse multipart form"))
		return
	}

	// parse form data
	req := dtos.PlaylistUpdateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		IsPublic:    r.FormValue("isPublic"),
	}

	// parse form file
	hasCover := true
	file, fileHeader, err := r.FormFile("cover")
	if err != nil {
		hasCover = false
	}

	if hasCover {
		// validate form file
		if err := storage.ValidateFile(fileHeader, storage.AllowedImageExtensions); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
	}

	// validate file if it needed
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

	a, err := h.pSvc.Update(ctx, userId, id, req, &file, fileHeader, hasCover)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			render.Status(r, http.StatusInternalServerError)
		} else if errors.Is(err, service.ErrExists) {
			render.Status(r, http.StatusConflict)
		} else if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return response
	render.JSON(w, r, a)
}

func (h *PlaylistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	// execute service function
	err := h.pSvc.Delete(ctx, userId, id)
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

func (h *PlaylistHandler) AddTrack(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")
	trackId := chi.URLParam(r, "trackId")

	err := h.pSvc.AddTrack(ctx, userId, id, trackId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else if errors.Is(err, service.ErrForbidden) {
			render.Status(r, http.StatusMethodNotAllowed)
			render.JSON(w, r, response.Error("you can't insert this track into the playlist"))
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *PlaylistHandler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")
	trackId := chi.URLParam(r, "trackId")

	err := h.pSvc.RemoveTrack(ctx, userId, id, trackId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusNoContent)
}

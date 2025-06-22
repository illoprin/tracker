package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
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

type AlbumHandler struct {
	aSvc *services.AlbumService
	tSvc *services.TrackService
	v    *validator.Validate
}

func NewAlbumHandler(s *services.AlbumService, t *services.TrackService) *AlbumHandler {
	v := validator.New()
	v.RegisterValidation("year", dtos.ValidateYear)
	v.RegisterValidation("type", dtos.ValidateType)
	v.RegisterValidation("genres", dtos.ValidateGenres)
	v.RegisterValidation("status", dtos.ValidateAlbumStatus)
	return &AlbumHandler{
		aSvc: s,
		tSvc: t,
		v:    v,
	}
}

func (h *AlbumHandler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	// start parse form
	err := r.ParseMultipartForm(storage.MaxFormSize)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to parse multipart form"))
		return
	}

	// parse body
	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse duration"))
		return
	}
	req := dtos.TrackCreateRequest{
		Name:     r.FormValue("name"),
		AlbumID:  id,
		Duration: duration,
		Genres:   strings.Split(r.FormValue("genres"), ","),
	}

	// validate body
	if !request.ValidateBody(w, r, h.v, req) {
		return
	}

	// get form file
	f, fh, err := r.FormFile("audio")
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to parse form file"))
		return
	}

	if err := storage.ValidateFile(fh, storage.AllowedAudioExtensions); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// execute service function
	t, err := h.tSvc.PushTrack(ctx, userId, req, &f, fh)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			render.Status(r, http.StatusForbidden)
		} else if errors.Is(err, service.ErrExists) {
			render.Status(r, http.StatusConflict)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}

		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return result
	render.JSON(w, r, t)
}

func (h *AlbumHandler) GetTracksByAlbumID(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	userRole := ctx.Value(middleware.UserRoleKey).(int)
	id := chi.URLParam(r, "id")

	// execute service function
	tracks, err := h.aSvc.GetTracks(ctx, userId, userRole, id)
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

func (h *AlbumHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	userRole := ctx.Value(middleware.UserRoleKey).(int)
	id := chi.URLParam(r, "id")

	// execute service function
	a, err := h.aSvc.GetByID(ctx, userId, userRole, id)
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

func (h *AlbumHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	albumId := chi.URLParam(r, "id")

	// parse form
	err := r.ParseMultipartForm(storage.MaxFormSize)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to parse form"))
		return
	}

	// decode form
	req := dtos.AlbumUpdateRequest{
		Name: r.FormValue("name"),
		Year: r.FormValue("year"),
		Type: r.FormValue("type"),
	}

	// validate request
	if !request.ValidateBody(w, r, h.v, req) {
		return
	}

	// get form file
	hasCover := true
	cover, coverHeader, err := r.FormFile("cover")
	if err != nil {
		hasCover = false
	}

	// validate form file
	if hasCover {
		if err := storage.ValidateFile(coverHeader, storage.AllowedImageExtensions); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
		}
	}

	err = h.aSvc.UpdateByID(
		ctx,
		userId,
		albumId,
		req,
		&cover,
		coverHeader,
		hasCover,
	)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else if errors.Is(err, service.ErrExists) {
			render.Status(r, http.StatusConflict)
		} else if errors.Is(err, service.ErrForbidden) {
			render.Status(r, http.StatusForbidden)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *AlbumHandler) Publish(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	albumId := chi.URLParam(r, "id")

	err := h.aSvc.Publish(ctx, userId, albumId)
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
	render.Status(r, http.StatusNoContent)
}

func (h *AlbumHandler) GetUnapproved(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()

	// get url params (limit, search, artistId)
	artistId := r.URL.Query().Get("artistId")
	query := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")

	// parse limit string
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// execute service function
	als, err := h.aSvc.GetUnapproved(ctx, limit, artistId, query)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return result
	render.JSON(w, r, als)
}

func (h *AlbumHandler) Moderate(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	albumId := chi.URLParam(r, "id")

	// decode json body
	var req dtos.AlbumModerationRequest
	if !request.DecodeJSONBody(w, r, &req) {
		return
	}

	// validate body
	if !request.ValidateBody(w, r, h.v, req) {
		return
	}

	// execute service function
	err := h.aSvc.Moderate(ctx, albumId, req)
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

func (h *AlbumHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

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

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

type ArtistHandler struct {
	aSvc  *services.ArtistService
	alSvc *services.AlbumService
	v     *validator.Validate
}

func NewArtistHandler(
	artistSvc *services.ArtistService,
	albumSvc *services.AlbumService,
) *ArtistHandler {
	v := validator.New()
	v.RegisterValidation("year", dtos.ValidateYear)
	v.RegisterValidation("type", dtos.ValidateType)
	v.RegisterValidation("genres", dtos.ValidateGenres)

	return &ArtistHandler{
		alSvc: albumSvc,
		aSvc:  artistSvc,
		v:     v,
	}
}

func (h *ArtistHandler) Create(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)

	err := r.ParseMultipartForm(storage.MaxFormSize)
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

	// parse avatar avatar
	avatar, avatarHeader, err := r.FormFile("avatar")
	hasAvatar := true
	if err != nil {
		hasAvatar = false
	}

	// validate file
	if hasAvatar {
		err = storage.ValidateFile(avatarHeader, storage.AllowedImageExtensions)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
	}

	// parse banner file
	banner, bannerHeader, err := r.FormFile("banner")
	hasBanner := true
	if err != nil {
		hasBanner = false
	}
	if hasBanner {
		err = storage.ValidateFile(bannerHeader, storage.AllowedImageExtensions)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
	}

	// execute service function
	a, err := h.aSvc.Create(
		ctx,
		userId,
		name,
		avatar,
		avatarHeader,
		hasAvatar,
		banner,
		bannerHeader,
		hasBanner,
	)
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

func (h *ArtistHandler) PushAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	artistId := chi.URLParam(r, "id")

	// parse form
	err := r.ParseMultipartForm(storage.MaxFormSize)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse multipart form"))
		return
	}

	// parse form data
	req := dtos.AlbumCreateRequest{
		ArtistID: artistId,
		Name:     r.FormValue("name"),
		Year:     r.FormValue("year"),
		Type:     r.FormValue("type"),
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

	a, err := h.alSvc.Create(ctx, userId, req, file, fileHeader, hasCover)
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

func (h *ArtistHandler) Like(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	// execute service function
	err := h.aSvc.Like(ctx, userId, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return result
	render.Status(r, http.StatusNoContent)
}

func (h *ArtistHandler) GetLiked(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)

	// execute service function
	a, err := h.aSvc.GetLiked(ctx, userId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return result
	render.JSON(w, r, a)
}

func (h *ArtistHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
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

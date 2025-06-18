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

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type AlbumHandler struct {
	aSvc *services.AlbumService
	v    *validator.Validate
}

func NewAlbumHandler(s *services.AlbumService) *AlbumHandler {
	v := validator.New()
	return &AlbumHandler{
		aSvc: s,
		v:    v,
	}
}

func (h *AlbumHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	req := dtos.AlbumCreateRequest{
		ArtistID: r.FormValue("artistId"),
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

	a, err := h.aSvc.Create(ctx, userId, req, file, fileHeader, hasCover)
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
func (h *AlbumHandler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}

func (h *AlbumHandler) GetTracksByAlbumID(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}
func (h *AlbumHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}
func (h *AlbumHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}
func (h *AlbumHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)

}

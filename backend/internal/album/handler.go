package album

import (
	"net/http"
	albumType "tracker-backend/internal/album/type"
	"tracker-backend/internal/auth"
	genreType "tracker-backend/internal/genre/type"
	"tracker-backend/internal/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type AlbumHandler struct {
	Service   *AlbumService
	Validator *validator.Validate
}

func NewAlbumHandler(s *AlbumService) *AlbumHandler {
	v := validator.New()
	v.RegisterValidation("status", albumType.ValidateStatus)
	v.RegisterValidation("year", albumType.ValidateYear)
	v.RegisterValidation("genres", genreType.ValidateGenres)

	return &AlbumHandler{
		Service:   s,
		Validator: v,
	}
}

func (h *AlbumHandler) Create(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)

	// decode json
	var req albumType.AlbumCreateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse form data"))
		return
	}

	// validate request
	if err := h.Validator.Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ValidationErrorsResp(err.(validator.ValidationErrors)))
		return
	}

	// execute service function
	album, err := h.Service.Create(ctx, userID, &req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// send response
	render.JSON(w, r, album.ToResponse())
}

func (h *AlbumHandler) Update(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)
	albumID := chi.URLParam(r, "id")

	// decode json
	var req albumType.AlbumUpdateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse form data"))
		return
	}

	// validate request
	if err := h.Validator.Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ValidationErrorsResp(err.(validator.ValidationErrors)))
		return
	}

	// execute service function
	album, err := h.Service.Update(ctx, userID, albumID, &req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// send response
	render.JSON(w, r, album.ToResponse())
}

func (h *AlbumHandler) UpdateCover(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)
	albumID := chi.URLParam(r, "id")

	// decode form data file
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("cover file required"))
		return
	}
	defer file.Close()

	// execute service function
	album, err := h.Service.UpdateCover(ctx, userID, albumID, &file, fileHeader)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	// send response
	render.JSON(w, r, album.ToResponse())
}

func (h *AlbumHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	albumID := chi.URLParam(r, "id")

	// execute service function
	album, err := h.Service.GetByID(ctx, albumID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error("album not found"))
		return
	}

	// send response
	render.JSON(w, r, album.ToResponse())
}

func (h *AlbumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)
	albumID := chi.URLParam(r, "id")

	// execute service function
	err := h.Service.Delete(ctx, userID, albumID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// send response
	render.Status(r, http.StatusNoContent)
}

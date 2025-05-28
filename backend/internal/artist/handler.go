package artist

import (
	"errors"
	"net/http"
	artistType "tracker-backend/internal/artist/type"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/pkg/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ArtistHandler struct {
	Service   *ArtistService
	Validator *validator.Validate
}

func NewArtistHandler(service *ArtistService) *ArtistHandler {
	v := validator.New()
	return &ArtistHandler{
		Service:   service,
		Validator: v,
	}
}

func (h *ArtistHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)

	var req artistType.CreateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ValidationErrorsResp(err.(validator.ValidationErrors)))
		return
	}

	artist, err := h.Service.Create(ctx, userID, req)
	if err != nil {
		if errors.Is(err, ErrNameTaken) {
			render.Status(r, http.StatusConflict)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, artist)
}

func (h *ArtistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)

	// get url ID param
	artistID := chi.URLParam(r, "id")
	if artistID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("missing artist id"))
		return
	}

	if err := h.Service.Delete(ctx, artistID, userID); err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *ArtistHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)

	// get url ID param
	artistID := chi.URLParam(r, "id")
	if artistID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("missing artist id"))
		return
	}

	var req artistType.UpdateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}

	if req.AvatarPath != nil {
		render.Status(r, http.StatusMethodNotAllowed)
		render.JSON(w, r, response.Error("use PUT /api/artist/{id}/avatar to update artist avatar"))
		return
	}

	artist, err := h.Service.Update(ctx, artistID, userID, req)
	// check update results
	if err != nil {
		if errors.Is(err, ErrNameTaken) {
			render.Status(r, http.StatusConflict)
		} else if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// send updated artist
	render.JSON(w, r, artist)
}

func (h *ArtistHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)

	// get url ID param
	artistID := chi.URLParam(r, "id")
	if artistID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("missing artist id"))
		return
	}

	// parse form data
	err := r.ParseMultipartForm(5 << 20) // 5MB
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse multipart form"))
	}

	// extract form file
	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to extract file"))
		return
	}

	// call update function
	artist, err := h.Service.UpdateAvatar(
		ctx, userID, artistID, &file, fileHeader,
	)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// send updated artist
	render.JSON(w, r, artist)
}

func (h *ArtistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get url ID param
	artistID := chi.URLParam(r, "id")
	if artistID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("missing artist id"))
		return
	}

	artist, err := h.Service.GetByID(ctx, artistID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, artist)
}

func (h *ArtistHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(auth.UserIDKey).(string)

	artists, err := h.Service.GetByUserID(ctx, userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to fetch artists"))
		return
	}

	render.JSON(w, r, artists)
}

package user

import (
	"errors"
	"net/http"
	"tracker-backend/internal/auth"
	"tracker-backend/internal/pkg/response"
	"tracker-backend/internal/pkg/service"
	userType "tracker-backend/internal/user/type"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Service   *UserService
	Validator *validator.Validate
}

func NewUserHandler(s *UserService) *UserHandler {
	v := validator.New()
	return &UserHandler{
		Service:   s,
		Validator: v,
	}
}

// POST /user
func (h *UserHandler) Register(
	w http.ResponseWriter, r *http.Request,
) {
	var req userType.RegisterRequest
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

	user, err := h.Service.Register(r.Context(), req)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user.ToResponse())
}

// POST /user/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userType.LoginRequest
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

	token, err := h.Service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else if errors.Is(err, service.ErrAccessDenied) {
			render.Status(r, http.StatusForbidden)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, map[string]string{"token": token})
}

// GET /user
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	user, err := h.Service.GetByID(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error("user not found"))
		return
	}

	// send response
	render.JSON(w, r, user.ToResponse())
}

// PUT /user
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var req userType.UpdateRequest
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

	// get 'Allow-Access' header
	aaHeader := r.Header.Get("Allow-Access")
	allowed := false
	if aaHeader == "1" {
		allowed = true
	}

	// check update results
	user, err := h.Service.Update(r.Context(), userID, req, allowed)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("role changing not allowed"))
			return
		} else if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error(err.Error()))
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
	}

	// send updated user
	render.JSON(w, r, user.ToResponse())
}

// DELETE /user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	err := h.Service.Delete(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusNoContent)
}

package handlers

import (
	"errors"
	"net/http"
	"tracker-backend/internal/domain/dtos"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/utils/request"
	"tracker-backend/internal/interfaces/rest/utils/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	s *services.AuthorizationService
	v *validator.Validate
}

func NewAuthHandler(s *services.AuthorizationService) *AuthHandler {
	return &AuthHandler{
		s: s,
		v: validator.New(),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// decode body
	var req dtos.LoginRequest
	if !request.DecodeJSONBody(w, r, &req) {
		return
	}
	// validate struct
	if !request.ValidateBody(w, r, h.v, req) {
		return
	}
	// execute service function
	token, err := h.s.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else if errors.Is(err, service.ErrInvalidPassword) {
			render.Status(r, http.StatusForbidden)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	// return result
	render.JSON(w, r, map[string]string{"token": token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// decode body
	var req dtos.RegisterRequest
	if !request.DecodeJSONBody(w, r, &req) {
		return
	}
	// validate struct
	if !request.ValidateBody(w, r, h.v, req) {
		return
	}
	// execute service function
	user, err := h.s.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrExists) {
			render.Status(r, http.StatusConflict)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return result
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	// decode body
	var req struct {
		Token string `json:"token"`
	}
	if !request.DecodeJSONBody(w, r, &req) {
		return
	}
	// validate body
	if len(req.Token) < 16 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid token"))
	}
	// execute service function
	_, _, _, err := h.s.Verify(r.Context(), req.Token)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			render.Status(r, http.StatusForbidden)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return result
	render.Status(r, http.StatusNoContent)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// decode body
	var req struct {
		Token string `json:"token"`
	}
	if !request.DecodeJSONBody(w, r, &req) {
		return
	}
	// validate body
	if len(req.Token) < 16 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid token"))
	}

	// execute service function
	token, err := h.s.Refresh(r.Context(), req.Token)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			render.Status(r, http.StatusForbidden)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	// return result
	render.JSON(w, r, map[string]string{"token": token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// decode body
	var req struct {
		Token string `json:"token"`
	}
	if !request.DecodeJSONBody(w, r, &req) {
		return
	}
	// validate body
	if len(req.Token) < 16 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid token"))
	}

	// execute service function
	err := h.s.Logout(r.Context(), req.Token)
	if err != nil {
		if !errors.Is(err, service.ErrInternal) {
			render.Status(r, http.StatusBadRequest)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	// return result
	render.Status(r, http.StatusNoContent)
}

package userHandler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"tracker-backend/internal/lib/handler/response"
	"tracker-backend/internal/middleware"
	userModel "tracker-backend/internal/user/model"
	userService "tracker-backend/internal/user/service"
)

type UserHandler struct {
	Service   *userService.UserService
	Validator *validator.Validate
}

func NewUserHandler(s *userService.UserService) *UserHandler {
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
	var req userModel.RegisterRequest
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

	fmt.Println(user)

	render.Status(r, http.StatusCreated)
}

// POST /user/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userModel.LoginRequest
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
		if err == userService.ErrNotFound {
			render.Status(r, http.StatusNotFound)
		} else if err == userService.ErrForbidden {
			render.Status(r, http.StatusForbidden)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, map[string]string{"token": token})
}

// GET /user
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	fmt.Println(userID)

	user, err := h.Service.GetByID(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error("user not found"))
		return
	}

	render.JSON(w, r, user)
}

// PUT /user
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req userModel.UpdateRequest
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

	user, err := h.Service.Update(r.Context(), userID, req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, user)
}

// DELETE /user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	err := h.Service.Delete(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.Status(r, http.StatusNoContent)
}

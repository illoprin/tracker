package handlers

import (
	"errors"
	"net/http"
	"tracker-backend/internal/domain/dtos"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/middleware"
	"tracker-backend/internal/interfaces/rest/utils/request"
	"tracker-backend/internal/interfaces/rest/utils/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	uSvc *services.UserService
	aSvc *services.AuthorizationService
	v    *validator.Validate
}

func NewUserHandler(
	userService *services.UserService,
	authService *services.AuthorizationService,
) *UserHandler {
	v := validator.New()
	return &UserHandler{
		uSvc: userService,
		aSvc: authService,
		v:    v,
	}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)

	// execute service function
	user, err := h.uSvc.GetByID(ctx, userId)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			render.Status(r, http.StatusInternalServerError)
		} else if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return response
	render.JSON(w, r, user)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	// decode body
	var req dtos.UserUpdateRequest
	if !request.DecodeJSONBody(w, r, &req) {
		return
	}
	// validate struct
	if !request.ValidateBody(w, r, h.v, req) {
		return
	}
	roleChangingAllowed := r.Header.Get("Allow-Access") == "1"
	// execute service function
	user, err := h.uSvc.UpdateByID(ctx, userId, req, roleChangingAllowed)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			render.Status(r, http.StatusInternalServerError)
		} else if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
		} else if errors.Is(err, service.ErrForbidden) {
			render.Status(r, http.StatusForbidden)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	// return response
	render.JSON(w, r, user)
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	// get context keys
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey).(string)
	// execute service function
	err := h.uSvc.DeleteByID(ctx, userId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	// invalidate session
	token := r.Header.Get("Authorization")[len("Bearer "):]
	err = h.aSvc.Logout(ctx, token)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	// return response
	render.Status(r, http.StatusNoContent)
}

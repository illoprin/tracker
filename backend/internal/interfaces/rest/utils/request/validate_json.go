package request

import (
	"net/http"
	"tracker-backend/internal/interfaces/rest/utils/response"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func ValidateBody(w http.ResponseWriter, r *http.Request, v *validator.Validate, req any) bool {
	if err := v.Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ValidationErrorsResp(err.(validator.ValidationErrors)))
		return false
	}
	return true
}

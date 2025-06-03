package request

import (
	"net/http"
	"tracker-backend/internal/interfaces/rest/utils/response"

	"github.com/go-chi/render"
)

// DecodeJSONBody returns true if body is valid
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := render.DecodeJSON(r.Body, target); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to parse form data"))
		return false
	}
	return true
}

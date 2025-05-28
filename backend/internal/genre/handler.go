package genre

import (
	"net/http"
	genreType "tracker-backend/internal/genre/type"

	"github.com/go-chi/render"
)

// GetAllGenres returns allowed genres
func GetAllGenres(w http.ResponseWriter, r *http.Request) {
	var res struct {
		Genres []string `json:"genres"`
	}
	res.Genres = genreType.AllowedGenres
	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, res)
}

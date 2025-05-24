package album

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AlbumHandler struct {
	Service   *AlbumService
	Validator *validator.Validate
}

func NewAlbumHandler(s *AlbumService) *AlbumHandler {
	v := validator.New()
	v.RegisterValidation("status", ValidateStatus)
	v.RegisterValidation("year", ValidateYear)

	return &AlbumHandler{
		Service:   s,
		Validator: v,
	}
}

func (h *AlbumHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *AlbumHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *AlbumHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h *AlbumHandler) GetByArtistID(w http.ResponseWriter, r *http.Request) {

}

func (h *AlbumHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

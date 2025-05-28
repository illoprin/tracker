package albumType

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type AlbumResponse struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	ArtistID  string   `json:"artistID"`
	Year      int      `json:"year"`
	CoverPath string   `json:"coverPath"`
	Genres    []string `json:"genres"`
	IsHidden  bool     `json:"isHidden"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"createdAt"`
}

type AlbumCreateRequest struct {
	Title    string   `json:"title" validate:"required,min=3,max=255"`
	Year     int      `json:"year" validate:"required,year"`
	Genres   []string `json:"genres" validate:"required,genres"`
	ArtistID string   `json:"artistID" validate:"required,uuid4"`
}

type AlbumUpdateRequest struct {
	Title  *string  `json:"title" validate:"omitempty,min=3,max=255"`
	Year   *int     `json:"year" validate:"omitempty,year"`
	Status *string  `json:"status" validate:"omitempty,status"`
	Genres []string `json:"genres" validate:"required,genres"`
}

var (
	statusValidation = map[string]bool{
		StatusModerated:    true,
		StatusOnModeration: true,
	}
)

func ValidateYear(fl validator.FieldLevel) bool {
	year := fl.Field().Interface().(int)
	return year < time.Now().Year()
}

func ValidateStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return statusValidation[status]
}

func (a *Album) ToResponse() AlbumResponse {
	return AlbumResponse{
		ID:        a.ID,
		Title:     a.Title,
		ArtistID:  a.ArtistID,
		Year:      a.Year,
		CoverPath: a.CoverPath,
		Genres:    a.Genres,
		Status:    a.Status,
		IsHidden:  a.IsHidden,
		CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

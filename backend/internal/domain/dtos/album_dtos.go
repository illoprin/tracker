package dtos

import (
	"strconv"
	"time"
	"tracker-backend/internal/domain/repository/schemas"

	"github.com/go-playground/validator/v10"
)

type AlbumCreateRequest struct {
	ArtistID string `json:"artistId" validate:"required,uuid4"`
	Name     string `json:"name" validate:"required,min=3"`
	Year     string `json:"year" validate:"required,year"`
	Type     string `json:"type" validate:"required,type"`
}

type AlbumUpdateRequest struct {
	Name string `json:"name" validate:"omitempty,min=3"`
	Year string `json:"year" validate:"omitempty,year"`
	Type string `json:"type" validate:"omitempty,type"`
}

type AlbumModerationRequest struct {
	Status  string `json:"status" validate:"required,status"`
	Comment string `json:"comment" validate:"required"`
}

type AlbumWithStats struct {
	schemas.Album
	TrackCount    int    `json:"count"`
	TotalDuration int    `json:"duration"`
	ArtistName    string `json:"artistName"`
	ArtistAvatar  string `json:"artistAvatar"`
}

func ValidateYear(fl validator.FieldLevel) bool {
	year := fl.Field().String()
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return false
	}
	return yearInt <= time.Now().Year()
}

var validTypes = map[string]bool{
	schemas.AlbumTypeAlbum:  true,
	schemas.AlbumTypeSingle: true,
}

func ValidateType(fl validator.FieldLevel) bool {
	return validTypes[fl.Field().String()]
}

var validStatuses = map[string]bool{
	schemas.AlbumStatusApproved: true,
	schemas.AlbumStatusRejected: true,
}

func ValidateAlbumStatus(fl validator.FieldLevel) bool {
	return validStatuses[fl.Field().String()]
}

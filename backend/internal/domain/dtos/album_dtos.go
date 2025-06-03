package dtos

type AlbumCreateRequest struct {
	ArtistID string `json:"artistId" validate:"required,uuid4"`
	Name     string `json:"name" validate:"required,min=3"`
	Year     int    `json:"year" validate:"required,year"`
	Type     string `json:"type" validate:"required,type"`
}

type AlbumUpdateRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=3"`
	Year     *int    `json:"year" validate:"omitempty,year"`
	Type     *string `json:"type" validate:"omitempty,type"`
	IsPublic *bool   `json:"isPublic"`
}

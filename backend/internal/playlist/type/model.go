package playlistType

type PlaylistCreateRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=255"`
	UserID    string `json:"userID" validate:"required,uuid4"`
	IsPublic  bool   `json:"isPublic" validate:"required"`
	IsDefault bool   `json:"isDefault" validate:"required"`
}

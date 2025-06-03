package dtos

type PlaylistCreateRequest struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
	IsPublic  bool   `json:"isPublic"`
}

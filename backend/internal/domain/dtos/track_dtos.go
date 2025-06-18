package dtos

type TrackCreateRequest struct {
	Name     string   `json:"name"`
	AlbumID  string   `json:"albumId"`
	Duration int      `json:"duration"`
	Genres   []string `json:"genres"`
}

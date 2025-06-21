package dtos

import "tracker-backend/internal/domain/repository/schemas"

type SearchResult struct {
	Tracks  []TrackResponse  `json:"tracks"`
	Albums  []AlbumWithStats `json:"albums"`
	Artists []schemas.Artist `json:"artists"`
}

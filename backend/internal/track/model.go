package track

// CreateTrackRequest represents a request to create a new track
type CreateTrackRequest struct {
	Title    string   `json:"title" validate:"required,min=1,max=128"`
	Genre    []string `json:"genre" validate:"required,min=1"`
	Duration int      `json:"duration" validate:"required,min=10"`
	AlbumID  string   `json:"albumId" validate:"uuid4"`
}

// TrackResponse represents a response with track information
type TrackResponse struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Duration  int      `json:"duration"`
	Genre     []string `json:"genre"`
	AudioFile string   `json:"audioFile"`
	AlbumID   string   `json:"albumId"`
	CreatedAt string   `json:"createdAt"`
}

// ToResponse converts Track to TrackResponse
func (t *Track) ToResponse() TrackResponse {
	return TrackResponse{
		ID:        t.ID,
		Title:     t.Title,
		Duration:  t.Duration,
		Genre:     t.Genre,
		AudioFile: t.AudioFile,
		AlbumID:   t.AlbumID,
		CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

package dtos

type TrackCreateRequest struct {
	Name     string   `json:"name" validate:"min=3"`
	AlbumID  string   `json:"albumId" validate:"uuid4"`
	Duration int      `json:"duration"`
	Genres   []string `json:"genres" validate:"genres"`
}

type TrackResponse struct {
	ID         string   `json:"id" bson:"id"`
	Name       string   `json:"name" bson:"name"`
	AlbumID    string   `json:"albumId" bson:"albumId"`
	ArtistID   string   `json:"artistId" bson:"artistId"`
	AlbumName  string   `json:"albumName" bson:"albumName"`
	ArtistName string   `json:"artistName" bson:"artistName"`
	Cover      string   `json:"cover" bson:"cover"`
	Duration   int      `json:"duration" bson:"duration"`
	Genres     []string `json:"genres" bson:"genres"`
}

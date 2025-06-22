package dtos

import "time"

type PlaylistCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	IsDefault   bool   `json:"isDefault"`
	Description string `json:"description" validate:"max=500"`
	IsPublic    bool   `json:"isPublic"`
}

type PlaylistUpdateRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Cover       *string `json:"cover"` // base64 или путь
	IsPublic    *bool   `json:"isPublic"`
}

type PlaylistResponse struct {
	ID          string    `bson:"id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Cover       string    `bson:"cover" json:"cover"`
	IsPublic    bool      `bson:"isPublic" json:"isPublic"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
	OwnerInfo   struct {
		ID   string `bson:"id" json:"id"`
		Name string `bson:"name" json:"name"`
	} `bson:"ownerInfo" json:"ownerInfo"`
}

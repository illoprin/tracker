package userModel

import (
	"time"
)

type User struct {
	ID               string    `json:"id,omitempty"`
	Login            string    `json:"login"`
	Email            string    `json:"email"`
	MyChoicePlaylist string    `json:"myChoicePlaylist"`
	CreatedAt        time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateRequest struct {
	Login *string `json:"login,omitempty" validate:"omitempty,min=3,max=32"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	// add possability to update password
	Password *string `json:"password,omitempty" validate:"omitempty,min=6,max=64"`
}

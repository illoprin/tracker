package dtos

type UserUpdateRequest struct {
	Login    *string `json:"login" validate:"omitempty,min=3"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=3"`
	Role     *int    `json:"role" validate:"omitempty"`
}

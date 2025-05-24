package user

import (
	"tracker-backend/internal/auth"

	"github.com/go-playground/validator/v10"
)

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
	// add possability to update role if 'Allow-Access header is set'
	Role *string `json:"role" validate:"omitempty,role"`
}

var validRoles = map[string]bool{
	auth.RoleAdmin:     true,
	auth.RoleModerator: true,
	auth.RoleCustomer:  true,
}

func ValidateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return validRoles[role]
}

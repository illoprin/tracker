package userType

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
	Role *int `json:"role" validate:"omitempty,role"`
}

type UserResponse struct {
	ID               string `json:"id"`
	Login            string `json:"login"`
	Email            string `json:"email"`
	PasswordHash     string `json:"passwordHash"`
	MyChoicePlaylist string `json:"myChoicePlaylist"`
	Role             int    `json:"role"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:               u.ID,
		Login:            u.Login,
		Email:            u.Email,
		PasswordHash:     u.PasswordHash,
		MyChoicePlaylist: u.MyChoicePlaylist,
		Role:             u.Role,
	}
}

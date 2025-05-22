package auth

import "context"

type AuthUser struct {
	ID    string `bson:"id"`
	Email string `bson:"email"`
	Role  string `bson:"role"`
}

const (
	RoleAdmin     = "Admin"
	RoleModerator = "Moderator"
	RoleCustomer  = "Customer"
)

type UserProvider interface {
	GetAuthDTOByID(ctx context.Context, id string, role string) (*AuthUser, error)
}

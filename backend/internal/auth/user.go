package auth

import "context"

type AuthUser struct {
	ID    string `bson:"id"`
	Email string `bson:"email"`
	Role  string `bson:"role"`
}

const (
	RoleAdmin     = 3
	RoleModerator = 2
	RoleCustomer  = 1
)

type UserProvider interface {
	GetAuthDTOByID(ctx context.Context, id string, role string) (*AuthUser, error)
}

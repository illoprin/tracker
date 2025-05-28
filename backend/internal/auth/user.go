package auth

import "context"

type AuthUser struct {
	ID    string `bson:"id"`
	Email string `bson:"email"`
	Role  int    `bson:"role"`
}

const (
	RoleAdmin     = 3
	RoleModerator = 2
	RoleCustomer  = 1
)

type UserProvider interface {
	GetAuthDTOByID(ctx context.Context, id string, role int) (*AuthUser, error)
}

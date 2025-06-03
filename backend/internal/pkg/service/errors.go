package service

import "errors"

var (
	ErrNotFound          = errors.New("requested resource not found")
	ErrForbidden         = errors.New("you have no rights to manage this resource")
	ErrExists            = errors.New("same resource exists")
	ErrInternal          = errors.New("internal server error")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrAlreadyAuthorized = errors.New("already authorized")
)

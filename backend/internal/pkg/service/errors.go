package service

import "errors"

var (
	ErrNotFound  = errors.New("requested resource not found")
	ErrForbidden = errors.New("you have no rights to manage this resource")
	ErrExists    = errors.New("same resource exists")
)

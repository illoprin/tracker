package service

import "errors"

var (
	ErrNotFound     = errors.New("requested resource not found")
	ErrAccessDenied = errors.New("user has no rights to access this resource")
	ErrUploadFailed = errors.New("file upload failed")
)

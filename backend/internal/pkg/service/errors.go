package service

import "errors"

var (
	ErrNotFound     = errors.New("requested resource not found")
	ErrAccessDenied = errors.New("user does not have access rights to this resource")
	ErrUploadFailed = errors.New("file upload failed")
)

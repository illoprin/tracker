package uploadfile

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

// ValidateFile checks file extension
func ValidateFile(
	fileHeader *multipart.FileHeader, allowedExtensions map[string]bool,
) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		return ErrInvalidFileType
	} else if fileHeader.Size > maxFileSize {
		return ErrFileTooLarge
	}

	return nil
}

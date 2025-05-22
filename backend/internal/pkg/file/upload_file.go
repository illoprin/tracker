package file

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file size exceeds 5MB")
)

const (
	maxFileSize = 5 << 20 // 5MB
)

// UploadFile saves file on server and returns path to file
func UploadFile(
	fileHeader *multipart.FileHeader,
	uploadDir string,
	allowedExt map[string]bool,
) (string, error) {
	// check size
	if fileHeader.Size > maxFileSize {
		return "", ErrFileTooLarge
	}

	// check extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExt[ext] {
		return "", ErrInvalidFileType
	}

	// create folder if it not exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// create unique file name
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, newFileName)

	// save file
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}

package storage

import (
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidFileType     = errors.New("invalid file type")
	ErrFileTooLarge        = errors.New("file size exceeds 5MB")
	AllowedImageExtensions = map[string]bool{
		".jpeg": true,
		".png":  true,
		".jpg":  true,
		".webp": true,
	}
	AllowedAudioExtensions = map[string]bool{
		".mp3": true,
		".wav": true,
		".m4a": true,
	}
)

const (
	maxFileSize = 30 << 20 // 30MB
)

// UploadFile saves file on server and returns full path to file
func UploadFile(
	fileHeader *multipart.FileHeader,
	file *multipart.File,
	uploadDir string, // full upload path
) (string, error) {
	// configure logger
	logger := slog.With(slog.String("function", "uploadfile.UploadFile"))

	// check size
	if fileHeader.Size > maxFileSize {
		return "", ErrFileTooLarge
	}

	// create folder if it not exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Warn("failed to create new directory", slog.String("error", err.Error()))
		return "", err
	}

	// create unique file name
	newFileName := uuid.New().String() + strings.ToLower(filepath.Ext(fileHeader.Filename))
	filePath := filepath.Join(uploadDir, newFileName)

	// open file
	src, err := fileHeader.Open()
	if err != nil {
		logger.Warn("failed to open file header", slog.String("error", err.Error()))
		return "", err
	}
	defer src.Close() // close file after saving

	// create file on server
	dst, err := os.Create(filePath)
	if err != nil {
		logger.Warn("failed to create new file", slog.String("error", err.Error()))
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		logger.Warn("failed to copy data to created file", slog.String("error", err.Error()))
		return "", err
	}

	return filePath, nil
}

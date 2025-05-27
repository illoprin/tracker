package uploadfile

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
	maxFileSize = 30 << 20  // 30MB
	BufferSize  = 32 * 1024 // 32KB
)

// UploadFile saves file on server and returns full path to file
func UploadFile(
	fileHeader *multipart.FileHeader,
	file *multipart.File,
	uploadDir string, // full upload path
	allowedExt map[string]bool, // allowed extensions
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

	// open file
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close() // close file after saving

	// create file on server
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

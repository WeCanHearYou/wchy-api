package blob

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
)

type Blob struct {
	Size        int64
	Content     []byte
	ContentType string
}

type StoreBlob struct {
	Key         string
	Content     []byte
	ContentType string
}

type RetrieveBlob struct {
	Key string

	Blob *Blob
}

type DeleteBlob struct {
	Key string
}

// ErrNotFound is returned when given blob is not found
var ErrNotFound = errors.New("Blob not found")

// ErrInvalidKeyFormat is returned when blob key is in invalid format
var ErrInvalidKeyFormat = errors.New("Blob key is in invalid format")

// SanitizeFileName replaces invalid characters from given filename
func SanitizeFileName(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	ext := filepath.Ext(fileName)
	if ext != "" {
		return slug.Make(fileName[0:len(fileName)-len(ext)]) + ext
	}
	return slug.Make(fileName)
}

// ValidateKey checks if key is is valid format
func ValidateKey(key string) error {
	if len(key) == 0 || len(key) > 512 || strings.Contains(key, " ") {
		return ErrInvalidKeyFormat
	}
	if strings.HasPrefix(key, "/") || strings.HasSuffix(key, "/") {
		return ErrInvalidKeyFormat
	}
	return nil
}

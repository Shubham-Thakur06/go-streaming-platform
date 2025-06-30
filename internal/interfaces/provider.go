package interfaces

import (
	"io"
)

type Provider interface {
	UploadFile(file io.Reader, filename string, contentType string) (string, error)
	DeleteFile(filename string) error
	GetFileURL(filename string) (string, error)
	GeneratePresignedURL(filename string, expiresIn int64) (string, error)
}

package repository

import (
	"mime/multipart"
	"time"
)

type IS3Repository interface {
	UploadFile(file *multipart.FileHeader, expir *time.Duration) (string, error)
	DeleteFile(key string) error
	DownloadFile(key string) ([]byte, error)
	UploadMultipleFiles(files []*multipart.FileHeader, expir *time.Duration) ([]string, error)
}

package service

import (
	"mime/multipart"
	"time"
)

type IS3Service interface {
	UploadFile(file *multipart.FileHeader, expir *time.Duration) (string, error)
	DeleteFile(key string) ([]byte, error)
}

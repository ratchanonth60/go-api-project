package service

import (
	"errors"
	"mime/multipart"
	In "project-api/internal/core/port/repository"
	"time"
)

type S3Service struct {
	UserRepo In.IUserRepository
	S3       In.IS3Repository
}

func NewS3Service(input S3Service) *S3Service {
	return &S3Service{
		S3:       input.S3,
		UserRepo: input.UserRepo,
	}
}
func (s *S3Service) UploadFile(file *multipart.FileHeader, expir *time.Duration) (string, error) {
	if expir == nil {
		return "", errors.New("expir is required")
	}
	url, err := s.S3.UploadFile(file, expir)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *S3Service) DeleteFile(key string) ([]byte, error) {
	return s.S3.DeleteFile(key)
}

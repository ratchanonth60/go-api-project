package service

import (
	"errors"
	"mime/multipart"
	"project-api/internal/core/common/utils"
	In "project-api/internal/core/port/repository"
	InS "project-api/internal/core/port/service"

	"time"

	"github.com/gofiber/fiber/v2"
)

type S3Service struct {
	UserRepo In.IUserRepository
	S3       In.IS3Repository
}

func NewS3Service(userRepo In.IUserRepository, s3Repo In.IS3Repository) InS.IS3Service {
	return &S3Service{
		S3:       s3Repo,
		UserRepo: userRepo,
	}
}
func (s *S3Service) UploadFile(c *fiber.Ctx, file *multipart.FileHeader, expir *time.Duration) (string, error) {
	if expir == nil {
		return "", errors.New("expir is required")
	}
	userIDFromContext, ok := utils.GetUserIDFromContext(c.UserContext()) // Renamed variable
	if !ok {
		return "", errors.New("authentication error")
	}
	_, err := s.UserRepo.GetById(c.Context(), userIDFromContext.UserID) // New variable for entity.User
	if err != nil {
		return "", err // Handle error from GetById, e.g., user not found
	}

	url, err := s.S3.UploadFile(file, expir)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *S3Service) DeleteFile(c *fiber.Ctx, key string) ([]byte, error) {
	return s.S3.DeleteFile(key)
}

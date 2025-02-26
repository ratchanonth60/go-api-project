package service

import (
	"mime/multipart"
	"project-api/internal/core/entity"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IS3Service interface {
	DeleteFile(c *fiber.Ctx, key string) error
	DownloadFile(c *fiber.Ctx, key string) ([]byte, *entity.File, error)
	UploadFile(c *fiber.Ctx, files []*multipart.FileHeader, expir *time.Duration) ([]string, error)
}

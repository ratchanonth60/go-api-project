package service

import (
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IS3Service interface {
	UploadFile(c *fiber.Ctx, file *multipart.FileHeader, expir *time.Duration) (string, error)
	DeleteFile(c *fiber.Ctx, key string) ([]byte, error)
}

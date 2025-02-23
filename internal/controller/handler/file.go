package controller

import (
	"project-api/internal/infra/aws"
	"project-api/internal/infra/config"
	"time"

	In "project-api/internal/core/port/service"

	"project-api/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/s3/v2"
)

type FileHeader struct {
	service In.IS3Service
}

func NewFileHandler() *FileHeader {
	s3Congfig := config.Config.GetS3Config()
	credential := config.Config.GetCredentials()
	repo := aws.New(s3.Config{
		Bucket:      s3Congfig.Bucket,
		Region:      s3Congfig.Region,
		Endpoint:    s3Congfig.Endpoint,
		Credentials: credential,
	})
	return &FileHeader{
		service: service.NewS3Service(
			service.S3Service{S3: repo}),
	}
}

func (f *FileHeader) UploadFile(c *fiber.Ctx) error {
	var expirt time.Duration = 0
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Failed to get file", "error": err.Error()})
	}

	fileURL, err := f.service.UploadFile(file, &expirt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to upload", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "File uploaded", "url": fileURL})
}

func (f *FileHeader) DeleteFile(c *fiber.Ctx) error {
	fileName := c.Params("fileName")
	fileURL, err := f.service.DeleteFile(fileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to download", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "File downloaded", "url": fileURL})
}

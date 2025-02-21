package controller

import (
	"project-api/internal/infra/aws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/s3/v2"
)

type FileHeader struct {
	*aws.StorageWrapper
}

func NewFileHandler() *FileHeader {
	return &FileHeader{
		aws.New(s3.Config{}),
	}
}

func (f *FileHeader) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Failed to get file", "error": err.Error()})
	}

	fileURL, err := f.Upload(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to upload", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "File uploaded", "url": fileURL})
}

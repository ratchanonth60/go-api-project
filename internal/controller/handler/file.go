package controller

import (
	"time"

	In "project-api/internal/core/port/service"

	"github.com/gofiber/fiber/v2"
)

type FileHeader struct {
	S3service   In.IS3Service
	UserService In.IUserService
}

func NewFileHandler(userService In.IUserService, s3Service In.IS3Service) *FileHeader {
	return &FileHeader{UserService: userService, S3service: s3Service}
}

func (f *FileHeader) UploadFile(c *fiber.Ctx) error {
	var expirt time.Duration = 0

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Failed to get file", "error": err.Error()})
	}

	fileURL, err := f.S3service.UploadFile(c, file, &expirt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to upload", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "File uploaded", "url": fileURL})
}

func (f *FileHeader) DeleteFile(c *fiber.Ctx) error {
	fileName := c.Params("fileName")

	fileURL, err := f.S3service.DeleteFile(c, fileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to download", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "File downloaded", "url": fileURL})
}

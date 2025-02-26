package controller

import (
	"fmt"
	"time"

	"project-api/internal/core/model/response"
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

	// รับไฟล์หลายไฟล์จาก field "files" (เปลี่ยนจาก "file" เพื่อความชัดเจน)
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusBadRequest,
			Msg:  "Failed to parse multipart form",
			Data: err.Error(),
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusBadRequest,
			Msg:  "Requires at least one file in 'files' field",
			Data: nil,
		})
	}

	// สมมติว่า S3Service.UploadFile รับ []*multipart.FileHeader และคืน []string
	fileURLs, err := f.S3service.UploadFile(c, files, &expirt)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusInternalServerError,
			Msg:  "Fail to UploadFile",
			Data: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccResponse{
		Msg:  fmt.Sprintf("Successfully uploaded %d file(s)", len(fileURLs)),
		Data: fileURLs,
	})
}

func (f *FileHeader) DeleteFile(c *fiber.Ctx) error {
	key := c.Params("key")
	if key != "" {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusBadRequest,
			Msg:  "Failed to delete file, key is required",
		})
	}

	if err := f.S3service.DeleteFile(c, key); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Code: fiber.StatusInternalServerError,
			Msg:  "Fail to delete file.",
			Data: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccResponse{
		Msg: "Successfully deleted file.",
	})
}

func (f *FileHeader) DownloadFile(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusBadRequest,
			Msg:  "Error: File key is required",
		})
	}

	// Download the file using S3service
	data, file, err := f.S3service.DownloadFile(c, key)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(response.ErrorResponse{
			Code: fiber.StatusInternalServerError,
			Msg:  "Error: Fail to download file.",
			Data: err.Error(),
		})
	}

	// Set appropriate headers for file download
	c.Set("Content-Type", file.FileType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))
	c.Set("Content-Length", fmt.Sprintf("%d", file.FileSize))

	// Send the file data
	return c.Send(data)
}

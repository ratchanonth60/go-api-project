package service

import (
	"errors"
	"mime/multipart"
	"project-api/internal/core/entity"
	In "project-api/internal/core/port/repository"
	InS "project-api/internal/core/port/service"
	"project-api/internal/infra/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// S3Service handles file operations with S3 and database
type S3Service struct {
	FileRepo In.IFileRepository
	S3       In.IS3Repository
}

// NewS3Service creates a new S3Service instance
func NewS3Service(fileRepo In.IFileRepository, s3Repo In.IS3Repository) InS.IS3Service {
	return &S3Service{
		S3:       s3Repo,
		FileRepo: fileRepo,
	}
}

func (s *S3Service) UploadFile(c *fiber.Ctx, files []*multipart.FileHeader, expir *time.Duration) ([]string, error) {
	if len(files) == 0 {
		logger.Error("No files provided for upload")
		return nil, errors.New("at least one file is required")
	}

	userID, err := s.getUserID(c)
	if err != nil {
		return nil, err
	}

	var urls []string
	keys := make([]string, len(files))
	for i, file := range files {
		if err := s.validateUploadInput(file, expir); err != nil {
			return nil, err
		}

		key := s.generateKey(file)
		url, err := s.uploadToS3(file, key, expir)
		if err != nil {
			// Cleanup ไฟล์ที่อัปโหลดไปแล้ว
			for _, uploadedKey := range keys[:i] {
				s.cleanupS3File(uploadedKey)
			}
			return nil, err
		}
		urls = append(urls, url)
		keys[i] = key
	}

	// บันทึก metadata ใน transaction เดียว
	tx := s.FileRepo.BeginTransaction(c.Context())
	if tx.Error != nil {
		logger.Error("Failed to start transaction", zap.Error(tx.Error))
		s.cleanupS3Files(keys) // Cleanup ถ้าเริ่ม transaction ไม่ได้
		return nil, errors.New("failed to start transaction")
	}

	for i, file := range files {
		if err := s.checkFileExists(c, tx, keys[i]); err != nil {
			tx.Rollback()
			s.cleanupS3Files(keys)
			return nil, err
		}

		newFile := s.createFileEntity(userID, file, keys[i], urls[i])
		if err := s.FileRepo.Create(c.Context(), newFile); err != nil {
			tx.Rollback()
			logger.Error("Failed to save file metadata",
				zap.String("key", keys[i]),
				zap.Error(err))
			s.cleanupS3Files(keys)
			return nil, errors.New("failed to save file metadata")
		}
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to commit transaction", zap.Error(err))
		s.cleanupS3Files(keys)
		return nil, errors.New("failed to commit transaction")
	}

	return urls, nil
}

// DeleteFile marks a file as deleted in the database and removes it from S3
func (s *S3Service) DeleteFile(c *fiber.Ctx, key string) error {
	userID, err := s.getUserID(c)
	if err != nil {
		return err
	}

	file, err := s.markFileAsDeleted(c, key, userID)
	if err != nil {
		return err
	}

	if err := s.S3.DeleteFile(file.FilePath); err != nil {
		return s.handleS3DeleteError(file.FilePath, err)
	}

	return nil
}

// DownloadFile retrieves a file from S3 after verifying ownership
func (s *S3Service) DownloadFile(c *fiber.Ctx, key string) ([]byte, *entity.File, error) {
	userID, err := s.getUserID(c)
	if err != nil {
		return nil, nil, err
	}

	file, err := s.verifyAndLockFile(c, key, userID)
	if err != nil {
		return nil, nil, err
	}

	data, err := s.S3.DownloadFile(file.FilePath)
	if err != nil {
		return nil, nil, s.handleS3DownloadError(err)
	}

	return data, file, nil
}

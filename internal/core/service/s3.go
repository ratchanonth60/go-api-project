package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"project-api/internal/core/common/utils"
	"project-api/internal/core/entity"
	In "project-api/internal/core/port/repository"
	InS "project-api/internal/core/port/service"
	"project-api/internal/infra/logger"

	"time"

	"github.com/gofiber/fiber/v2"
)

type S3Service struct {
	FileRepo In.IFileRepository
	S3       In.IS3Repository
}

func NewS3Service(fileRepo In.IFileRepository, s3Repo In.IS3Repository) InS.IS3Service {
	return &S3Service{
		S3:       s3Repo,
		FileRepo: fileRepo,
	}
}

func (s *S3Service) UploadFile(c *fiber.Ctx, file *multipart.FileHeader, expir *time.Duration) (string, error) {
	if expir == nil {
		logger.Error("Expiration duration is required")
		return "", errors.New("expiration duration is required")
	}

	// Validate file
	if file.Size > 10<<20 { // 10MB limit, adjust as needed
		logger.Error("File size exceeds 10MB limit")
		return "", errors.New("file size exceeds 10MB limit")
	}

	userIDFromContext, ok := utils.GetUserIDFromContext(c.UserContext())
	if !ok {
		logger.Error(" Context error")
		return "", errors.New("authentication error")
	}

	// Generate a unique key (e.g., using UUID or timestamp to prevent overwrites)
	// Upload to S3 first (outside transaction due to S3's nature)
	key := fmt.Sprintf("file/%s/%s", file.Header.Get("Content-Type"), file.Filename)

	filePath, err := s.S3.UploadFile(file, expir)
	if err != nil {
		logger.Error("Failed to upload file to S3: " + err.Error())
		return "", errors.New("failed to upload file to S3: " + err.Error())
	}

	// Start transaction for database operation
	tx := s.FileRepo.BeginTransaction(c.Context())
	if tx.Error != nil {
		// Cleanup S3 file on transaction failure
		s.S3.DeleteFile(key) // Ignore error for simplicity; log in production
		return "", errors.New("failed to start transaction: " + tx.Error.Error())
	}

	var existingFile entity.File
	if err := s.FileRepo.FindByKey(c.Context(), key, &existingFile); err == nil && !existingFile.IsDeleted {
		tx.Rollback()
		s.S3.DeleteFile(key) // Cleanup new upload
		return "", errors.New("file already exists for this user")
	}

	// Create file entity
	newFile := &entity.File{
		UserID:   userIDFromContext.UserID,
		FileName: file.Filename,
		FileSize: file.Size,
		FileType: file.Header.Get("Content-Type"),
		FilePath: key,
		UrlPath:  filePath,
	}

	// Save to database
	if err := s.FileRepo.Create(c.Context(), newFile); err != nil {
		tx.Rollback()
		// Cleanup S3 file on DB failure
		if delErr := s.S3.DeleteFile(key); delErr != nil {
			// Log this error in production; file is orphaned
			logger.Error("Failed to delete file from S3: " + delErr.Error())
		}
		return "", errors.New("failed to save file metadata: " + err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		s.S3.DeleteFile(key) // Cleanup
		logger.Error("Failed to commit transaction: " + err.Error())
		return "", errors.New("failed to commit transaction: " + err.Error())
	}

	return filePath, nil
}

func (s *S3Service) DeleteFile(c *fiber.Ctx, key string) error {
	userIDFromContext, ok := utils.GetUserIDFromContext(c.UserContext())
	if !ok {
		return errors.New("authentication error")
	}

	// Start transaction
	tx := s.FileRepo.BeginTransaction(c.Context())
	if tx.Error != nil {
		return errors.New("failed to start transaction: " + tx.Error.Error())
	}

	// Find file with lock to prevent concurrent modifications
	var file entity.File
	if err := s.FileRepo.FindByKeyForUpdate(c.Context(), key, &file); err != nil {
		tx.Rollback()
		return errors.New("file not found or already deleted")
	}

	// Check ownership
	if file.UserID != userIDFromContext.UserID {
		tx.Rollback()
		return errors.New("unauthorized: you do not own this file")
	}

	// Idempotency check
	if file.IsDeleted {
		tx.Rollback()
		return nil
	}

	// Mark file as deleted in DB first
	file.IsDeleted = true
	if err := s.FileRepo.Update(c.Context(), &file); err != nil {
		tx.Rollback()
		return errors.New("failed to mark file as deleted: " + err.Error())
	}

	// Commit DB transaction before S3 operation
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("failed to commit transaction: " + err.Error())
	}

	// Delete from S3 (outside transaction)

	if err := s.S3.DeleteFile(file.FilePath); err != nil {
		return errors.New("file marked as deleted in DB but failed to delete from S3: " + err.Error())
	}

	return nil
}

func (s *S3Service) DownloadFile(c *fiber.Ctx, key string) ([]byte, *entity.File, error) {
	userIDFromContext, ok := utils.GetUserIDFromContext(c.UserContext())
	if !ok {
		logger.Error("Context error")
		return nil, nil, errors.New("authentication error")
	}

	// Start transaction
	tx := s.FileRepo.BeginTransaction(c.Context())
	if tx.Error != nil {
		logger.Error("Failed to start transaction: " + tx.Error.Error())
		return nil, nil, errors.New("failed to start transaction: " + tx.Error.Error())
	}
	defer tx.Rollback() // Rollback if not committed

	// Find file with lock
	var file entity.File
	if err := s.FileRepo.FindByKeyForUpdate(c.Context(), key, &file); err != nil {
		logger.Error("File not found: " + err.Error())
		return nil, nil, errors.New("file not found")
	}

	// Check ownership and deletion status
	if file.UserID != userIDFromContext.UserID {
		logger.Error("Unauthorized access attempt")
		return nil, nil, errors.New("unauthorized: you do not own this file")
	}
	if file.IsDeleted {
		logger.Error("File already deleted")
		return nil, nil, errors.New("file has been deleted")
	}

	// Download from S3
	data, err := s.S3.DownloadFile(file.FilePath)
	if err != nil {
		logger.Error("Failed to download file from S3: " + err.Error())
		return nil, nil, errors.New("failed to download file from S3: " + err.Error())
	}

	// Commit transaction (no changes to DB, but maintains consistency)
	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit transaction: " + err.Error())
		return nil, nil, errors.New("failed to commit transaction: " + err.Error())
	}

	return data, &file, nil
}

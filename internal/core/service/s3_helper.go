package service

import (
	"errors"
	"mime/multipart"
	"project-api/internal/core/common/utils"
	"project-api/internal/core/entity"
	"project-api/internal/infra/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// validateUploadInput checks input constraints for file upload
func (s *S3Service) validateUploadInput(file *multipart.FileHeader, expir *time.Duration) error {
	if expir == nil {
		logger.Error("Expiration duration is required")
		return errors.New("expiration duration is required")
	}

	const maxFileSize = 10 << 20 // 10MB
	if file.Size > maxFileSize {
		logger.Error("File size exceeds limit",
			zap.Int64("size", file.Size),
			zap.String("filename", file.Filename),
			zap.Int64("maxSize", maxFileSize))
		return errors.New("file size exceeds 10MB limit")
	}
	return nil
}

// getUserID retrieves the authenticated user's ID from context
func (s *S3Service) getUserID(c *fiber.Ctx) (uint, error) {
	userIDFromContext, ok := utils.GetUserIDFromContext(c.UserContext())
	if !ok {
		logger.Error("Context error: unable to retrieve user ID")
		return 0, errors.New("authentication error")
	}
	return userIDFromContext.UserID, nil
}

// generateKey creates a unique key for S3 storage
func (s *S3Service) generateKey(file *multipart.FileHeader) string {
	return "file/" + file.Header.Get("Content-Type") + "/" + file.Filename
}

// uploadToS3 uploads a file to S3 and returns the URL
func (s *S3Service) uploadToS3(file *multipart.FileHeader, key string, expir *time.Duration) (string, error) {
	url, err := s.S3.UploadFile(file, expir)
	if err != nil {
		logger.Error("Failed to upload file to S3",
			zap.String("key", key),
			zap.Error(err))
		return "", errors.New("failed to upload file to S3")
	}
	return url, nil
}

// saveFileMetadata saves file metadata to the database within a transaction
func (s *S3Service) saveFileMetadata(c *fiber.Ctx, userID uint, file *multipart.FileHeader, key, url string) error {
	tx := s.FileRepo.BeginTransaction(c.Context())
	if tx.Error != nil {
		logger.Error("Failed to start transaction",
			zap.Error(tx.Error))
		return errors.New("failed to start transaction")
	}

	if err := s.checkFileExists(c, tx, key); err != nil {
		tx.Rollback()
		return err
	}

	newFile := s.createFileEntity(userID, file, key, url)
	if err := s.FileRepo.Create(c.Context(), newFile); err != nil {
		tx.Rollback()
		logger.Error("Failed to save file metadata",
			zap.String("key", key),
			zap.Error(err))
		return errors.New("failed to save file metadata")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to commit transaction",
			zap.String("key", key),
			zap.Error(err))
		return errors.New("failed to commit transaction")
	}
	return nil
}

// checkFileExists checks if a file already exists and is not deleted
func (s *S3Service) checkFileExists(c *fiber.Ctx, tx *gorm.DB, key string) error {
	var existingFile entity.File
	if err := s.FileRepo.FindByKey(c.Context(), key, &existingFile); err == nil && !existingFile.IsDeleted {
		logger.Error("File already exists",
			zap.String("key", key))
		return errors.New("file already exists")
	}
	return nil
}

// createFileEntity constructs a new File entity
func (s *S3Service) createFileEntity(userID uint, file *multipart.FileHeader, key, url string) *entity.File {
	return &entity.File{
		UserID:   userID,
		FileName: file.Filename,
		FileSize: file.Size,
		FileType: file.Header.Get("Content-Type"),
		FilePath: key,
		UrlPath:  url,
	}
}

// cleanupS3File deletes a single file from S3 if an error occurs
func (s *S3Service) cleanupS3File(key string) {
	if err := s.S3.DeleteFile(key); err != nil {
		logger.Error("Failed to delete file from S3",
			zap.String("key", key),
			zap.Error(err))
	}
}

// cleanupS3Files deletes multiple files from S3 if an error occurs
func (s *S3Service) cleanupS3Files(keys []string) {
	for _, key := range keys {
		s.cleanupS3File(key)
	}
}

// markFileAsDeleted updates the file status to deleted within a transaction
func (s *S3Service) markFileAsDeleted(c *fiber.Ctx, key string, userID uint) (*entity.File, error) {
	tx := s.FileRepo.BeginTransaction(c.Context())
	if tx.Error != nil {
		logger.Error("Failed to start transaction",
			zap.Error(tx.Error))
		return nil, errors.New("failed to start transaction")
	}

	var file entity.File
	if err := s.FileRepo.FindByKeyForUpdate(c.Context(), key, &file); err != nil {
		tx.Rollback()
		logger.Error("File not found or already deleted",
			zap.String("key", key),
			zap.Error(err))
		return nil, errors.New("file not found or already deleted")
	}

	if file.UserID != userID {
		tx.Rollback()
		logger.Error("Unauthorized: user does not own this file",
			zap.String("key", key),
			zap.Uint("userID", userID))
		return nil, errors.New("unauthorized: you do not own this file")
	}

	if file.IsDeleted {
		tx.Rollback()
		return &file, nil // Idempotent: already deleted
	}

	file.IsDeleted = true
	if err := s.FileRepo.Update(c.Context(), &file); err != nil {
		tx.Rollback()
		logger.Error("Failed to mark file as deleted",
			zap.String("key", key),
			zap.Error(err))
		return nil, errors.New("failed to mark file as deleted")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to commit transaction",
			zap.String("key", key),
			zap.Error(err))
		return nil, errors.New("failed to commit transaction")
	}

	return &file, nil
}

// verifyAndLockFile locks and verifies the file for download
func (s *S3Service) verifyAndLockFile(c *fiber.Ctx, key string, userID uint) (*entity.File, error) {
	tx := s.FileRepo.BeginTransaction(c.Context())
	if tx.Error != nil {
		logger.Error("Failed to start transaction",
			zap.Error(tx.Error))
		return nil, errors.New("failed to start transaction")
	}
	defer tx.Rollback()

	var file entity.File
	if err := s.FileRepo.FindByKeyForUpdate(c.Context(), key, &file); err != nil {
		logger.Error("File not found",
			zap.String("key", key),
			zap.Error(err))
		return nil, errors.New("file not found")
	}

	if file.UserID != userID {
		logger.Error("Unauthorized access attempt",
			zap.String("key", key),
			zap.Uint("userID", userID))
		return nil, errors.New("unauthorized: you do not own this file")
	}

	if file.IsDeleted {
		logger.Error("File already deleted",
			zap.String("key", key))
		return nil, errors.New("file has been deleted")
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit transaction",
			zap.String("key", key),
			zap.Error(err))
		return nil, errors.New("failed to commit transaction")
	}

	return &file, nil
}

// handleS3DeleteError logs and formats S3 deletion errors
func (s *S3Service) handleS3DeleteError(key string, err error) error {
	logger.Error("Failed to delete file from S3",
		zap.String("key", key),
		zap.Error(err))
	return errors.New("file marked as deleted in DB but failed to delete from S3")
}

// handleS3DownloadError logs and formats S3 download errors
func (s *S3Service) handleS3DownloadError(err error) error {
	logger.Error("Failed to download file from S3",
		zap.Error(err))
	return errors.New("failed to download file from S3")
}

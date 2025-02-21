package aws

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"project-api/internal/infra/config"
	"time"

	"github.com/gofiber/storage/s3/v2"
)

type StorageWrapper struct {
	*s3.Storage
}

func New(config s3.Config) *StorageWrapper {
	return &StorageWrapper{
		s3.New(config),
	}
}
func (s *StorageWrapper) Upload(file *multipart.FileHeader) (string, error) {
	// เปิดไฟล์
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// อ่านไฟล์เป็น []byte
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, src); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	data := buf.Bytes()

	// สร้าง key สำหรับจัดเก็บไฟล์ใน S3
	key := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)

	// อัปโหลดไฟล์ไปยัง S3
	err = s.Set(key, data, 0) // TTL = 0 (ไฟล์จะไม่หมดอายุ)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	config := config.S3.GetS3Config()
	// URL ของไฟล์ที่อัปโหลด
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.Bucket, config.Region, key)
	return fileURL, nil
}

// Download ดึงไฟล์จาก S3
func (s *StorageWrapper) Download(key string) ([]byte, error) {
	data, err := s.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %v", err)
	}
	return data, nil
}

// Delete ลบไฟล์จาก S3
func (s *StorageWrapper) Delete(key string) error {
	err := s.Storage.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}
	return nil
}

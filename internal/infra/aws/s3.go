package aws

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"project-api/internal/infra/config"
	"sync"
	"time"

	"project-api/internal/core/port/repository"

	"github.com/gofiber/storage/s3/v2"
)

type StorageWrapper struct {
	*s3.Storage
}

func New(config s3.Config) repository.IS3Repository {
	return &StorageWrapper{
		s3.New(config),
	}
}

func (s *StorageWrapper) UploadFile(file *multipart.FileHeader, expir *time.Duration) (string, error) {
	// เปิดไฟล์
	config := config.Config.GetS3Config()

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
	key := fmt.Sprintf("file/%s/%s", file.Header.Get("Content-Type"), file.Filename)

	// อัปโหลดไฟล์ไปยัง S3
	if expir == nil {
		expir = &DefaultExpiry
	}
	err = s.Set(key, data, *expir) // TTL = 0 (ไฟล์จะไม่หมดอายุ)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// URL ของไฟล์ที่อัปโหลด
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.Bucket, config.Region, key)

	return fileURL, nil
}

func (s *StorageWrapper) UploadMultipleFiles(files []*multipart.FileHeader, expir *time.Duration) ([]string, error) {
	var (
		fileURLs     []string
		uploadErrors []error
		mu           sync.Mutex
		wg           sync.WaitGroup
	)

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			fileURL, err := s.UploadFile(file, expir)
			mu.Lock()
			if err != nil {
				uploadErrors = append(uploadErrors, err)
			} else {
				fileURLs = append(fileURLs, fileURL)
			}
			mu.Unlock()
		}(file)
	}
	wg.Wait()

	if len(uploadErrors) > 0 {
		return fileURLs, fmt.Errorf("failed to upload %d file(s): %v", len(uploadErrors), uploadErrors)
	}
	return fileURLs, nil
}

// Download ดึงไฟล์จาก S3
func (s *StorageWrapper) DeleteFile(key string) error {
	err := s.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to download file from S3: %v", err)
	}
	return nil
}

func (s *StorageWrapper) DownloadFile(key string) ([]byte, error) {
	// Get file from S3
	data, err := s.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %v", err)
	}

	if data == nil {
		return nil, fmt.Errorf("file not found in S3")
	}

	return data, nil
}

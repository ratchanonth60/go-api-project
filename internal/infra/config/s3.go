package config

import (
	"github.com/gofiber/storage/s3/v2"
)

func (s *AppConfig) GetS3Config() *s3.Config {
	return &s3.Config{
		Bucket:   s.S3.Bucket,
		Endpoint: s.S3.Endpoint,
		Region:   s.S3.Region,
	}
}

func (s *AppConfig) GetCredentials() s3.Credentials {
	return s3.Credentials{
		AccessKey:       s.Credentials.AccessKey,
		SecretAccessKey: s.Credentials.SecretKey,
	}
}

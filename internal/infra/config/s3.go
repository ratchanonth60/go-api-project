package config

import "github.com/gofiber/storage/s3/v2"

type S3Config struct {
	S3 struct {
		Region   string `yaml:"region" env:"AWS_REGION"`
		Bucket   string `yaml:"bucket" env:"AWS_BUCKET"`
		Endpoint string `yaml:"endpoint" env:"AWS_ENDPOINT"`
	} `yaml:"s3"`
	Credentials struct {
		AccessKey string `yaml:"access_key" env:"AWS_ACCESS_KEY_ID"`
		SecretKey string `yaml:"secret_key" env:"AWS_SECRET_ACCESS_KEY"`
	} `yaml:"credentials"`
}

var S3 *S3Config

func (s *S3Config) GetS3Config() *s3.Config {
	return &s3.Config{
		Bucket:   s.S3.Bucket,
		Endpoint: s.S3.Endpoint,
		Region:   s.S3.Endpoint,
	}
}

func (s *S3Config) GetCredentials() s3.Credentials {
	return s3.Credentials{
		AccessKey:       s.Credentials.AccessKey,
		SecretAccessKey: s.Credentials.SecretKey,
	}
}

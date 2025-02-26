package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func (s *AppConfig) GetSQSConfig() *aws.Config {
	return &aws.Config{
		Region:   aws.String(s.SQS.Region),
		Endpoint: aws.String(s.SQS.Endpoint),
	}
}

func (s *AppConfig) GetCredentialSQS() *credentials.Credentials {
	return credentials.NewStaticCredentials(s.SQS.AccessKey, s.SQS.SecretKey, "")
}

package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func (s *AppConfig) GetSESConfig() *aws.Config {
	return &aws.Config{
		Region:   aws.String(s.SES.Region),
		Endpoint: aws.String(s.SES.Endpoint),
	}
}

func (s *AppConfig) GetCredentialSES() *credentials.Credentials {
	return credentials.NewStaticCredentials(s.SES.AccessKey, s.SES.SecretKey, "")
}

package task

import (
	"project-api/internal/infra/aws"
	"project-api/internal/infra/google"
)

func TaskSendConfirmationEmail(toEmail string, token string, name string, host string) error {
	return aws.SendConfirmationEmail(toEmail, token, name, host)
}

func TaskSendConfirmationEmailSMTP(toEmail string, token string, name string, host string) error {
	return google.SendConfirmationEmailSMTP(toEmail, token, name, host)
}

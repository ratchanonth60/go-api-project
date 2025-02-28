package task

import (
	"project-api/internal/infra/aws"
)

func TaskSendConfirmationEmail(toEmail string, token string, name string, host string) error {
	return aws.SendConfirmationEmail(toEmail, token, name, host)
}

func TaskSendResetPasswordEmail(toEmail string, token string, name string, host string) error {
	return aws.SendResetPasswordEmail(toEmail, token, name, host)
}

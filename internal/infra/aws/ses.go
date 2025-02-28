package aws

import (
	"bytes"
	"fmt"
	"html/template"
	"project-api/internal/infra"
	"project-api/internal/infra/config"
	"project-api/internal/infra/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

func SendConfirmationEmail(toEmail string, token string, name string, host string) error {
	awsConfig := config.Config.GetSESConfig()
	awsCredential := config.Config.GetCredentialSES()
	// สร้าง AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      awsConfig.Region,
		Credentials: awsCredential,
	})
	if err != nil {
		logger.Error("Failed to create AWS session for SES", zap.Error(err))
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	// สร้าง SES client
	sesClient := ses.New(sess)

	// โหลดและ render เทมเพลต HTML
	tmpl, err := template.ParseFiles("templates/email_confirmation.html")
	if err != nil {
		logger.Error("Failed to parse email template", zap.Error(err))
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	data := infra.EmailData{
		Name:  name,
		Token: token,
		Host:  host,
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		logger.Error("Failed to render email template", zap.Error(err))
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// ส่งอีเมลผ่าน SES
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(toEmail)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Confirm Your Email Address"),
			},
		},
		Source: aws.String(config.Config.SES.From),
	}

	_, err = sesClient.SendEmail(input)
	if err != nil {
		logger.Error("Failed to send email via SES", zap.Error(err), zap.String("to", toEmail))
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	logger.Info("Confirmation email sent via SES", zap.String("to", toEmail))
	return nil
}

func SendResetPasswordEmail(toEmail string, token string, name string, host string) error {
	awsConfig := config.Config.GetSESConfig()
	awsCredential := config.Config.GetCredentialSES()

	// สร้าง AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      awsConfig.Region,
		Credentials: awsCredential,
	})
	if err != nil {
		logger.Error("Failed to create AWS session for SES", zap.Error(err))
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	// สร้าง SES client
	sesClient := ses.New(sess)

	// โหลดและ render เทมเพลต HTML
	tmpl, err := template.ParseFiles("templates/email_reset_password.html")
	if err != nil {
		logger.Error("Failed to parse reset password email template", zap.Error(err))
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	data := infra.EmailData{
		Name:  name,
		Token: token,
		Host:  host,
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		logger.Error("Failed to render reset password email template", zap.Error(err))
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// ส่งอีเมลผ่าน SES
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(toEmail)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Reset Your Password"),
			},
		},
		Source: aws.String(config.Config.SES.From),
	}

	_, err = sesClient.SendEmail(input)
	if err != nil {
		logger.Error("Failed to send reset password email via SES", zap.Error(err), zap.String("to", toEmail))
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	logger.Info("Reset password email sent via SES", zap.String("to", toEmail))
	return nil
}

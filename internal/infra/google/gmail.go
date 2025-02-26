package google

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"project-api/internal/infra"
	"project-api/internal/infra/config"
	"project-api/internal/infra/logger"

	"go.uber.org/zap"
)

func SendConfirmationEmailSMTP(toEmail, token, name, host string) error {
	// โหลด SMTP settings จาก config
	from := config.Config.GmailSMTP.From
	password := config.Config.GmailSMTP.Password
	smtpHost := config.Config.GmailSMTP.SMTPHost
	smtpPort := config.Config.GmailSMTP.SMTPPort

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

	// สร้างอีเมล message
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: Confirm Your Email Address\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", toEmail, from, body.String()))

	// SMTP authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// ส่งอีเมล
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, msg)
	fmt.Println(smtpHost + ":" + smtpPort)
	if err != nil {
		logger.Error("Failed to send email via Gmail SMTP", zap.Error(err), zap.String("to", toEmail))
		return fmt.Errorf("failed to send email via Gmail SMTP: %w", err)
	}

	logger.Info("Confirmation email sent via Gmail SMTP", zap.String("to", toEmail))
	return nil
}

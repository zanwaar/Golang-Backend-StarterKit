package utils

import (
	"fmt"
	"golang-backend/config"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(toEmail, code string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "noreply@golang-backend.com")
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", "Email Verification Code")
	mailer.SetBody("text/html", fmt.Sprintf("Your verification code is: <b>%s</b>", code))

	dialer := gomail.NewDialer(
		config.AppConfig.SMTPHost,
		config.AppConfig.SMTPPort,
		config.AppConfig.SMTPUser,
		config.AppConfig.SMTPPass,
	)

	return dialer.DialAndSend(mailer)
}

func SendResetPasswordEmail(toEmail, code string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "noreply@golang-backend.com")
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", "Reset Password Code")
	mailer.SetBody("text/html", fmt.Sprintf("Your reset password code is: <b>%s</b>", code))

	dialer := gomail.NewDialer(
		config.AppConfig.SMTPHost,
		config.AppConfig.SMTPPort,
		config.AppConfig.SMTPUser,
		config.AppConfig.SMTPPass,
	)

	return dialer.DialAndSend(mailer)
}

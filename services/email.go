package services

import (
	"log"
	"net/smtp"
)

// EmailService define a interface para o envio de e-mails
type EmailService interface {
	SendEmail(to, subject, body string) error
}

// emailService encapsula as configurações do SMTP
type emailService struct {
	fromEmail         string
	fromEmailSMTP     string
	fromEmailPassword string
	SMTPAddrress      string
}

type EmailRequestBody struct {
	ToAddr  string `json:"to_addr"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewEmailService(fromEmail, fromEmailSMTP, fromEmailPassword, SMTPAddrress string) EmailService {
	return &emailService{
		fromEmail:         fromEmail,
		fromEmailSMTP:     fromEmailSMTP,
		fromEmailPassword: fromEmailPassword,
		SMTPAddrress:      SMTPAddrress,
	}
}

func (e *emailService) SendEmail(to string, subject string, htmlBody string) error {
	auth := smtp.PlainAuth(
		"",
		e.fromEmail,
		e.fromEmailPassword,
		e.fromEmailSMTP,
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	message := "Subject: " + subject + "\n" + headers + "\n\n" + htmlBody
	err := smtp.SendMail(
		e.SMTPAddrress,
		auth,
		e.fromEmail,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		log.Print(err)
		return err
	}

	return err
}

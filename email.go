package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"

)

type EmailSender struct {
	SMTPServer   string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
}

func NewEmailSender() *EmailSender{
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("FROM_EMAIL")
	
	
	if smtpServer == "" || smtpPort == "" || smtpUser == "" || smtpPassword == "" || fromEmail == "" {
		log.Fatal("Email variables not found")
	}
	
	emailSender := &EmailSender{
		SMTPServer:   smtpServer,
		SMTPPort:     smtpPort,
		SMTPUser:     smtpUser,
		SMTPPassword: smtpPassword,
		FromEmail:    fromEmail,
	}

	return emailSender
}

func (eSender *EmailSender) SendWelcomeEmail(to, name string) error {
	auth := smtp.PlainAuth("", eSender.SMTPUser, eSender.SMTPPassword, eSender.SMTPServer)

	tmpl, err := template.ParseFiles("email.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer

	body.Write([]byte("From: " + eSender.FromEmail + "\r\n"))
	body.Write([]byte("To: " + to + "\r\n"))
	body.Write([]byte("Subject: Welcome\r\n"))
	body.Write([]byte("MIME-Version: 1.0\r\n"))
	body.Write([]byte("Content-Type: text/html; charset=\"UTF-8\"\r\n"))
	body.Write([]byte("\r\n"))

	err = tmpl.Execute(&body, struct {
		Name    string
		Email   string
		Message string
	}{
		Name:    name,
		Email:   to,
		Message: "You're now a part of something beautiful!",
	})

	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	err = smtp.SendMail(
		eSender.SMTPServer+":"+eSender.SMTPPort,
		auth,
		eSender.FromEmail,
		[]string{to},
		body.Bytes(),
	)

	if err != nil {
		log.Printf("SMTP error details: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("Email sent successfully")
	return nil
}

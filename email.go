package main

import (
	"fmt"
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

func (eSender *EmailSender) SendWelcomeEmail(to, name string) error{
	auth := smtp.PlainAuth("",eSender.SMTPUser,eSender.SMTPPassword,eSender.SMTPServer)

	subject := "Welcome to Our Service!"
    body := fmt.Sprintf("Hi %s,\n\nThank you for signing up with our service!\n\nBest regards,\nThe Team", name)
    
    msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", 
        eSender.FromEmail, 
        to, 
        subject, 
        body)


	err := smtp.SendMail(
		eSender.SMTPServer+":"+eSender.SMTPPort,
		auth,
		eSender.FromEmail,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
        log.Printf("SMTP error details: %v", err)
        return fmt.Errorf("failed to send email: %w", err)
    }
    
    log.Println("Email sent successfully")
    return nil
}
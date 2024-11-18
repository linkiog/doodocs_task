package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/smtp"
	"os"
	"strings"
)

func SendFileByEmail(file io.Reader, header *multipart.FileHeader, emails string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	fmt.Printf("SMTP_HOST: %s\n", smtpHost)
	fmt.Printf("SMTP_PORT: %s\n", smtpPort)
	fmt.Printf("SMTP_USER: %s\n", smtpUser)
	fmt.Printf("SMTP_PASS: %s\n", smtpPass)

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP configuration is missing")
	}

	emailList := strings.Split(emails, ",")
	if len(emailList) == 0 {
		return fmt.Errorf("No valid emails provided")
	}

	subject := "Here is your file"
	body := "Please find the attached file."
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", smtpUser, strings.Join(emailList, ","), subject, body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, emailList, []byte(message))
	if err != nil {
		return fmt.Errorf("Failed to send email: %v", err)
	}

	return nil
}

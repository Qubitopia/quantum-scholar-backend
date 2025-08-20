package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendEmailToNewUser(to string, Name string, magicLink string) {
	// Load environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("FROM_EMAIL")

	// Read HTML template from file
	templatePath := "mail/newMail.html"
	bodyTemplateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		log.Println("Failed to read newMail.html:", err)
		return
	}
	bodyTemplate := string(bodyTemplateBytes)

	// Email content
	subject := "Subject: Welcome to Quantum Scholar by Qubitopia\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(bodyTemplate, Name, magicLink)
	msg := []byte(subject + mime + body)

	// Auth
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Failed to send email:", err)
	}

	fmt.Println("✅ Email sent successfully.")
}

func SendEmailToOldUser(to string, Name string, magicLink string) {
	// Load environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("FROM_EMAIL")

	// Read HTML template from file
	templatePath := "mail/oldMail.html"
	bodyTemplateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		log.Println("Failed to read oldMail.html:", err)
		return
	}
	bodyTemplate := string(bodyTemplateBytes)

	// Email content
	subject := fmt.Sprintf("Subject: Welcome back, %s — your secure login link inside\r\n", Name)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(bodyTemplate, Name, magicLink)
	msg := []byte(subject + mime + body)

	// Auth
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Failed to send email:", err)
	} else {
		log.Println("✅ Email sent successfully.")
	}
}

func SendEmailInvoiceForQSCoinsPurchase(to string, Name string, order_id string, coins_amount string, currency string, rate string) {
	// Load environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("FROM_EMAIL")

	// Read HTML template from file
	templatePath := "mail/qsCoinsPurchaseInvoice.html"
	bodyTemplateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		log.Println("Failed to read qsCoinsPurchaseInvoice.html:", err)
		return
	}
	bodyTemplate := string(bodyTemplateBytes)

	// Email content
	subject := fmt.Sprintf("Subject: %s — Purchase of %s QS Coins\r\n", order_id, coins_amount)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(bodyTemplate, Name, coins_amount, order_id, coins_amount, currency, rate, currency, rate)
	msg := []byte(subject + mime + body)

	// Auth
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Failed to send email:", err)
	} else {
		log.Println("✅ Email sent successfully.")
	}
}

func SendEmailNotificationOfUserLogin(to string, Name string, timestamp string, ipAddress string, userAgent string) {
	// Load environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("FROM_EMAIL")

	// Read HTML template from file
	templatePath := "mail/userLoginNotification.html"
	bodyTemplateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		log.Println("Failed to read userLoginNotification.html:", err)
		return
	}
	bodyTemplate := string(bodyTemplateBytes)

	// Email content
	subject := fmt.Sprintf("Subject: New Login Attempt Detected for %s\r\n", to)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(bodyTemplate, Name, to, timestamp, ipAddress, userAgent)
	msg := []byte(subject + mime + body)

	// Auth
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Failed to send email:", err)
	} else {
		log.Println("✅ Email sent successfully.")
	}
}

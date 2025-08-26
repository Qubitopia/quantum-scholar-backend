package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/Qubitopia/quantum-scholar-backend/database"
)

var (
	newUserTemplate  string
	oldUserTemplate  string
	invoiceTemplate  string
	newLoginTemplate string
	auth             smtp.Auth
)

func LoadEmailTemplates() {
	// Load New User Email Template
	newUserBytes, err := os.ReadFile("mail/newMail.html")
	if err != nil {
		log.Fatalf("Failed to read newMail.html: %v", err)
	}
	newUserTemplate = string(newUserBytes)

	// Load Old User Email Template
	oldUserBytes, err := os.ReadFile("mail/oldMail.html")
	if err != nil {
		log.Fatalf("Failed to read oldMail.html: %v", err)
	}
	oldUserTemplate = string(oldUserBytes)

	// Load Invoice Email Template
	invoiceBytes, err := os.ReadFile("mail/qsCoinsPurchaseInvoice.html")
	if err != nil {
		log.Fatalf("Failed to read qsCoinsPurchaseInvoice.html: %v", err)
	}
	invoiceTemplate = string(invoiceBytes)

	// Load New Login Notification Email Template
	newLoginBytes, err := os.ReadFile("mail/userLoginNotification.html")
	if err != nil {
		log.Fatalf("Failed to read userLoginNotification.html: %v", err)
	}
	newLoginTemplate = string(newLoginBytes)
}

func InitEmail() {
	auth = smtp.PlainAuth("", database.SMTP_USERNAME, database.SMTP_PASSWORD, database.SMTP_HOST)
}

func sendEmail(to string, subject string, body string) error {
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	msg := []byte(subject + mime + body)

	// Send email
	err := smtp.SendMail(database.SMTP_HOST+":"+database.SMTP_PORT, auth, database.FROM_EMAIL, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}

func SendEmailToNewUser(to string, Name string, magicLink string) error {
	// Email content
	subject := "Subject: " + "Welcome to Quantum Scholar by Qubitopia" + "\r\n"
	body := fmt.Sprintf(newUserTemplate, Name, database.BASE_URL, magicLink, database.BASE_URL, database.BASE_URL)

	// Send email
	err := sendEmail(to, subject, body)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}
	log.Println("✅ Email sent successfully.")
	return nil
}

func SendEmailToOldUser(to string, Name string, magicLink string) error {
	// Email content
	subject := fmt.Sprintf("Subject: Welcome back, %s — your secure login link inside\r\n", Name)
	body := fmt.Sprintf(oldUserTemplate, Name, database.BASE_URL, magicLink, database.BASE_URL, database.BASE_URL)

	// Send email
	err := sendEmail(to, subject, body)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}
	log.Println("✅ Email sent successfully.")
	return nil
}

func SendEmailInvoiceForQSCoinsPurchase(to string, Name string, order_id string, coins_amount string, currency string, rate string) error {
	// Email content
	subject := fmt.Sprintf("Subject: ORDER-%s — Purchase of %s QS Coins\r\n", order_id, coins_amount)
	body := fmt.Sprintf(invoiceTemplate, Name, coins_amount, order_id, coins_amount, currency, rate, currency, rate, database.BASE_URL, database.BASE_URL, database.BASE_URL)

	// Send email
	err := sendEmail(to, subject, body)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}
	log.Println("✅ Email sent successfully.")
	return nil
}

func SendEmailNotificationOfUserLogin(to string, Name string, timestamp string, ipAddress string, userAgent string) error {
	// Email content
	subject := fmt.Sprintf("Subject: New Login Attempt Detected for %s\r\n", to)
	body := fmt.Sprintf(newLoginTemplate, Name, to, timestamp, ipAddress, userAgent, database.BASE_URL, database.BASE_URL, database.BASE_URL)

	// Send email
	err := sendEmail(to, subject, body)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}
	log.Println("✅ Email sent successfully.")
	return nil
}

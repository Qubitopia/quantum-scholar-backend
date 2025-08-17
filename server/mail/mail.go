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

	bodyTemplate := `<html>
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f4f7;
      margin: 0;
      padding: 0;
    }
    .container {
      max-width: 600px;
      background: #ffffff;
      margin: 40px auto;
      padding: 30px;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.05);
    }
    h2 {
      color: #333;
    }
    p {
      font-size: 16px;
      color: #555;
      line-height: 1.6;
    }
    .button-container {
      text-align: center;
      margin: 30px 0;
    }
    .button {
      background-color: #007BFF;
      color: white !important;
      padding: 14px 30px;
      text-decoration: none;
      border-radius: 6px;
      font-size: 16px;
      display: inline-block;
      font-weight: bold;
    }
    .footer {
      font-size: 12px;
      color: #999;
      text-align: center;
      margin-top: 40px;
    }
    .footer a {
      color: #007BFF;
      text-decoration: none;
    }
    @media (max-width: 600px) {
      .container {
        padding: 20px;
        margin: 20px;
      }
      .button {
        width: 100%%;
        box-sizing: border-box;
      }
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Hello %s,</h2>
    <p>Thank you for using <strong>Quantum Scholar</strong> by <strong>Qubitopia</strong>! We're thrilled to have you with us.</p>
    <p>To securely log in to your account, click the magic link below. This link is valid for the next <strong>15 minutes</strong>.</p>
    <div class="button-container">
      <a href="%s" class="button">Login to Quantum Scholar</a>
    </div>
    <p>If you didn't request this login, you can safely ignore this email or contact our support team.</p>
    <div class="footer">
      <p>You are receiving this email because you signed up for Quantum Scholar.<br />
        If you'd like to stop receiving these emails, <a href="https://quantumscholar.pages.dev/mail/unsubscribe">unsubscribe here</a>.
      </p>
      <p>Qubitopia Inc. | India | <a href="https://quantumscholar.pages.dev/privacypolicy">Privacy Policy</a></p>
    </div>
  </div>
</body>
</html>
`

	// Email content
	subject := "Subject: Welcome to Quantum Scholar by Qubitopia\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(bodyTemplate, Name, magicLink)
	msg := []byte(subject + mime + body)

	// Auth
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
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

	bodyTemplate := `<html>
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f4f7;
      margin: 0;
      padding: 0;
    }
    .container {
      max-width: 600px;
      background: #ffffff;
      margin: 40px auto;
      padding: 30px;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.05);
    }
    h2 {
      color: #333;
    }
    p {
      font-size: 16px;
      color: #555;
      line-height: 1.6;
    }
    .button-container {
      text-align: center;
      margin: 30px 0;
    }
    .button {
      background-color: #007BFF;
      color: white !important;
      padding: 14px 30px;
      text-decoration: none;
      border-radius: 6px;
      font-size: 16px;
      display: inline-block;
      font-weight: bold;
    }
    .footer {
      font-size: 12px;
      color: #999;
      text-align: center;
      margin-top: 40px;
    }
    .footer a {
      color: #007BFF;
      text-decoration: none;
    }
    @media (max-width: 600px) {
      .container {
        padding: 20px;
        margin: 20px;
      }
      .button {
        width: 100%%;
        box-sizing: border-box;
      }
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Welcome back, %s 👋</h2>
    <p>It's great to see you again on <strong>Quantum Scholar</strong> by <strong>Qubitopia</strong>! We're excited to continue this journey with you.</p>
    <p>For secure access to your account, click the link below. This link will remain valid for the next <strong>15 minutes</strong>.</p>
    <div class="button-container">
      <a href="%s" class="button">Login to Quantum Scholar</a>
    </div>
    <p>If you didn't request this login, please ignore this email or contact our support team immediately.</p>
    <div class="footer">
      <p>You are receiving this email because you signed up for Quantum Scholar.<br />
        If you'd like to stop receiving these emails, <a href="https://quantumscholar.pages.dev/mail/unsubscribe">unsubscribe here</a>.
      </p>
      <p>Qubitopia Inc. | India | <a href="https://quantumscholar.pages.dev/privacypolicy">Privacy Policy</a></p>
    </div>
  </div>
</body>
</html>
`

	// Email content
	subject := fmt.Sprintf("Subject: Welcome back, %s — your secure login link inside\r\n", Name)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(bodyTemplate, Name, magicLink)
	msg := []byte(subject + mime + body)

	// Auth
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Failed to send email:", err)
	} else {
		log.Println("✅ Email sent successfully.")
	}
}

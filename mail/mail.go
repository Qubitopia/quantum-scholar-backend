package mail

import (
	"fmt"
	"log"
	"net/smtp"

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
	newUserTemplate = `<html>
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
      <a href="%s%s" class="button">Login to Quantum Scholar</a>
    </div>
    <p>If you didn't request this login, you can safely ignore this email or contact our support team.</p>
    <div class="footer">
      <p>You are receiving this email because you signed up for Quantum Scholar.<br />
        If you'd like to stop receiving these emails, <a href="%s/mail/unsubscribe">unsubscribe here</a>.
      </p>
      <p>Qubitopia Inc. | India | <a href="%s/privacypolicy">Privacy Policy</a></p>
    </div>
  </div>
</body>
</html>`

	// Load Old User Email Template
	oldUserTemplate = `<html>
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
      <a href="%s%s" class="button">Login to Quantum Scholar</a>
    </div>
    <p>If you didn't request this login, please ignore this email or contact our support team immediately.</p>
    <div class="footer">
      <p>You are receiving this email because you signed up for Quantum Scholar.<br />
        If you'd like to stop receiving these emails, <a href="%s/mail/unsubscribe">unsubscribe here</a>.
      </p>
      <p>Qubitopia Inc. | India | <a href="%s/privacypolicy">Privacy Policy</a></p>
    </div>
  </div>
</body>
</html>`

	// Load Invoice Email Template
	invoiceTemplate = `<html>
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
    .invoice-box {
      margin-top: 20px;
      border: 1px solid #eee;
      border-radius: 6px;
      padding: 20px;
      background-color: #fafafa;
    }
    .invoice-box table {
      width: 100%%;
      border-collapse: collapse;
    }
    .invoice-box table td {
      padding: 10px;
      font-size: 15px;
      color: #444;
    }
    .invoice-box table tr.total td {
      border-top: 2px solid #007BFF;
      font-weight: bold;
      color: #222;
      font-size: 16px;
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
    <h2>Payment Invoice 💳</h2>
    <p>Hi %s,</p>
    <p>Thank you for purchasing %s QS Coins 🪙 on <strong>Quantum Scholar</strong>! Here are your order details:</p>

    <div class="invoice-box">
      <table>
        <tr>
          <td><strong>Order ID</strong></td>
          <td>%s</td>
        </tr>
        <tr>
          <td><strong>Purchase</strong></td>
          <td>%s QS Coins 🪙</td>
        </tr>
        <tr>
          <td><strong>Amount</strong></td>
          <td>%s %s</td>
        </tr>
        <tr>
          <td><strong>Status</strong></td>
          <td>Paid</td>
        </tr>
        <tr class="total">
          <td>Total</td>
          <td><strong>%s %s</strong></td>
        </tr>
      </table>
    </div>

    <div class="button-container">
      <a href="%s/billing" class="button">View Order Details</a>
    </div>

    <p>If you have any questions about this invoice, please reach out to our support team.</p>

    <div class="footer">
      <p>You are receiving this email because you purchased QS Coins on Quantum Scholar.<br />
        If you did not authorize this payment, please <a href="%s/support">contact support</a>.
      </p>
      <p>Qubitopia Inc. | India | <a href="%s/privacypolicy">Privacy Policy</a></p>
    </div>
  </div>
</body>
</html>`

	// Load New Login Notification Email Template
	newLoginTemplate = `<html>
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
      color: #d9534f;
    }
    p {
      font-size: 16px;
      color: #555;
      line-height: 1.6;
    }
    .details-box {
      margin-top: 20px;
      border: 1px solid #eee;
      border-radius: 6px;
      padding: 20px;
      background-color: #fafafa;
    }
    .details-box table {
      width: 100%%;
      border-collapse: collapse;
    }
    .details-box table td {
      padding: 10px;
      font-size: 15px;
      color: #444;
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
    <h2>🔐 New Login Attempt Detected</h2>
    <p>Hello %s,</p>
    <p>We noticed a login attempt to your <strong>Quantum Scholar</strong> account with the following details:</p>

    <div class="details-box">
      <table>
        <tr>
          <td><strong>Email</strong></td>
          <td>%s</td> <!-- User's email -->
        </tr>
        <tr>
          <td><strong>Time</strong></td>
          <td>%s</td> <!-- Timestamp -->
        </tr>
        <tr>
          <td><strong>IP Address</strong></td>
          <td>%s</td> <!-- Client IP -->
        </tr>
        <tr>
          <td><strong>Device / Browser</strong></td>
          <td>%s</td> <!-- User Agent -->
        </tr>
      </table>
    </div>

    <p>If this was you, you can safely ignore this email.</p>
    <p>If this wasn't you, we recommend you to contact our support team.</p>

    <div class="button-container">
      <a href="%s/support" class="button">Contact Support</a>
    </div>

    <div class="footer">
      <p>You are receiving this email from <strong>Quantum Scholar</strong> because of your registered account activity.<br />
        If you need assistance, please <a href="%s/support">contact support</a>.
      </p>
      <p>Qubitopia Inc. | India | <a href="%s/privacypolicy">Privacy Policy</a></p>
    </div>
  </div>
</body>
</html>`
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

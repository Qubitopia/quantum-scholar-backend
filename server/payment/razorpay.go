package payment

import (
	"log"
	"os"

	razorpay "github.com/razorpay/razorpay-go"
)

var RazorpayClient *razorpay.Client

func InitRazorpayClient() {
	key := os.Getenv("RZP_KEY_ID")
	secret := os.Getenv("RZP_KEY_SECRET")
	if key == "" || secret == "" {
		log.Fatal("Razorpay credentials not set in environment")
	}
	RazorpayClient = razorpay.NewClient(key, secret)
}

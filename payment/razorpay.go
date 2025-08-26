package payment

import (
	"github.com/Qubitopia/quantum-scholar-backend/database"
	razorpay "github.com/razorpay/razorpay-go"
)

var RazorpayClient *razorpay.Client

func InitRazorpayClient() {
	RazorpayClient = razorpay.NewClient(database.RZP_KEY_ID, database.RZP_KEY_SECRET)
}

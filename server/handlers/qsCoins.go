package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
	utils "github.com/razorpay/razorpay-go/utils"
)

type PurchaseQSCoinsINRRequest struct {
	QScoins uint64 `json:"qscoins" binding:"required,min=1"`
}

type PurchaseQSCoinsUSDRequest struct {
	QScoins uint64 `json:"qscoins" binding:"required,min=1"`
}

type RazorpayVerifyRequest struct {
	RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
	RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" binding:"required"`
}

func PurchaseQSCoinsINR(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	examiner, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	var req PurchaseQSCoinsINRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Razorpay credentials from environment
	key := os.Getenv("RZP_KEY_ID")
	secret := os.Getenv("RZP_KEY_SECRET")
	if key == "" || secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay credentials not set in environment"})
		return
	}

	// Create Razorpay client
	client := razorpay.NewClient(key, secret)

	// Step 1: Create payment record with pending status, no RazorpayPaymentID yet
	payment := models.PaymentTable{
		UserID:           examiner.ID,
		Amount:           int32(req.QScoins),
		Currency:         "INR",
		QSCoinsPurchased: int64(req.QScoins),
		PaymentStatus:    false,
		DateTime:         time.Now(),
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record", "details": err.Error()})
		return
	}

	// Step 2: Create Razorpay order using DB-generated OrderID as receipt
	receipt := "ORDER-" + fmt.Sprint(payment.OrderID)
	orderData := map[string]interface{}{
		"amount":          req.QScoins * 100, // INR to paise
		"currency":        "INR",
		"receipt":         receipt,
		"payment_capture": 1,
	}
	order, err := client.Order.Create(orderData, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order", "details": err.Error()})
		return
	}

	// Step 3: Update payment record with RazorpayPaymentID
	payment.RazorpayPaymentID = order["id"].(string)
	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment record", "details": err.Error()})
		return
	}

	// Return order info to frontend for payment processing
	c.JSON(http.StatusOK, gin.H{
		// "order":      order,
		"payment_id": payment.RazorpayPaymentID,
		"order_id":   payment.OrderID,
		"amount":     req.QScoins,
		"currency":   "INR",
	})
}

func PurchaseQSCoinsUSD(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	examiner, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	var req PurchaseQSCoinsINRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Razorpay credentials from environment
	key := os.Getenv("RZP_KEY_ID")
	secret := os.Getenv("RZP_KEY_SECRET")
	if key == "" || secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay credentials not set in environment"})
		return
	}

	// Create Razorpay client
	client := razorpay.NewClient(key, secret)

	// Step 1: Create payment record with pending status, no RazorpayPaymentID yet
	payment := models.PaymentTable{
		UserID:           examiner.ID,
		Amount:           int32(((req.QScoins * 100) / 75)),
		Currency:         "USD",
		QSCoinsPurchased: int64(req.QScoins),
		PaymentStatus:    false,
		DateTime:         time.Now(),
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record", "details": err.Error()})
		return
	}

	// Step 2: Create Razorpay order using DB-generated OrderID as receipt
	receipt := "ORDER-" + fmt.Sprint(payment.OrderID)
	orderData := map[string]interface{}{
		"amount":          ((req.QScoins * 100) / 75), // Tokens to USD
		"currency":        "USD",
		"receipt":         receipt,
		"payment_capture": 1,
	}
	order, err := client.Order.Create(orderData, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order", "details": err.Error()})
		return
	}

	// Step 3: Update payment record with RazorpayPaymentID
	payment.RazorpayPaymentID = order["id"].(string)
	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment record", "details": err.Error()})
		return
	}

	// Return order info to frontend for payment processing
	c.JSON(http.StatusOK, gin.H{
		// "order":      order,
		"payment_id": payment.RazorpayPaymentID,
		"order_id":   payment.OrderID,
		"amount":     ((req.QScoins * 100) / 75),
		"currency":   "USD",
	})
}

// Handler to verify payment and update coins
func VerifyRazorpayPayment(c *gin.Context) {
	var req RazorpayVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secret := os.Getenv("RZP_KEY_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay secret not set"})
		return
	}

	params := map[string]interface{}{
		"razorpay_order_id":   req.RazorpayOrderID,
		"razorpay_payment_id": req.RazorpayPaymentID,
	}

	// Verify signature
	if !utils.VerifyPaymentSignature(params, req.RazorpaySignature, secret) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment signature"})
		return
	}

	// Find payment record by RazorpayPaymentID
	var payment models.PaymentTable
	if err := database.DB.Where("razorpay_payment_id = ?", req.RazorpayOrderID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment record not found"})
		return
	}

	// Update payment status
	payment.PaymentStatus = true
	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	// Update user's QS coins
	var user models.User
	if err := database.DB.First(&user, payment.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	user.QSCoins += int64(payment.QSCoinsPurchased)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user coins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment verified and coins added"})
}

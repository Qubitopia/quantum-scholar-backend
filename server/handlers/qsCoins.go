package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/models"
	"github.com/Qubitopia/QuantumScholar/server/payment"

	"github.com/gin-gonic/gin"

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

	// Use global RazorpayClient
	if payment.RazorpayClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay client not initialized"})
		return
	}
	client := payment.RazorpayClient

	// Step 1: Create payment record with pending status, no RazorpayPaymentID yet
	payment := models.PaymentTable{
		UserID:           examiner.ID,
		Amount:           int32(req.QScoins),
		Currency:         "INR",
		QSCoinsPurchased: int64(req.QScoins),
		PaymentStatus:    "pending",
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

	// Step 3: Update payment record with RazorpayOrderID
	payment.RazorpayOrderID = order["id"].(string)
	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment record", "details": err.Error()})
		return
	}

	// Return order info to frontend for payment processing
	c.JSON(http.StatusOK, gin.H{
		// "order":      order,
		"razorpay_order_id": payment.RazorpayOrderID,
		"order_id":          payment.OrderID,
		"amount":            req.QScoins,
		"currency":          "INR",
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

	var req PurchaseQSCoinsUSDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use global RazorpayClient
	if payment.RazorpayClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay client not initialized"})
		return
	}
	client := payment.RazorpayClient

	amountInUSD := req.QScoins / 75 // 1 USD = 75 QSCoins

	// Step 1: Create payment record with pending status, no RazorpayPaymentID yet
	payment := models.PaymentTable{
		UserID:           examiner.ID,
		Amount:           int32(amountInUSD),
		Currency:         "USD",
		QSCoinsPurchased: int64(req.QScoins),
		PaymentStatus:    "pending",
		DateTime:         time.Now(),
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record", "details": err.Error()})
		return
	}

	// Step 2: Create Razorpay order using DB-generated OrderID as receipt
	receipt := "ORDER-" + fmt.Sprint(payment.OrderID)
	orderData := map[string]interface{}{
		"amount":          amountInUSD * 100, // USD to paise
		"currency":        "USD",
		"receipt":         receipt,
		"payment_capture": 1,
	}
	order, err := client.Order.Create(orderData, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order", "details": err.Error()})
		return
	}

	// Step 3: Update payment record with RazorpayOrderID
	payment.RazorpayOrderID = order["id"].(string)
	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment record", "details": err.Error()})
		return
	}

	// Return order info to frontend for payment processing
	c.JSON(http.StatusOK, gin.H{
		// "order":      order,
		"razorpay_order_id": payment.RazorpayOrderID,
		"order_id":          payment.OrderID,
		"amount":            amountInUSD * 100,
		"currency":          "USD",
	})
}

// Handler to verify payment and update coins
func VerifyRazorpayPayment(c *gin.Context) {
	var req RazorpayVerifyRequest
	log.Println(req)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find payment record by Razorpay order id (should be in RazorpayOrderID field)
	var payment models.PaymentTable
	if err := database.DB.Where("razorpay_order_id = ?", req.RazorpayOrderID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment record not found"})
		return
	}

	// Store the payload fields in DB (RazorpayPaymentID, RazorpaySignature)
	payment.RazorpayPaymentID = req.RazorpayPaymentID
	payment.RazorpaySignature = req.RazorpaySignature
	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment record with Razorpay details"})
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

	// Update payment status
	payment.PaymentStatus = "completed"
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

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/models"
	"github.com/Qubitopia/QuantumScholar/server/payment"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm/clause"

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

	// Use transaction and row-level locking to prevent race conditions
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var payment models.PaymentTable
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("razorpay_order_id = ?", req.RazorpayOrderID).First(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment record not found", "details": err.Error()})
		return
	}

	if payment.PaymentStatus == "completed" {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"message": "Payment already processed"})
		return
	}

	if payment.PaymentStatus != "pending" {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Payment is not pending nor completed"})
		return
	}

	// Store the payload fields in DB (RazorpayPaymentID, RazorpaySignature)
	payment.RazorpayPaymentID = req.RazorpayPaymentID
	payment.RazorpaySignature = req.RazorpaySignature

	secret := os.Getenv("RZP_KEY_SECRET")
	if secret == "" {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay secret not set"})
		return
	}

	params := map[string]interface{}{
		"razorpay_order_id":   req.RazorpayOrderID,
		"razorpay_payment_id": req.RazorpayPaymentID,
	}

	// Verify signature
	if !utils.VerifyPaymentSignature(params, req.RazorpaySignature, secret) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment signature"})
		return
	}

	// Update payment status
	payment.PaymentStatus = "completed"
	if err := tx.Save(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status", "details": err.Error()})
		return
	}

	// Update user's QS coins
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, payment.UserID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found", "details": err.Error()})
		return
	}
	user.QSCoins += int64(payment.QSCoinsPurchased)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user coins", "details": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Payment verified and coins added"})
}

// RazorpayWebhookHandler handles Razorpay webhook events for order.paid webhook
func RazorpayWebhookHandler(c *gin.Context) {
	// 1) Read raw body for signature verification
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "failed to read body")
		return
	}
	// Restore body for potential later use
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 2) Extract signature header
	signature := c.GetHeader("X-Razorpay-Signature")
	if signature == "" {
		c.String(http.StatusBadRequest, "missing X-Razorpay-Signature")
		return
	}

	// 3) Load webhook secret
	secret := os.Getenv("RZP_WEBHOOK_SECRET")
	if secret == "" {
		c.String(http.StatusInternalServerError, "webhook secret not configured")
		return
	}

	// 4) Verify signature
	if !utils.VerifyWebhookSignature(string(bodyBytes), signature, secret) {
		c.String(http.StatusUnauthorized, "signature verification failed")
		return
	}

	// 5) Parse JSON into a generic map
	var root map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &root); err != nil {
		c.String(http.StatusBadRequest, "invalid json")
		return
	}

	// 6) Extract required fields using safe map navigation
	get := func(m map[string]interface{}, keys ...string) (interface{}, bool) {
		var cur interface{} = m
		for _, k := range keys {
			asMap, ok := cur.(map[string]interface{})
			if !ok {
				return nil, false
			}
			v, ok := asMap[k]
			if !ok {
				return nil, false
			}
			cur = v
		}
		return cur, true
	}

	// Event check for order.paid
	if v, ok := root["event"]; ok {
		if s, _ := v.(string); s != "order.paid" {
			c.String(http.StatusOK, "ignored")
			return
		}
	}

	// Preferred order.entity fields for payment status and order_id
	var razorpay_order_id, razorpay_payment_id, payment_status string

	if v, ok := get(root, "payload", "order", "entity", "id"); ok {
		if s, ok := v.(string); ok {
			razorpay_order_id = s
		}
	}
	if v, ok := get(root, "payload", "order", "entity", "status"); ok {
		if s, ok := v.(string); ok {
			payment_status = s
		}
	}
	if v, ok := get(root, "payload", "payment", "entity", "id"); ok {
		if s, ok := v.(string); ok {
			razorpay_payment_id = s
		}
	}

	// 7) If payment_status is 'paid', update PaymentStatus and user's QSCoins if pending, using transaction and row lock
	if payment_status == "paid" && razorpay_order_id != "" {
		tx := database.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		var payment models.PaymentTable
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("razorpay_order_id = ?", razorpay_order_id).First(&payment).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment record not found", "details": err.Error()})
			return
		}
		if payment.PaymentStatus == "completed" {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"message": "Payment already processed"})
			return
		}
		if payment.PaymentStatus != "pending" {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Payment is not pending nor completed"})
			return
		}
		// Only process if pending
		payment.PaymentStatus = "completed"
		payment.RazorpayPaymentID = razorpay_payment_id
		if err := tx.Save(&payment).Error; err != nil {
			log.Printf("Failed to update payment status for order_id=%s: %v", razorpay_order_id, err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status", "details": err.Error()})
			return
		}
		var user models.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, payment.UserID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found", "details": err.Error()})
			return
		}
		user.QSCoins += int64(payment.QSCoinsPurchased)
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user coins", "details": err.Error()})
			return
		}
		tx.Commit()
		log.Println("Webhook processed successfully")
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

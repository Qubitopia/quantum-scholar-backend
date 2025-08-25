package handlers

import (
	"net/http"
	"time"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/models"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.(models.User),
	})
}

func UpdateProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser := user.(models.User)

	var updateData struct {
		Name        string `json:"name"`
		PublicEmail string `json:"public_email"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.Name != "" {
		currentUser.Name = updateData.Name
	}

	if updateData.PublicEmail != "" {
		currentUser.PublicEmail = updateData.PublicEmail
	}

	if err := database.DB.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    currentUser,
	})
}

func GetAllOrdersByUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser := user.(models.User)

	// Only select necessary fields from the database
	type OrderPartial struct {
		OrderID          uint32    `json:"order_id"`
		Amount           int32     `json:"amount"`
		Currency         string    `json:"currency"`
		QSCoinsPurchased int64     `json:"qs_coins_purchased"`
		PaymentStatus    string    `json:"payment_status"`
		DateTime         time.Time `json:"date_time"`
	}

	var orders []OrderPartial
	if err := database.DB.Model(&models.PaymentTable{}).
		Where("user_id = ?", currentUser.ID).
		Select("order_id, amount, currency, qs_coins_purchased, payment_status, date_time").
		Scan(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

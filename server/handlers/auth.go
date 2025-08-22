package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/mail"
	"github.com/Qubitopia/QuantumScholar/server/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

type AuthResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token,omitempty"`
	User    *models.User `json:"user,omitempty"`
}

func generateMagicToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateJWT(userID uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	return token.SignedString([]byte(database.JWT_SECRET))
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find or create user
	var user models.User
	result := database.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		// Create new user if not exists
		user = models.User{
			Email:       req.Email,
			PublicEmail: req.Email,
			Name:        req.Email,
			IsActive:    true,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Generate magic token
	token, err := generateMagicToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate magic token"})
		return
	}

	// Parse expiry duration
	expiryDuration, err := time.ParseDuration(database.MAGIC_LINK_EXPIRY)
	if err != nil {
		expiryDuration = 15 * time.Minute // Default to 15 minutes
	}

	// Create magic link record
	magicLink := models.MagicLink{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	if err := database.DB.Create(&magicLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create magic link"})
		return
	}

	if user.Name == user.Email {
		// Mail to new user or who has not updated their details
		magicLinkURL := fmt.Sprintf("/authNewUser/verify?token=%s", token)
		mail.SendEmailToNewUser(user.Email, user.Name, magicLinkURL)
	} else {
		// Mail to old user
		magicLinkURL := fmt.Sprintf("/auth/verify?token=%s", token)
		mail.SendEmailToOldUser(user.Email, user.Name, magicLinkURL)
	}

	c.JSON(http.StatusOK, AuthResponse{
		Message: "Magic link sent to your email",
	})
}

func VerifyMagicLink(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find magic link
	var magicLink models.MagicLink
	result := database.DB.Preload("User").Where("token = ? AND expires_at > ?",
		req.Token, time.Now()).First(&magicLink)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired magic link"})
		return
	}

	// Delete magic link after successful verification
	database.DB.Delete(&magicLink)

	// Generate JWT token
	jwtToken, err := generateJWT(magicLink.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Send login notification email
	timestamp := time.Now().Format("2006-01-02 15:04:05 MST")
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	mail.SendEmailNotificationOfUserLogin(magicLink.User.Email, magicLink.User.Name, timestamp, ipAddress, userAgent)

	c.JSON(http.StatusOK, AuthResponse{
		Message: "Login successful",
		Token:   jwtToken,
		User:    &magicLink.User,
	})
}

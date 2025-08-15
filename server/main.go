package main

import (
	"log"
	"strings"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/handlers"
	"github.com/Qubitopia/QuantumScholar/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	database.Connect()
	database.Migrate()

	// Initialize Gin router
	r := gin.Default()
	r.Use(gin.Recovery())
	r.TrustedPlatform = gin.PlatformCloudflare

	// Allowed origins
	allowedOrigins := map[string]bool{
		"https://quantumscholar.pages.dev": true,
		"http://localhost:3000":            true,
		"https://localhost:3000":           true,
		"http://127.0.0.1:3000":            true,
		"https://127.0.0.1:3000":           true,
	}

	// CORS middleware
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		} else if origin != "" {
			// Optionally, block disallowed origins explicitly
			// c.AbortWithStatus(403); return
			// Or just omit CORS headers; browser will block
			c.Header("Vary", "Origin")
		}

		if strings.EqualFold(c.Request.Method, "OPTIONS") {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to QuantumScholar API. Visit https://github.com/Qubitopia/QuantumScholar for more information."})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	auth := r.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/verify", handlers.VerifyMagicLink)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Profile
		api.GET("/profile", handlers.GetProfile)
		api.PUT("/profile", handlers.UpdateProfile)

		// QS Coins purchase and verification
		api.POST("/purchase-qscoins-inr", handlers.PurchaseQSCoinsINR)
		api.POST("/purchase-qscoins-usd", handlers.PurchaseQSCoinsUSD)
		api.POST("/verify-razorpay-payment", handlers.VerifyRazorpayPayment)

		// Test
		api.POST("/test/create", handlers.CreateNewTest)
		api.PUT("/test/update", handlers.UpdateTest)

	}

	// Start server
	log.Printf("Server starting on port set in varible PORT in .env")
	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

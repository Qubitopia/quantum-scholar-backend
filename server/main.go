package main

import (
	"log"
	"os"

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

	// CORS middleware (if needed)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
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
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

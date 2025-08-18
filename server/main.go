package main

import (
	"log"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/handlers"
	"github.com/Qubitopia/QuantumScholar/server/middleware"
	"github.com/Qubitopia/QuantumScholar/server/payment"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to PostgreSQL
	database.ConnectPgsql()
	database.MigratePgsql()

	// Connect to Redis
	database.ConnectRedis()

	// Initialize Razorpay client
	payment.InitRazorpayClient()

	// Initialize Gin router
	r := gin.Default()
	r.Use(gin.Recovery())
	r.TrustedPlatform = gin.PlatformCloudflare

	// Rate limiting middleware
	r.Use(middleware.RateLimitMiddleware())

	// CORS middleware
	r.Use(middleware.CORSMiddleware())

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

		// Orders
		api.GET("/orders", handlers.GetAllOrdersByUser)

		// QS Coins purchase and verification
		api.POST("/purchase-qscoins-inr", handlers.PurchaseQSCoinsINR)
		api.POST("/purchase-qscoins-usd", handlers.PurchaseQSCoinsUSD)
		api.POST("/verify-razorpay-payment", handlers.VerifyRazorpayPayment)

		// Test
		api.POST("/test/create", handlers.CreateNewTest)
		api.PUT("/test/update-que-ans", handlers.UpdateQuestionsAndAnswersInTest)
		api.GET("/test", handlers.GetAllTestsCreatedByUser)

	}

	// Start server
	log.Printf("Server starting on port set in varible PORT in .env")
	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"log"

	"github.com/Qubitopia/quantum-scholar-backend/database"
	"github.com/Qubitopia/quantum-scholar-backend/handlers"
	"github.com/Qubitopia/quantum-scholar-backend/mail"
	"github.com/Qubitopia/quantum-scholar-backend/middleware"
	"github.com/Qubitopia/quantum-scholar-backend/payment"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load environment variables into global variables
	database.LoadEnvVariables()

	// Connect to PostgreSQL
	database.ConnectPgsql()
	database.MigratePgsql()

	// Connect to Redis
	database.ConnectRedis()

	// Initialize Cloudflare R2 (S3) client
	if err := database.InitR2Client(); err != nil {
		log.Fatal("Failed to initialize R2 client:", err)
	}

	// Initialize Razorpay client
	payment.InitRazorpayClient()

	// Initialize email
	mail.LoadEmailTemplates()
	mail.InitEmail()

	// Initialize Gin router
	r := gin.Default()
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
		api.GET("/test/:id", handlers.GetTestByID)
		api.PUT("/test/add-candidates", handlers.AddCandidatesToTest)

		// Image upload
		api.POST("/upload-image", handlers.UploadImage)
		api.GET("/image-url/:imagename", handlers.GetImageURL)

	}

	webhook := r.Group("/webhook")
	{
		webhook.POST("/razorpay", handlers.RazorpayWebhookHandler)
	}

	// Start server
	log.Printf("Server starting on port set in variable PORT in .env")
	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var AllowedOrigins = map[string]bool{
	"https://quantumscholar.pages.dev":     true,
	"https://dev.quantumscholar.pages.dev": true,
	"http://localhost:3000":                true,
	"https://localhost:3000":               true,
	"http://127.0.0.1:3000":                true,
	"https://127.0.0.1:3000":               true,
	"http://localhost:5500":                true,
	"https://localhost:5500":               true,
	"http://127.0.0.1:5500":                true,
	"https://127.0.0.1:5500":               true,
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if AllowedOrigins[origin] {
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
	}
}

package middleware

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	rateLimit := 20
	rateLimitStr := os.Getenv("API_RATE_LIMIT_PER_MINUTE")
	if rateLimitStr != "" {
		if v, err := strconv.Atoi(rateLimitStr); err == nil {
			rateLimit = v
		}
	}
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := "rl:" + ip
		count, err := database.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Rate limiter error"})
			return
		}
		if count == 1 {
			database.RedisClient.Expire(ctx, key, time.Minute)
		}
		if int(count) > rateLimit {
			c.AbortWithStatusJSON(429, gin.H{"error": "Too many requests. Please try again later."})
			return
		}
		c.Next()
	}
}

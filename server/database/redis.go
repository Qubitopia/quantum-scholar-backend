package database

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisAddr := os.Getenv("DB_HOST")
	if redisAddr == "" {
		redisAddr = "redis:6379" // fallback default
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		redisAddr = "redis:" + port
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}

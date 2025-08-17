package database

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	redisAddr := redisHost + ":" + redisPort

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test connection
	if err := RedisClient.Ping(RedisClient.Context()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis: " + err.Error())
	}
}

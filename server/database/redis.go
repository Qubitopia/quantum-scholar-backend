package database

import (
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis() {
	// Use global variables from environmentVariable.go
	redisAddr := REDIS_HOST + ":" + REDIS_PORT

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test connection
	if err := RedisClient.Ping(RedisClient.Context()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis: " + err.Error())
	}
}

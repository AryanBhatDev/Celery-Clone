package main

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	redisConnString := os.Getenv("REDIS_URL")

	if redisConnString == "" {
		log.Fatal("Redis url not found")
	}

	opt, err := redis.ParseURL(redisConnString)

	if err != nil {
		log.Fatal("error while parsing url", err)
	}

	client := redis.NewClient(opt)
	return client
}

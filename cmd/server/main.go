package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	api "github.com/harrymitchinson/dynamic-inventory-svc/pkg/api"
	"go.uber.org/zap"
)

// List of default variables
var (
	RedisHost     string = "localhost"
	RedisPort     string = "6379"
	RedisUsername string = ""
	RedisPassword string = ""
	HTTPPort      string = "8080"
)

func main() {
	if v, ok := os.LookupEnv("REDIS_HOST"); ok {
		RedisHost = v
	}
	if v, ok := os.LookupEnv("REDIS_PORT"); ok {
		RedisPort = v
	}
	if v, ok := os.LookupEnv("REDIS_USERNAME"); ok {
		RedisUsername = v
	}
	if v, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		RedisPassword = v
	}
	if v, ok := os.LookupEnv("HTTP_PORT"); ok {
		HTTPPort = v
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", RedisHost, RedisPort),
		Username: RedisUsername,
		Password: RedisPassword,
	})
	defer redis.Close()
	if err = redis.Ping(redis.Context()).Err(); err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}

	// Setup the API server.
	if err := api.Setup(&api.Builder{
		Redis:  redis,
		Logger: logger,
	}).Run(fmt.Sprintf(":%s", HTTPPort)); err != nil {
		logger.Fatal("api server failed", zap.Error(err))
	}
}

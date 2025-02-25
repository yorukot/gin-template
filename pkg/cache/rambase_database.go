package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/yorukot/go-template/pkg/logger"
)

var (
	RedisClient  *redis.Client
	RedisLimiter *redis.Client
)

// InitializeRedis initializes the Redis client.
func init() {
	// Read Redis database number.
	limiterDB := 1
	normalDB := 0

	host := os.Getenv("CACHE_HOST")
	port := os.Getenv("CACHE_PORT")
	if host == "" || port == "" {
		logger.Log.Fatal("CACHE_HOST or CACHE_PORT is not set")
	}

	url := fmt.Sprintf("%s:%s", host, port)

	// Set Redis options.
	options := &redis.Options{
		Addr:         url,
		Password:     os.Getenv("CACHE_PASSWORD"),
		DB:           normalDB,
		PoolSize:     10,
		MinIdleConns: 1,
	}

	RedisClient = redis.NewClient(options)

	options = &redis.Options{
		Addr:         url,
		Password:     os.Getenv("CACHE_PASSWORD"),
		DB:           limiterDB,
		PoolSize:     10,
		MinIdleConns: 1,
	}

	RedisLimiter = redis.NewClient(options)

	// Test connection
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	if err := RedisLimiter.Ping(context.Background()).Err(); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
}

package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"create-paste-service/model"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient        *redis.Client
	ctx                = context.Background()
	RedisCachePasteTTL time.Duration
)

// InitRedis initializes the Redis client connection
func InitRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = "" // No password by default
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         redisHost,
		Password:     redisPassword,
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 20,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// Get TTL from environment in seconds
	ttlStr := os.Getenv("REDIS_CACHE_PASTE_TTL_SECONDS")
	if ttlStr != "" {
		ttlSec, err := strconv.Atoi(ttlStr)
		if err != nil {
			log.Printf("Invalid REDIS_CACHE_PASTE_TTL_SECONDS: %s. Using default 3600s (1h)", ttlStr)
			RedisCachePasteTTL = 3600 * time.Second
		} else {
			RedisCachePasteTTL = time.Duration(ttlSec) * time.Second
		}
	} else {
		RedisCachePasteTTL = 3600 * time.Second // Default: 1 hour
	}

	// Test Redis connection
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return err
	}

	log.Println("Successfully connected to Redis")
	return nil
}

// CloseRedis closes the Redis client connection
func CloseRedis() {
	if RedisClient != nil {
		_ = RedisClient.Close()
	}
}

// CachePaste caches a paste in Redis with proper TTL
func CachePaste(paste *model.Paste) error {
	pasteJSON, err := json.Marshal(paste)
	if err != nil {
		log.Printf("Error marshaling paste for Redis: %v", err)
		return err
	}

	var expiration time.Duration
	if paste.ExpiresAt != nil {
		expiration = time.Until(*paste.ExpiresAt)
		if expiration <= 0 {
			// Don't cache expired pastes
			return nil
		}
	} else {
		expiration = RedisCachePasteTTL
	}

	err = RedisClient.Set(ctx, "paste:"+paste.ID, pasteJSON, expiration).Err()
	if err != nil {
		log.Printf("Error caching paste ID %s: %v", paste.ID, err)
		return err
	}

	log.Printf("Successfully cached paste ID: %s", paste.ID)
	return nil
}

// GetPaste retrieves a paste from Redis cache
func GetPaste(id string) (*model.Paste, error) {
	pasteJSON, err := RedisClient.Get(ctx, "paste:"+id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		log.Printf("Error retrieving paste ID %s: %v", id, err)
		return nil, err
	}

	var paste model.Paste
	err = json.Unmarshal([]byte(pasteJSON), &paste)
	if err != nil {
		log.Printf("Error unmarshaling cached paste ID %s: %v", id, err)
		return nil, err
	}

	return &paste, nil
}

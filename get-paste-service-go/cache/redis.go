package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"get-paste-service/model"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient        *redis.Client
	ctx                = context.Background()
	RedisCachePasteTTL time.Duration
)

// InitRedis initializes the Redis client
func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         redisHost,
		Password:     redisPassword,
		DB:           0,
		PoolSize:     500,
		MinIdleConns: 100,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// TTL setting from env in seconds
	ttlStr := os.Getenv("REDIS_CACHE_PASTE_TTL_SECONDS")
	if ttlStr != "" {
		ttlSec, err := strconv.Atoi(ttlStr)
		if err != nil {
			log.Printf("Invalid REDIS_CACHE_PASTE_TTL_SECONDS: %s. Using default 3600s", ttlStr)
			RedisCachePasteTTL = 3600 * time.Second
		} else {
			RedisCachePasteTTL = time.Duration(ttlSec) * time.Second
		}
	} else {
		RedisCachePasteTTL = 3600 * time.Second // Default TTL: 1 hour
	}

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Could not connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
}

// GetPasteFromCache attempts to retrieve a paste from Redis
func GetPasteFromCache(pasteID string) (*model.Paste, error) {
	val, err := RedisClient.Get(ctx, "paste:"+pasteID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		log.Printf("Error getting paste ID %s from Redis: %v", pasteID, err)
		return nil, err
	}

	var paste model.Paste
	if err := json.Unmarshal([]byte(val), &paste); err != nil {
		log.Printf("Error unmarshaling paste ID %s from Redis: %v", pasteID, err)
		return nil, err
	}

	return &paste, nil
}

// SetPasteToCache stores a paste in Redis with appropriate TTL
func SetPasteToCache(paste *model.Paste) error {
	ttl := RedisCachePasteTTL

	if paste.ExpiresAt != nil {
		expiresIn := time.Until(*paste.ExpiresAt)
		if expiresIn > 0 {
			ttl = expiresIn
		} else {
			return nil // Don't cache expired pastes
		}
	}

	pasteJSON, err := json.Marshal(paste)
	if err != nil {
		log.Printf("Error marshaling paste ID %s for Redis: %v", paste.ID, err)
		return err
	}

	err = RedisClient.Set(ctx, "paste:"+paste.ID, pasteJSON, ttl).Err()
	if err != nil {
		log.Printf("Error setting paste ID %s to Redis: %v", paste.ID, err)
		return err
	}

	log.Printf("Cached paste ID %s to Redis with TTL %s", paste.ID, ttl)
	return nil
}

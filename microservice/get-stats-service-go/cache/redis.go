package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"get-stats-service/model"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient             *redis.Client
	ctx                     = context.Background()
	TTLOtherMonth           time.Duration
	defaultTTLOtherMonthSec = 60
)

// InitRedis initializes Redis client and TTL values
func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Could not connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}

	// Load TTLs from environment
	TTLOtherMonth = loadTTLFromEnv("REDIS_CACHE_STATS_TTL_OTHER_MONTH", defaultTTLOtherMonthSec)
}

func loadTTLFromEnv(envKey string, defaultSeconds int) time.Duration {
	val := os.Getenv(envKey)
	if val == "" {
		return time.Duration(defaultSeconds) * time.Second
	}

	sec, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Invalid %s: %s. Using default %ds", envKey, val, defaultSeconds)
		return time.Duration(defaultSeconds) * time.Second
	}

	return time.Duration(sec) * time.Second
}

// GetMonthlyStatsFromCache retrieves monthly stats from Redis
func GetMonthlyStatsFromCache(year, month int) (*model.MonthlyStats, error) {
	key := fmt.Sprintf("stats:monthly:%d-%02d", year, month)

	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		log.Printf("Error getting stats from cache: %v", err)
		return nil, err
	}

	var stats model.MonthlyStats
	if err := json.Unmarshal([]byte(val), &stats); err != nil {
		log.Printf("Error unmarshaling cached stats: %v", err)
		return nil, err
	}

	log.Printf("Cache hit: stats for %d-%02d", year, month)
	return &stats, nil
}

// SetMonthlyStatsToCache stores monthly stats in Redis
func SetMonthlyStatsToCache(year, month int, stats *model.MonthlyStats) error {
	key := fmt.Sprintf("stats:monthly:%d-%02d", year, month)

	ttl := TTLOtherMonth

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Error marshaling stats for %d-%02d: %v", year, month, err)
		return err
	}

	err = RedisClient.Set(ctx, key, statsJSON, ttl).Err()
	if err != nil {
		log.Printf("Error caching stats for %d-%02d: %v", year, month, err)
		return err
	}

	log.Printf("Cached stats for %d-%02d with TTL %s", year, month, ttl)
	return nil
}

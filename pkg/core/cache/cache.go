package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/config"
	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Remove(ctx context.Context, key string) error
	CountKeys(ctx context.Context) (int64, error)
	ClearPattern(ctx context.Context, pattern string) (int64, error)
	Close() error
}

// RedisCacheService implements CacheService using Redis
type RedisCacheService struct {
	client *redis.Client
}

// NewRedisCacheService creates a new instance of RedisCacheService
func NewRedisCacheService(ctx context.Context) (*RedisCacheService, error) {
	// Get Redis server address and password from environment variables
	if config.GlobalConfig.IsRedis {
		redisURI := config.GlobalConfig.RedisURI
		if redisURI == "" {
			redisURI = "localhost:6379" // Default Redis server address
		}
		redisPassword := config.GlobalConfig.RedisPassword
		if redisPassword == "" {
			redisPassword = "" // No password if not provided
		}

		redisDB := config.GlobalConfig.RedisDB
		if redisDB == -1 {
			redisDB = 0
		}

		// Create a new Redis client
		client := redis.NewClient(&redis.Options{
			Addr:     redisURI,
			Password: redisPassword,
			DB:       redisDB,
		})

		// Ping the Redis server to ensure connectivity
		if err := client.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("failed to ping Redis server: %w", err)
		}

		return &RedisCacheService{
			client: client,
		}, nil
	} else {
		return nil, nil
	}
}

// Get retrieves value from cache by key
func (svc *RedisCacheService) Get(ctx context.Context, key string) (string, error) {
	if config.GlobalConfig.IsRedis {
		val, err := svc.client.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return "", nil // Cache miss
			}
			return "", fmt.Errorf("failed to get value from cache: %w", err)
		}
		return val, nil
	} else {
		return "", nil
	}
}

// Set sets value in cache with specified key
func (svc *RedisCacheService) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	if config.GlobalConfig.IsRedis {
		err := svc.client.Set(ctx, key, value, expiration).Err()
		if err != nil {
			return fmt.Errorf("failed to set value in cache: %w", err)
		}
		return nil
	} else {
		return nil
	}
}

// Remove implements CacheService.
func (svc *RedisCacheService) Remove(ctx context.Context, key string) error {
	// Use context with timeout to prevent blocking indefinitely
	if config.GlobalConfig.IsRedis {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err := svc.client.Del(ctx, key).Err()
		if err != nil {
			return fmt.Errorf("failed to remove value from cache: %w", err)
		}

		return nil
	} else {
		return nil
	}
}

// CountKeys counts the number of keys in the Redis cache
func (svc *RedisCacheService) CountKeys(ctx context.Context) (int64, error) {
	// Use SCAN command to iterate over keys in the cache
	if config.GlobalConfig.IsRedis {
		var cursor uint64 = 0
		var keysCount int64 = 0

		for {
			// Scan keys with cursor and pattern
			keys, nextCursor, err := svc.client.Scan(ctx, cursor, "*", 100).Result()
			if err != nil {
				return 0, fmt.Errorf("failed to scan keys in cache: %w", err)
			}

			// Increment keys count
			keysCount += int64(len(keys))

			// Update cursor for next iteration
			cursor = nextCursor

			// Break if iteration is complete
			if cursor == 0 {
				break
			}
		}

		return keysCount, nil
	} else {
		return 0, nil
	}
}

func (svc *RedisCacheService) ClearPattern(ctx context.Context, pattern string) (int64, error) {
	// Use context with timeout to prevent blocking indefinitely
	if !config.GlobalConfig.IsRedis {
		return 0, nil
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(config.GlobalConfig.RedisExp)*time.Second)
	defer cancel()

	// Use SCAN command to iterate over keys in the cache matching the pattern
	var cursor uint64 = 0
	var deletedKeysCount int64 = 0

	for {
		// Scan keys with cursor and pattern
		keys, nextCursor, err := svc.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return deletedKeysCount, fmt.Errorf("failed to scan keys in cache: %w", err)
		}

		// Delete keys matching the pattern
		if len(keys) > 0 {
			deletedCount, err := svc.client.Del(ctx, keys...).Result()
			if err != nil {
				return deletedKeysCount, fmt.Errorf("failed to delete keys in cache: %w", err)
			}
			deletedKeysCount += deletedCount
		}

		// Update cursor for next iteration
		cursor = nextCursor

		// Break if iteration is complete
		if cursor == 0 {
			break
		}
	}

	return deletedKeysCount, nil
}

// Close closes the Redis client
func (svc *RedisCacheService) Close() error {
	if config.GlobalConfig.IsRedis {
		return svc.client.Close()
	} else {
		return nil
	}
}

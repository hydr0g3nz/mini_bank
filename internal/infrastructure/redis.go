package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}
type CacheConfig struct {
	Host     string
	Port     int
	Password string
	Db       int
}

// NewRedisClient creates a new Redis client instance
func NewRedisClient(cfg CacheConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:       cfg.Password,
		DB:             cfg.Db,
		MaxActiveConns: 0,
	})
	fmt.Println("redis", fmt.Sprintf("%s:%d  %s", cfg.Host, cfg.Port, cfg.Password))
	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %w", err))
	}

	return &RedisClient{client: client}
}

// Set stores a value with expiration
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key does not exist: %s", key)
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	return json.Unmarshal(data, dest)
}

// Delete removes a key
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// HashSet stores a hash field
func (r *RedisClient) HashSet(ctx context.Context, key, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.HSet(ctx, key, field, data).Err()
}

// HashGet retrieves a hash field
func (r *RedisClient) HashGet(ctx context.Context, key, field string, dest interface{}) error {
	data, err := r.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return fmt.Errorf("failed to get hash field: %w", err)
	}

	return json.Unmarshal(data, dest)
}

// SetNX sets a value if the key doesn't exist (useful for distributed locks)
func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.SetNX(ctx, key, data, expiration).Result()
}

// Incr increments a key's value
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

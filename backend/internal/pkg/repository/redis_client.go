package repository

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

func NewRedisClient(
	ctx context.Context, host string, port string,
) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	// create connection
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// check connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	fmt.Println("redis connection established")

	// return structure
	return &RedisClient{
		client: rdb,
	}, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) SetJSON(
	ctx context.Context,
	key string, value interface{},
	ttl time.Duration,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisClient) SetString(
	ctx context.Context,
	key string, value string,
	ttl time.Duration,
) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// target must be a pointer to structure
func (r *RedisClient) GetJSON(
	ctx context.Context, key string, target any,
) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), target); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetJSONArray(
	ctx context.Context, key string, target any,
) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), target); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) Delete(
	ctx context.Context, key string,
) error {
	return r.client.Del(ctx, key).Err()
}

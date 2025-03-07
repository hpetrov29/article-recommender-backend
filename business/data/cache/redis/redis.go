package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hpetrov29/resttemplate/business/data/cache"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	c *redis.Client
}

func (rc *RedisClient) Open(ctx context.Context, cfg cache.Config) error {
	rc.c = redis.NewClient(&redis.Options{
        Addr:     cfg.Host,
        Password: cfg.Password, // no password set
        DB:       cfg.DbName,  // use default DB
    })

	_, err := rc.c.Ping(ctx).Result()
    if err != nil {
        return fmt.Errorf("could not connect to Redis: %v", err)
    }

    return nil
}

func (rc *RedisClient) StatusCheck(ctx context.Context) error {
	_, err := rc.c.Ping(ctx).Result()
    if err != nil {
        return fmt.Errorf("could not connect to Redis: %v", err)
    }

    return nil
}

func (rc *RedisClient) Close() error {
	if err := rc.c.Close(); err != nil {
		return fmt.Errorf("error closing Redis connection: %w", err)
	}
	return nil
}

func (rc *RedisClient) Set(ctx context.Context, key string, value []byte) error {
	if err := rc.c.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("error inserting key-value pair with key '%s' into Redis: %w", key, err)
	}
	return nil
}

func (rc *RedisClient) SetWithTTL(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := rc.c.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("error inserting key-value pair with key '%s' into Redis with TTL '%s': %w", key, ttl, err)
	}
	return nil
}

func (rc *RedisClient) GetNonFatal(ctx context.Context, key string) ([]byte, bool, error) {
	value, err := rc.c.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, fmt.Errorf("error retrieving key '%s' from Redis: %w", key, err)
	}

	return []byte(value), true, nil
}

func (rc *RedisClient) GetFatal(ctx context.Context, key string) ([]byte, error) {
	value, err := rc.c.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key '%s' not found in Redis: %w", key, err)
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving key '%s' from Redis: %w", key, err)
	}

	return []byte(value), nil
}
package redis

import (
	"context"
	"errors"
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

func (rc *RedisClient) HSetWithTTL(ctx context.Context, key string, pairs map[string]interface{}, ttl time.Duration) error {
	if (ttl <= 0) { 
		return errors.New("expiry must be a positive value")
	}

	err := rc.c.HSet(ctx, key, pairs).Err()
	if err != nil {
		return fmt.Errorf("failed to set hash fields for key %s: %w", key, err)
	}

	err = rc.c.Expire(ctx, key, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set TTL for key %s: %w", key, err)
	}

	return nil
}

func (rc *RedisClient) HVals(ctx context.Context, key string, field string) ([]string, error) {
	values, err := rc.c.HVals(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get values from hash for key %s: %w", key, err)
	}
	return values, nil
}

func (rc *RedisClient) HSetField(ctx context.Context, key string, field string, value interface{}) error {
	err := rc.c.HSet(ctx, key, field, value).Err()
	if err != nil {
		return fmt.Errorf("failed to set field %s in hash for key %s: %w", field, key, err)
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
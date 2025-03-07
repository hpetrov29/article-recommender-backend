package cache

import (
	"context"
	"time"
)

type Config struct {
	Password 	string
	Host     	string
	DbName 	 	int
}

type Cache interface {
	Open(ctx context.Context, cfg Config) error
	StatusCheck(ctx context.Context) error
	Close() error
	Set(ctx context.Context, key string, value []byte) error
	SetWithTTL(ctx context.Context, key string, value []byte, ttl time.Duration) error
	GetNonFatal(ctx context.Context, key string) ([]byte, bool, error)
	GetFatal(ctx context.Context, key string) ([]byte, error)
}
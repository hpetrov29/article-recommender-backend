package dbnosql

import "context"

type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxOpenConns int
}

type NOSQLDB interface {
	Open(cfg Config) error
	StatusCheck(ctx context.Context) error
	Close() error
	GetRepository(string) NOSQLDBrepo
}

type NOSQLDBrepo interface {
	Insert(ctx context.Context, record string) error
	Query(ctx context.Context, record string) (string, error)
	Delete(ctx context.Context, record string) error
}
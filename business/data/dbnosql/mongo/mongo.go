package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxOpenConns int
}

func Open(cfg Config) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", cfg.User, cfg.Password, cfg.Host)

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(uint64(cfg.MaxOpenConns))

	client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

	return client, nil
}

// StatusCheck returns nil if it can successfully talk to the MongoDB database.
// It returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, c *mongo.Client) error {
	// Ping the MongoDB server to check connectivity.
	err := c.Ping(ctx, nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			time.Sleep(5 * time.Second)
			return StatusCheck(ctx, c)
		}
		return err
	}
	return nil
}

func Close(client *mongo.Client) error {
    if client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := client.Disconnect(ctx); err != nil {
            return fmt.Errorf("error closing MongoDB connection: %w", err)
		}
		return nil
    }
	return errors.New("mongoDB client was a null pointer")
}
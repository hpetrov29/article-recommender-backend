package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hpetrov29/resttemplate/business/data/dbnosql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------------------------------------------------------------------------
// MongoClient implements the NOSQLDB interface

type MongoClient struct {
	c *mongo.Database
}

func (mc *MongoClient) Open(cfg dbnosql.Config) (error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", cfg.User, cfg.Password, cfg.Host)

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(uint64(cfg.MaxOpenConns))

	client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %w", err)
    }
	mc.c = client.Database(cfg.Name)
	return nil
}

// StatusCheck returns nil if it can successfully talk to the MongoDB database.
// It returns a non-nil error otherwise.
func (mc *MongoClient) StatusCheck(ctx context.Context) error {
	// Ping the MongoDB server to check connectivity.
	err := mc.c.Client().Ping(ctx, nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			time.Sleep(5 * time.Second)
			return mc.StatusCheck(ctx)
		}
		return err
	}
	return nil
}

func (mc *MongoClient) Close() error {
    if mc.c != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := mc.c.Client().Disconnect(ctx); err != nil {
            return fmt.Errorf("error closing MongoDB connection: %w", err)
		}
		return nil
    }
	return errors.New("mongoDB client is not initialized")
}

// -------------------------------------------------------------------------
// mongoRepository implements the interface required by the core api logic

type MongoRepository struct {
    collection *mongo.Collection
}

func (mc *MongoClient) GetRepository(collectionName string) dbnosql.NOSQLDBrepo {
	collection := mc.c.Collection(collectionName)
    return &MongoRepository{collection}
}

// Insert inserts a new record into the MongoDB collection.
func (r *MongoRepository) Insert(ctx context.Context, record interface{}) error {
    _, err := r.collection.InsertOne(ctx, record)
	if err != nil {
        return fmt.Errorf("failed to insert record in mongoDB: %w", err)
	}
    return nil
}

// QueryById retrieves a record with the specified id from the MongoDB collection.
func (r *MongoRepository) QueryById(ctx context.Context, id int64, data any) error {
	if data == nil {
        return errors.New("(*MongoRepository) QueryById expects data to be a non-nil pointer")
    }

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(data)
    if errors.Is(err, mongo.ErrNoDocuments) {
        return fmt.Errorf("document not found in mongoDB: %w", err)
    }

	return err
}

// Delete deletes a record from the MongoDB collection.
func (r *MongoRepository) Delete(ctx context.Context, id uint64) error {
    res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err 
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("document with id %d not found in mongoDB", id)
	}

    return nil
}
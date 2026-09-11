package repository

import (
	"context"
	"fmt"
	"time"

	"log-aggregator/internal/config"
	"log-aggregator/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository interface {
	InsertBatch(ctx context.Context, logs []model.Log) error
	Close(ctx context.Context) error
}

type MongoRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoRepository(ctx context.Context , cfg *config.Config) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}
	db := client.Database(cfg.MongoDBName)
	collection := db.Collection(cfg.MongoCol)

	return &MongoRepository{
		client:     client,
		db:         db,
		collection: collection,
	}, nil

}

func (r *MongoRepository) InsertBatch(ctx context.Context, logs []model.Log) error {
	if len(logs) == 0 {
		return nil
	}

	documents := make([]interface{}, len(logs))
	for i, log := range logs {
		documents[i] = log
	}
	_, err := r.collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert batch logs: %w", err)
	}

	return nil
}

func (r *MongoRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}

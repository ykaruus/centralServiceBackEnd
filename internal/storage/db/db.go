package db

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ClientDB struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewClientDB(url string, dbName string) (*ClientDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(url))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}
	slog.Info("connected to database", "layer", "db-storage")
	return &ClientDB{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

func (c *ClientDB) LoadCollection(collname string) *mongo.Collection {
	return c.database.Collection(collname)
}

func (c *ClientDB) CreateUniqueIndexes(coll *mongo.Collection, indexModels []mongo.IndexModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	_, err := coll.Indexes().CreateMany(ctx, indexModels)

	if err != nil {
		return err
	}

	return err
}

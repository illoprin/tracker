package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoClient(ctx context.Context, uri, dbName string) (*MongoClient, error) {
	// create options
	clientOptions := options.Client().ApplyURI(uri)

	// connect
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	// check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	fmt.Println("mongodb connection established")

	return &MongoClient{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

func (c *MongoClient) Disconnect(ctx context.Context) error {
	if err := c.Client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

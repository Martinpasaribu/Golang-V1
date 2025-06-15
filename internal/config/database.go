package config

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client     *mongo.Client
	DB         *mongo.Database
	MongoCtx   context.Context
)

func InitMongoDB() error {
	MongoCtx = context.Background()
	
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(MongoCtx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(MongoCtx, nil)
	if err != nil {
		return err
	}

	Client = client
	DB = client.Database("projectdb")
	return nil
}

func CloseMongoDB() {
	if Client != nil {
		Client.Disconnect(MongoCtx)
	}
}
package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	DB       *mongo.Database
	MongoCtx context.Context
)

func InitMongoDB() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found. Using system environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	// Fallback to default if not set
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	
	if dbName == "" {
		dbName = "Server-go"
		log.Println("Using default DB name: Server-go")
	}

	MongoCtx = context.Background()
	
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(MongoCtx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return err
	}

	err = client.Ping(MongoCtx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return err
	}

	Client = client
	DB = client.Database(dbName)
	log.Printf("✅ Connected to MongoDB! Database: %s", dbName)
	return nil
}

func CloseMongoDB() {
	if Client != nil {
		if err := Client.Disconnect(MongoCtx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		} else {
			log.Println("🚪 MongoDB connection closed")
		}
	}
}

func GetDB() *mongo.Database {
	if DB == nil {
		log.Println("⚠️ Warning: DB is nil! Make sure InitMongoDB was called")
	}
	return DB
}
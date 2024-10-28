package database

import (
	"context"
	"log"
	"time"

	"pymouse/pymouse/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB Contexts
var (
	tdContext = context.TODO()
	bgContext = context.Background()
)

// MongoDB Client
var MongoClient *mongo.Client
var Database *mongo.Database

// InitDB initializes the MongoDB connection
func InitDB() {
	log.Println("\033[0;33mConnecting with MongoDB...\033[0m")
	clientOptions := options.Client().ApplyURI(config.DatabaseURI)
	client, err := mongo.Connect(bgContext, clientOptions)
	if err != nil {
		log.Fatalf("\033[0;31m[MongoDB][Connect][Error]: %v\033[0m", err)
	}

	ctx, cancel := context.WithTimeout(bgContext, 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("\033[0;31m[MongoDB][Ping][Error]: %v\033[0m", err)
	}

	MongoClient = client
	Database = MongoClient.Database("PyMouse")
	log.Println("\033[0;32mConnected to MongoDB successfully!\033[0m")
}

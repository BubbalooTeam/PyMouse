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
	log.Println("\033[0;34mConnecting with MongoDB...\033")
	clientOptions := options.Client().ApplyURI(config.DatabaseURI)
	client, err := mongo.Connect(bgContext, clientOptions)
	if err != nil {
		log.Fatalf("[MongoDB][Connect][Error]: %v", err)
	}

	ctx, cancel := context.WithTimeout(bgContext, 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("[MongoDB][Ping][Error]: %v", err)
	}

	MongoClient = client
	Database = client.Database("PyMouse")
	log.Println("Connected to MongoDB successfully!")
}

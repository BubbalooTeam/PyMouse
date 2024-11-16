package database

import (
	"context"
	"time"

	"pymouse/pymouse/config"

	"github.com/sirupsen/logrus"
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
	logrus.Info("Connecting with MongoDB...")
	clientOptions := options.Client().ApplyURI(config.DatabaseURI)
	client, err := mongo.Connect(bgContext, clientOptions)
	if err != nil {
		logrus.Errorf("[MongoDB][Connect][Error]: %v", err)
	}

	ctx, cancel := context.WithTimeout(bgContext, 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		logrus.Errorf("[MongoDB][Ping][Error]: %v", err)
	}

	MongoClient = client
	Database = MongoClient.Database("PyMouse")
	logrus.Info("Connected to MongoDB successfully!")
}

func CloseDB() {
	if err := MongoClient.Disconnect(bgContext); err != nil {
		logrus.Errorf("Failed to close MongoDB: %v", err)
		return
	}
}

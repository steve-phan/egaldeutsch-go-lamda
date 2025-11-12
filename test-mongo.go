package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Get environment variables
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE")

	fmt.Printf("Testing MongoDB Connection...\n")
	fmt.Printf("URI: %s\n", uri)
	fmt.Printf("Database: %s\n", dbName)

	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	if dbName == "" {
		dbName = "egaldeutsch"
		fmt.Printf("Using default database: %s\n", dbName)
	}

	// Configure client options
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(10)
	clientOptions.SetMinPoolSize(0)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)
	clientOptions.SetServerSelectionTimeout(5 * time.Second)
	clientOptions.SetConnectTimeout(10 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Attempting to connect...")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Testing ping...")
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("✅ MongoDB connection successful!")

	// Test database access
	database := client.Database(dbName)
	collections, err := database.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}

	fmt.Printf("✅ Database accessible! Collections found: %v\n", collections)
}

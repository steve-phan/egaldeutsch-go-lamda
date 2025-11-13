package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
	once     sync.Once
	connErr  error
)

// Connect initializes MongoDB connection with connection pooling optimized for serverless
func Connect() error {
	once.Do(func() {
		connErr = initializeConnection()
	})
	return connErr
}

// GetClient returns the MongoDB client, initializing if necessary
func GetClient() (*mongo.Client, error) {
	if err := Connect(); err != nil {
		return nil, err
	}
	return Client, nil
}

// GetDatabase returns the MongoDB database, initializing if necessary
func GetDatabase() (*mongo.Database, error) {
	if err := Connect(); err != nil {
		return nil, err
	}
	return Database, nil
}

func initializeConnection() error {
	// Load .env file if it exists (for local development)
	// This will be ignored in production where env vars are set directly
	err := godotenv.Load()
	if err != nil {
		// Don't fail if .env file doesn't exist - this is normal in production
		fmt.Printf("No .env file found or failed to load (this is normal in production): %v\n", err)
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "egaldeutsch"
	}

	// Configure client options optimized for serverless
	clientOptions := options.Client().ApplyURI(uri)

	// Serverless-optimized connection pool settings
	clientOptions.SetMaxPoolSize(10)                         // Limit concurrent connections
	clientOptions.SetMinPoolSize(0)                          // No minimum connections (serverless friendly)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)       // Close idle connections quickly
	clientOptions.SetServerSelectionTimeout(5 * time.Second) // Quick timeout for server selection
	clientOptions.SetConnectTimeout(10 * time.Second)        // Connection timeout
	clientOptions.SetSocketTimeout(30 * time.Second)         // Socket timeout

	// Read preference for better performance
	clientOptions.SetReadPreference(readpref.Primary())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to ensure connectivity
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		client.Disconnect(ctx)
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	Client = client
	Database = client.Database(dbName)

	return nil
}

// Disconnect closes MongoDB connection
func Disconnect() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return Client.Disconnect(ctx)
	}
	return nil
}

// IsHealthy checks if the MongoDB connection is healthy
func IsHealthy(ctx context.Context) error {
	if Client == nil {
		return fmt.Errorf("client not initialized")
	}

	// Quick ping with timeout
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return Client.Ping(pingCtx, readpref.Primary())
}

// EnsureConnection ensures a healthy connection exists, reconnecting if necessary
func EnsureConnection() error {
	if Client == nil {
		return Connect()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := IsHealthy(ctx); err != nil {
		// Connection is unhealthy, reset and reconnect
		Client = nil
		Database = nil
		once = sync.Once{} // Reset the once flag
		return Connect()
	}

	return nil
}

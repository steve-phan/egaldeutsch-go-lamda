package services

import (
	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollections returns the database collections used by stories-management
func GetCollections() (*mongo.Collection, error) {
	return db.GetCollection(db.Collections.Stories)
}

package services

import (
	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollections returns the database collections used by stories-management
func GetCollections() (*mongo.Collection, error) {
	if err := db.EnsureConnection(); err != nil {
		return nil, err
	}

	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	storiesCollection := database.Collection("stories")

	return storiesCollection, nil
}

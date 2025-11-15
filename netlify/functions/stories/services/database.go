package services

import (
	"fmt"

	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollections ensures database connection and returns collections
func GetCollections() (*mongo.Collection, *mongo.Collection, error) {
	if err := db.EnsureConnection(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	database, err := db.GetDatabase()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get database: %w", err)
	}

	storiesCollection := database.Collection("stories")
	questionsCollection := database.Collection("questions")

	return storiesCollection, questionsCollection, nil
}

// GetStoriesCollection returns only the stories collection
func GetStoriesCollection() (*mongo.Collection, error) {
	storiesCollection, _, err := GetCollections()
	return storiesCollection, err
}

// GetQuestionsCollection returns only the questions collection
func GetQuestionsCollection() (*mongo.Collection, error) {
	_, questionsCollection, err := GetCollections()
	return questionsCollection, err
}

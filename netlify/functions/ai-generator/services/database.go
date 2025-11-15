package services

import (
	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollections returns the database collections used by ai-generator
func GetCollections() (*mongo.Collection, *mongo.Collection, *mongo.Collection, error) {
	if err := db.EnsureConnection(); err != nil {
		return nil, nil, nil, err
	}

	database, err := db.GetDatabase()
	if err != nil {
		return nil, nil, nil, err
	}

	storiesCollection := database.Collection("stories")
	questionsCollection := database.Collection("questions")
	quizzesCollection := database.Collection("quizzes")

	return storiesCollection, questionsCollection, quizzesCollection, nil
}

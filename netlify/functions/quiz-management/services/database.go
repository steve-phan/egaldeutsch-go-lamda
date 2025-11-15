package services

import (
	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollections returns the database collections used by quiz-management
func GetCollections() (*mongo.Collection, *mongo.Collection, *mongo.Collection, *mongo.Collection, *mongo.Collection, error) {
	if err := db.EnsureConnection(); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	database, err := db.GetDatabase()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	quizzesCollection := database.Collection("quizzes")
	storiesCollection := database.Collection("stories")
	questionsCollection := database.Collection("questions")
	submissionsCollection := database.Collection("submissions")
	usersCollection := database.Collection("users")

	return quizzesCollection, storiesCollection, questionsCollection, submissionsCollection, usersCollection, nil
}

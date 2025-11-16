package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionNames struct {
	Stories             string
	Questions           string
	Submissions         string
	Users               string
	Quizzes             string
	PasswordResetTokens string
}

var Collections = CollectionNames{
	Stories:             "stories",
	Questions:           "questions",
	Submissions:         "submissions",
	Users:               "users",
	Quizzes:             "quizzes",
	PasswordResetTokens: "password_reset_tokens",
}

// GetCollection returns a signle collection with connection management
func GetCollection(collectionName string) (*mongo.Collection, error) {
	if err := EnsureConnection(); err != nil {
		return nil, fmt.Errorf("failed  to ensure database connection: %w", err)
	}
	database, err := GetDatabase()

	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	return database.Collection(collectionName), nil
}

// GetCollections returns mutiple collections with connection management
// Pass collection names as variadic arguments

func GetCollections(collectionsNames ...string) ([]*mongo.Collection, error) {
	if err := EnsureConnection(); err != nil {
		return nil, fmt.Errorf("failed to ensure database connection: %w", err)
	}

	database, err := GetDatabase()

	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	collections := make([]*mongo.Collection, len(collectionsNames))
	for i, name := range collectionsNames {
		collections[i] = database.Collection(name)
	}

	return collections, nil

}

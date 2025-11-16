package services

import (
	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollections ensures database connection and returns all quiz-related collections
func GetCollections() ([]*mongo.Collection, error) {

	return db.GetCollections(db.Collections.Stories, db.Collections.Questions, db.Collections.Submissions)
}

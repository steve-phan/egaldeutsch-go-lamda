package services

import (
	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetPasswordResetTokenCollection() (*mongo.Collection, error) {
	return db.GetCollection(db.Collections.PasswordResetTokens)
}

func GetUserCollection() (*mongo.Collection, error) {
	return db.GetCollection(db.Collections.Users)
}

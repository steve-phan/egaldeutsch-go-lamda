package services

import (
	"context"
	"fmt"
	"time"

	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId"`
	Token     string             `bson:"token"`
	CreatedAt time.Time          `bson:"createdAt"`
	ExpiresAt time.Time          `bson:"expiresAt"`
	Used      bool               `bson:"used"`
}

// CreatePasswordResetToken creates a new password reset token
func CreatePasswordResetToken(userID primitive.ObjectID) (string, error) {
	collection := db.Database.Collection("password_reset_tokens")

	// Generate token (simple implementation)
	token := fmt.Sprintf("%d_%s", time.Now().Unix(), userID.Hex())

	resetToken := PasswordResetToken{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		Used:      false,
	}

	_, err := collection.InsertOne(context.TODO(), resetToken)
	if err != nil {
		return "", fmt.Errorf("failed to create password reset token: %w", err)
	}

	return token, nil
}

// ValidateResetToken validates and returns the user ID for a reset token
func ValidateResetToken(token string) (primitive.ObjectID, error) {
	collection := db.Database.Collection("password_reset_tokens")

	var resetToken PasswordResetToken
	err := collection.FindOne(context.TODO(), bson.M{
		"token":     token,
		"used":      false,
		"expiresAt": bson.M{"$gt": time.Now()},
	}).Decode(&resetToken)

	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid or expired token")
	}

	return resetToken.UserID, nil
}

// MarkTokenAsUsed marks a reset token as used
func MarkTokenAsUsed(token string) error {
	collection := db.Database.Collection("password_reset_tokens")

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"token": token},
		bson.M{"$set": bson.M{"used": true}},
	)

	return err
}

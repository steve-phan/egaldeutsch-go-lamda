package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/user-management/services"
	"egaldeutsch-serverless/netlify/functions/user-management/types"
	"egaldeutsch-serverless/pkg/notification"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// ForgotPassword handles password reset request
func ForgotPassword(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req types.ForgotPasswordRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid request format"), nil
	}

	// Find user by email
	collection, _ := services.GetUserCollection()
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)

	// Always return success to prevent email enumeration
	if err != nil {
		return response.SuccessJSONWithDefault(200, nil, "If the email exists, a password reset link has been sent")
	}

	// Check if user is active
	if user.Status != models.UserStatusActive {
		return response.SuccessJSONWithDefault(200, nil, "If the email exists, a password reset link has been sent")
	}

	// Generate password reset token
	token, err := services.CreatePasswordResetToken(user.ID)
	if err != nil {
		log.Printf("Failed to create password reset token for user %s: %v", user.ID.Hex(), err)
		return response.SuccessJSONWithDefault(200, nil, "If the email exists, a password reset link has been sent")
	}

	// Send password reset email directly (no goroutine needed in Lambda)
	log.Printf("Attempting to send password reset email to: %s", user.Email)
	if err := services.SendPasswordResetEmail(user.Email, user.Name, token); err != nil {
		log.Printf("Failed to send password reset email to %s: %v", user.Email, err)
		// Note: We don't fail the request if email fails to prevent enumeration
	} else {
		log.Printf("Successfully sent password reset email to: %s", user.Email)
	}
	return response.SuccessJSONWithDefault(200, nil, "If the email exists, a password reset link has been sent")
}

// ResetPassword handles password reset with token
func ResetPassword(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req types.ResetPasswordRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid request format"), nil
	}

	// Validate password length
	if len(req.NewPassword) < 8 {
		return response.SimpleErrorWithDefault(400, "Password must be at least 8 characters long"), nil
	}

	// Validate reset token
	userID, err := services.ValidateResetToken(req.Token)
	if err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid or expired reset token"), nil
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return response.SimpleErrorWithDefault(500, "Failed to hash password"), nil
	}

	// Update user password
	collection, _ := services.GetUserCollection()
	var user models.User
	err = collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"passwordHash": string(hashedPassword),
			"updatedAt":    time.Now(),
		}},
	).Decode(&user)

	if err != nil {
		return response.SimpleErrorWithDefault(500, "Failed to update password"), nil
	}

	// Mark token as used
	if err := services.MarkTokenAsUsed(req.Token); err != nil {
		log.Printf("Failed to mark reset token as used: %v", err)
	}

	// Invalidate all existing sessions for this user
	sessionCollection := db.Database.Collection("sessions")
	_, err = sessionCollection.DeleteMany(context.TODO(), bson.M{"userId": userID})
	if err != nil {
		log.Printf("Failed to delete user sessions after password reset: %v", err)
	}

	// Send password changed email (we can skip this for now since it's not in our requirements)
	// TODO: Add password changed email template if needed

	// Create notification for password change
	if err := notification.CreateNotification(
		userID,
		notification.NotificationTypePasswordChanged,
		"Password Changed",
		"Your password has been successfully changed.",
		"",
		nil,
	); err != nil {
		log.Printf("Failed to create password change notification: %v", err)
	}

	return response.SuccessJSONWithDefault(200, nil, "Password has been successfully reset")
}

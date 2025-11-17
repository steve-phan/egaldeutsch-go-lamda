package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/user-management/services"
	"egaldeutsch-serverless/netlify/functions/user-management/types"
	"egaldeutsch-serverless/pkg/auth"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser handles user authentication
func LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Login attempt - starting authentication process")

	var loginReq types.UserLoginRequest
	if err := json.Unmarshal([]byte(request.Body), &loginReq); err != nil {
		log.Printf("Login error - invalid request format: %v", err)
		return response.SimpleErrorWithDefault(400, "Invalid request format"), nil
	}

	log.Printf("Login attempt for username/email: %s", loginReq.Username)

	// Find user by username or email
	collection, _ := services.GetUserCollection()
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"username": loginReq.Username},
			{"email": loginReq.Username},
		},
	}).Decode(&user)

	if err != nil {
		log.Printf("Login error - user not found: %v", err)
		return response.SimpleErrorWithDefault(401, "Invalid credentials"), nil
	}

	log.Printf("User found: %s (ID: %s, Status: %s)", user.Username, user.ID.Hex(), user.Status)

	// Check if user is active
	if user.Status != models.UserStatusActive {
		log.Printf("Login error - user account not active: %s (Status: %s)", user.Username, user.Status)
		return response.SimpleErrorWithDefault(403, "Account is not active"), nil
	}

	// Verify password
	log.Printf("Verifying password for user: %s", user.Username)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		log.Printf("Login error - password verification failed for user %s: %v", user.Username, err)
		return response.SimpleErrorWithDefault(401, "Invalid credentials"), nil
	}

	log.Printf("Password verified successfully for user: %s", user.Username)

	// Generate JWT token
	log.Printf("Generating JWT token for user: %s", user.Username)
	token, err := auth.GenerateJWT(&user, 24) // 24 hour token
	if err != nil {
		log.Printf("Login error - failed to generate JWT for user %s: %v", user.Username, err)
		return response.SimpleErrorWithDefault(500, "Failed to generate authentication token"), nil
	}

	log.Printf("JWT token generated successfully for user: %s", user.Username)

	// Update user's last login
	now := time.Now()
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		}},
	)
	if err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.ID.Hex(), err)
	}
	user.LastLoginAt = &now

	// Return authentication response
	userResponse := types.UserResponse{
		ID:          user.ID.Hex(),
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Role:        string(user.Role),
		Status:      string(user.Status),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	authResponse := types.AuthTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      userResponse,
	}

	return response.JSONWithDefault(200, authResponse)
}

package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/user-management/services"
	"egaldeutsch-serverless/netlify/functions/user-management/types"
	"egaldeutsch-serverless/pkg/notification"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
func RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var regReq types.UserRegistrationRequest
	if err := json.Unmarshal([]byte(request.Body), &regReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid request format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Check if username or email already exists
	collection, _ := services.GetUserCollection()
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"username": regReq.Username},
			{"email": regReq.Email},
		},
	}).Decode(&existingUser)

	if err == nil {
		return response.SimpleError(409, "Username or email already exists"), nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.SimpleError(500, "Failed to hash password"), nil
	}

	// Create new user
	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Username:     regReq.Username,
		Email:        regReq.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    regReq.FirstName,
		LastName:     regReq.LastName,
		Name:         regReq.FirstName + " " + regReq.LastName,
		Role:         models.RoleCreator, // Default role, admin can upgrade
		Status:       models.UserStatusActive,
		IsActive:     true, // Backward compatibility
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user to database
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return response.SimpleError(500, "Failed to create user"), nil
	}

	// Send welcome email directly (no goroutine needed in Lambda)
	if err := services.SendWelcomeEmail(newUser.Email, newUser.Name); err != nil {
		log.Printf("Failed to send welcome email to %s: %v", newUser.Email, err)
		// Note: We don't fail the registration if email fails
	}

	// Create notification for user registration
	if err := notification.CreateNotification(
		newUser.ID,
		notification.NotificationTypeUserRegistered,
		"Welcome to EgalDeutsch!",
		"Your account has been successfully created. Start learning German today!",
		"/stories",
		nil,
	); err != nil {
		log.Printf("Failed to create registration notification: %v", err)
	}

	// Return user response (without password)
	userResponse := types.UserResponse{
		ID:          newUser.ID.Hex(),
		Username:    newUser.Username,
		Email:       newUser.Email,
		FirstName:   newUser.FirstName,
		LastName:    newUser.LastName,
		Role:        string(newUser.Role),
		Status:      string(newUser.Status),
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		LastLoginAt: newUser.LastLoginAt,
	}

	return response.JSON(201, userResponse)
}

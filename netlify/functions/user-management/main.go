package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserRegistrationRequest represents user registration data
type UserRegistrationRequest struct {
	Username      string `json:"username" validate:"required,min=3,max=50"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8,max=128"`
	FirstName     string `json:"firstName" validate:"required,min=2,max=50"`
	LastName      string `json:"lastName" validate:"required,min=2,max=50"`
	PreferredRole string `json:"preferredRole" validate:"required,oneof=creator reviewer"`
}

// UserLoginRequest represents login credentials
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents user data in responses (without sensitive info)
type UserResponse struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// AuthTokenResponse represents authentication response
type AuthTokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expiresAt"`
	User      UserResponse `json:"user"`
}

// Session represents a user session
type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Token     string             `bson:"token" json:"token"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	IPAddress string             `bson:"ipAddress" json:"ipAddress"`
	UserAgent string             `bson:"userAgent" json:"userAgent"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Database connection failed"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer db.Disconnect()

	// Handle CORS preflight requests
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Access-Control-Max-Age":       "86400",
			},
		}, nil
	}

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "POST":
		if strings.Contains(request.Path, "/register") {
			return registerUser(request)
		}
		if strings.Contains(request.Path, "/login") {
			return loginUser(request)
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error": "Endpoint not found"}`,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		}, nil
	case "GET":
		if strings.Contains(request.Path, "/profile") {
			return getUserProfile(request)
		}
		return listUsers(request)
	case "PUT":
		return updateUserProfile(request)
	case "DELETE":
		if strings.Contains(request.Path, "/logout") {
			return logoutUser(request)
		}
		return deleteUser(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       `{"error": "Method not allowed"}`,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		}, nil
	}
}

// registerUser handles user registration
func registerUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var regReq UserRegistrationRequest
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
	collection := db.Database.Collection("users")
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"username": regReq.Username},
			{"email": regReq.Email},
		},
	}).Decode(&existingUser)

	if err == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 409,
			Body:       `{"error": "Username or email already exists"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to hash password"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
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
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to create user"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Return user response (without password)
	userResponse := UserResponse{
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

	response, err := json.Marshal(userResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to encode response"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(response),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// loginUser handles user authentication
func loginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginReq UserLoginRequest
	if err := json.Unmarshal([]byte(request.Body), &loginReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid request format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Find user by username or email
	collection := db.Database.Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"username": loginReq.Username},
			{"email": loginReq.Username},
		},
	}).Decode(&user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Invalid credentials"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Check if user is active
	if user.Status != models.UserStatusActive {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       `{"error": "Account is not active"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Invalid credentials"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Generate session token
	token, err := generateSessionToken()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to generate token"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Create session
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hour session
	session := Session{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		IPAddress: getClientIP(request),
		UserAgent: request.Headers["User-Agent"],
	}

	// Save session to database
	sessionCollection := db.Database.Collection("sessions")
	_, err = sessionCollection.InsertOne(context.TODO(), session)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to create session"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

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
	userResponse := UserResponse{
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

	authResponse := AuthTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      userResponse,
	}

	response, err := json.Marshal(authResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to encode response"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// generateSessionToken creates a secure random token for session management
func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// getClientIP extracts client IP from request headers
func getClientIP(request events.APIGatewayProxyRequest) string {
	if ip := request.Headers["X-Forwarded-For"]; ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := request.Headers["X-Real-IP"]; ip != "" {
		return ip
	}
	return request.RequestContext.Identity.SourceIP
}

// validateSession validates a session token and returns the user
func validateSession(request events.APIGatewayProxyRequest) (*models.User, error) {
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"] // Check lowercase
	}

	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	token := parts[1]

	// Find session in database
	sessionCollection := db.Database.Collection("sessions")
	var session Session
	err := sessionCollection.FindOne(context.TODO(), bson.M{
		"token":     token,
		"expiresAt": bson.M{"$gt": time.Now()},
	}).Decode(&session)

	if err != nil {
		return nil, fmt.Errorf("invalid or expired token")
	}

	// Get user from session
	userCollection := db.Database.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": session.UserID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if user.Status != models.UserStatusActive {
		return nil, fmt.Errorf("user account is not active")
	}

	return &user, nil
}

// getUserProfile returns the current user's profile
func getUserProfile(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session
	user, err := validateSession(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Unauthorized"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Return user profile
	userResponse := UserResponse{
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

	response, err := json.Marshal(userResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to encode response"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func listUsers(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session and check admin role
	user, err := validateSession(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Unauthorized"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Check if user has admin role
	if user.Role != models.RoleAdmin && user.Role != models.RoleReviewer {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       `{"error": "Insufficient permissions"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Get query parameters for filtering
	queryParams := request.QueryStringParameters
	filter := bson.M{}

	if role := queryParams["role"]; role != "" {
		filter["role"] = role
	}
	if status := queryParams["status"]; status != "" {
		filter["status"] = status
	}

	// Find users
	collection := db.Database.Collection("users")
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to fetch users"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to decode users"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Convert to response format (without passwords)
	var userResponses []UserResponse
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:          u.ID.Hex(),
			Username:    u.Username,
			Email:       u.Email,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Role:        string(u.Role),
			Status:      string(u.Status),
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
			LastLoginAt: u.LastLoginAt,
		})
	}

	response, err := json.Marshal(userResponses)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to encode response"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// UserUpdateRequest represents user profile update data
type UserUpdateRequest struct {
	FirstName   *string `json:"firstName,omitempty"`
	LastName    *string `json:"lastName,omitempty"`
	Email       *string `json:"email,omitempty"`
	OldPassword *string `json:"oldPassword,omitempty"`
	NewPassword *string `json:"newPassword,omitempty"`
}

func updateUserProfile(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session
	user, err := validateSession(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Unauthorized"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	var updateReq UserUpdateRequest
	if err := json.Unmarshal([]byte(request.Body), &updateReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid request format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	collection := db.Database.Collection("users")
	updateFields := bson.M{"updatedAt": time.Now()}

	// Update basic profile fields
	if updateReq.FirstName != nil {
		updateFields["firstName"] = *updateReq.FirstName
	}
	if updateReq.LastName != nil {
		updateFields["lastName"] = *updateReq.LastName
	}

	// Update name field when first or last name changes
	if updateReq.FirstName != nil || updateReq.LastName != nil {
		firstName := user.FirstName
		lastName := user.LastName
		if updateReq.FirstName != nil {
			firstName = *updateReq.FirstName
		}
		if updateReq.LastName != nil {
			lastName = *updateReq.LastName
		}
		updateFields["name"] = firstName + " " + lastName
	}

	// Handle email update (check for uniqueness)
	if updateReq.Email != nil && *updateReq.Email != user.Email {
		var existingUser models.User
		err := collection.FindOne(context.TODO(), bson.M{"email": *updateReq.Email}).Decode(&existingUser)
		if err == nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 409,
				Body:       `{"error": "Email already in use"}`,
				Headers: map[string]string{
					"Content-Type":                "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
		}
		updateFields["email"] = *updateReq.Email
	}

	// Handle password update
	if updateReq.OldPassword != nil && updateReq.NewPassword != nil {
		// Verify old password
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*updateReq.OldPassword))
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `{"error": "Current password is incorrect"}`,
				Headers: map[string]string{
					"Content-Type":                "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateReq.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"error": "Failed to hash new password"}`,
				Headers: map[string]string{
					"Content-Type":                "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
		}
		updateFields["passwordHash"] = string(hashedPassword)
	}

	// Update user in database
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to update user"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Fetch updated user
	var updatedUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": user.ID}).Decode(&updatedUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to fetch updated user"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Return updated user profile
	userResponse := UserResponse{
		ID:          updatedUser.ID.Hex(),
		Username:    updatedUser.Username,
		Email:       updatedUser.Email,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		Role:        string(updatedUser.Role),
		Status:      string(updatedUser.Status),
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		LastLoginAt: updatedUser.LastLoginAt,
	}

	response, err := json.Marshal(userResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to encode response"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func logoutUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"] // Check lowercase
	}

	if authHeader == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "No authorization header"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid authorization header format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	token := parts[1]

	// Delete session from database
	sessionCollection := db.Database.Collection("sessions")
	_, err := sessionCollection.DeleteOne(context.TODO(), bson.M{"token": token})
	if err != nil {
		// Even if deletion fails, we should still return success
		// as the client side logout should still work
		log.Printf("Failed to delete session: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message": "Logged out successfully"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func deleteUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session and check admin role
	user, err := validateSession(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Unauthorized"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Only admins can delete users
	if user.Role != models.RoleAdmin {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       `{"error": "Insufficient permissions"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Get user ID from path parameters
	userID := request.PathParameters["id"]
	if userID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "User ID is required"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid user ID format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Prevent admin from deleting themselves
	if objectID == user.ID {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Cannot delete your own account"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Check if user exists
	collection := db.Database.Collection("users")
	var targetUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&targetUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error": "User not found"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Delete user's sessions first
	sessionCollection := db.Database.Collection("sessions")
	_, err = sessionCollection.DeleteMany(context.TODO(), bson.M{"userId": objectID})
	if err != nil {
		log.Printf("Failed to delete user sessions: %v", err)
	}

	// Instead of hard delete, we'll soft delete by setting status to suspended
	// This preserves content relationships and audit trail
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"status":    models.UserStatusSuspended,
			"isActive":  false,
			"updatedAt": time.Now(),
		}},
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to delete user"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message": "User deleted successfully"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}

package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
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
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
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

// Stub functions for remaining endpoints - to be implemented
func getUserProfile(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 501,
		Body:       `{"error": "Not implemented yet"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func listUsers(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 501,
		Body:       `{"error": "Not implemented yet"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func updateUserProfile(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 501,
		Body:       `{"error": "Not implemented yet"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func logoutUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 501,
		Body:       `{"error": "Not implemented yet"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func deleteUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 501,
		Body:       `{"error": "Not implemented yet"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}

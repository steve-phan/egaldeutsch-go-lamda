package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// GenerateSessionToken creates a secure random token for session management
func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GetClientIP extracts client IP from request headers
func GetClientIP(request events.APIGatewayProxyRequest) string {
	if ip := request.Headers["X-Forwarded-For"]; ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := request.Headers["X-Real-IP"]; ip != "" {
		return ip
	}
	return request.RequestContext.Identity.SourceIP
}

// ValidateSession validates a session token and returns the user
func ValidateSession(request events.APIGatewayProxyRequest) (*models.User, error) {
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

// ValidateSessionWithRole validates session and checks if user has required role
func ValidateSessionWithRole(request events.APIGatewayProxyRequest, allowedRoles ...models.UserRole) (*models.User, error) {
	user, err := ValidateSession(request)
	if err != nil {
		return nil, err
	}

	// Check if user has one of the allowed roles
	for _, role := range allowedRoles {
		if user.Role == role {
			return user, nil
		}
	}

	return nil, fmt.Errorf("insufficient permissions")
}

// CreateSession creates a new session for a user
func CreateSession(userID primitive.ObjectID, request events.APIGatewayProxyRequest) (*Session, string, error) {
	token, err := GenerateSessionToken()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour) // 24 hour session
	session := Session{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		IPAddress: GetClientIP(request),
		UserAgent: request.Headers["User-Agent"],
	}

	// Save session to database
	sessionCollection := db.Database.Collection("sessions")
	_, err = sessionCollection.InsertOne(context.TODO(), session)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	return &session, token, nil
}

// DeleteSession deletes a session by token
func DeleteSession(token string) error {
	sessionCollection := db.Database.Collection("sessions")
	_, err := sessionCollection.DeleteOne(context.TODO(), bson.M{"token": token})
	return err
}

// ExtractToken extracts the token from authorization header
func ExtractToken(request events.APIGatewayProxyRequest) (string, error) {
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"]
	}

	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

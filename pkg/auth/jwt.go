package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"egaldeutsch-serverless/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JWTClaims represents the JWT payload
type JWTClaims struct {
	UserID    string          `json:"userId"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	IssuedAt  int64           `json:"iat"`
	ExpiresAt int64           `json:"exp"`
}

// JWT secret key - should be loaded from environment
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback for development - in production, this MUST be set
		secret = "your-super-secret-jwt-key-change-in-production"
	}
	return secret
}

// GenerateJWT creates a JWT token for a user
func GenerateJWT(user *models.User, expirationHours int) (string, error) {
	if expirationHours == 0 {
		expirationHours = 24 // Default 24 hours
	}

	now := time.Now()
	claims := JWTClaims{
		UserID:    user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Duration(expirationHours) * time.Hour).Unix(),
	}

	// Create header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	// Encode header
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Encode payload
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create signature
	message := headerEncoded + "." + payloadEncoded
	signature := createHMAC(message, getJWTSecret())
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	// Combine all parts
	token := message + "." + signatureEncoded

	return token, nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(token string) (*JWTClaims, error) {
	// Split token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signatureEncoded := parts[2]

	// Verify signature
	message := headerEncoded + "." + payloadEncoded
	expectedSignature := createHMAC(message, getJWTSecret())
	expectedSignatureEncoded := base64.RawURLEncoding.EncodeToString(expectedSignature)

	if signatureEncoded != expectedSignatureEncoded {
		return nil, fmt.Errorf("invalid token signature")
	}

	// Decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Check expiration
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token has expired")
	}

	return &claims, nil
}

// ValidateJWTFromRequest extracts and validates JWT from request headers
func ValidateJWTFromRequest(authHeader string) (*JWTClaims, error) {
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	token := parts[1]
	return ValidateJWT(token)
}

// GetUserFromClaims creates a basic User object from JWT claims
func GetUserFromClaims(claims *JWTClaims) (*models.User, error) {
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in claims: %w", err)
	}

	return &models.User{
		ID:       userID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
		Status:   models.UserStatusActive, // Assumed active if token is valid
	}, nil
}

// createHMAC creates an HMAC signature
func createHMAC(message, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return h.Sum(nil)
}

// RefreshToken generates a new token from an existing valid token
func RefreshToken(oldToken string) (string, error) {
	claims, err := ValidateJWT(oldToken)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Create a user object from claims
	user, err := GetUserFromClaims(claims)
	if err != nil {
		return "", err
	}

	// Generate new token with fresh expiration
	return GenerateJWT(user, 24)
}

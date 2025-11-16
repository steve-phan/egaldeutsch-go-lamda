package middleware

import (
	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/pkg/auth"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// ValidateJWT validates JWT from request and returns user claims
func ValidateJWT(request events.APIGatewayProxyRequest) (*auth.JWTClaims, error) {
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"] // Check lowercase
	}

	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	claims, err := auth.ValidateJWTFromRequest(authHeader)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

// ValidateJWTWithRole validates JWT and checks if user has required role
func ValidateJWTWithRole(request events.APIGatewayProxyRequest, allowedRoles ...models.UserRole) (*auth.JWTClaims, error) {
	claims, err := ValidateJWT(request)
	if err != nil {
		return nil, err
	}

	// Check if user has one of the allowed roles
	for _, role := range allowedRoles {
		if claims.Role == role {
			return claims, nil
		}
	}

	return nil, fmt.Errorf("insufficient permissions")
}

// RequireAuth is a middleware that validates JWT and returns error response if invalid
func RequireAuth(request events.APIGatewayProxyRequest) (*auth.JWTClaims, *events.APIGatewayProxyResponse) {
	claims, err := ValidateJWT(request)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		errorResponse := events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: fmt.Sprintf(`{"success": false, "error": "Unauthorized: %s"}`, err.Error()),
		}
		return nil, &errorResponse
	}

	return claims, nil
}

// RequireRole is a middleware that validates JWT and checks role
func RequireRole(request events.APIGatewayProxyRequest, allowedRoles ...models.UserRole) (*auth.JWTClaims, *events.APIGatewayProxyResponse) {
	claims, err := ValidateJWTWithRole(request, allowedRoles...)
	if err != nil {
		log.Printf("Authorization failed: %v", err)
		errorResponse := events.APIGatewayProxyResponse{
			StatusCode: http.StatusForbidden,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: fmt.Sprintf(`{"success": false, "error": "Forbidden: %s"}`, err.Error()),
		}
		return nil, &errorResponse
	}

	return claims, nil
}

// Helper function to extract user from claims
func GetUserFromClaims(claims *auth.JWTClaims) (*models.User, error) {
	return auth.GetUserFromClaims(claims)
}

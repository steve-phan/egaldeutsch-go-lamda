package auth

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

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

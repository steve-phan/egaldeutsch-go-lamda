package handlers

import (
	"log"

	"egaldeutsch-serverless/pkg/auth"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
)

// LogoutUser handles user logout by invalidating the session
func LogoutUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract token using auth package
	token, err := auth.ExtractToken(request)
	if err != nil {
		return response.SimpleError(400, "No valid authorization token"), nil
	}

	// Delete session using auth package
	if err := auth.DeleteSession(token); err != nil {
		log.Printf("Failed to delete session: %v", err)
	}

	return response.SuccessJSON(200, nil, "Logged out successfully")
}

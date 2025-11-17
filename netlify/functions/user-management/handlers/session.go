package handlers

import (
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
)

// LogoutUser handles user logout
// With JWT, logout is handled client-side by removing the token
// This endpoint validates the token and returns success
func LogoutUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate JWT to ensure it's a valid authenticated request
	_, errResponse := middleware.RequireAuth(request)
	if errResponse != nil {
		return *errResponse, nil
	}

	// With JWT, no server-side action needed
	// Client will delete the token from localStorage
	return response.SuccessJSONWithDefault(200, nil, "Logged out successfully")
}

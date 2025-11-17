package services

import (
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
)

// ErrorResponse creates a standardized error response with CORS headers
func ErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return response.SimpleErrorWithDefault(statusCode, message)
}

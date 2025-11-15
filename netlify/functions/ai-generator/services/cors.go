package services

import (
	"github.com/aws/aws-lambda-go/events"
)

// GetCORSHeaders returns standard CORS headers
func GetCORSHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Methods":     "POST, OPTIONS",
		"Access-Control-Allow-Credentials": "true",
	}
}

// HandleCORSOptions handles OPTIONS requests for CORS preflight
func HandleCORSOptions() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    GetCORSHeaders(),
	}
}

// ErrorResponse creates a standardized error response with CORS headers
func ErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       `{"error": "` + message + `"}`,
		Headers:    GetCORSHeaders(),
	}
}

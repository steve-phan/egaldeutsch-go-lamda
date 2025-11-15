package services

import (
	"github.com/aws/aws-lambda-go/events"
)

// GetCORSHeaders returns standard CORS headers for quiz API responses
func GetCORSHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Methods":     "GET, POST, OPTIONS",
		"Access-Control-Allow-Credentials": "true",
	}
}

// HandleCORSOptions handles OPTIONS preflight requests
func HandleCORSOptions() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    GetCORSHeaders(),
	}
}

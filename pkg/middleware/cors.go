package middleware

import "github.com/aws/aws-lambda-go/events"

// GetCORSHeaders returns standard CORS headers for API responses
func GetCORSHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		"Access-Control-Allow-Credentials": "true",
		"Content-Type":                     "application/json",
	}
}

// HandleCORS handles OPTIONS preflight requests
func HandleCORS(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, bool) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    GetCORSHeaders(),
		}, true
	}
	return events.APIGatewayProxyResponse{}, false
}

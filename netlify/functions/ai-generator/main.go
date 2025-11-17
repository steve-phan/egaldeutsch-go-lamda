package main

import (
	"egaldeutsch-serverless/netlify/functions/ai-generator/handlers"
	"egaldeutsch-serverless/netlify/functions/ai-generator/services"
	"egaldeutsch-serverless/pkg/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

// init loads environment variables from .env file
func init() {
	// Try to load .env file from project root (3 levels up from function directory)
	// This works for local development
	if err := godotenv.Load("../../../.env"); err != nil {
		// If that fails, try loading from current directory
		if err := godotenv.Load(".env"); err != nil {
			// Silently continue - in production, env vars should be set directly
			// The "No .env file found" message is expected and normal
		}
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle OPTIONS request
	if request.HTTPMethod == "OPTIONS" {
		return middleware.HandleCORSOptions(middleware.PublicAPI), nil
	}

	// Only POST is allowed
	if request.HTTPMethod != "POST" {
		return services.ErrorResponse(405, "Method not allowed"), nil
	}

	// Route to main generation handler
	return handlers.GenerateContent(request)
}

func main() {
	lambda.Start(handler)
}

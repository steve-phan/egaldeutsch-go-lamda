package main

import (
	"log"
	"os"
	"strings"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/netlify/functions/user-management/handlers"
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

// init loads environment variables from .env file in development
func init() {
	log.Printf("Initializing user-management function")

	// Try multiple paths for .env file - Netlify dev changes working directory
	paths := []string{
		"../../../.env", // From function directory to project root
		"../../.env",    // Alternative path
		".env",          // Current directory
		"./.env",        // Explicit current directory
	}

	envLoaded := false
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Successfully loaded .env file from: %s", path)
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		log.Printf("No .env file loaded (this is normal in production)")
	}

	// Log some key environment variables (without exposing sensitive data)
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DATABASE")

	log.Printf("Environment check - MONGODB_URI present: %t", mongoURI != "")
	log.Printf("Environment check - MONGODB_DATABASE: %s", mongoDBName)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		return response.SimpleErrorWithDefault(500, "Database connection failed"), nil
	}
	defer db.Disconnect()

	// Handle CORS preflight requests
	if corsResponse, handled := middleware.HandleCORS(request); handled {
		return corsResponse, nil
	}

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "POST":
		if strings.Contains(request.Path, "/register") {
			return handlers.RegisterUser(request)
		}
		if strings.Contains(request.Path, "/login") {
			return handlers.LoginUser(request)
		}
		if strings.Contains(request.Path, "/forgot-password") {
			return handlers.ForgotPassword(request)
		}
		if strings.Contains(request.Path, "/reset-password") {
			return handlers.ResetPassword(request)
		}
		return response.SimpleErrorWithDefault(404, "Endpoint not found"), nil
	case "GET":
		if strings.Contains(request.Path, "/profile") {
			return handlers.GetUserProfile(request)
		}
		return handlers.ListUsers(request)
	case "PUT":
		return handlers.UpdateUserProfile(request)
	case "DELETE":
		return handlers.DeleteUser(request)
	default:
		return response.SimpleErrorWithDefault(405, "Method not allowed"), nil
	}
}

func main() {
	lambda.Start(handler)
}

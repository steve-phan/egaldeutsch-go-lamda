package main

import (
	"log"
	"os"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/netlify/functions/stories-management/handlers"
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

// init loads environment variables from .env file in development
func init() {
	log.Printf("Initializing stories-management function")

	// Try multiple paths for .env file
	paths := []string{
		"../../../.env",
		"../../.env",
		".env",
		"./.env",
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

	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DATABASE")

	log.Printf("Environment check - MONGODB_URI present: %t", mongoURI != "")
	log.Printf("Environment check - MONGODB_DATABASE: %s", mongoDBName)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		return response.DatabaseError("Failed to connect to database"), nil
	}
	defer db.Disconnect()

	// Handle CORS preflight
	if corsResponse, handled := middleware.HandleCORS(request); handled {
		return corsResponse, nil
	}

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "GET":
		// Check if requesting by slug or ID
		pathParam := request.PathParameters["id"]
		if pathParam != "" {
			// Try to determine if it's a slug or ID
			// MongoDB ObjectIDs are exactly 24 hex characters
			if len(pathParam) == 24 && isHexString(pathParam) {
				// Looks like an ID
				return handlers.GetStory(request)
			} else {
				// Treat as slug
				return handlers.GetStoryBySlug(request)
			}
		}
		return handlers.ListStories(request)
	case "POST":
		return handlers.CreateStory(request)
	case "PUT":
		return handlers.UpdateStory(request)
	case "PATCH":
		return handlers.UpdateStoryStatus(request)
	case "DELETE":
		return handlers.DeleteStory(request)
	default:
		return response.SimpleErrorWithDefault(405, "Method not allowed"), nil
	}
}

// isHexString checks if a string contains only hexadecimal characters
func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func main() {
	lambda.Start(handler)
}

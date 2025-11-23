package main

import (
	"context"
	"net/http"
	"strings"

	"egaldeutsch-serverless/netlify/functions/stories/handlers"
	"egaldeutsch-serverless/pkg/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Handle OPTIONS request
	if req.HTTPMethod == "OPTIONS" {
		return middleware.HandleCORSOptions(middleware.PublicAPI), nil
	}

	// Route requests
	method := req.HTTPMethod

	// Check for story ID in query parameters
	storyID := ""
	if id, exists := req.QueryStringParameters["id"]; exists && id != "" {
		storyID = id
	}

	// Check for slug in query parameters
	slug := ""
	if s, exists := req.QueryStringParameters["slug"]; exists && s != "" {
		slug = s
	}

	// Try different path parameter keys
	if s, exists := req.PathParameters["proxy"]; exists && s != "" {
		slug = s
	} else if s, exists := req.PathParameters["slug"]; exists && s != "" {
		slug = s
	} else {
		// Extract slug from path manually
		// Path format: /.netlify/functions/stories/my-slug
		path := req.Path
		parts := strings.Split(strings.Trim(path, "/"), "/")

		// Look for the slug after "stories"
		for i, part := range parts {
			if part == "stories" && i+1 < len(parts) {
				slug = parts[i+1]
				break
			}
		}
	}

	storyStatus := ""
	if status, exists := req.QueryStringParameters["status"]; exists && status != "" {
		storyStatus = status
	}

	switch {
	case method == "GET" && slug != "":
		// Request by slug: /stories/my-story-slug
		return handlers.GetStoryBySlug(ctx, slug)
	case method == "GET" && storyID != "":
		// Request by ID: /stories?id=123
		return handlers.GetStoryByID(ctx, storyID)
	case method == "GET":
		// List all stories: /stories or /stories?status=published
		return handlers.GetAllStories(ctx, storyStatus)
	case method == "POST":
		return handlers.CreateStory(ctx, req)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    middleware.GetPublicCORSHeaders(),
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}

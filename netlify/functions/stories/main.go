package main

import (
	"context"
	"net/http"

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
	storyStatus := ""
	if status, exists := req.QueryStringParameters["status"]; exists && status != "" {
		storyStatus = status
	}

	switch {
	case method == "GET" && storyID == "":
		return handlers.GetAllStories(ctx, storyStatus)
	case method == "GET" && storyID != "":
		return handlers.GetStoryByID(ctx, storyID)
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

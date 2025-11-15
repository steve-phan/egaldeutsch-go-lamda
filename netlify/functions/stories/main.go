package main

import (
	"context"
	"net/http"

	"egaldeutsch-serverless/netlify/functions/stories/handlers"
	"egaldeutsch-serverless/netlify/functions/stories/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Handle OPTIONS request
	if req.HTTPMethod == "OPTIONS" {
		return services.HandleCORSOptions(), nil
	}

	// Route requests
	method := req.HTTPMethod

	// Check for story ID in query parameters
	storyID := ""
	if id, exists := req.QueryStringParameters["id"]; exists && id != "" {
		storyID = id
	}

	switch {
	case method == "GET" && storyID == "":
		return handlers.GetAllStories(ctx)
	case method == "GET" && storyID != "":
		return handlers.GetStoryByID(ctx, storyID)
	case method == "POST":
		return handlers.CreateStory(ctx, req)
	default:
		headers := services.GetCORSHeaders()
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    headers,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}

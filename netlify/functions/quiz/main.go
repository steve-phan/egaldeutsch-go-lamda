package main

import (
	"context"
	"net/http"

	"egaldeutsch-serverless/netlify/functions/quiz/handlers"
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

	// Get story ID from query parameters
	storyID := ""
	if id, exists := req.QueryStringParameters["story_id"]; exists && id != "" {
		storyID = id
	}

	// Check for submit action
	action := ""
	if act, exists := req.QueryStringParameters["action"]; exists && act != "" {
		action = act
	}

	switch {
	case method == "GET" && storyID != "":
		return handlers.GetQuiz(ctx, storyID)
	case method == "POST" && storyID != "" && action == "submit":
		return handlers.SubmitQuiz(ctx, storyID, req)
	default:
		headers := middleware.GetPublicCORSHeaders()
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    headers,
			Body:       `{"success": false, "error": "Method not allowed or invalid route"}`,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}

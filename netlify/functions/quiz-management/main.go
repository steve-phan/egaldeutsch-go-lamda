package main

import (
	"strings"

	"egaldeutsch-serverless/netlify/functions/quiz-management/handlers"
	"egaldeutsch-serverless/netlify/functions/quiz-management/services"
	"egaldeutsch-serverless/pkg/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle OPTIONS request
	if request.HTTPMethod == "OPTIONS" {
		return middleware.HandleCORSOptions(middleware.AuthenticatedAPI), nil
	}

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "GET":
		if request.PathParameters["id"] != "" {
			// Check if it's a submissions endpoint
			if strings.Contains(request.Path, "/submissions") {
				return handlers.GetQuizSubmissions(request)
			}
			return handlers.GetQuiz(request)
		}
		return handlers.ListQuizzes(request)
	case "POST":
		if request.PathParameters["auto"] != "" {
			return handlers.AutoGenerateQuiz(request)
		}
		// Check if it's a submit endpoint
		if strings.Contains(request.Path, "/submit") {
			return handlers.SubmitQuiz(request)
		}
		return handlers.CreateQuiz(request)
	case "PUT":
		return handlers.UpdateQuiz(request)
	case "DELETE":
		return handlers.DeleteQuiz(request)
	default:
		return services.ErrorResponse(405, "Method not allowed"), nil
	}
}

func main() {
	lambda.Start(handler)
}

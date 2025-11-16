package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Set CORS headers
	headers := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Methods":     "GET, POST, OPTIONS",
		"Access-Control-Allow-Credentials": "true",
	}

	// Handle OPTIONS request
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
		}, nil
	}

	// Route requests
	path := req.Path
	method := req.HTTPMethod

	// Extract story ID from path
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	storyID := ""
	if len(pathParts) > 1 {
		storyID = pathParts[1]
	}

	switch {
	case method == "GET" && storyID != "":
		return getQuestionsByStoryID(ctx, storyID, headers)
	case method == "POST":
		return createQuestion(ctx, req, headers)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    headers,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func getQuestionsByStoryID(ctx context.Context, storyID string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	id, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       `{"error": "Invalid story ID"}`,
		}, nil
	}

	questionCollection, _ := db.GetCollection(db.Collections.Questions)

	cursor, err := questionCollection.Find(ctx, bson.M{"storyId": id})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "%v"}`, err),
		}, nil
	}
	defer cursor.Close(ctx)

	var questions []models.Question
	if err = cursor.All(ctx, &questions); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "%v"}`, err),
		}, nil
	}

	jsonData, err := json.Marshal(questions)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "%v"}`, err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonData),
	}, nil
}

func createQuestion(ctx context.Context, req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var question models.Question
	if err := json.Unmarshal([]byte(req.Body), &question); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "Invalid request body: %v"}`, err),
		}, nil
	}

	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()
	questionCollection, _ := db.GetCollection(db.Collections.Questions)

	_, err := questionCollection.InsertOne(ctx, question)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "%v"}`, err),
		}, nil
	}

	jsonData, err := json.Marshal(question)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "%v"}`, err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    headers,
		Body:       string(jsonData),
	}, nil
}

func main() {
	lambda.Start(handler)
}

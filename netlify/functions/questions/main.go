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
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle CORS preflight
	if corsResponse, handled := middleware.HandlePublicCORS(req); handled {
		return corsResponse, nil
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
		return getQuestionsByStoryID(ctx, storyID)
	case method == "POST":
		return createQuestion(ctx, req)
	default:
		return response.SimpleError(http.StatusMethodNotAllowed, "Method not allowed", middleware.PublicAPI), nil
	}
}

func getQuestionsByStoryID(ctx context.Context, storyID string) (events.APIGatewayProxyResponse, error) {
	id, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return response.SimpleError(http.StatusBadRequest, "Invalid story ID", middleware.PublicAPI), nil
	}

	questionCollection, _ := db.GetCollection(db.Collections.Questions)

	cursor, err := questionCollection.Find(ctx, bson.M{"storyId": id})
	if err != nil {
		return response.SimpleError(http.StatusInternalServerError, err.Error(), middleware.PublicAPI), nil
	}
	defer cursor.Close(ctx)

	var questions []models.Question
	if err = cursor.All(ctx, &questions); err != nil {
		return response.SimpleError(http.StatusInternalServerError, err.Error(), middleware.PublicAPI), nil
	}

	return response.JSON(http.StatusOK, questions, middleware.PublicAPI)
}

func createQuestion(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var question models.Question
	if err := json.Unmarshal([]byte(req.Body), &question); err != nil {
		return response.SimpleError(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err), middleware.PublicAPI), nil
	}

	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()
	questionCollection, _ := db.GetCollection(db.Collections.Questions)

	_, err := questionCollection.InsertOne(ctx, question)
	if err != nil {
		return response.SimpleError(http.StatusInternalServerError, err.Error(), middleware.PublicAPI), nil
	}

	return response.JSON(http.StatusCreated, question, middleware.PublicAPI)
}

func main() {
	lambda.Start(handler)
}

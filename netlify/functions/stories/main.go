package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var storiesCollection *mongo.Collection
var questionsCollection *mongo.Collection

func init() {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		fmt.Printf("Error connecting to MongoDB: %v\n", err)
		os.Exit(1)
	}

	storiesCollection = db.Database.Collection("stories")
	questionsCollection = db.Database.Collection("questions")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Set CORS headers
	headers := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
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

	// Extract story ID from path if present
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	storyID := ""
	if len(pathParts) > 1 {
		storyID = pathParts[1]
	}

	switch {
	case method == "GET" && storyID == "":
		return getAllStories(ctx, headers)
	case method == "GET" && storyID != "":
		return getStoryByID(ctx, storyID, headers)
	case method == "POST":
		return createStory(ctx, req, headers)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    headers,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func getAllStories(ctx context.Context, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	// Only fetch active stories
	filter := bson.M{"isActive": true}
	cursor, err := storiesCollection.Find(ctx, filter)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to fetch stories"}`,
		}, nil
	}
	defer cursor.Close(ctx)

	var stories []models.Story
	if err = cursor.All(ctx, &stories); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to decode stories"}`,
		}, nil
	}

	// If no stories found, return empty array
	if stories == nil {
		stories = []models.Story{}
	}

	response := map[string]interface{}{
		"success": true,
		"data":    stories,
		"message": fmt.Sprintf("Found %d stories", len(stories)),
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to serialize response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonData),
	}, nil
}

func getStoryByID(ctx context.Context, storyID string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	id, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       `{"success": false, "error": "Invalid story ID format"}`,
		}, nil
	}

	var story models.Story
	filter := bson.M{"_id": id, "isActive": true}
	err = storiesCollection.FindOne(ctx, filter).Decode(&story)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Headers:    headers,
				Body:       `{"success": false, "error": "Story not found or not active"}`,
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to fetch story"}`,
		}, nil
	}

	response := map[string]interface{}{
		"success": true,
		"data":    story,
		"message": "Story retrieved successfully",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to serialize response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonData),
	}, nil
}

func createStory(ctx context.Context, req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var story models.Story
	if err := json.Unmarshal([]byte(req.Body), &story); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "Invalid request body: %v"}`, err),
		}, nil
	}

	story.ID = primitive.NewObjectID()
	story.CreatedAt = time.Now()
	story.UpdatedAt = time.Now()

	_, err := storiesCollection.InsertOne(ctx, story)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "%v"}`, err),
		}, nil
	}

	jsonData, err := json.Marshal(story)
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

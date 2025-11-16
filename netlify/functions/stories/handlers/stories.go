package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/stories/services"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllStories retrieves all published stories
func GetAllStories(ctx context.Context, status string) (events.APIGatewayProxyResponse, error) {
	headers := services.GetCORSHeaders()

	storiesCollection, err := services.GetStoriesCollection()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Database connection failed"}`,
		}, nil
	}
	filterStatus := ""
	filter := bson.M{}

	if status != "" {
		filterStatus = status
		filter = bson.M{"status": filterStatus}
	}

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
		"data":    stories, // need to change the field to the stories
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

// GetStoryByID retrieves a specific story by its ID
func GetStoryByID(ctx context.Context, storyID string) (events.APIGatewayProxyResponse, error) {
	headers := services.GetCORSHeaders()

	storiesCollection, err := services.GetStoriesCollection()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Database connection failed"}`,
		}, nil
	}

	id, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       `{"success": false, "error": "Invalid story ID format"}`,
		}, nil
	}

	var story models.Story
	filter := bson.M{"_id": id, "status": "published"}
	err = storiesCollection.FindOne(ctx, filter).Decode(&story)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Headers:    headers,
				Body:       `{"success": false, "error": "Story not found or not published"}`,
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

// CreateStory creates a new story
func CreateStory(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := services.GetCORSHeaders()

	storiesCollection, err := services.GetStoriesCollection()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Database connection failed"}`,
		}, nil
	}

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

	_, err = storiesCollection.InsertOne(ctx, story)
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

package handlers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/stories-management/services"
	"egaldeutsch-serverless/netlify/functions/stories-management/types"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateStory creates a new story in draft status
func CreateStory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var storyReq types.StoryRequest
	if err := json.Unmarshal([]byte(request.Body), &storyReq); err != nil {
		return services.ErrorResponse(400, "Invalid request body"), nil
	}

	// Get database collection
	storiesCollection, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Calculate word count and reading time
	wordCount := len(strings.Fields(storyReq.Content))
	readingTime := wordCount / 200 // Assume 200 words per minute

	// Create story with metadata
	userID := primitive.NewObjectID() // TODO: Get from auth context
	now := time.Now()

	story := &models.Story{
		ID:          primitive.NewObjectID(),
		Title:       storyReq.Title,
		Content:     storyReq.Content,
		Level:       storyReq.Level,
		WordCount:   wordCount,
		ReadingTime: readingTime,
		Topics:      storyReq.Topics,
		Summary:     storyReq.Summary,
		Vocabulary:  storyReq.Vocabulary,
		ContentMetadata: models.ContentMetadata{
			Status:    models.StatusDraft,
			CreatedBy: userID,
			CreatedAt: now,
			UpdatedAt: now,
			Version:   1,
		},
		IsAIQuestionsGenerated: false,
	}

	// Insert story into database
	_, err = storiesCollection.InsertOne(context.Background(), story)
	if err != nil {
		return services.ErrorResponse(500, "Failed to create story"), nil
	}

	// Send email notification for new story
	notificationService := services.NewNotificationService()
	if err := notificationService.SendNewStoryNotification(story.ID.Hex(), story.Title, story.Level); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	response := convertToStoryResponse(*story)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers:    services.GetCORSHeaders(),
	}, nil
}

// GetStory retrieves a single story by ID
func GetStory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	storyID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid story ID"), nil
	}

	// Get database collection
	storiesCollection, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	var story models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&story)
	if err != nil {
		return services.ErrorResponse(404, "Story not found"), nil
	}

	response := convertToStoryResponse(story)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    services.GetCORSHeaders(),
	}, nil
}

// ListStories retrieves all stories with optional filtering and pagination
func ListStories(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get database collection
	storiesCollection, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Parse query parameters
	queryParams := request.QueryStringParameters

	// Build filter
	filter := bson.M{}

	if level := queryParams["level"]; level != "" {
		filter["level"] = level
	}

	if status := queryParams["status"]; status != "" {
		filter["status"] = status
	}

	if topic := queryParams["topic"]; topic != "" {
		filter["topics"] = bson.M{"$in": []string{topic}}
	}

	// Parse pagination parameters
	page := 1
	if p := queryParams["page"]; p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 10
	if l := queryParams["limit"]; l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	skip := (page - 1) * limit

	// Get total count
	total, err := storiesCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return services.ErrorResponse(500, "Failed to count stories"), nil
	}

	// Set up find options
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}}) // Most recent first
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Execute query
	cursor, err := storiesCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return services.ErrorResponse(500, "Failed to fetch stories"), nil
	}
	defer cursor.Close(context.Background())

	var stories []models.Story
	if err = cursor.All(context.Background(), &stories); err != nil {
		return services.ErrorResponse(500, "Failed to decode stories"), nil
	}

	// Convert to response format
	var storyResponses []types.StoryResponse
	for _, story := range stories {
		storyResponses = append(storyResponses, convertToStoryResponse(story))
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	response := types.ListStoriesResponse{
		Stories:    storyResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    services.GetCORSHeaders(),
	}, nil
}

// DeleteStory deletes a story
func DeleteStory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	storyID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid story ID"), nil
	}

	// Get database collection
	storiesCollection, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	result, err := storiesCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return services.ErrorResponse(500, "Failed to delete story"), nil
	}

	if result.DeletedCount == 0 {
		return services.ErrorResponse(404, "Story not found"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
		Headers:    services.GetCORSHeaders(),
	}, nil
}

// convertToStoryResponse converts a models.Story to types.StoryResponse
func convertToStoryResponse(story models.Story) types.StoryResponse {
	return types.StoryResponse{
		ID:                     story.ID.Hex(),
		Title:                  story.Title,
		Content:                story.Content,
		Level:                  story.Level,
		WordCount:              story.WordCount,
		ReadingTime:            story.ReadingTime,
		Topics:                 story.Topics,
		Summary:                story.Summary,
		Vocabulary:             story.Vocabulary,
		Status:                 story.Status,
		IsAIQuestionsGenerated: story.IsAIQuestionsGenerated,
		CreatedBy:              story.CreatedBy.Hex(),
		CreatedAt:              story.CreatedAt,
		UpdatedAt:              story.UpdatedAt,
		Version:                story.Version,
	}
}

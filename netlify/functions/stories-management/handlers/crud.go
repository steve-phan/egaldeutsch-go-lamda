package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/stories-management/services"
	"egaldeutsch-serverless/netlify/functions/stories-management/types"
	"egaldeutsch-serverless/pkg/middleware"

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

	// Generate unique slug from title
	baseSlug := models.GenerateSlug(storyReq.Title)
	slug := baseSlug
	counter := 1

	// Check for slug uniqueness and append number if needed
	for {
		var existingStory models.Story
		err = storiesCollection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&existingStory)
		if err != nil {
			// Slug doesn't exist, we can use it
			break
		}
		// Slug exists, try with a number suffix
		slug = baseSlug + "-" + strconv.Itoa(counter)
		counter++
	}

	// Create story with metadata
	userID := primitive.NewObjectID() // TODO: Get from auth context
	now := time.Now()

	story := &models.Story{
		ID:          primitive.NewObjectID(),
		Title:       storyReq.Title,
		Slug:        slug,
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
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
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
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// GetStoryBySlug retrieves a single story by its slug
func GetStoryBySlug(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract slug from path - handle both PathParameters and manual extraction
	var slug string
	if len(request.PathParameters) > 0 && request.PathParameters["slug"] != "" {
		slug = request.PathParameters["slug"]
	} else {
		// Extract from path manually
		parts := strings.Split(strings.Trim(request.Path, "/"), "/")
		if len(parts) < 2 {
			return services.ErrorResponse(400, "Story slug is required"), nil
		}
		slug = parts[len(parts)-1]
	}

	if slug == "" {
		return services.ErrorResponse(400, "Story slug is required"), nil
	}

	// Get database collection
	storiesCollection, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	var story models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&story)
	if err != nil {
		return services.ErrorResponse(404, "Story not found"), nil
	}

	response := convertToStoryResponse(story)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
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
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
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
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// convertToStoryResponse converts a models.Story to types.StoryResponse
func convertToStoryResponse(story models.Story) types.StoryResponse {
	return types.StoryResponse{
		ID:                     story.ID.Hex(),
		Title:                  story.Title,
		Slug:                   story.Slug,
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

// UpdateStory updates an existing story
func UpdateStory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract story ID from path
	var storyID string
	if len(request.PathParameters) > 0 && request.PathParameters["id"] != "" {
		storyID = request.PathParameters["id"]
	} else {
		// Extract from path manually
		parts := strings.Split(strings.Trim(request.Path, "/"), "/")
		if len(parts) < 2 {
			return services.ErrorResponse(400, "Story ID is required"), nil
		}
		storyID = parts[len(parts)-1]
	}

	objectID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid story ID"), nil
	}

	var storyReq types.StoryRequest
	if err := json.Unmarshal([]byte(request.Body), &storyReq); err != nil {
		return services.ErrorResponse(400, "Invalid request body"), nil
	}

	// Get database collection
	storiesCollection, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Fetch existing story first
	var existingStory models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&existingStory)
	if err != nil {
		return services.ErrorResponse(404, "Story not found"), nil
	}

	// Calculate word count and reading time
	wordCount := len(strings.Fields(storyReq.Content))
	readingTime := wordCount / 200

	// Check if title changed and regenerate slug if needed
	slug := existingStory.Slug
	if storyReq.Title != existingStory.Title {
		baseSlug := models.GenerateSlug(storyReq.Title)
		slug = baseSlug
		counter := 1

		// Check for slug uniqueness (excluding current story)
		for {
			var checkStory models.Story
			err = storiesCollection.FindOne(context.Background(), bson.M{
				"slug": slug,
				"_id":  bson.M{"$ne": objectID},
			}).Decode(&checkStory)
			if err != nil {
				// Slug doesn't exist or error, we can use it
				break
			}
			// Slug exists, try with a number suffix
			slug = baseSlug + "-" + strconv.Itoa(counter)
			counter++
		}
	}

	// Update story in database
	update := bson.M{
		"$set": bson.M{
			"title":       storyReq.Title,
			"slug":        slug,
			"content":     storyReq.Content,
			"level":       storyReq.Level,
			"wordCount":   wordCount,
			"readingTime": readingTime,
			"topics":      storyReq.Topics,
			"summary":     storyReq.Summary,
			"vocabulary":  storyReq.Vocabulary,
			"updatedAt":   time.Now(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	result := storiesCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedStory models.Story
	if err := result.Decode(&updatedStory); err != nil {
		return services.ErrorResponse(404, "Story not found"), nil
	}

	response := convertToStoryResponse(updatedStory)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// UpdateStoryStatus handles status transitions (submit for review, approve, reject, etc.)
func UpdateStoryStatus(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract story ID from path
	var storyID string

	// First try path parameters
	if id, exists := request.PathParameters["id"]; exists && id != "" {
		storyID = id
	} else {
		// Try to parse from path manually
		pathParts := strings.Split(strings.Trim(request.Path, "/"), "/")
		if len(pathParts) >= 4 {
			storyID = pathParts[3] // Index 3 should be the story ID
		}
	}

	if storyID == "" {
		return services.ErrorResponse(400, "Story ID is required"), nil
	}

	objectID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid story ID format"), nil
	}

	var statusReq types.StatusUpdateRequest
	if err := json.Unmarshal([]byte(request.Body), &statusReq); err != nil {
		return services.ErrorResponse(400, "Invalid request body"), nil
	}

	// Get database collection
	storiesCollection, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Get current story to validate status transition
	var currentStory models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&currentStory)
	if err != nil {
		return services.ErrorResponse(404, "Story not found"), nil
	}

	// Validate status transition
	if !currentStory.Status.CanTransitionTo(statusReq.Status) {
		return services.ErrorResponse(400, fmt.Sprintf("Cannot transition from %s to %s", currentStory.Status, statusReq.Status)), nil
	}

	// Prepare update
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":    statusReq.Status,
			"updatedAt": now,
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	// Set appropriate timestamp based on status
	userID := primitive.NewObjectID() // TODO: Get from auth context
	switch statusReq.Status {
	case models.StatusPreview:
		// Story submitted for review
	case models.StatusReady:
		update["$set"].(bson.M)["reviewedBy"] = userID
		update["$set"].(bson.M)["reviewedAt"] = now
		update["$set"].(bson.M)["approvedAt"] = now
	case models.StatusPublished:
		update["$set"].(bson.M)["activatedAt"] = now
	case models.StatusDraft:
		// Story sent back to draft
		update["$set"].(bson.M)["reviewedBy"] = userID
		update["$set"].(bson.M)["reviewedAt"] = now
	}

	// Add comment if provided
	if statusReq.Comment != "" {
		comment := models.ReviewComment{
			ID:         primitive.NewObjectID(),
			ReviewerID: userID,
			Comment:    statusReq.Comment,
			CreatedAt:  now,
			Type:       getCommentType(statusReq.Status),
		}
		update["$push"] = bson.M{"comments": comment}
	}

	result := storiesCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedStory models.Story
	if err := result.Decode(&updatedStory); err != nil {
		return services.ErrorResponse(500, "Failed to update story status"), nil
	}

	response := convertToStoryResponse(updatedStory)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// getCommentType returns the comment type based on status
func getCommentType(status models.ContentStatus) string {
	switch status {
	case models.StatusReady:
		return "approval_note"
	case models.StatusDraft:
		return "revision_note"
	case models.StatusPublished:
		return "publication_note"
	default:
		return "feedback"
	}
}

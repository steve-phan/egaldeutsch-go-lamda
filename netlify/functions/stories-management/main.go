package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/notification"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StoryRequest represents the request structure for creating/updating stories
type StoryRequest struct {
	Title      string                  `json:"title"`
	Content    string                  `json:"content"`
	Level      string                  `json:"level"`
	Topics     []string                `json:"topics"`
	Summary    string                  `json:"summary"`
	Vocabulary []models.VocabularyWord `json:"vocabulary"`
}

// StoryResponse represents the response structure for stories
type StoryResponse struct {
	ID                     string                  `json:"id"`
	Title                  string                  `json:"title"`
	Content                string                  `json:"content"`
	Level                  string                  `json:"level"`
	WordCount              int                     `json:"wordCount"`
	ReadingTime            int                     `json:"readingTime"`
	Topics                 []string                `json:"topics"`
	Summary                string                  `json:"summary"`
	Vocabulary             []models.VocabularyWord `json:"vocabulary"`
	Status                 models.ContentStatus    `json:"status"`
	IsAIQuestionsGenerated bool                    `json:"isAIQuestionsGenerated"`
	CreatedBy              string                  `json:"createdBy"`
	CreatedAt              time.Time               `json:"createdAt"`
	UpdatedAt              time.Time               `json:"updatedAt"`
	Version                int                     `json:"version"`
}

// StatusUpdateRequest represents status transition requests
type StatusUpdateRequest struct {
	Status  models.ContentStatus `json:"status"`
	Comment string               `json:"comment,omitempty"`
}

// ListStoriesResponse represents the response for listing stories
type ListStoriesResponse struct {
	Stories    []StoryResponse `json:"stories"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"totalPages"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		return response.DatabaseError("Failed to connect to database"), nil
	}
	defer db.Disconnect()

	// Handle CORS preflight
	if corsResponse, handled := middleware.HandleCORS(request); handled {
		return corsResponse, nil
	}

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "GET":
		if request.PathParameters["id"] != "" {
			return getStory(request)
		}
		return listStories(request)
	case "POST":
		return createStory(request)
	case "PUT":
		return updateStory(request)
	case "PATCH":
		return updateStoryStatus(request)
	case "DELETE":
		return deleteStory(request)
	default:
		return response.SimpleErrorWithDefault(405, "Method not allowed"), nil
	}
}

// createStory creates a new story in draft status
func createStory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate JWT and get user claims
	claims, errResponse := middleware.RequireAuth(request)
	if errResponse != nil {
		return *errResponse, nil
	}

	// Get user from claims
	user, err := middleware.GetUserFromClaims(claims)
	if err != nil {
		return response.SimpleErrorWithDefault(500, "Failed to get user information"), nil
	}

	// Check if user has permission to create stories (Creator, Reviewer, or Admin)
	if user.Role != models.RoleCreator && user.Role != models.RoleReviewer && user.Role != models.RoleAdmin {
		return response.SimpleErrorWithDefault(403, "Insufficient permissions to create stories"), nil
	}

	var storyReq StoryRequest
	if err := json.Unmarshal([]byte(request.Body), &storyReq); err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid request body"), nil
	}

	// Validate the story
	story := &models.Story{
		Title:      storyReq.Title,
		Content:    storyReq.Content,
		Level:      storyReq.Level,
		Topics:     storyReq.Topics,
		Summary:    storyReq.Summary,
		Vocabulary: storyReq.Vocabulary,
	}

	if err := story.Validate(); err != nil {
		return response.ValidationError(fmt.Sprintf("Validation error: %v", err)), nil
	}

	// Calculate word count and reading time
	story.WordCount = len([]rune(story.Content)) / 5 // Rough word count estimation
	story.ReadingTime = story.WordCount / 200        // 200 words per minute

	// Use authenticated user ID
	now := time.Now()

	// Set up story with draft status
	story.ID = primitive.NewObjectID()
	story.ContentMetadata = models.ContentMetadata{
		Status:    models.StatusDraft,
		CreatedBy: user.ID,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}

	// Insert story into database
	collection := db.Database.Collection("stories")
	_, insertErr := collection.InsertOne(context.Background(), story)
	if insertErr != nil {
		return response.DatabaseError("Failed to create story"), nil
	}

	// Convert to response format
	response := convertToStoryResponse(*story)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// getStory retrieves a single story by ID
func getStory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	storyID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid story ID"), nil
	}

	collection, _ := db.GetCollection(db.Collections.Stories)
	var story models.Story

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&story)
	if err != nil {
		return response.NotFoundError("Story"), nil
	}

	response := convertToStoryResponse(story)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// listStories retrieves stories with filtering and pagination
func listStories(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse query parameters
	queryParams := request.QueryStringParameters

	page := 1
	limit := 10

	if p := queryParams["page"]; p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := queryParams["limit"]; l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// Build filter
	filter := bson.M{}

	if status := queryParams["status"]; status != "" {
		filter["status"] = status
	}

	if level := queryParams["level"]; level != "" {
		filter["level"] = level
	}

	if topic := queryParams["topic"]; topic != "" {
		filter["topics"] = bson.M{"$in": []string{topic}}
	}

	if search := queryParams["search"]; search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"summary": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	collection, _ := db.GetCollection(db.Collections.Stories)

	// Get total count
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return response.DatabaseError("Failed to count stories"), nil
	}

	// Get stories with pagination
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"createdAt": -1})

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return response.DatabaseError("Failed to retrieve stories"), nil
	}
	defer cursor.Close(context.Background())

	var stories []models.Story
	if err = cursor.All(context.Background(), &stories); err != nil {
		return response.DatabaseError("Failed to decode stories"), nil
	}

	// Convert to response format
	storyResponses := make([]StoryResponse, len(stories))
	for i, story := range stories {
		storyResponses[i] = convertToStoryResponse(story)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := ListStoriesResponse{
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

// updateStory updates an existing story
func updateStory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	storyID := request.PathParameters["id"]
	log.Printf("requestid is , %s", storyID)
	objectID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid story ID"), nil
	}

	var storyReq StoryRequest
	if err := json.Unmarshal([]byte(request.Body), &storyReq); err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid request body"), nil
	}

	// Validate the updated story
	story := &models.Story{
		Title:      storyReq.Title,
		Content:    storyReq.Content,
		Level:      storyReq.Level,
		Topics:     storyReq.Topics,
		Summary:    storyReq.Summary,
		Vocabulary: storyReq.Vocabulary,
	}

	if err := story.Validate(); err != nil {
		return response.ValidationError(fmt.Sprintf("Validation error: %v", err)), nil
	}

	// Calculate updated word count and reading time
	story.WordCount = len([]rune(story.Content)) / 5
	story.ReadingTime = story.WordCount / 200

	// Update the story in database
	collection, _ := db.GetCollection(db.Collections.Stories)
	update := bson.M{
		"$set": bson.M{
			"title":       story.Title,
			"content":     story.Content,
			"level":       story.Level,
			"wordCount":   story.WordCount,
			"readingTime": story.ReadingTime,
			"topics":      story.Topics,
			"summary":     story.Summary,
			"vocabulary":  story.Vocabulary,
			"updatedAt":   time.Now(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	result := collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedStory models.Story
	if err := result.Decode(&updatedStory); err != nil {
		return response.NotFoundError("Story"), nil
	}

	response := convertToStoryResponse(updatedStory)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// updateStoryStatus handles status transitions (submit for review, approve, reject, etc.)
func updateStoryStatus(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Debug: log the request path and parameters
	fmt.Printf("Request path: %s\n", request.Path)
	fmt.Printf("Path parameters: %+v\n", request.PathParameters)
	fmt.Printf("Query parameters: %+v\n", request.QueryStringParameters)

	// Try to get story ID from different sources
	var storyID string

	// First try path parameters
	if id, exists := request.PathParameters["id"]; exists && id != "" {
		storyID = id
	} else if id, exists := request.QueryStringParameters["id"]; exists && id != "" {
		// Try query parameters
		storyID = id
	} else {
		// Try to parse from path manually
		// Path might be like: /.netlify/functions/stories-management/69159f73a354489a4d203f11/status
		pathParts := strings.Split(strings.Trim(request.Path, "/"), "/")
		if len(pathParts) >= 4 {
			storyID = pathParts[3] // Index 3 should be the story ID
		}
	}

	fmt.Printf("Extracted story ID: %s\n", storyID)

	if storyID == "" {
		return response.SimpleErrorWithDefault(400, "Story ID is required"), nil
	}

	objectID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid story ID format"), nil
	}

	var statusReq StatusUpdateRequest
	if err := json.Unmarshal([]byte(request.Body), &statusReq); err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid request body"), nil
	}

	// Get current story to validate status transition
	collection, _ := db.GetCollection(db.Collections.Stories)
	var currentStory models.Story
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&currentStory)
	if err != nil {
		return response.NotFoundError("Story"), nil
	}

	// Validate status transition
	if !currentStory.Status.CanTransitionTo(statusReq.Status) {
		return response.SimpleErrorWithDefault(400, fmt.Sprintf("Cannot transition from %s to %s", currentStory.Status, statusReq.Status)), nil
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
	userID := primitive.NewObjectID() // Mock user ID
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

	result := collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedStory models.Story
	if err := result.Decode(&updatedStory); err != nil {
		return response.DatabaseError("Failed to update story status"), nil
	}

	// Send notifications based on status change
	switch statusReq.Status {
	case models.StatusPublished:
		// Notify all users when story is published
		go func() {
			if err := notification.NotifyStoryPublished(updatedStory.ID, updatedStory.Title); err != nil {
				log.Printf("Failed to send publication notifications for story %s: %v", updatedStory.ID.Hex(), err)
			}
		}()

		// Send email notifications to subscribers
		go func() {
			if err := sendNewStoryEmailNotification(updatedStory.ID.Hex(), updatedStory.Title, updatedStory.Level); err != nil {
				log.Printf("Failed to send email notifications for story %s: %v", updatedStory.ID.Hex(), err)
			}
		}()
	case models.StatusPreview:
		// Notify admins/reviewers when story is submitted for review
		go func() {
			if err := notification.NotifyStorySubmitted(currentStory.CreatedBy, updatedStory.ID, updatedStory.Title); err != nil {
				log.Printf("Failed to send submission notifications for story %s: %v", updatedStory.ID.Hex(), err)
			}
		}()
	case models.StatusReady, models.StatusDraft:
		// Notify creator about status change
		message := statusReq.Comment
		if message == "" {
			if statusReq.Status == models.StatusReady {
				message = "Your story has been approved and is ready to publish!"
			} else {
				message = "Your story has been sent back for revision."
			}
		}
		go func() {
			if err := notification.NotifyStoryStatusChange(currentStory.CreatedBy, updatedStory.ID, updatedStory.Title, string(statusReq.Status), message); err != nil {
				log.Printf("Failed to send status change notification for story %s: %v", updatedStory.ID.Hex(), err)
			}
		}()
	}

	response := convertToStoryResponse(updatedStory)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// deleteStory soft deletes a story (archives it)
func deleteStory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	storyID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid story ID"), nil
	}

	collection, _ := db.GetCollection(db.Collections.Stories)
	// In simplified workflow, delete means permanent removal
	result, err := collection.DeleteOne(
		context.Background(),
		bson.M{"_id": objectID},
	)

	if err != nil {
		return response.DatabaseError("Failed to delete story"), nil
	}

	if result.DeletedCount == 0 {
		return response.NotFoundError("Story"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// Helper functions

func convertToStoryResponse(story models.Story) StoryResponse {
	return StoryResponse{
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

func errorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	body := fmt.Sprintf(`{"error": "%s"}`, message)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// sendNewStoryEmailNotification sends email notifications for new published stories
func sendNewStoryEmailNotification(storyID, storyTitle, storyLevel string) error {
	emailServiceURL := getEmailServiceURL()

	payload := map[string]interface{}{
		"storyId":    storyID,
		"storyTitle": storyTitle,
		"storyLevel": storyLevel,
		"maxUsers":   10, // Limit for trial account
	}

	return callEmailService(emailServiceURL+"/send-new-story-notification", payload)
}

// callEmailService makes an HTTP call to the email service
func callEmailService(url string, payload map[string]interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to call email service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("email service returned status %d", resp.StatusCode)
	}

	return nil
}

// getEmailServiceURL returns the email service URL
func getEmailServiceURL() string {
	if baseURL := os.Getenv("NETLIFY_FUNCTIONS_URL"); baseURL != "" {
		return baseURL + "/email-service"
	}
	// Default for local development
	return "http://localhost:8888/.netlify/functions/email-service"
}

func main() {
	lambda.Start(handler)
}

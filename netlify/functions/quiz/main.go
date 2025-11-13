package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// getCollections ensures database connection and returns collections
func getCollections() (*mongo.Collection, *mongo.Collection, *mongo.Collection, error) {
	if err := db.EnsureConnection(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	database, err := db.GetDatabase()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get database: %w", err)
	}

	storiesCollection := database.Collection("stories")
	questionsCollection := database.Collection("questions")
	submissionsCollection := database.Collection("submissions")

	return storiesCollection, questionsCollection, submissionsCollection, nil
}

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
		return getQuiz(ctx, storyID, headers)
	case method == "POST" && storyID != "" && action == "submit":
		return submitQuiz(ctx, storyID, req, headers)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    headers,
			Body:       `{"success": false, "error": "Method not allowed or invalid route"}`,
		}, nil
	}
}

func getQuiz(ctx context.Context, storyID string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	storiesCollection, questionsCollection, _, err := getCollections()
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

	// Get story
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

	// Get questions for this story, ordered by order field
	filter = bson.M{"storyId": id}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "order", Value: 1}}) // Sort by order field ascending
	questionsCursor, err := questionsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to fetch questions"}`,
		}, nil
	}
	defer questionsCursor.Close(ctx)

	var questions []models.Question
	if err = questionsCursor.All(ctx, &questions); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to decode questions"}`,
		}, nil
	}

	// Sort questions by order if not empty
	if len(questions) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    headers,
			Body:       `{"success": false, "error": "No questions available for this story. Please generate AI questions first to create a quiz."}`,
		}, nil
	}

	// Remove correct answer from response (for security)
	for i := range questions {
		questions[i].CorrectAnswer = -1 // Hide correct answer
	}

	quiz := models.Quiz{
		StoryID:   story.ID,
		Story:     &story,
		Questions: questions,
	}

	response := map[string]interface{}{
		"success": true,
		"data":    quiz,
		"message": fmt.Sprintf("Quiz loaded with %d questions", len(questions)),
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

func submitQuiz(ctx context.Context, storyID string, req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	_, questionsCollection, submissionsCollection, err := getCollections()
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
	var submissionRequest struct {
		Answers []int `json:"answers"`
	}

	if err := json.Unmarshal([]byte(req.Body), &submissionRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       `{"success": false, "error": "Invalid request body"}`,
		}, nil
	}

	// Get questions with correct answers (ordered by order field)
	filter := bson.M{"storyId": id}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "order", Value: 1}}) // Sort by order field ascending
	var questionsCursor *mongo.Cursor
	questionsCursor, err = questionsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to fetch questions"}`,
		}, nil
	}
	defer questionsCursor.Close(ctx)

	var questions []models.Question
	if err = questionsCursor.All(ctx, &questions); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "%v"}`, err),
		}, nil
	}

	// Validate answers length
	if len(submissionRequest.Answers) != len(questions) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       `{"success": false, "error": "Answer count does not match question count"}`,
		}, nil
	}

	// Calculate score and track correct answers
	score := 0
	correctAnswers := make([]bool, len(questions))
	totalPoints := 0
	earnedPoints := 0

	for i, answer := range submissionRequest.Answers {
		if i < len(questions) {
			totalPoints += questions[i].Points
			if answer == questions[i].CorrectAnswer {
				score++
				correctAnswers[i] = true
				earnedPoints += questions[i].Points
			}
		}
	}

	percentage := float64(score) / float64(len(questions)) * 100
	passed := percentage >= 70.0 // 70% passing grade

	// Save submission
	submission := models.QuizSubmission{
		ID:             primitive.NewObjectID(),
		StoryID:        id,
		Answers:        submissionRequest.Answers,
		Score:          score,
		TotalQuestions: len(questions),
		SubmittedAt:    time.Now(),
	}

	_, err = submissionsCollection.InsertOne(ctx, submission)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to save submission"}`,
		}, nil
	}

	// Prepare detailed response
	result := map[string]interface{}{
		"id":             submission.ID.Hex(),
		"storyId":        storyID,
		"answers":        submissionRequest.Answers,
		"score":          score,
		"totalQuestions": len(questions),
		"percentage":     percentage,
		"passed":         passed,
		"correctAnswers": correctAnswers,
		"earnedPoints":   earnedPoints,
		"totalPoints":    totalPoints,
		"submittedAt":    submission.SubmittedAt,
	}

	response := map[string]interface{}{
		"success": true,
		"data":    result,
		"message": fmt.Sprintf("Quiz submitted successfully. Score: %d/%d (%.1f%%)", score, len(questions), percentage),
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

func main() {
	lambda.Start(handler)
}

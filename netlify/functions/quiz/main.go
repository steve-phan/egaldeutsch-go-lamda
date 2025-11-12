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
var submissionsCollection *mongo.Collection

func init() {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		fmt.Printf("Error connecting to MongoDB: %v\n", err)
		os.Exit(1)
	}

	storiesCollection = db.Database.Collection("stories")
	questionsCollection = db.Database.Collection("questions")
	submissionsCollection = db.Database.Collection("quiz_submissions")
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
	path := req.Path
	method := req.HTTPMethod

	// Extract story ID from path
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	storyID := ""
	action := ""
	if len(pathParts) > 1 {
		storyID = pathParts[1]
	}
	if len(pathParts) > 2 {
		action = pathParts[2]
	}

	switch {
	case method == "GET" && storyID != "" && action == "":
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

	// Get questions for this story, ordered by order field
	cursor, err := questionsCollection.Find(ctx, bson.M{"storyId": id}, )
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to fetch questions"}`,
		}, nil
	}
	defer cursor.Close(ctx)

	var questions []models.Question
	if err = cursor.All(ctx, &questions); err != nil {
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
			Body:       `{"success": false, "error": "No questions found for this story"}`,
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

func submitQuiz(ctx context.Context, req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var submissionRequest struct {
		StoryID primitive.ObjectID `json:"storyId"`
		Answers []int              `json:"answers"`
	}

	if err := json.Unmarshal([]byte(req.Body), &submissionRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "Invalid request body: %v"}`, err),
		}, nil
	}

	// Get questions with correct answers
	cursor, err := questionsCollection.Find(ctx, bson.M{"storyId": submissionRequest.StoryID})
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

	// Calculate score
	score := 0
	for i, answer := range submissionRequest.Answers {
		if i < len(questions) && answer == questions[i].CorrectAnswer {
			score++
		}
	}

	// Save submission
	submission := models.QuizSubmission{
		ID:              primitive.NewObjectID(),
		StoryID:         submissionRequest.StoryID,
		Answers:         submissionRequest.Answers,
		Score:           score,
		TotalQuestions:  len(questions),
		SubmittedAt:     time.Now(),
	}

	_, err = submissionsCollection.InsertOne(ctx, submission)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "%v"}`, err),
		}, nil
	}

	// Prepare response with results
	type QuizResult struct {
		Score           int     `json:"score"`
		TotalQuestions  int     `json:"totalQuestions"`
		Percentage      float64 `json:"percentage"`
		SubmittedAt     time.Time `json:"submittedAt"`
	}

	result := QuizResult{
		Score:          score,
		TotalQuestions: len(questions),
		Percentage:     float64(score) / float64(len(questions)) * 100,
		SubmittedAt:    submission.SubmittedAt,
	}

	jsonData, err := json.Marshal(result)
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

func main() {
	lambda.Start(handler)
}

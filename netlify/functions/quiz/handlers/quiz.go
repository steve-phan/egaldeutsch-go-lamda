package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/quiz/services"
	"egaldeutsch-serverless/netlify/functions/quiz/types"
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetQuiz retrieves a quiz for a specific story
func GetQuiz(ctx context.Context, storyID string) (events.APIGatewayProxyResponse, error) {
	headers := middleware.GetCORSHeaders()

	collections, err := services.GetCollections()
	if err != nil {
		return response.DatabaseError("Failed to connect to database"), nil
	}

	id, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return response.BadRequestError("Invalid story ID format"), nil
	}

	// Get story
	var story models.Story
	filter := bson.M{"_id": id, "status": "published"}
	err = collections[0].FindOne(ctx, filter).Decode(&story)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFoundError("Story"), nil
		}
		return response.DatabaseError("Failed to fetch story"), nil
	}

	// Get questions for this story, ordered by order field
	filter = bson.M{"storyId": id}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "order", Value: 1}}) // Sort by order field ascending
	questionsCursor, err := collections[1].Find(ctx, filter, findOptions)
	if err != nil {
		return response.DatabaseError("Failed to fetch questions"), nil
	}
	defer questionsCursor.Close(ctx)

	var questions []models.Question
	if err = questionsCursor.All(ctx, &questions); err != nil {
		return response.DatabaseError("Failed to decode questions"), nil
	}

	// Check if questions exist
	if len(questions) == 0 {
		return response.NotFoundError("Story"), nil
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

	resp := map[string]interface{}{
		"success": true,
		"data":    quiz,
		"message": fmt.Sprintf("Quiz loaded with %d questions", len(questions)),
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		return response.DatabaseError("Failed to serialize response"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonData),
	}, nil
}

// SubmitQuiz handles quiz submission and scoring
func SubmitQuiz(ctx context.Context, storyID string, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := middleware.GetCORSHeaders()

	collections, err := services.GetCollections()

	if err != nil {
		return response.DatabaseError("Failed to connect to database"), nil
	}

	id, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return response.BadRequestError("Invalid story ID format"), nil
	}

	var submissionRequest types.QuizSubmissionRequest
	if err := json.Unmarshal([]byte(req.Body), &submissionRequest); err != nil {
		return response.BadRequestError("Invalid request body"), nil
	}

	// Get questions with correct answers (ordered by order field)
	filter := bson.M{"storyId": id}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "order", Value: 1}}) // Sort by order field ascending
	var questionsCursor *mongo.Cursor
	questionsCursor, err = collections[1].Find(ctx, filter, findOptions)
	if err != nil {
		return response.DatabaseError("Failed to fetch questions"), nil
	}
	defer questionsCursor.Close(ctx)

	var questions []models.Question
	if err = questionsCursor.All(ctx, &questions); err != nil {
		return response.DatabaseError(fmt.Sprintf(`{"error": "%v"}`, err)), nil
	}

	// Validate answers length
	if len(submissionRequest.Answers) != len(questions) {
		return response.BadRequestError("Answer count does not match question count"), nil
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

	_, err = collections[2].InsertOne(ctx, submission)
	if err != nil {
		return response.DatabaseError("Failed to save submission"), nil
	}

	// Prepare detailed response using the types
	result := types.QuizResult{
		ID:             submission.ID.Hex(),
		StoryID:        storyID,
		Answers:        submissionRequest.Answers,
		Score:          score,
		TotalQuestions: len(questions),
		Percentage:     percentage,
		Passed:         passed,
		CorrectAnswers: correctAnswers,
		EarnedPoints:   earnedPoints,
		TotalPoints:    totalPoints,
		SubmittedAt:    submission.SubmittedAt,
	}

	resp := map[string]interface{}{
		"success": true,
		"data":    result,
		"message": fmt.Sprintf("Quiz submitted successfully. Score: %d/%d (%.1f%%)", score, len(questions), percentage),
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		return response.DatabaseError("Failed to serialize response"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonData),
	}, nil
}

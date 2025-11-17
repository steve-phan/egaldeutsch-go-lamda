package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/pkg/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// QuestionRequest represents the request structure for creating/updating questions
type QuestionRequest struct {
	StoryID       string   `json:"storyId"`
	Question      string   `json:"question"`
	QuestionType  string   `json:"questionType"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correctAnswer"`
	Explanation   string   `json:"explanation"`
	Points        int      `json:"points"`
	Order         int      `json:"order"` // Accept as number from frontend
	Difficulty    string   `json:"difficulty"`
}

// QuestionResponse represents the response structure for questions
type QuestionResponse struct {
	ID            string               `json:"id"`
	StoryID       string               `json:"storyId"`
	Question      string               `json:"question"`
	QuestionType  string               `json:"questionType"`
	Options       []string             `json:"options"`
	CorrectAnswer int                  `json:"correctAnswer"`
	Explanation   string               `json:"explanation"`
	Points        int                  `json:"points"`
	Order         int                  `json:"order"`
	Difficulty    string               `json:"difficulty"`
	Status        models.ContentStatus `json:"status"`
	CreatedBy     string               `json:"createdBy"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
	Version       int                  `json:"version"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "GET":
		if request.PathParameters["id"] != "" {
			return getQuestion(request)
		}
		return listQuestions(request)
	case "POST":
		return createQuestion(request)
	case "PUT":
		return updateQuestion(request)
	case "DELETE":
		return deleteQuestion(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       `{"error": "Method not allowed"}`,
			Headers:    middleware.GetAuthenticatedCORSHeaders(),
		}, nil
	}
}

// createQuestion creates a new question in draft status
func createQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var questionReq QuestionRequest
	if err := json.Unmarshal([]byte(request.Body), &questionReq); err != nil {
		return errorResponse(400, "Invalid request body")
	}

	// Validate story exists
	storyID, err := primitive.ObjectIDFromHex(questionReq.StoryID)
	if err != nil {
		return errorResponse(400, "Invalid story ID")
	}

	storiesCollection, _ := db.GetCollection(db.Collections.Stories)
	var story models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": storyID}).Decode(&story)
	if err != nil {
		return errorResponse(404, "Story not found")
	}

	// Create and validate the question
	question := &models.Question{
		StoryID:       storyID,
		Question:      questionReq.Question,
		QuestionType:  questionReq.QuestionType,
		Options:       questionReq.Options,
		CorrectAnswer: questionReq.CorrectAnswer,
		Explanation:   questionReq.Explanation,
		Points:        questionReq.Points,
		Order:         questionReq.Order,
		Difficulty:    questionReq.Difficulty,
	}

	if err := question.Validate(); err != nil {
		return errorResponse(400, fmt.Sprintf("Validation error: %v", err))
	}

	// Get user ID from context
	userID := primitive.NewObjectID()
	now := time.Now()

	// Set up question with draft status
	question.ID = primitive.NewObjectID()
	question.ContentMetadata = models.ContentMetadata{
		Status:    models.StatusDraft,
		CreatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}

	// Insert question into database
	questionCollection, _ := db.GetCollection(db.Collections.Questions)
	_, err = questionCollection.InsertOne(context.Background(), question)
	if err != nil {
		return errorResponse(500, "Failed to create question")
	}

	response := convertToQuestionResponse(*question)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

func getQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	questionID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return errorResponse(400, "Invalid question ID")
	}

	questionCollection, _ := db.GetCollection(db.Collections.Questions)
	var question models.Question

	err = questionCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&question)
	if err != nil {
		return errorResponse(404, "Question not found")
	}

	response := convertToQuestionResponse(question)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

func listQuestions(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	filter := bson.M{}

	if storyID := queryParams["storyId"]; storyID != "" {
		objectID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return errorResponse(400, "Invalid story ID")
		}
		filter["storyId"] = objectID
	}

	questionCollection, _ := db.GetCollection(db.Collections.Questions)

	total, err := questionCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return errorResponse(500, "Failed to count questions")
	}

	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"order": 1})

	cursor, err := questionCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return errorResponse(500, "Failed to retrieve questions")
	}
	defer cursor.Close(context.Background())

	var questions []models.Question
	if err = cursor.All(context.Background(), &questions); err != nil {
		return errorResponse(500, "Failed to decode questions")
	}

	questionResponses := make([]QuestionResponse, len(questions))
	for i, question := range questions {
		questionResponses[i] = convertToQuestionResponse(question)
	}

	response := map[string]interface{}{
		"questions": questionResponses,
		"total":     total,
		"page":      page,
		"limit":     limit,
	}

	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

func updateQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	questionID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return errorResponse(400, "Invalid question ID")
	}

	var questionReq QuestionRequest
	if err := json.Unmarshal([]byte(request.Body), &questionReq); err != nil {
		return errorResponse(400, "Invalid request body")
	}

	// Validate the updated question
	question := &models.Question{
		Question:      questionReq.Question,
		QuestionType:  questionReq.QuestionType,
		Options:       questionReq.Options,
		CorrectAnswer: questionReq.CorrectAnswer,
		Explanation:   questionReq.Explanation,
		Points:        questionReq.Points,
		Order:         questionReq.Order,
		Difficulty:    questionReq.Difficulty,
	}

	if err := question.Validate(); err != nil {
		return errorResponse(400, fmt.Sprintf("Validation error: %v", err))
	}

	// Update the question in database
	questionCollection, _ := db.GetCollection(db.Collections.Questions)
	update := bson.M{
		"$set": bson.M{
			"question":      question.Question,
			"questionType":  question.QuestionType,
			"options":       question.Options,
			"correctAnswer": question.CorrectAnswer,
			"explanation":   question.Explanation,
			"points":        question.Points,
			"order":         question.Order,
			"difficulty":    question.Difficulty,
			"updatedAt":     time.Now(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	result := questionCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedQuestion models.Question
	if err := result.Decode(&updatedQuestion); err != nil {
		return errorResponse(404, "Question not found")
	}

	response := convertToQuestionResponse(updatedQuestion)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

func deleteQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	questionID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return errorResponse(400, "Invalid question ID")
	}

	questionCollection, _ := db.GetCollection(db.Collections.Questions)
	// In simplified workflow, delete means permanent removal
	result, err := questionCollection.DeleteOne(
		context.Background(),
		bson.M{"_id": objectID},
	)

	if err != nil {
		return errorResponse(500, "Failed to delete question")
	}

	if result.DeletedCount == 0 {
		return errorResponse(404, "Question not found")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

func convertToQuestionResponse(question models.Question) QuestionResponse {
	return QuestionResponse{
		ID:            question.ID.Hex(),
		StoryID:       question.StoryID.Hex(),
		Question:      question.Question,
		QuestionType:  question.QuestionType,
		Options:       question.Options,
		CorrectAnswer: question.CorrectAnswer,
		Explanation:   question.Explanation,
		Points:        question.Points,
		Order:         question.Order,
		Difficulty:    question.Difficulty,
		Status:        question.Status,
		CreatedBy:     question.CreatedBy.Hex(),
		CreatedAt:     question.CreatedAt,
		UpdatedAt:     question.UpdatedAt,
		Version:       question.Version,
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

func main() {
	lambda.Start(handler)
}

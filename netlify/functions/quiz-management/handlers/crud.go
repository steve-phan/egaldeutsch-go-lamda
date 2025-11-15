package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/quiz-management/services"
	"egaldeutsch-serverless/netlify/functions/quiz-management/types"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateQuiz creates a new quiz in draft status
func CreateQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var quizReq types.QuizRequest
	if err := json.Unmarshal([]byte(request.Body), &quizReq); err != nil {
		return services.ErrorResponse(400, "Invalid request body"), nil
	}

	// Get collections
	quizzesCollection, storiesCollection, _, _, _, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Validate story exists
	storyID, err := primitive.ObjectIDFromHex(quizReq.StoryID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid story ID"), nil
	}

	var story models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": storyID}).Decode(&story)
	if err != nil {
		return services.ErrorResponse(404, "Story not found"), nil
	}

	// Convert question IDs to ObjectIDs
	var questionObjectIDs []primitive.ObjectID
	for _, qidStr := range quizReq.QuestionIDs {
		qid, err := primitive.ObjectIDFromHex(qidStr)
		if err != nil {
			return services.ErrorResponse(400, fmt.Sprintf("Invalid question ID: %s", qidStr)), nil
		}
		questionObjectIDs = append(questionObjectIDs, qid)
	}

	// Create and validate the quiz
	quiz := &models.Quiz{
		StoryID:        storyID,
		Title:          quizReq.Title,
		Description:    quizReq.Description,
		QuestionIDs:    questionObjectIDs,
		TotalQuestions: quizReq.TotalQuestions,
		PassingScore:   quizReq.PassingScore,
		QuizType:       quizReq.QuizType,
	}

	if err := quiz.Validate(); err != nil {
		return services.ErrorResponse(400, fmt.Sprintf("Validation error: %v", err)), nil
	}

	// Get user ID from context
	userID := primitive.NewObjectID()
	now := time.Now()

	// Set up quiz with draft status
	quiz.ID = primitive.NewObjectID()
	quiz.ContentMetadata = models.ContentMetadata{
		Status:    models.StatusDraft,
		CreatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}

	// Insert quiz into database
	_, err = quizzesCollection.InsertOne(context.Background(), quiz)
	if err != nil {
		return services.ErrorResponse(500, "Failed to create quiz"), nil
	}

	response := convertToQuizResponse(*quiz)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers:    services.GetCORSHeaders(),
	}, nil
}

// GetQuiz retrieves a single quiz by ID
func GetQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	quizID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid quiz ID"), nil
	}

	// Get collections
	quizzesCollection, _, _, _, _, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	var quiz models.Quiz
	err = quizzesCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&quiz)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return services.ErrorResponse(404, "Quiz not found"), nil
		}
		return services.ErrorResponse(500, "Failed to fetch quiz"), nil
	}

	response := convertToQuizResponse(quiz)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    services.GetCORSHeaders(),
	}, nil
}

// ListQuizzes retrieves all quizzes with optional filtering
func ListQuizzes(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get collections
	quizzesCollection, _, _, _, _, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Build filter based on query parameters
	filter := bson.M{}

	// Filter by story ID if provided
	if storyID := request.QueryStringParameters["storyId"]; storyID != "" {
		objectID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return services.ErrorResponse(400, "Invalid story ID"), nil
		}
		filter["storyId"] = objectID
	}

	// Filter by status if provided
	if status := request.QueryStringParameters["status"]; status != "" {
		filter["status"] = status
	}

	// Set up sorting and pagination
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}}) // Most recent first

	cursor, err := quizzesCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return services.ErrorResponse(500, "Failed to fetch quizzes"), nil
	}
	defer cursor.Close(context.Background())

	var quizzes []models.Quiz
	if err = cursor.All(context.Background(), &quizzes); err != nil {
		return services.ErrorResponse(500, "Failed to decode quizzes"), nil
	}

	// Convert to response format
	var quizResponses []types.QuizResponse
	for _, quiz := range quizzes {
		quizResponses = append(quizResponses, convertToQuizResponse(quiz))
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"quizzes": quizResponses,
		"total":   len(quizResponses),
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    services.GetCORSHeaders(),
	}, nil
}

// UpdateQuiz updates an existing quiz
func UpdateQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Implementation would go here - placeholder for now
	return services.ErrorResponse(501, "Update quiz not implemented yet"), nil
}

// DeleteQuiz deletes a quiz
func DeleteQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Implementation would go here - placeholder for now
	return services.ErrorResponse(501, "Delete quiz not implemented yet"), nil
}

// convertToQuizResponse converts a models.Quiz to types.QuizResponse
func convertToQuizResponse(quiz models.Quiz) types.QuizResponse {
	questionIDs := make([]string, len(quiz.QuestionIDs))
	for i, id := range quiz.QuestionIDs {
		questionIDs[i] = id.Hex()
	}

	return types.QuizResponse{
		ID:             quiz.ID.Hex(),
		StoryID:        quiz.StoryID.Hex(),
		Title:          quiz.Title,
		Description:    quiz.Description,
		QuestionIDs:    questionIDs,
		TotalQuestions: quiz.TotalQuestions,
		PassingScore:   quiz.PassingScore,
		QuizType:       quiz.QuizType,
		Status:         quiz.Status,
		CreatedBy:      quiz.CreatedBy.Hex(),
		CreatedAt:      quiz.CreatedAt,
		UpdatedAt:      quiz.UpdatedAt,
		Version:        quiz.Version,
	}
}

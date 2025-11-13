package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuizRequest represents the request structure for creating/updating quizzes
type QuizRequest struct {
	StoryID        string   `json:"storyId"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	QuestionIDs    []string `json:"questionIds,omitempty"`
	TotalQuestions int      `json:"totalQuestions"`
	PassingScore   int      `json:"passingScore"`
	QuizType       string   `json:"quizType"`
}

// QuizResponse represents the response structure for quizzes
type QuizResponse struct {
	ID             string               `json:"id"`
	StoryID        string               `json:"storyId"`
	Title          string               `json:"title"`
	Description    string               `json:"description"`
	QuestionIDs    []string             `json:"questionIds"`
	TotalQuestions int                  `json:"totalQuestions"`
	PassingScore   int                  `json:"passingScore"`
	QuizType       string               `json:"quizType"`
	Status         models.ContentStatus `json:"status"`
	CreatedBy      string               `json:"createdBy"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
	Version        int                  `json:"version"`
}

// AutoGenerateRequest represents request for auto-generating quiz
type AutoGenerateRequest struct {
	StoryID       string   `json:"storyId"`
	QuestionCount int      `json:"questionCount"`
	DifficultyMix string   `json:"difficultyMix"` // "easy", "medium", "hard", "mixed"
	QuestionTypes []string `json:"questionTypes"` // ["comprehension", "vocabulary", "grammar"]
}

// QuizSubmissionRequest represents a quiz submission
type QuizSubmissionRequest struct {
	Answers   map[string]int `json:"answers" validate:"required"`
	StartTime time.Time      `json:"startTime" validate:"required"`
	EndTime   time.Time      `json:"endTime" validate:"required"`
}

// QuizSubmission represents a completed quiz submission
type QuizSubmission struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	QuizID         primitive.ObjectID `json:"quizId" bson:"quizId"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	Score          int                `json:"score" bson:"score"`
	TotalQuestions int                `json:"totalQuestions" bson:"totalQuestions"`
	CorrectAnswers int                `json:"correctAnswers" bson:"correctAnswers"`
	Percentage     float64            `json:"percentage" bson:"percentage"`
	Passed         bool               `json:"passed" bson:"passed"`
	TimeSpent      int                `json:"timeSpent" bson:"timeSpent"` // in seconds
	SubmittedAt    time.Time          `json:"submittedAt" bson:"submittedAt"`
	Results        []QuestionResult   `json:"results" bson:"results"`
}

// QuestionResult represents the result for a single question
type QuestionResult struct {
	QuestionID     primitive.ObjectID `json:"questionId" bson:"questionId"`
	Question       string             `json:"question" bson:"question"`
	SelectedAnswer int                `json:"selectedAnswer" bson:"selectedAnswer"`
	CorrectAnswer  int                `json:"correctAnswer" bson:"correctAnswer"`
	IsCorrect      bool               `json:"isCorrect" bson:"isCorrect"`
	Points         int                `json:"points" bson:"points"`
	Explanation    string             `json:"explanation" bson:"explanation"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Database connection failed"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer db.Disconnect()

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "GET":
		if request.PathParameters["id"] != "" {
			// Check if it's a submissions endpoint
			if strings.Contains(request.Path, "/submissions") {
				return getQuizSubmissions(request)
			}
			return getQuiz(request)
		}
		return listQuizzes(request)
	case "POST":
		if request.PathParameters["auto"] != "" {
			return autoGenerateQuiz(request)
		}
		// Check if it's a submit endpoint
		if strings.Contains(request.Path, "/submit") {
			return submitQuiz(request)
		}
		return createQuiz(request)
	case "PUT":
		return updateQuiz(request)
	case "DELETE":
		return deleteQuiz(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       `{"error": "Method not allowed"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
}

// createQuiz creates a new quiz in draft status
func createQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var quizReq QuizRequest
	if err := json.Unmarshal([]byte(request.Body), &quizReq); err != nil {
		return errorResponse(400, "Invalid request body")
	}

	// Validate story exists
	storyID, err := primitive.ObjectIDFromHex(quizReq.StoryID)
	if err != nil {
		return errorResponse(400, "Invalid story ID")
	}

	storiesCollection := db.Database.Collection("stories")
	var story models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": storyID}).Decode(&story)
	if err != nil {
		return errorResponse(404, "Story not found")
	}

	// Convert question IDs to ObjectIDs
	var questionObjectIDs []primitive.ObjectID
	for _, qidStr := range quizReq.QuestionIDs {
		qid, err := primitive.ObjectIDFromHex(qidStr)
		if err != nil {
			return errorResponse(400, fmt.Sprintf("Invalid question ID: %s", qidStr))
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
		return errorResponse(400, fmt.Sprintf("Validation error: %v", err))
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
	collection := db.Database.Collection("quizzes")
	_, err = collection.InsertOne(context.Background(), quiz)
	if err != nil {
		return errorResponse(500, "Failed to create quiz")
	}

	response := convertToQuizResponse(*quiz)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func getQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	quizID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return errorResponse(400, "Invalid quiz ID")
	}

	collection := db.Database.Collection("quizzes")
	var quiz models.Quiz

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&quiz)
	if err != nil {
		return errorResponse(404, "Quiz not found")
	}

	response := convertToQuizResponse(quiz)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func listQuizzes(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := request.QueryStringParameters

	// Build filter
	filter := bson.M{}

	if storyID := queryParams["storyId"]; storyID != "" {
		objectID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return errorResponse(400, "Invalid story ID")
		}
		filter["storyId"] = objectID
	}

	if status := queryParams["status"]; status != "" {
		filter["status"] = status
	}

	if quizType := queryParams["type"]; quizType != "" {
		filter["quizType"] = quizType
	}

	collection := db.Database.Collection("quizzes")

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return errorResponse(500, "Failed to retrieve quizzes")
	}
	defer cursor.Close(context.Background())

	var quizzes []models.Quiz
	if err = cursor.All(context.Background(), &quizzes); err != nil {
		return errorResponse(500, "Failed to decode quizzes")
	}

	quizResponses := make([]QuizResponse, len(quizzes))
	for i, quiz := range quizzes {
		quizResponses[i] = convertToQuizResponse(quiz)
	}

	response := map[string]interface{}{
		"quizzes": quizResponses,
		"total":   len(quizzes),
	}

	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func updateQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return errorResponse(501, "Update quiz not implemented yet")
}

func deleteQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	quizID := request.PathParameters["id"]
	objectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return errorResponse(400, "Invalid quiz ID")
	}

	collection := db.Database.Collection("quizzes")
	// In simplified workflow, delete means permanent removal
	result, err := collection.DeleteOne(
		context.Background(),
		bson.M{"_id": objectID},
	)

	if err != nil {
		return errorResponse(500, "Failed to delete quiz")
	}

	if result.DeletedCount == 0 {
		return errorResponse(404, "Quiz not found")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// autoGenerateQuiz automatically generates a quiz from story questions
func autoGenerateQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var autoReq AutoGenerateRequest
	if err := json.Unmarshal([]byte(request.Body), &autoReq); err != nil {
		return errorResponse(400, "Invalid request body")
	}

	// Validate story exists
	storyID, err := primitive.ObjectIDFromHex(autoReq.StoryID)
	if err != nil {
		return errorResponse(400, "Invalid story ID")
	}

	storiesCollection := db.Database.Collection("stories")
	var story models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": storyID}).Decode(&story)
	if err != nil {
		return errorResponse(404, "Story not found")
	}

	// Build filter for questions
	filter := bson.M{
		"storyId": storyID,
		"status":  models.StatusPublished,
	}

	if len(autoReq.QuestionTypes) > 0 {
		filter["questionType"] = bson.M{"$in": autoReq.QuestionTypes}
	}

	if autoReq.DifficultyMix != "mixed" && autoReq.DifficultyMix != "" {
		filter["difficulty"] = autoReq.DifficultyMix
	}

	// Get available questions
	questionsCollection := db.Database.Collection("questions")
	cursor, err := questionsCollection.Find(context.Background(), filter)
	if err != nil {
		return errorResponse(500, "Failed to retrieve questions")
	}
	defer cursor.Close(context.Background())

	var questions []models.Question
	if err = cursor.All(context.Background(), &questions); err != nil {
		return errorResponse(500, "Failed to decode questions")
	}

	if len(questions) < autoReq.QuestionCount {
		return errorResponse(400, fmt.Sprintf("Not enough questions available. Found %d, requested %d", len(questions), autoReq.QuestionCount))
	}

	// Select questions based on criteria
	selectedQuestions := selectQuestionsForQuiz(questions, autoReq.QuestionCount, autoReq.DifficultyMix)

	var questionIDs []primitive.ObjectID
	for _, q := range selectedQuestions {
		questionIDs = append(questionIDs, q.ID)
	}

	// Create the quiz
	userID := primitive.NewObjectID()
	now := time.Now()

	quiz := &models.Quiz{
		ID:             primitive.NewObjectID(),
		StoryID:        storyID,
		Title:          fmt.Sprintf("Auto-Generated Quiz: %s", story.Title),
		Description:    fmt.Sprintf("Automatically generated quiz with %d questions from %s", autoReq.QuestionCount, story.Title),
		QuestionIDs:    questionIDs,
		TotalQuestions: autoReq.QuestionCount,
		PassingScore:   70, // Default passing score
		QuizType:       "auto_generated",
		ContentMetadata: models.ContentMetadata{
			Status:    models.StatusDraft,
			CreatedBy: userID,
			CreatedAt: now,
			UpdatedAt: now,
			Version:   1,
		},
	}

	// Insert quiz into database
	collection := db.Database.Collection("quizzes")
	_, err = collection.InsertOne(context.Background(), quiz)
	if err != nil {
		return errorResponse(500, "Failed to create auto-generated quiz")
	}

	response := convertToQuizResponse(*quiz)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func selectQuestionsForQuiz(questions []models.Question, count int, difficultyMix string) []models.Question {
	if difficultyMix == "mixed" {
		// Try to get a balanced mix of difficulties
		easy := make([]models.Question, 0)
		medium := make([]models.Question, 0)
		hard := make([]models.Question, 0)

		for _, q := range questions {
			switch q.Difficulty {
			case "easy":
				easy = append(easy, q)
			case "medium":
				medium = append(medium, q)
			case "hard":
				hard = append(hard, q)
			}
		}

		// Aim for 40% easy, 40% medium, 20% hard
		easyCount := count * 2 / 5
		mediumCount := count * 2 / 5
		hardCount := count - easyCount - mediumCount

		selected := make([]models.Question, 0, count)

		// Add easy questions
		for i := 0; i < easyCount && i < len(easy); i++ {
			selected = append(selected, easy[i])
		}

		// Add medium questions
		for i := 0; i < mediumCount && i < len(medium); i++ {
			selected = append(selected, medium[i])
		}

		// Add hard questions
		for i := 0; i < hardCount && i < len(hard); i++ {
			selected = append(selected, hard[i])
		}

		// Fill remaining slots if needed
		remaining := count - len(selected)
		allRemaining := make([]models.Question, 0)
		for _, q := range questions {
			found := false
			for _, s := range selected {
				if s.ID == q.ID {
					found = true
					break
				}
			}
			if !found {
				allRemaining = append(allRemaining, q)
			}
		}

		for i := 0; i < remaining && i < len(allRemaining); i++ {
			selected = append(selected, allRemaining[i])
		}

		return selected
	}

	// For non-mixed, just take the first 'count' questions
	if count > len(questions) {
		return questions
	}
	return questions[:count]
}

// Helper functions
func convertToQuizResponse(quiz models.Quiz) QuizResponse {
	questionIDs := make([]string, len(quiz.QuestionIDs))
	for i, id := range quiz.QuestionIDs {
		questionIDs[i] = id.Hex()
	}

	return QuizResponse{
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

func errorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	body := fmt.Sprintf(`{"error": "%s"}`, message)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// submitQuiz handles quiz submission
func submitQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	quizID := request.PathParameters["id"]
	if quizID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Quiz ID is required"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Parse submission request
	var submissionReq QuizSubmissionRequest
	if err := json.Unmarshal([]byte(request.Body), &submissionReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid request format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Convert quiz ID to ObjectID
	quizObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid quiz ID format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Get quiz from database
	var quiz models.Quiz
	collection := db.Database.Collection("quizzes")
	err = collection.FindOne(context.TODO(), bson.M{"_id": quizObjectID}).Decode(&quiz)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error": "Quiz not found"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Check if quiz is published and available
	if quiz.Status != models.StatusPublished {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Quiz is not published"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Get questions for the quiz
	questionCollection := db.Database.Collection("questions")
	questionObjectIDs := make([]primitive.ObjectID, len(quiz.QuestionIDs))
	for i, qID := range quiz.QuestionIDs {
		questionObjectIDs[i] = qID
	}

	cursor, err := questionCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": questionObjectIDs}})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to retrieve questions"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer cursor.Close(context.TODO())

	var questions []models.Question
	if err = cursor.All(context.TODO(), &questions); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to decode questions"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Calculate score and results
	var results []QuestionResult
	correctAnswers := 0
	totalPoints := 0

	for _, question := range questions {
		questionIDStr := question.ID.Hex()
		selectedAnswer, exists := submissionReq.Answers[questionIDStr]

		isCorrect := exists && selectedAnswer == question.CorrectAnswer
		points := 0
		if isCorrect {
			correctAnswers++
			points = 1 // Each question worth 1 point
		}
		totalPoints += points

		results = append(results, QuestionResult{
			QuestionID:     question.ID,
			Question:       question.Question,
			SelectedAnswer: selectedAnswer,
			CorrectAnswer:  question.CorrectAnswer,
			IsCorrect:      isCorrect,
			Points:         points,
			Explanation:    question.Explanation,
		})
	}

	// Calculate percentage and pass/fail
	percentage := float64(correctAnswers) / float64(len(questions)) * 100
	passed := percentage >= float64(quiz.PassingScore)
	timeSpent := int(submissionReq.EndTime.Sub(submissionReq.StartTime).Seconds())

	// Create submission record
	submission := QuizSubmission{
		ID:             primitive.NewObjectID(),
		QuizID:         quizObjectID,
		UserID:         primitive.NewObjectID(), // TODO: Get from authentication
		Score:          totalPoints,
		TotalQuestions: len(questions),
		CorrectAnswers: correctAnswers,
		Percentage:     percentage,
		Passed:         passed,
		TimeSpent:      timeSpent,
		SubmittedAt:    time.Now(),
		Results:        results,
	}

	// Save submission to database
	submissionCollection := db.Database.Collection("quiz_submissions")
	_, err = submissionCollection.InsertOne(context.TODO(), submission)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to save submission"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Return response
	response, err := json.Marshal(submission)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to encode response"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(response),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

// getQuizSubmissions retrieves all submissions for a quiz
func getQuizSubmissions(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	quizID := request.PathParameters["id"]
	if quizID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Quiz ID is required"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Convert quiz ID to ObjectID
	quizObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid quiz ID format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Get submissions from database
	collection := db.Database.Collection("quiz_submissions")
	cursor, err := collection.Find(context.TODO(), bson.M{"quizId": quizObjectID})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to retrieve submissions"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer cursor.Close(context.TODO())

	var submissions []QuizSubmission
	if err = cursor.All(context.TODO(), &submissions); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to decode submissions"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Return response
	response := map[string]interface{}{
		"submissions": submissions,
		"total":       len(submissions),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to encode response"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}

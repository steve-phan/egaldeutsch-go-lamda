package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/quiz-management/services"
	"egaldeutsch-serverless/netlify/functions/quiz-management/types"
	"egaldeutsch-serverless/pkg/middleware"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AutoGenerateQuiz automatically generates a quiz based on story and criteria
func AutoGenerateQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var autoReq types.AutoGenerateRequest
	if err := json.Unmarshal([]byte(request.Body), &autoReq); err != nil {
		return services.ErrorResponse(400, "Invalid request body"), nil
	}

	// Get collections
	quizzesCollection, storiesCollection, questionsCollection, _, _, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Validate story exists
	storyID, err := primitive.ObjectIDFromHex(autoReq.StoryID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid story ID"), nil
	}

	var story models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": storyID}).Decode(&story)
	if err != nil {
		return services.ErrorResponse(404, "Story not found"), nil
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
	cursor, err := questionsCollection.Find(context.Background(), filter)
	if err != nil {
		return services.ErrorResponse(500, "Failed to retrieve questions"), nil
	}
	defer cursor.Close(context.Background())

	var questions []models.Question
	if err = cursor.All(context.Background(), &questions); err != nil {
		return services.ErrorResponse(500, "Failed to decode questions"), nil
	}

	if len(questions) < autoReq.QuestionCount {
		return services.ErrorResponse(400, fmt.Sprintf("Not enough questions available. Found %d, requested %d", len(questions), autoReq.QuestionCount)), nil
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
	_, err = quizzesCollection.InsertOne(context.Background(), quiz)
	if err != nil {
		return services.ErrorResponse(500, "Failed to create auto-generated quiz"), nil
	}

	response := convertToQuizResponse(*quiz)
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// SubmitQuiz processes a quiz submission and calculates results
func SubmitQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	quizID := request.PathParameters["id"]
	if quizID == "" {
		return services.ErrorResponse(400, "Quiz ID is required"), nil
	}

	// Parse submission request
	var submissionReq types.QuizSubmissionRequest
	if err := json.Unmarshal([]byte(request.Body), &submissionReq); err != nil {
		return services.ErrorResponse(400, "Invalid request format"), nil
	}

	// Get collections
	quizzesCollection, _, questionsCollection, submissionsCollection, _, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Convert quiz ID to ObjectID
	quizObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid quiz ID format"), nil
	}

	// Get quiz from database
	var quiz models.Quiz
	err = quizzesCollection.FindOne(context.TODO(), bson.M{"_id": quizObjectID}).Decode(&quiz)
	if err != nil {
		return services.ErrorResponse(404, "Quiz not found"), nil
	}

	// Check if quiz is published and available
	if quiz.Status != models.StatusPublished {
		return services.ErrorResponse(400, "Quiz is not available for submission"), nil
	}

	// Get questions for this quiz
	questionFilter := bson.M{"_id": bson.M{"$in": quiz.QuestionIDs}}
	questionCursor, err := questionsCollection.Find(context.Background(), questionFilter)
	if err != nil {
		return services.ErrorResponse(500, "Failed to retrieve quiz questions"), nil
	}
	defer questionCursor.Close(context.Background())

	var questions []models.Question
	if err = questionCursor.All(context.Background(), &questions); err != nil {
		return services.ErrorResponse(500, "Failed to decode questions"), nil
	}

	// Create a map for easy question lookup
	questionMap := make(map[string]models.Question)
	for _, q := range questions {
		questionMap[q.ID.Hex()] = q
	}

	// Calculate results
	correctAnswers := 0
	totalQuestions := len(questions)
	totalPoints := 0
	earnedPoints := 0
	results := make([]types.QuestionResult, 0, totalQuestions)

	for questionIDStr, selectedAnswer := range submissionReq.Answers {
		if question, exists := questionMap[questionIDStr]; exists {
			totalPoints += question.Points
			isCorrect := selectedAnswer == question.CorrectAnswer

			if isCorrect {
				correctAnswers++
				earnedPoints += question.Points
			}

			result := types.QuestionResult{
				QuestionID:     question.ID,
				Question:       question.Question,
				SelectedAnswer: selectedAnswer,
				CorrectAnswer:  question.CorrectAnswer,
				IsCorrect:      isCorrect,
				Points:         question.Points,
				Explanation:    question.Explanation,
			}
			results = append(results, result)
		}
	}

	// Calculate percentage and pass status
	percentage := float64(correctAnswers) / float64(totalQuestions) * 100
	passed := percentage >= float64(quiz.PassingScore)
	timeSpent := int(submissionReq.EndTime.Sub(submissionReq.StartTime).Seconds())

	// Create submission record
	submission := types.QuizSubmission{
		ID:             primitive.NewObjectID(),
		QuizID:         quizObjectID,
		UserID:         primitive.NewObjectID(), // Should get from auth context
		Score:          earnedPoints,
		TotalQuestions: totalQuestions,
		CorrectAnswers: correctAnswers,
		Percentage:     percentage,
		Passed:         passed,
		TimeSpent:      timeSpent,
		SubmittedAt:    time.Now(),
		Results:        results,
	}

	// Save submission to database
	_, err = submissionsCollection.InsertOne(context.Background(), submission)
	if err != nil {
		return services.ErrorResponse(500, "Failed to save quiz submission"), nil
	}

	responseBody, _ := json.Marshal(submission)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// GetQuizSubmissions retrieves submissions for a specific quiz
func GetQuizSubmissions(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	quizID := request.PathParameters["id"]
	if quizID == "" {
		return services.ErrorResponse(400, "Quiz ID is required"), nil
	}

	// Get collections
	_, _, _, submissionsCollection, _, err := services.GetCollections()
	if err != nil {
		return services.ErrorResponse(500, "Database connection failed"), nil
	}

	// Convert quiz ID to ObjectID
	quizObjectID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return services.ErrorResponse(400, "Invalid quiz ID format"), nil
	}

	// Query submissions for this quiz
	filter := bson.M{"quizId": quizObjectID}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "submittedAt", Value: -1}}) // Most recent first

	cursor, err := submissionsCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return services.ErrorResponse(500, "Failed to retrieve submissions"), nil
	}
	defer cursor.Close(context.Background())

	var submissions []types.QuizSubmission
	if err = cursor.All(context.Background(), &submissions); err != nil {
		return services.ErrorResponse(500, "Failed to decode submissions"), nil
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"submissions": submissions,
		"total":       len(submissions),
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    middleware.GetAuthenticatedCORSHeaders(),
	}, nil
}

// selectQuestionsForQuiz selects questions based on criteria and difficulty mix
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

		// Fill remaining slots with any available questions
		remaining := count - len(selected)
		if remaining > 0 {
			allRemaining := make([]models.Question, 0)
			usedIDs := make(map[primitive.ObjectID]bool)

			for _, q := range selected {
				usedIDs[q.ID] = true
			}

			for _, q := range questions {
				if !usedIDs[q.ID] {
					allRemaining = append(allRemaining, q)
				}
			}

			for i := 0; i < remaining && i < len(allRemaining); i++ {
				selected = append(selected, allRemaining[i])
			}
		}

		return selected
	}

	// For non-mixed difficulty, just take the first 'count' questions
	if count > len(questions) {
		count = len(questions)
	}

	return questions[:count]
}

package types

import (
	"time"

	"egaldeutsch-serverless/models"

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

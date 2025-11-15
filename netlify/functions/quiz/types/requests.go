package types

import "time"

// QuizSubmissionRequest represents the request body for quiz submission
type QuizSubmissionRequest struct {
	Answers []int `json:"answers"`
}

// QuizResult represents the detailed result of a quiz submission
type QuizResult struct {
	ID             string    `json:"id"`
	StoryID        string    `json:"storyId"`
	Answers        []int     `json:"answers"`
	Score          int       `json:"score"`
	TotalQuestions int       `json:"totalQuestions"`
	Percentage     float64   `json:"percentage"`
	Passed         bool      `json:"passed"`
	CorrectAnswers []bool    `json:"correctAnswers"`
	EarnedPoints   int       `json:"earnedPoints"`
	TotalPoints    int       `json:"totalPoints"`
	SubmittedAt    time.Time `json:"submittedAt"`
}

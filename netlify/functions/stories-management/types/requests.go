package types

import (
	"time"

	"egaldeutsch-serverless/models"
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

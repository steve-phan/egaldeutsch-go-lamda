package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// QuestionGenerated represents a generated question from OpenAI
type QuestionGenerated struct {
	Question      string   `json:"question"`
	QuestionType  string   `json:"questionType"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correctAnswer"`
	Explanation   string   `json:"explanation"`
	Difficulty    string   `json:"difficulty"`
	Points        int      `json:"points"`
	GermanConcept string   `json:"germanConcept"`
}

// VocabularyExtracted represents vocabulary extracted from story
type VocabularyExtracted struct {
	German   string `json:"german"`
	English  string `json:"english"`
	WordType string `json:"wordType"`
	Article  string `json:"article,omitempty"`
	Context  string `json:"context"`
}

// QuestionGenerationResponse represents the OpenAI response for questions
type QuestionGenerationResponse struct {
	Questions           []QuestionGenerated   `json:"questions"`
	VocabularyExtracted []VocabularyExtracted `json:"vocabularyExtracted"`
}

// QuizGenerationResponse represents the OpenAI response for quiz metadata
type QuizGenerationResponse struct {
	Title                    string   `json:"title"`
	Description              string   `json:"description"`
	QuizType                 string   `json:"quizType"`
	EstimatedTime            int      `json:"estimatedTime"`
	PassingScore             int      `json:"passingScore"`
	RecommendedQuestionTypes []string `json:"recommendedQuestionTypes"`
	DifficultyDistribution   struct {
		Easy   int `json:"easy"`
		Medium int `json:"medium"`
		Hard   int `json:"hard"`
	} `json:"difficultyDistribution"`
}

// GenerationResult represents the API response
type GenerationResult struct {
	Success        bool     `json:"success"`
	Message        string   `json:"message"`
	QuestionsCount int      `json:"questionsCount,omitempty"`
	QuizID         string   `json:"quizId,omitempty"`
	QuestionIDs    []string `json:"questionIds,omitempty"`
}

// GenerationRequest represents the request parameters
type GenerationRequest struct {
	StoryID primitive.ObjectID `json:"storyId"`
	Type    string             `json:"type"` // "questions", "quiz", or "both"
}

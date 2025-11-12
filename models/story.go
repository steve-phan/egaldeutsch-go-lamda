package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Story represents a German learning story
type Story struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Content     string             `bson:"content" json:"content"`
	Level       string             `bson:"level" json:"level"` // A1, A2, B1, B2, C1, C2
	WordCount   int                `bson:"wordCount" json:"wordCount"`
	ReadingTime int                `bson:"readingTime" json:"readingTime"` // Estimated reading time in minutes
	Topic       string             `bson:"topic" json:"topic"`             // e.g., "Family", "Travel", "Food"
	Vocabulary  []VocabularyWord   `bson:"vocabulary" json:"vocabulary"`   // Key vocabulary from the story
	Summary     string             `bson:"summary" json:"summary"`         // Brief story summary
	IsActive    bool               `bson:"isActive" json:"isActive"`       // Whether story is published
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// VocabularyWord represents a key German word from the story
type VocabularyWord struct {
	German   string `bson:"german" json:"german"`
	English  string `bson:"english" json:"english"`
	WordType string `bson:"wordType" json:"wordType"` // "noun", "verb", "adjective", etc.
	Article  string `bson:"article" json:"article"`   // "der", "die", "das" for nouns
}

// Question represents a quiz question for a story
type Question struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StoryID       primitive.ObjectID `bson:"storyId" json:"storyId"`
	Question      string             `bson:"question" json:"question"`
	QuestionType  string             `bson:"questionType" json:"questionType"`   // "comprehension", "vocabulary", "grammar"
	Options       []string           `bson:"options" json:"options"`             // Array of answer options
	CorrectAnswer int                `bson:"correctAnswer" json:"correctAnswer"` // Index of correct answer in options array
	Explanation   string             `bson:"explanation" json:"explanation"`
	Points        int                `bson:"points" json:"points"` // Points awarded for correct answer
	Order         int                `bson:"order" json:"order"`   // Question order in quiz (1-10)
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}

// Quiz represents a complete quiz for a story
type Quiz struct {
	StoryID   primitive.ObjectID `bson:"storyId" json:"storyId"`
	Story     *Story             `bson:"story,omitempty" json:"story,omitempty"`
	Questions []Question         `bson:"questions" json:"questions"`
}

// QuizSubmission represents a user's quiz submission
type QuizSubmission struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StoryID        primitive.ObjectID `bson:"storyId" json:"storyId"`
	Answers        []int              `bson:"answers" json:"answers"` // Array of selected answer indices
	Score          int                `bson:"score" json:"score"`     // Number of correct answers
	TotalQuestions int                `bson:"totalQuestions" json:"totalQuestions"`
	SubmittedAt    time.Time          `bson:"submittedAt" json:"submittedAt"`
}

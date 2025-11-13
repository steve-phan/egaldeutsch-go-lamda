package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContentStatus represents the status of content in the approval workflow
type ContentStatus string

const (
	StatusDraft         ContentStatus = "draft"
	StatusPendingReview ContentStatus = "pending_review"
	StatusApproved      ContentStatus = "approved"
	StatusActive        ContentStatus = "active"
	StatusInactive      ContentStatus = "inactive"
	StatusRejected      ContentStatus = "rejected"
	StatusArchived      ContentStatus = "archived"
)

// ReviewComment represents a comment in the review process
type ReviewComment struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	ReviewerID primitive.ObjectID `bson:"reviewerId" json:"reviewerId"`
	Comment    string             `bson:"comment" json:"comment"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	Type       string             `bson:"type" json:"type"` // "feedback", "rejection_reason", "approval_note"
}

// ContentMetadata contains common fields for content management
type ContentMetadata struct {
	Status      ContentStatus       `bson:"status" json:"status"`
	CreatedBy   primitive.ObjectID  `bson:"createdBy" json:"createdBy"`
	CreatedAt   time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updatedAt" json:"updatedAt"`
	ReviewedBy  *primitive.ObjectID `bson:"reviewedBy,omitempty" json:"reviewedBy,omitempty"`
	ReviewedAt  *time.Time          `bson:"reviewedAt,omitempty" json:"reviewedAt,omitempty"`
	ApprovedAt  *time.Time          `bson:"approvedAt,omitempty" json:"approvedAt,omitempty"`
	ActivatedAt *time.Time          `bson:"activatedAt,omitempty" json:"activatedAt,omitempty"`
	Comments    []ReviewComment     `bson:"comments,omitempty" json:"comments,omitempty"`
	Version     int                 `bson:"version" json:"version"`
}

// Story represents a German learning story
type Story struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Content     string             `bson:"content" json:"content"`
	Level       string             `bson:"level" json:"level"` // A1, A2, B1, B2, C1, C2
	WordCount   int                `bson:"wordCount" json:"wordCount"`
	ReadingTime int                `bson:"readingTime" json:"readingTime"` // Estimated reading time in minutes
	Topics      []string           `bson:"topics" json:"topics"`           // e.g., ["Family", "Travel", "Food"]
	Vocabulary  []VocabularyWord   `bson:"vocabulary" json:"vocabulary"`   // Key vocabulary from the story
	Summary     string             `bson:"summary" json:"summary"`         // Brief story summary

	// Content Management Fields
	ContentMetadata `bson:",inline"`

	// Deprecated: Use Status field instead
	IsActive bool `bson:"isActive" json:"isActive"` // For backward compatibility
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
	Points        int                `bson:"points" json:"points"`         // Points awarded for correct answer
	Order         int                `bson:"order" json:"order"`           // Question order in quiz (1-10)
	Difficulty    string             `bson:"difficulty" json:"difficulty"` // "easy", "medium", "hard"

	// Content Management Fields
	ContentMetadata `bson:",inline"`
}

// Quiz represents a complete quiz for a story
type Quiz struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	StoryID        primitive.ObjectID   `bson:"storyId" json:"storyId"`
	Title          string               `bson:"title" json:"title"`
	Description    string               `bson:"description" json:"description"`
	QuestionIDs    []primitive.ObjectID `bson:"questionIds" json:"questionIds"` // References to questions
	TotalQuestions int                  `bson:"totalQuestions" json:"totalQuestions"`
	TotalPoints    int                  `bson:"totalPoints" json:"totalPoints"`
	TimeLimit      int                  `bson:"timeLimit" json:"timeLimit"`       // Time limit in minutes (0 = no limit)
	PassingScore   int                  `bson:"passingScore" json:"passingScore"` // Minimum score to pass (percentage)
	QuizType       string               `bson:"quizType" json:"quizType"`         // "auto_generated", "manual", "mixed"

	// Content Management Fields
	ContentMetadata `bson:",inline"`

	// Populated fields (not stored in DB)
	Story     *Story     `bson:"-" json:"story,omitempty"`
	Questions []Question `bson:"-" json:"questions,omitempty"`
}

// QuizSubmission represents a user's quiz submission
type QuizSubmission struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuizID         primitive.ObjectID `bson:"quizId" json:"quizId"`
	StoryID        primitive.ObjectID `bson:"storyId" json:"storyId"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`   // User who submitted
	Username       string             `bson:"username,omitempty" json:"username,omitempty"` // Denormalized for leaderboard
	Answers        []int              `bson:"answers" json:"answers"` // Array of selected answer indices
	Score          int                `bson:"score" json:"score"`     // Number of correct answers
	TotalQuestions int                `bson:"totalQuestions" json:"totalQuestions"`
	TotalPoints    int                `bson:"totalPoints" json:"totalPoints"`
	EarnedPoints   int                `bson:"earnedPoints" json:"earnedPoints"`
	Percentage     float64            `bson:"percentage" json:"percentage"`
	Passed         bool               `bson:"passed" json:"passed"`
	TimeSpent      int                `bson:"timeSpent" json:"timeSpent"` // Time spent in seconds
	SubmittedAt    time.Time          `bson:"submittedAt" json:"submittedAt"`
}

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	UserID          primitive.ObjectID `bson:"userId" json:"userId"`
	Username        string             `bson:"username" json:"username"`
	TotalPoints     int                `bson:"totalPoints" json:"totalPoints"`
	QuizzesTaken    int                `bson:"quizzesTaken" json:"quizzesTaken"`
	QuizzesPassed   int                `bson:"quizzesPassed" json:"quizzesPassed"`
	AverageScore    float64            `bson:"averageScore" json:"averageScore"`
	Rank            int                `bson:"-" json:"rank"` // Computed field
	LastSubmittedAt time.Time          `bson:"lastSubmittedAt" json:"lastSubmittedAt"`
}

// User represents a system user (creators, reviewers, admins)
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash" json:"-"` // Never include in JSON
	FirstName    string             `bson:"firstName" json:"firstName"`
	LastName     string             `bson:"lastName" json:"lastName"`
	Name         string             `bson:"name" json:"name"` // Computed from FirstName + LastName
	Role         UserRole           `bson:"role" json:"role"`
	Status       UserStatus         `bson:"status" json:"status"`
	IsActive     bool               `bson:"isActive" json:"isActive"` // Backward compatibility
	LastLoginAt  *time.Time         `bson:"lastLoginAt,omitempty" json:"lastLoginAt,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// UserRole represents user permission levels
type UserRole string

const (
	RoleCreator  UserRole = "creator"
	RoleReviewer UserRole = "reviewer"
	RoleAdmin    UserRole = "admin"
)

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending   UserStatus = "pending"
)

// Dashboard represents dashboard statistics
type DashboardStats struct {
	TotalStories     int `json:"totalStories"`
	PendingStories   int `json:"pendingStories"`
	ActiveStories    int `json:"activeStories"`
	TotalQuestions   int `json:"totalQuestions"`
	PendingQuestions int `json:"pendingQuestions"`
	ActiveQuestions  int `json:"activeQuestions"`
	TotalQuizzes     int `json:"totalQuizzes"`
	ActiveQuizzes    int `json:"activeQuizzes"`
	TotalSubmissions int `json:"totalSubmissions"`
}

// ContentListFilter represents filters for content listing
type ContentListFilter struct {
	Status    []ContentStatus     `json:"status,omitempty"`
	Level     []string            `json:"level,omitempty"`
	Topics    []string            `json:"topics,omitempty"`
	CreatedBy *primitive.ObjectID `json:"createdBy,omitempty"`
	DateFrom  *time.Time          `json:"dateFrom,omitempty"`
	DateTo    *time.Time          `json:"dateTo,omitempty"`
	Search    string              `json:"search,omitempty"`
	Page      int                 `json:"page"`
	Limit     int                 `json:"limit"`
	SortBy    string              `json:"sortBy"`    // "createdAt", "updatedAt", "title"
	SortOrder int                 `json:"sortOrder"` // 1 for ascending, -1 for descending
}

// Status transition validation methods

// CanTransitionTo checks if content can transition to target status
func (s ContentStatus) CanTransitionTo(target ContentStatus) bool {
	validTransitions := map[ContentStatus][]ContentStatus{
		StatusDraft:         {StatusPendingReview, StatusArchived},
		StatusPendingReview: {StatusApproved, StatusRejected, StatusArchived},
		StatusApproved:      {StatusActive, StatusInactive, StatusArchived},
		StatusActive:        {StatusInactive, StatusArchived},
		StatusInactive:      {StatusActive, StatusArchived},
		StatusRejected:      {StatusDraft, StatusArchived},
		StatusArchived:      {}, // Cannot transition from archived
	}

	allowed, exists := validTransitions[s]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowed {
		if allowedStatus == target {
			return true
		}
	}
	return false
}

// Validation methods for Story
func (s *Story) Validate() error {
	if len(s.Title) < 10 || len(s.Title) > 200 {
		return fmt.Errorf("title must be between 10 and 200 characters")
	}

	if len(s.Content) < 100 || len(s.Content) > 15000 {
		return fmt.Errorf("content must be between 100 and 15000 characters")
	}

	validLevels := map[string]bool{"A1": true, "A2": true, "B1": true, "B2": true, "C1": true, "C2": true}
	if !validLevels[s.Level] {
		return fmt.Errorf("level must be one of: A1, A2, B1, B2, C1, C2")
	}

	if len(s.Topics) == 0 || len(s.Topics) > 5 {
		return fmt.Errorf("must have between 1 and 5 topics")
	}

	return nil
}

// Validation methods for Question
func (q *Question) Validate() error {
	if len(q.Question) < 10 || len(q.Question) > 300 {
		return fmt.Errorf("question text must be between 10 and 300 characters")
	}

	if len(q.Options) < 2 || len(q.Options) > 6 {
		return fmt.Errorf("must have between 2 and 6 options")
	}

	if q.CorrectAnswer < 0 || q.CorrectAnswer >= len(q.Options) {
		return fmt.Errorf("correct answer index is invalid")
	}

	if q.Points < 1 || q.Points > 10 {
		return fmt.Errorf("points must be between 1 and 10")
	}

	for i, option := range q.Options {
		if len(option) < 1 || len(option) > 100 {
			return fmt.Errorf("option %d must be between 1 and 100 characters", i+1)
		}
	}

	validTypes := map[string]bool{"comprehension": true, "vocabulary": true, "grammar": true}
	if !validTypes[q.QuestionType] {
		return fmt.Errorf("question type must be one of: comprehension, vocabulary, grammar")
	}

	return nil
}

// Validation methods for Quiz
func (qz *Quiz) Validate() error {
	if len(qz.Title) < 5 || len(qz.Title) > 200 {
		return fmt.Errorf("quiz title must be between 5 and 200 characters")
	}

	if qz.TotalQuestions < 5 || qz.TotalQuestions > 20 {
		return fmt.Errorf("quiz must have between 5 and 20 questions")
	}

	if qz.PassingScore < 1 || qz.PassingScore > 100 {
		return fmt.Errorf("passing score must be between 1 and 100 percent")
	}

	validTypes := map[string]bool{"auto_generated": true, "manual": true, "mixed": true}
	if !validTypes[qz.QuizType] {
		return fmt.Errorf("quiz type must be one of: auto_generated, manual, mixed")
	}

	return nil
}

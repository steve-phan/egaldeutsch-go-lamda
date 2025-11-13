package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	Title                    string `json:"title"`
	Description              string `json:"description"`
	QuizType                 string `json:"quizType"`
	EstimatedTime            int    `json:"estimatedTime"`
	PassingScore             int    `json:"passingScore"`
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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		return errorResponse(500, "Database connection failed")
	}
	defer db.Disconnect()

	// Only POST is allowed
	if request.HTTPMethod != "POST" {
		return errorResponse(405, "Method not allowed")
	}

	// Get query parameters
	generationType := request.QueryStringParameters["type"]
	storyIDStr := request.QueryStringParameters["story_id"]

	if generationType == "" || storyIDStr == "" {
		return errorResponse(400, "Missing required parameters: type and story_id")
	}

	// Validate generation type
	validTypes := map[string]bool{"questions": true, "quiz": true, "both": true}
	if !validTypes[generationType] {
		return errorResponse(400, "Invalid type. Must be: questions, quiz, or both")
	}

	// Parse and validate story ID
	storyID, err := primitive.ObjectIDFromHex(storyIDStr)
	if err != nil {
		return errorResponse(400, "Invalid story ID")
	}

	// Fetch the story
	story, err := getStory(storyID)
	if err != nil {
		return errorResponse(404, fmt.Sprintf("Story not found: %v", err))
	}

	// Check if story is published
	if story.Status != models.StatusPublished {
		return errorResponse(400, "Story must be published to generate content")
	}

	// Check for OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return errorResponse(500, "OpenAI API key not configured")
	}

	// Route to appropriate generation function
	var result GenerationResult
	switch generationType {
	case "questions":
		result, err = generateQuestions(story, apiKey)
	case "quiz":
		result, err = generateQuiz(story, apiKey)
	case "both":
		result, err = generateBoth(story, apiKey)
	}

	if err != nil {
		return errorResponse(500, fmt.Sprintf("Generation failed: %v", err))
	}

	responseBody, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func getStory(storyID primitive.ObjectID) (*models.Story, error) {
	collection := db.Database.Collection("stories")
	var story models.Story
	err := collection.FindOne(context.Background(), bson.M{"_id": storyID}).Decode(&story)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

func generateQuestions(story *models.Story, apiKey string) (GenerationResult, error) {
	// Build prompt for question generation with JSON format instruction
	prompt := buildQuestionPrompt(story) + `

Please respond with valid JSON in exactly this format:
{
  "questions": [
    {
      "question": "Question text here",
      "questionType": "comprehension|vocabulary|grammar",
      "options": ["option1", "option2", "option3", "option4"],
      "correctAnswer": 0,
      "explanation": "Explanation here",
      "difficulty": "easy|medium|hard",
      "points": 5,
      "germanConcept": "Concept being tested"
    }
  ],
  "vocabularyExtracted": [
    {
      "german": "word",
      "english": "translation",
      "wordType": "noun",
      "article": "der",
      "context": "context in story"
    }
  ]
}`

	// Call OpenAI API with simple JSON mode
	client := openai.NewClient(option.WithAPIKey(apiKey))
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	jsonObjectPtr := shared.NewResponseFormatJSONObjectParam()
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &jsonObjectPtr,
		},
	}

	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return GenerationResult{}, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(completion.Choices) == 0 {
		return GenerationResult{}, fmt.Errorf("no response from OpenAI")
	}

	// Parse the response
	var response QuestionGenerationResponse
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &response); err != nil {
		return GenerationResult{}, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	// Store questions in database
	questionIDs, err := storeQuestions(story.ID, response.Questions)
	if err != nil {
		return GenerationResult{}, fmt.Errorf("failed to store questions: %w", err)
	}

	return GenerationResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully generated %d questions", len(questionIDs)),
		QuestionsCount: len(questionIDs),
		QuestionIDs:    questionIDs,
	}, nil
}

func generateQuiz(story *models.Story, apiKey string) (GenerationResult, error) {
	// Build prompt for quiz generation with JSON format instruction
	prompt := buildQuizPrompt(story) + `

Please respond with valid JSON in exactly this format:
{
  "title": "Quiz title here",
  "description": "Quiz description",
  "quizType": "comprehension|vocabulary|mixed",
  "estimatedTime": 15,
  "passingScore": 70,
  "recommendedQuestionTypes": ["comprehension", "vocabulary"],
  "difficultyDistribution": {
    "easy": 3,
    "medium": 5,
    "hard": 2
  }
}`

	client := openai.NewClient(option.WithAPIKey(apiKey))
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	jsonObjectPtr := shared.NewResponseFormatJSONObjectParam()
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &jsonObjectPtr,
		},
	}

	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return GenerationResult{}, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(completion.Choices) == 0 {
		return GenerationResult{}, fmt.Errorf("no response from OpenAI")
	}

	// Parse the response
	var response QuizGenerationResponse
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &response); err != nil {
		return GenerationResult{}, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	// Get existing questions for this story to create quiz
	existingQuestions, err := getStoryQuestions(story.ID)
	if err != nil {
		return GenerationResult{}, fmt.Errorf("failed to fetch existing questions: %w", err)
	}

	if len(existingQuestions) == 0 {
		return GenerationResult{}, fmt.Errorf("no questions available for quiz. Generate questions first")
	}

	// Store quiz in database
	quizID, err := storeQuiz(story.ID, response, existingQuestions)
	if err != nil {
		return GenerationResult{}, fmt.Errorf("failed to store quiz: %w", err)
	}

	return GenerationResult{
		Success: true,
		Message: fmt.Sprintf("Successfully generated quiz: %s", response.Title),
		QuizID:  quizID,
	}, nil
}

func generateBoth(story *models.Story, apiKey string) (GenerationResult, error) {
	// First generate questions
	questionsResult, err := generateQuestions(story, apiKey)
	if err != nil {
		return GenerationResult{}, fmt.Errorf("failed to generate questions: %w", err)
	}

	// Then generate quiz
	quizResult, err := generateQuiz(story, apiKey)
	if err != nil {
		return GenerationResult{}, fmt.Errorf("failed to generate quiz: %w", err)
	}

	return GenerationResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully generated %d questions and quiz", questionsResult.QuestionsCount),
		QuestionsCount: questionsResult.QuestionsCount,
		QuestionIDs:    questionsResult.QuestionIDs,
		QuizID:         quizResult.QuizID,
	}, nil
}

func buildQuestionPrompt(story *models.Story) string {
	vocabStr := formatVocabulary(story.Vocabulary)
	topicsStr := strings.Join(story.Topics, ", ")

	return fmt.Sprintf(`You are an expert German language teacher creating assessment materials for CEFR level %s.

STORY CONTENT:
Title: %s
Level: %s
Content: %s
Topics: %s
Existing Vocabulary: %s
Word Count: %d

GENERATION REQUIREMENTS:
- Create 8-12 diverse questions covering the story content
- Ensure questions test genuine comprehension, not just recall
- Vocabulary questions should use words actually present in the story
- Grammar questions should focus on structures used in the text
- All questions must be answerable from the story content
- Difficulty should match CEFR level %s
- Include cultural context when relevant
- Provide clear, helpful explanations

QUESTION DISTRIBUTION:
- 40%% comprehension (main ideas, details, inference)
- 35%% vocabulary (definitions, usage, context)
- 25%% grammar (structures from the text)

DIFFICULTY DISTRIBUTION:
- 30%% easy questions
- 50%% medium questions
- 20%% hard questions

QUALITY STANDARDS:
- Questions should be unambiguous with clear correct answers
- Distractors (wrong options) should be plausible but clearly incorrect
- German grammar terminology should be accurate
- Explanations should aid learning, not just confirm correctness
- Each question should have exactly 4 options
- Points should range from 5-10 based on difficulty

Please generate high-quality questions following these guidelines.`,
		story.Level, story.Title, story.Level, story.Content, topicsStr, vocabStr, story.WordCount, story.Level)
}

func buildQuizPrompt(story *models.Story) string {
	topicsStr := strings.Join(story.Topics, ", ")

	return fmt.Sprintf(`You are an expert German language teacher creating a quiz for CEFR level %s learners.

STORY CONTENT:
Title: %s
Level: %s
Topics: %s
Word Count: %d
Reading Time: %d minutes

QUIZ GENERATION REQUIREMENTS:
Create a complete quiz structure with:
- Engaging title related to story content
- Clear description explaining what students will be assessed on
- Appropriate time estimate for completion (5-30 minutes)
- Recommended question type distribution
- Balanced difficulty progression
- Passing score between 50-90%%

The quiz should be appropriate for %s level learners and assess comprehension of the story "%s".

Consider the story topics (%s) when creating the quiz metadata.`,
		story.Level, story.Title, story.Level, topicsStr, story.WordCount, story.ReadingTime,
		story.Level, story.Title, topicsStr)
}

func formatVocabulary(vocabulary []models.VocabularyWord) string {
	if len(vocabulary) == 0 {
		return "No vocabulary provided"
	}

	var vocabStrs []string
	for _, v := range vocabulary {
		if v.Article != "" {
			vocabStrs = append(vocabStrs, fmt.Sprintf("%s %s (%s) - %s", v.Article, v.German, v.WordType, v.English))
		} else {
			vocabStrs = append(vocabStrs, fmt.Sprintf("%s (%s) - %s", v.German, v.WordType, v.English))
		}
	}
	return strings.Join(vocabStrs, ", ")
}

func getAIGeneratorUserID() primitive.ObjectID {
	// Check if there's a special AI generator user in the system
	// For now, create a consistent ObjectID for AI-generated content
	collection := db.Database.Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"username": "ai-generator"}).Decode(&user)
	if err == nil {
		return user.ID
	}

	// Create a deterministic ObjectID for AI generator
	// This ensures consistency across invocations
	aiGenID, _ := primitive.ObjectIDFromHex("000000000000000000000001")
	return aiGenID
}

func storeQuestions(storyID primitive.ObjectID, questions []QuestionGenerated) ([]string, error) {
	collection := db.Database.Collection("questions")
	questionIDs := make([]string, 0, len(questions))
	now := time.Now()
	aiUserID := getAIGeneratorUserID()

	for i, q := range questions {
		// Map difficulty to points if not in valid range
		points := q.Points
		if points < 1 || points > 10 {
			switch q.Difficulty {
			case "easy":
				points = 5
			case "medium":
				points = 7
			case "hard":
				points = 10
			default:
				points = 7
			}
		}

		question := models.Question{
			ID:            primitive.NewObjectID(),
			StoryID:       storyID,
			Question:      q.Question,
			QuestionType:  q.QuestionType,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
			Explanation:   q.Explanation,
			Points:        points,
			Order:         i + 1,
			Difficulty:    q.Difficulty,
			ContentMetadata: models.ContentMetadata{
				Status:    models.StatusDraft,
				CreatedBy: aiUserID,
				CreatedAt: now,
				UpdatedAt: now,
				Version:   1,
			},
		}

		// Validate before inserting
		if err := question.Validate(); err != nil {
			return nil, fmt.Errorf("question %d validation failed: %w", i+1, err)
		}

		_, err := collection.InsertOne(context.Background(), question)
		if err != nil {
			return nil, fmt.Errorf("failed to insert question %d: %w", i+1, err)
		}

		questionIDs = append(questionIDs, question.ID.Hex())
	}

	return questionIDs, nil
}

func getStoryQuestions(storyID primitive.ObjectID) ([]models.Question, error) {
	collection := db.Database.Collection("questions")
	cursor, err := collection.Find(context.Background(), bson.M{
		"storyId": storyID,
		"status":  bson.M{"$in": []models.ContentStatus{models.StatusDraft, models.StatusPreview, models.StatusReady, models.StatusPublished}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var questions []models.Question
	if err := cursor.All(context.Background(), &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

func storeQuiz(storyID primitive.ObjectID, quizData QuizGenerationResponse, questions []models.Question) (string, error) {
	collection := db.Database.Collection("quizzes")
	now := time.Now()
	aiUserID := getAIGeneratorUserID()

	// Extract question IDs
	questionIDs := make([]primitive.ObjectID, 0, len(questions))
	totalPoints := 0
	for _, q := range questions {
		questionIDs = append(questionIDs, q.ID)
		totalPoints += q.Points
	}

	// Calculate time limit based on estimated time
	timeLimit := quizData.EstimatedTime

	quiz := models.Quiz{
		ID:             primitive.NewObjectID(),
		StoryID:        storyID,
		Title:          quizData.Title,
		Description:    quizData.Description,
		QuestionIDs:    questionIDs,
		TotalQuestions: len(questionIDs),
		TotalPoints:    totalPoints,
		TimeLimit:      timeLimit,
		PassingScore:   quizData.PassingScore,
		QuizType:       "auto_generated",
		ContentMetadata: models.ContentMetadata{
			Status:    models.StatusDraft,
			CreatedBy: aiUserID,
			CreatedAt: now,
			UpdatedAt: now,
			Version:   1,
		},
	}

	// Validate before inserting
	if err := quiz.Validate(); err != nil {
		return "", fmt.Errorf("quiz validation failed: %w", err)
	}

	_, err := collection.InsertOne(context.Background(), quiz)
	if err != nil {
		return "", fmt.Errorf("failed to insert quiz: %w", err)
	}

	return quiz.ID.Hex(), nil
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

func main() {
	lambda.Start(handler)
}

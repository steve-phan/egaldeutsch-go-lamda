package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/ai-generator/services"
	"egaldeutsch-serverless/netlify/functions/ai-generator/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateContent handles the main AI generation request
func GenerateContent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	handlerStartTime := time.Now()
	fmt.Printf("🚀 AI Generator handler started at %v\n", handlerStartTime)

	// Get query parameters
	generationType := request.QueryStringParameters["type"]
	storyIDStr := request.QueryStringParameters["story_id"]

	if generationType == "" || storyIDStr == "" {
		return services.ErrorResponse(400, "Missing required parameters: type and story_id"), nil
	}

	// Validate generation type
	validTypes := map[string]bool{"questions": true, "quiz": true, "both": true}
	if !validTypes[generationType] {
		return services.ErrorResponse(400, "Invalid type. Must be: questions, quiz, or both"), nil
	}

	// Parse and validate story ID
	storyID, err := primitive.ObjectIDFromHex(storyIDStr)
	if err != nil {
		return services.ErrorResponse(400, "Invalid story ID"), nil
	}

	// Fetch the story
	story, err := getStory(storyID)
	if err != nil {
		return services.ErrorResponse(404, fmt.Sprintf("Story not found: %v", err)), nil
	}

	// Check if AI questions have already been generated for questions or both
	if (generationType == "questions" || generationType == "both") && story.IsAIQuestionsGenerated {
		return services.ErrorResponse(409, "AI questions have already been generated for this story"), nil
	}

	// Get OpenAI API key
	apiKey := services.GetOpenAIAPIKey()
	if apiKey == "" {
		return services.ErrorResponse(500, "OpenAI API key not configured"), nil
	}

	var result types.GenerationResult

	// Generate based on type
	switch generationType {
	case "questions":
		result, err = generateQuestions(story, apiKey)
	case "quiz":
		result, err = generateQuiz(story, apiKey)
	case "both":
		result, err = generateBoth(story, apiKey)
	}

	if err != nil {
		return services.ErrorResponse(500, fmt.Sprintf("Generation failed: %v", err)), nil
	}

	responseBody, _ := json.Marshal(result)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    services.GetCORSHeaders(),
	}, nil
}

// getStory fetches a story from the database
func getStory(storyID primitive.ObjectID) (*models.Story, error) {
	storiesCollection, _, _, err := services.GetCollections()
	if err != nil {
		return nil, err
	}

	var story models.Story
	err = storiesCollection.FindOne(context.Background(), bson.M{"_id": storyID}).Decode(&story)
	return &story, err
}

// generateQuestions generates questions for a story
func generateQuestions(story *models.Story, apiKey string) (types.GenerationResult, error) {
	fmt.Printf("🔍 Generating questions for story: %s\n", story.Title)

	// Create OpenAI client
	client := openai.NewClient(option.WithAPIKey(apiKey))

	// Build prompt
	prompt := services.BuildQuestionPrompt(story)

	// Make OpenAI request
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

	completion, err := client.Chat.Completions.New(context.Background(), params)
	if err != nil {
		return types.GenerationResult{}, fmt.Errorf("OpenAI API error: %w", err)
	}

	fmt.Printf("✅ OpenAI response received: %d characters\n", len(completion.Choices[0].Message.Content))

	// Parse response
	var questionResponse types.QuestionGenerationResponse
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &questionResponse); err != nil {
		return types.GenerationResult{}, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	fmt.Printf("📝 Generated %d questions\n", len(questionResponse.Questions))

	// Store questions in database
	questionIDs, err := storeQuestions(story.ID, questionResponse.Questions)
	if err != nil {
		return types.GenerationResult{}, fmt.Errorf("failed to store questions: %w", err)
	}

	// Mark story as having AI questions generated
	if err := markStoryAIQuestionsGenerated(story.ID); err != nil {
		return types.GenerationResult{}, fmt.Errorf("failed to update story status: %w", err)
	}

	return types.GenerationResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully generated %d questions for story '%s'", len(questionResponse.Questions), story.Title),
		QuestionsCount: len(questionResponse.Questions),
		QuestionIDs:    questionIDs,
	}, nil
}

// generateQuiz generates quiz metadata for a story
func generateQuiz(story *models.Story, apiKey string) (types.GenerationResult, error) {
	fmt.Printf("🎯 Generating quiz for story: %s\n", story.Title)

	// Create OpenAI client
	client := openai.NewClient(option.WithAPIKey(apiKey))

	// Build prompt
	prompt := services.BuildQuizPrompt(story)

	// Make OpenAI request
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

	completion, err := client.Chat.Completions.New(context.Background(), params)
	if err != nil {
		return types.GenerationResult{}, fmt.Errorf("OpenAI API error: %w", err)
	}

	// Parse response
	var quizResponse types.QuizGenerationResponse
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &quizResponse); err != nil {
		return types.GenerationResult{}, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	// Get existing questions for the story
	questions, err := getStoryQuestions(story.ID)
	if err != nil {
		return types.GenerationResult{}, fmt.Errorf("failed to get story questions: %w", err)
	}

	if len(questions) == 0 {
		return types.GenerationResult{}, fmt.Errorf("no questions found for story. Please generate questions first")
	}

	// Store quiz in database
	quizID, err := storeQuiz(story.ID, quizResponse, questions)
	if err != nil {
		return types.GenerationResult{}, fmt.Errorf("failed to store quiz: %w", err)
	}

	return types.GenerationResult{
		Success: true,
		Message: fmt.Sprintf("Successfully generated quiz '%s' for story '%s'", quizResponse.Title, story.Title),
		QuizID:  quizID,
	}, nil
}

// generateBoth generates both questions and quiz for a story
func generateBoth(story *models.Story, apiKey string) (types.GenerationResult, error) {
	fmt.Printf("🔄 Generating questions and quiz for story: %s\n", story.Title)

	// First generate questions
	questionsResult, err := generateQuestions(story, apiKey)
	if err != nil {
		return types.GenerationResult{}, err
	}

	// Wait a moment to avoid rate limiting
	time.Sleep(1 * time.Second)

	// Then generate quiz
	quizResult, err := generateQuiz(story, apiKey)
	if err != nil {
		return types.GenerationResult{}, err
	}

	return types.GenerationResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully generated %d questions and quiz for story '%s'", questionsResult.QuestionsCount, story.Title),
		QuestionsCount: questionsResult.QuestionsCount,
		QuizID:         quizResult.QuizID,
		QuestionIDs:    questionsResult.QuestionIDs,
	}, nil
}

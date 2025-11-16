package handlers

import (
	"context"
	"fmt"
	"time"

	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/ai-generator/services"
	"egaldeutsch-serverless/netlify/functions/ai-generator/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// storeQuestions stores generated questions in the database
func storeQuestions(storyID primitive.ObjectID, questions []types.QuestionGenerated) ([]string, error) {
	_, questionsCollection, _, err := services.GetCollections()
	if err != nil {
		return nil, err
	}

	var questionIDs []string
	userID := getAIGeneratorUserID()
	now := time.Now()

	for i, q := range questions {
		question := models.Question{
			ID:            primitive.NewObjectID(),
			StoryID:       storyID,
			Question:      q.Question,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
			Explanation:   q.Explanation,
			QuestionType:  q.QuestionType,
			Difficulty:    q.Difficulty,
			Points:        q.Points,
			Order:         i + 1,
			ContentMetadata: models.ContentMetadata{
				Status:    models.StatusPublished,
				CreatedBy: userID,
				CreatedAt: now,
				UpdatedAt: now,
				Version:   1,
			},
		}

		if _, err := questionsCollection.InsertOne(context.Background(), question); err != nil {
			return nil, fmt.Errorf("failed to insert question %d: %w", i+1, err)
		}

		questionIDs = append(questionIDs, question.ID.Hex())
	}

	return questionIDs, nil
}

// getStoryQuestions retrieves all questions for a story
func getStoryQuestions(storyID primitive.ObjectID) ([]models.Question, error) {
	_, questionsCollection, _, err := services.GetCollections()
	if err != nil {
		return nil, err
	}

	filter := bson.M{"storyId": storyID, "status": models.StatusPublished}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "order", Value: 1}})

	cursor, err := questionsCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var questions []models.Question
	if err = cursor.All(context.Background(), &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

// storeQuiz stores generated quiz metadata in the database
func storeQuiz(storyID primitive.ObjectID, quizData types.QuizGenerationResponse, questions []models.Question) (string, error) {
	_, _, quizzesCollection, err := services.GetCollections()
	if err != nil {
		return "", err
	}

	// Extract question IDs
	var questionIDs []primitive.ObjectID
	for _, q := range questions {
		questionIDs = append(questionIDs, q.ID)
	}

	userID := getAIGeneratorUserID()
	now := time.Now()

	quiz := models.Quiz{
		ID:             primitive.NewObjectID(),
		StoryID:        storyID,
		Title:          quizData.Title,
		Description:    quizData.Description,
		QuestionIDs:    questionIDs,
		TotalQuestions: len(questionIDs),
		PassingScore:   quizData.PassingScore,
		QuizType:       quizData.QuizType,
		ContentMetadata: models.ContentMetadata{
			Status:    models.StatusPublished,
			CreatedBy: userID,
			CreatedAt: now,
			UpdatedAt: now,
			Version:   1,
		},
	}

	if _, err := quizzesCollection.InsertOne(context.Background(), quiz); err != nil {
		return "", fmt.Errorf("failed to insert quiz: %w", err)
	}

	fmt.Printf("✅ Stored quiz: %s\n", quiz.Title)
	return quiz.ID.Hex(), nil
}

// markStoryAIQuestionsGenerated marks a story as having AI questions generated
func markStoryAIQuestionsGenerated(storyID primitive.ObjectID) error {
	storiesCollection, _, _, err := services.GetCollections()
	if err != nil {
		return err
	}

	filter := bson.M{"_id": storyID}
	update := bson.M{
		"$set": bson.M{
			"isAIQuestionsGenerated": true,
			"updatedAt":              time.Now(),
		},
	}

	result, err := storiesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("story not found")
	}

	fmt.Printf("✅ Marked story as having AI questions generated\n")
	return nil
}

// getAIGeneratorUserID returns a consistent user ID for AI-generated content
func getAIGeneratorUserID() primitive.ObjectID {
	// Use a fixed ObjectID for AI-generated content
	// This represents the "AI Generator" system user
	aiUserIDHex := "507f1f77bcf86cd799439011"
	aiUserID, err := primitive.ObjectIDFromHex(aiUserIDHex)
	if err != nil {
		// Fallback to a new ObjectID if parsing fails
		return primitive.NewObjectID()
	}
	return aiUserID
}

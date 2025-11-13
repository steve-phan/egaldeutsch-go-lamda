package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer db.Disconnect()

	// Create sample stories and questions
	if err := seedData(); err != nil {
		log.Fatalf("Error seeding data: %v", err)
	}

	fmt.Println("Sample data seeded successfully!")
}

func seedData() error {
	ctx := context.Background()

	// Sample Story 1: A1 Level
	now := time.Now()
	adminID := primitive.NewObjectID() // Would be actual admin user ID in real system

	story1 := models.Story{
		ID:          primitive.NewObjectID(),
		Title:       "Meine Familie",
		Content:     "Hallo! Ich heiße Anna. Ich bin 25 Jahre alt. Ich komme aus Deutschland. Ich wohne in Berlin.\n\nMeine Familie ist klein. Mein Vater heißt Thomas. Er ist 50 Jahre alt und arbeitet als Lehrer. Meine Mutter heißt Maria. Sie ist 48 Jahre alt und ist Ärztin.\n\nIch habe einen Bruder. Er heißt Max und ist 23 Jahre alt. Max studiert Informatik an der Universität. Wir verstehen uns sehr gut.\n\nAm Wochenende essen wir zusammen. Wir kochen oft deutsche Gerichte. Meine Mutter macht die beste Schnitzel! Wir reden viel und lachen zusammen.\n\nMeine Familie ist sehr wichtig für mich. Wir helfen uns immer.",
		Level:       "A1",
		WordCount:   120,
		ReadingTime: 2,
		Topics:      []string{"Familie", "Alltag"},
		Summary:     "Anna stellt ihre kleine Familie vor: ihre Eltern Thomas und Maria, und ihren Bruder Max.",
		Vocabulary: []models.VocabularyWord{
			{German: "Familie", English: "family", WordType: "noun", Article: "die"},
			{German: "Vater", English: "father", WordType: "noun", Article: "der"},
			{German: "Mutter", English: "mother", WordType: "noun", Article: "die"},
			{German: "Bruder", English: "brother", WordType: "noun", Article: "der"},
			{German: "arbeiten", English: "to work", WordType: "verb"},
			{German: "studieren", English: "to study", WordType: "verb"},
			{German: "wichtig", English: "important", WordType: "adjective"},
		},
		ContentMetadata: models.ContentMetadata{
			Status:      models.StatusPublished,
			CreatedBy:   adminID,
			CreatedAt:   now,
			UpdatedAt:   now,
			ReviewedBy:  &adminID,
			ReviewedAt:  &now,
			ApprovedAt:  &now,
			ActivatedAt: &now,
			Version:     1,
		},
	}

	// Insert story
	storiesCollection := db.Database.Collection("stories")
	_, err := storiesCollection.InsertOne(ctx, story1)
	if err != nil {
		return fmt.Errorf("failed to insert story: %w", err)
	}

	// Sample Questions for Story 1
	questionMetadata := models.ContentMetadata{
		Status:      models.StatusPublished,
		CreatedBy:   adminID,
		CreatedAt:   now,
		UpdatedAt:   now,
		ReviewedBy:  &adminID,
		ReviewedAt:  &now,
		ApprovedAt:  &now,
		ActivatedAt: &now,
		Version:     1,
	}

	questions1 := []models.Question{
		{
			ID:              primitive.NewObjectID(),
			StoryID:         story1.ID,
			Question:        "Wie heißt die Hauptperson?",
			QuestionType:    "comprehension",
			Options:         []string{"Maria", "Anna", "Max", "Thomas"},
			CorrectAnswer:   1,
			Explanation:     "Im ersten Satz sagt sie: 'Ich heiße Anna.'",
			Points:          10,
			Order:           1,
			Difficulty:      "easy",
			ContentMetadata: questionMetadata,
		},
		{
			ID:              primitive.NewObjectID(),
			StoryID:         story1.ID,
			Question:        "Wie alt ist Anna?",
			QuestionType:    "comprehension",
			Options:         []string{"23", "25", "48", "50"},
			CorrectAnswer:   1,
			Explanation:     "Anna sagt: 'Ich bin 25 Jahre alt.'",
			Points:          10,
			Order:           2,
			Difficulty:      "easy",
			ContentMetadata: questionMetadata,
		},
		{
			ID:              primitive.NewObjectID(),
			StoryID:         story1.ID,
			Question:        "Was macht Annas Vater beruflich?",
			QuestionType:    "comprehension",
			Options:         []string{"Arzt", "Lehrer", "Student", "Informatiker"},
			CorrectAnswer:   1,
			Explanation:     "Der Text sagt: 'Mein Vater heißt Thomas. Er ist 50 Jahre alt und arbeitet als Lehrer.'",
			Points:          10,
			Order:           3,
			Difficulty:      "medium",
			ContentMetadata: questionMetadata,
		},
		{
			ID:              primitive.NewObjectID(),
			StoryID:         story1.ID,
			Question:        "Was studiert Max?",
			QuestionType:    "comprehension",
			Options:         []string{"Medizin", "Informatik", "Deutsch", "Pädagogik"},
			CorrectAnswer:   1,
			Explanation:     "Der Text erklärt: 'Max studiert Informatik an der Universität.'",
			Points:          10,
			Order:           4,
			Difficulty:      "medium",
			ContentMetadata: questionMetadata,
		},
		{
			ID:              primitive.NewObjectID(),
			StoryID:         story1.ID,
			Question:        "Was bedeutet 'Familie' auf Englisch?",
			QuestionType:    "vocabulary",
			Options:         []string{"friend", "family", "house", "work"},
			CorrectAnswer:   1,
			Explanation:     "'Familie' bedeutet 'family' auf Englisch.",
			Points:          10,
			Order:           5,
			Difficulty:      "easy",
			ContentMetadata: questionMetadata,
		},
	}

	// Insert questions
	questionsCollection := db.Database.Collection("questions")
	for _, question := range questions1 {
		_, err := questionsCollection.InsertOne(ctx, question)
		if err != nil {
			return fmt.Errorf("failed to insert question: %w", err)
		}
	}

	// Create a quiz for the story
	quiz1 := models.Quiz{
		ID:              primitive.NewObjectID(),
		StoryID:         story1.ID,
		Title:           "Quiz: Meine Familie",
		Description:     "Test your understanding of Anna's family story",
		QuestionIDs:     []primitive.ObjectID{},
		TotalQuestions:  5,
		PassingScore:    70,
		QuizType:        "auto_generated",
		ContentMetadata: questionMetadata,
	}

	// Get question IDs for the quiz
	for _, q := range questions1 {
		quiz1.QuestionIDs = append(quiz1.QuestionIDs, q.ID)
	}

	// Insert quiz
	quizCollection := db.Database.Collection("quizzes")
	_, err = quizCollection.InsertOne(ctx, quiz1)
	if err != nil {
		return fmt.Errorf("failed to insert quiz: %w", err)
	}

	return nil
}

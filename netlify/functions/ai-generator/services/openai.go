package services

import (
	"fmt"
	"os"

	"egaldeutsch-serverless/models"
)

// GetOpenAIAPIKey returns the OpenAI API key from environment
func GetOpenAIAPIKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

// BuildQuestionPrompt creates the prompt for question generation
func BuildQuestionPrompt(story *models.Story) string {
	vocabularyStr := FormatVocabulary(story.Vocabulary)

	return fmt.Sprintf(`You are a German language learning expert. Based on the following German story, generate diverse questions for German learners.

Story Details:
Title: %s
Level: %s
Content: %s

Key Vocabulary:
%s

Generate exactly 2-4 questions of different types. For each question, extract relevant German vocabulary that appears in the story.

Return your response as valid JSON in this exact format:
{
  "questions": [
    {
      "question": "What does the main character do in the story?",
      "questionType": "comprehension",
      "options": ["Option A", "Option B", "Option C", "Option D"],
      "correctAnswer": 0,
      "explanation": "Brief explanation of the correct answer",
      "difficulty": "easy",
      "points": 1,
      "germanConcept": "Relevant German grammar/vocabulary concept being tested"
    }
  ],
  "vocabularyExtracted": [
    {
      "german": "das Haus",
      "english": "house",
      "wordType": "noun",
      "article": "das",
      "context": "Context from the story where this word appears"
    }
  ]
}

Question types to include:
- comprehension: About story content and meaning
- vocabulary: About German words and their meanings
- grammar: About German grammar rules, cases, verb conjugations
- culture: About German culture references in the story

Difficulty levels: easy (basic vocabulary/grammar), medium (intermediate concepts), hard (complex grammar/cultural nuances)
Points: 1 for easy, 2 for medium, 3 for hard questions`,
		story.Title, story.Level, story.Content, vocabularyStr)
}

// BuildQuizPrompt creates the prompt for quiz metadata generation
func BuildQuizPrompt(story *models.Story) string {
	return fmt.Sprintf(`You are a German language learning expert. Based on the following German story, generate quiz metadata.

Story Details:
Title: %s
Level: %s
Content: %s

Generate quiz metadata that would be appropriate for this story level and content.

Return your response as valid JSON in this exact format:
{
  "title": "Quiz title based on the story",
  "description": "Brief description of what this quiz tests",
  "quizType": "comprehension_and_vocabulary",
  "estimatedTime": 15,
  "passingScore": 70,
  "recommendedQuestionTypes": ["comprehension", "vocabulary", "grammar"],
  "difficultyDistribution": {
    "easy": 4,
    "medium": 4,
    "hard": 2
  }
}

Guidelines:
- Title should be engaging and relevant to the story
- Description should explain what skills are being tested
- Estimated time in minutes (typically 10-20 minutes)
- Passing score as percentage (typically 60-80%%)
- Difficulty distribution should total 10 questions`,
		story.Title, story.Level, story.Content)
}

// FormatVocabulary formats vocabulary words for the prompt
func FormatVocabulary(vocabulary []models.VocabularyWord) string {
	if len(vocabulary) == 0 {
		return "No vocabulary provided"
	}

	result := ""
	for _, word := range vocabulary {
		if word.Article != "" {
			result += fmt.Sprintf("- %s %s (%s) - %s\n", word.Article, word.German, word.WordType, word.English)
		} else {
			result += fmt.Sprintf("- %s (%s) - %s\n", word.German, word.WordType, word.English)
		}
	}
	return result
}

package main

import (
	"context"
	"log"
	"strconv"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	log.Println("Starting slug migration for existing stories...")

	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Disconnect()

	log.Println("Connected to database")

	// Get stories collection
	collection := db.Database.Collection("stories")

	// Find all stories without a slug or with empty slug
	filter := bson.M{
		"$or": []bson.M{
			{"slug": bson.M{"$exists": false}},
			{"slug": ""},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to fetch stories: %v", err)
	}
	defer cursor.Close(context.Background())

	var stories []models.Story
	if err = cursor.All(context.Background(), &stories); err != nil {
		log.Fatalf("Failed to decode stories: %v", err)
	}

	log.Printf("Found %d stories without slugs", len(stories))

	if len(stories) == 0 {
		log.Println("No stories need slug migration. Exiting.")
		return
	}

	// Process each story
	updatedCount := 0
	failedCount := 0

	for i, story := range stories {
		log.Printf("Processing story %d/%d: %s", i+1, len(stories), story.Title)

		// Generate base slug
		baseSlug := models.GenerateSlug(story.Title)
		slug := baseSlug
		counter := 1

		// Check for uniqueness
		for {
			var existingStory models.Story
			err := collection.FindOne(context.Background(), bson.M{
				"slug": slug,
				"_id":  bson.M{"$ne": story.ID},
			}).Decode(&existingStory)

			if err == mongo.ErrNoDocuments {
				// Slug is unique, we can use it
				break
			} else if err != nil {
				log.Printf("  Error checking slug uniqueness: %v", err)
				failedCount++
				continue
			}

			// Slug already exists, try with number suffix
			slug = baseSlug + "-" + strconv.Itoa(counter)
			counter++
		}

		// Update story with slug
		result, err := collection.UpdateOne(
			context.Background(),
			bson.M{"_id": story.ID},
			bson.M{"$set": bson.M{"slug": slug}},
		)

		if err != nil {
			log.Printf("  ❌ Failed to update story: %v", err)
			failedCount++
		} else if result.ModifiedCount > 0 {
			log.Printf("  ✅ Updated with slug: %s", slug)
			updatedCount++
		} else {
			log.Printf("  ⚠️  No changes made")
		}
	}

	// log.Println("\n" + "=".repeat(50))
	log.Printf("Migration complete!")
	log.Printf("Total stories: %d", len(stories))
	log.Printf("Successfully updated: %d", updatedCount)
	log.Printf("Failed: %d", failedCount)
	// log.Println("=".repeat(50))
}

// repeat is a helper function to repeat a string n times
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

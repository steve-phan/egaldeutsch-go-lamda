package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Set CORS headers
	headers := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Methods":     "GET, OPTIONS",
		"Access-Control-Allow-Credentials": "true",
	}

	// Handle OPTIONS request
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
		}, nil
	}

	// Only allow GET requests
	if req.HTTPMethod != "GET" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    headers,
			Body:       `{"success": false, "error": "Method not allowed"}`,
		}, nil
	}

	return getLeaderboard(ctx, req, headers)
}

func getLeaderboard(ctx context.Context, req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	submissionsCollection, err := db.GetCollection(db.Collections.Submissions)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Database connection failed"}`,
		}, nil
	}

	// Parse query parameters
	limit := 100 // Default limit
	if limitParam, ok := req.QueryStringParameters["limit"]; ok {
		fmt.Sscanf(limitParam, "%d", &limit)
		if limit > 500 {
			limit = 500 // Max limit
		}
		if limit < 1 {
			limit = 10
		}
	}

	// Aggregate leaderboard data
	// Group by userId, calculate total points, quizzes taken, passed, and average score
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"userId": bson.M{"$exists": true},
			},
		},
		{
			"$group": bson.M{
				"_id": "$userId",
				"totalPoints": bson.M{
					"$sum": "$earnedPoints",
				},
				"quizzesTaken": bson.M{
					"$sum": 1,
				},
				"quizzesPassed": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							"$passed",
							1,
							0,
						},
					},
				},
				"averageScore": bson.M{
					"$avg": "$percentage",
				},
				"lastSubmittedAt": bson.M{
					"$max": "$submittedAt",
				},
				"username": bson.M{
					"$first": "$username",
				},
			},
		},
		{
			"$sort": bson.M{
				"totalPoints":     -1, // Sort by total points descending
				"lastSubmittedAt": -1, // Then by most recent activity
			},
		},
		{
			"$limit": int64(limit),
		},
		{
			"$project": bson.M{
				"_id":             0,
				"userId":          "$_id",
				"username":        1,
				"totalPoints":     1,
				"quizzesTaken":    1,
				"quizzesPassed":   1,
				"averageScore":    1,
				"lastSubmittedAt": 1,
			},
		},
	}

	cursor, err := submissionsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"success": false, "error": "Failed to aggregate leaderboard: %v"}`, err),
		}, nil
	}
	defer cursor.Close(ctx)

	var entries []models.LeaderboardEntry
	if err = cursor.All(ctx, &entries); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to decode leaderboard"}`,
		}, nil
	}

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	// Get total number of participants for statistics
	totalParticipants, err := submissionsCollection.Distinct(ctx, "userId", bson.M{
		"userId": bson.M{"$exists": true},
	})
	if err != nil {
		totalParticipants = []interface{}{} // Default to empty
	}

	// Calculate total quizzes submitted
	totalQuizzes, err := submissionsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		totalQuizzes = 0
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"entries":           entries,
			"totalParticipants": len(totalParticipants),
			"totalQuizzes":      totalQuizzes,
			"generatedAt":       time.Now(),
		},
		"message": fmt.Sprintf("Leaderboard retrieved with %d entries", len(entries)),
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"success": false, "error": "Failed to serialize response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonData),
	}, nil
}

func main() {
	lambda.Start(handler)
}

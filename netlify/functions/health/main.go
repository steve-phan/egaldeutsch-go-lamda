package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/pkg/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := middleware.GetPublicCORSHeaders()

	// Handle CORS preflight
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	if req.HTTPMethod != "GET" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    headers,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}

	healthResponse := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services:  make(map[string]string),
		Version:   "1.0.0",
	}

	// Check MongoDB connection
	mongoStatus := checkMongoDB(ctx)
	healthResponse.Services["mongodb"] = mongoStatus

	// Determine overall status
	if mongoStatus != "healthy" {
		healthResponse.Status = "unhealthy"
	}

	// Set appropriate HTTP status code
	statusCode := http.StatusOK
	if healthResponse.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	jsonData, err := json.Marshal(healthResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"error": "Failed to marshal response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(jsonData),
	}, nil
}

func checkMongoDB(ctx context.Context) string {
	// Create a timeout context for the health check
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Ensure we have a healthy connection
	if err := db.EnsureConnection(); err != nil {
		return "connection_failed"
	}

	// Test database health
	if err := db.IsHealthy(healthCtx); err != nil {
		return "ping_failed"
	}

	// Get database instance
	database, err := db.GetDatabase()
	if err != nil {
		return "database_not_initialized"
	}

	// Test a simple operation - list collections
	_, err = database.ListCollectionNames(healthCtx, map[string]interface{}{})
	if err != nil {
		return "database_access_failed"
	}

	return "healthy"
}

func main() {
	lambda.Start(handler)
}

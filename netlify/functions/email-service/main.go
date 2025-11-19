package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/pkg/email"
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// init loads environment variables from .env file in development
func init() {
	log.Printf("Initializing email service function")

	// Try multiple paths for .env file - Netlify dev changes working directory
	paths := []string{
		"../../../.env", // From function directory to project root
		"../../.env",    // Alternative path
		".env",          // Current directory
		"./.env",        // Explicit current directory
	}

	envLoaded := false
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Successfully loaded .env file from: %s", path)
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		log.Printf("No .env file loaded (this is normal in production)")
	}

	// Log environment variables (without exposing sensitive data)
	emailProvider := os.Getenv("EMAIL_PROVIDER")
	emailAPIKey := os.Getenv("EMAIL_API_KEY")
	emailFrom := os.Getenv("EMAIL_FROM")
	emailFromName := os.Getenv("EMAIL_FROM_NAME")

	log.Printf("Environment check - EMAIL_PROVIDER: %s", emailProvider)
	log.Printf("Environment check - EMAIL_API_KEY present: %t", emailAPIKey != "")
	log.Printf("Environment check - EMAIL_FROM: %s", emailFrom)
	log.Printf("Environment check - EMAIL_FROM_NAME: %s", emailFromName)
}

// EmailRequest represents email sending request
type EmailRequest struct {
	Type     string            `json:"type"`               // "welcome", "password_reset", "new_story"
	To       []string          `json:"to"`                 // Recipients
	Data     map[string]string `json:"data"`               // Template data
	Subject  string            `json:"subject,omitempty"`  // Optional custom subject
	Priority string            `json:"priority,omitempty"` // "high", "normal", "low"
}

// EmailResponse represents email sending response
type EmailResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	SentCount int    `json:"sentCount,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Email service handler called - Method: %s, Path: %s", request.HTTPMethod, request.Path)

	// Handle CORS preflight requests
	if corsResponse, handled := middleware.HandleAuthenticatedCORS(request); handled {
		log.Printf("CORS preflight handled")
		return corsResponse, nil
	}

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "POST":
		// Check specific routes first, then fallback to generic ones
		if strings.Contains(request.Path, "/send-password-reset") {
			log.Printf("Routing to password reset email handler")
			return sendPasswordResetEmail(request)
		}
		if strings.Contains(request.Path, "/send-welcome") {
			return sendWelcomeEmail(request)
		}
		if strings.Contains(request.Path, "/send-new-story-notification") {
			return sendNewStoryNotificationWithDB(request)
		}
		if strings.Contains(request.Path, "/send") {
			return sendEmail(request)
		}
		return response.SimpleErrorWithDefault(404, "Endpoint not found"), nil
	case "GET":
		if strings.Contains(request.Path, "/config") {
			return getEmailConfig(request)
		}
		return response.SimpleErrorWithDefault(404, "Endpoint not found"), nil
	default:
		return response.SimpleErrorWithDefault(405, "Method not allowed"), nil
	}
}

// sendEmail handles generic email sending
func sendEmail(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var emailReq EmailRequest
	if err := json.Unmarshal([]byte(request.Body), &emailReq); err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid request body"), nil
	}

	// Initialize email service
	emailService, err := email.NewServiceFromEnv()
	if err != nil {
		log.Printf("Failed to initialize email service: %v", err)
		return response.SimpleErrorWithDefault(500, "Email service initialization failed"), nil
	}

	var sentCount int
	var emailErr error

	switch emailReq.Type {
	case "welcome":
		for _, recipient := range emailReq.To {
			userName := emailReq.Data["userName"]
			if userName == "" {
				userName = "New User"
			}
			if err := emailService.SendWelcomeEmail(recipient, userName); err != nil {
				emailErr = err
				break
			}
			sentCount++
		}
	case "password_reset":
		for _, recipient := range emailReq.To {
			userName := emailReq.Data["userName"]
			resetToken := emailReq.Data["resetToken"]
			if userName == "" || resetToken == "" {
				emailErr = fmt.Errorf("userName and resetToken are required for password reset emails")
				break
			}
			if err := emailService.SendPasswordResetEmail(recipient, userName, resetToken); err != nil {
				emailErr = err
				break
			}
			sentCount++
		}
	case "new_story":
		storyTitle := emailReq.Data["storyTitle"]
		storyLevel := emailReq.Data["storyLevel"]
		storyID := emailReq.Data["storyID"]
		if storyTitle == "" || storyLevel == "" || storyID == "" {
			emailErr = fmt.Errorf("storyTitle, storyLevel, and storyID are required for new story emails")
		} else {
			if err := emailService.SendNewStoryNotification(emailReq.To, storyTitle, storyLevel, storyID); err != nil {
				emailErr = err
			} else {
				sentCount = len(emailReq.To)
			}
		}
	default:
		return response.SimpleErrorWithDefault(400, "Unsupported email type"), nil
	}

	if emailErr != nil {
		log.Printf("Email sending failed: %v", emailErr)
		return response.SimpleErrorWithDefault(500, fmt.Sprintf("Failed to send email: %v", emailErr)), nil
	}

	responseData := EmailResponse{
		Success:   true,
		Message:   "Email sent successfully",
		SentCount: sentCount,
		Provider:  emailService.GetConfig().Provider,
	}

	return response.JSONWithDefault(200, responseData)
}

// sendWelcomeEmail handles welcome email sending
func sendWelcomeEmail(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req struct {
		Email    string `json:"email"`
		UserName string `json:"userName"`
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid request body"), nil
	}

	emailService, err := email.NewServiceFromEnv()
	if err != nil {
		log.Printf("Failed to initialize email service: %v", err)
		return response.SimpleErrorWithDefault(500, "Email service initialization failed"), nil
	}

	if err := emailService.SendWelcomeEmail(req.Email, req.UserName); err != nil {
		log.Printf("Failed to send welcome email: %v", err)
		return response.SimpleErrorWithDefault(500, "Failed to send welcome email"), nil
	}

	responseData := EmailResponse{
		Success:   true,
		Message:   "Welcome email sent successfully",
		SentCount: 1,
		Provider:  emailService.GetConfig().Provider,
	}

	return response.JSONWithDefault(200, responseData)
}

// sendPasswordResetEmail handles password reset email sending
func sendPasswordResetEmail(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Password reset email request received")

	var req struct {
		Email      string `json:"email"`
		UserName   string `json:"userName"`
		ResetToken string `json:"resetToken"`
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Printf("Failed to unmarshal password reset request: %v", err)
		return response.SimpleErrorWithDefault(400, "Invalid request body"), nil
	}

	log.Printf("Password reset request parsed - Email: %s, UserName: %s", req.Email, req.UserName)

	log.Printf("Initializing email service...")
	emailService, err := email.NewServiceFromEnv()
	if err != nil {
		log.Printf("Failed to initialize email service: %v", err)
		return response.SimpleErrorWithDefault(500, "Email service initialization failed"), nil
	}
	log.Printf("Email service initialized successfully")

	log.Printf("Sending password reset email to: %s", req.Email)
	if err := emailService.SendPasswordResetEmail(req.Email, req.UserName, req.ResetToken); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return response.SimpleErrorWithDefault(500, "Failed to send password reset email"), nil
	}
	log.Printf("Password reset email sent successfully to: %s", req.Email)

	responseData := EmailResponse{
		Success:   true,
		Message:   "Password reset email sent successfully",
		SentCount: 1,
		Provider:  emailService.GetConfig().Provider,
	}

	return response.JSONWithDefault(200, responseData)
}

// sendNewStoryNotification handles new story notification emails
func sendNewStoryNotification(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req struct {
		StoryID    string `json:"storyId"`
		StoryTitle string `json:"storyTitle"`
		StoryLevel string `json:"storyLevel"`
		MaxUsers   int    `json:"maxUsers,omitempty"` // Optional limit
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.SimpleErrorWithDefault(400, "Invalid request body"), nil
	}

	// Get subscriber emails from database
	subscriberEmails, err := getSubscriberEmails(req.MaxUsers)
	if err != nil {
		log.Printf("Failed to get subscriber emails: %v", err)
		return response.SimpleErrorWithDefault(500, "Failed to get subscriber emails"), nil
	}

	if len(subscriberEmails) == 0 {
		return response.JSONWithDefault(200, EmailResponse{
			Success:   true,
			Message:   "No subscribers to notify",
			SentCount: 0,
		})
	}

	emailService, err := email.NewServiceFromEnv()
	if err != nil {
		log.Printf("Failed to initialize email service: %v", err)
		return response.SimpleErrorWithDefault(500, "Email service initialization failed"), nil
	}

	if err := emailService.SendNewStoryNotification(subscriberEmails, req.StoryTitle, req.StoryLevel, req.StoryID); err != nil {
		log.Printf("Failed to send new story notification: %v", err)
		return response.SimpleErrorWithDefault(500, "Failed to send new story notification"), nil
	}

	responseData := EmailResponse{
		Success:   true,
		Message:   "New story notification sent successfully",
		SentCount: len(subscriberEmails),
		Provider:  emailService.GetConfig().Provider,
	}

	return response.JSONWithDefault(200, responseData)
}

// sendNewStoryNotificationWithDB handles new story notification with database connection
func sendNewStoryNotificationWithDB(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB only for this function that needs it
	if err := db.Connect(); err != nil {
		log.Printf("Database connection failed: %v", err)
		return response.SimpleErrorWithDefault(500, "Database connection failed"), nil
	}
	defer db.Disconnect()

	// Call the original function
	return sendNewStoryNotification(request)
}

// getEmailConfig returns masked email configuration
func getEmailConfig(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	emailService, err := email.NewServiceFromEnv()
	if err != nil {
		return response.SimpleErrorWithDefault(500, "Email service initialization failed"), nil
	}

	config := emailService.GetConfig()
	return response.JSONWithDefault(200, config)
}

// getSubscriberEmails retrieves subscriber emails from the database
func getSubscriberEmails(maxUsers int) ([]string, error) {
	collection := db.Database.Collection("users")

	// Set default limit if not specified
	if maxUsers == 0 {
		maxUsers = 10 // Default for trial accounts
	}

	// Find active users who want to receive notifications
	// You can add a field like "emailNotifications: true" to User model
	filter := bson.M{
		"role": bson.M{"$in": []string{"user", "creator"}}, // Active users
		// Add email notification preference when implemented
		// "emailNotifications": true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	// Extract emails and apply limit
	var emails []string
	for i, user := range users {
		if i >= maxUsers {
			break
		}
		if user.Email != "" {
			emails = append(emails, user.Email)
		}
	}

	return emails, nil
}

func main() {
	lambda.Start(handler)
}

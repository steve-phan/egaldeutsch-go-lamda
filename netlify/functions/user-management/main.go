package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/pkg/auth"
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/notification"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserRegistrationRequest represents user registration data
type UserRegistrationRequest struct {
	Username      string `json:"username" validate:"required,min=3,max=50"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8,max=128"`
	FirstName     string `json:"firstName" validate:"required,min=2,max=50"`
	LastName      string `json:"lastName" validate:"required,min=2,max=50"`
	PreferredRole string `json:"preferredRole" validate:"required,oneof=creator reviewer"`
}

// UserLoginRequest represents login credentials
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents user data in responses (without sensitive info)
type UserResponse struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// AuthTokenResponse represents authentication response
type AuthTokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expiresAt"`
	User      UserResponse `json:"user"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	if err := db.Connect(); err != nil {
		return response.SimpleError(500, "Database connection failed"), nil
	}
	defer db.Disconnect()

	// Handle CORS preflight requests
	if corsResponse, handled := middleware.HandleCORS(request); handled {
		return corsResponse, nil
	}

	// Route based on HTTP method and path
	switch request.HTTPMethod {
	case "POST":
		if strings.Contains(request.Path, "/register") {
			return registerUser(request)
		}
		if strings.Contains(request.Path, "/login") {
			return loginUser(request)
		}
		if strings.Contains(request.Path, "/forgot-password") {
			return forgotPassword(request)
		}
		if strings.Contains(request.Path, "/reset-password") {
			return resetPassword(request)
		}
		return response.SimpleError(404, "Endpoint not found"), nil
	case "GET":
		if strings.Contains(request.Path, "/profile") {
			return getUserProfile(request)
		}
		return listUsers(request)
	case "PUT":
		return updateUserProfile(request)
	case "DELETE":
		if strings.Contains(request.Path, "/logout") {
			return logoutUser(request)
		}
		return deleteUser(request)
	default:
		return response.SimpleError(405, "Method not allowed"), nil
	}
}

// registerUser handles user registration
func registerUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var regReq UserRegistrationRequest
	if err := json.Unmarshal([]byte(request.Body), &regReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid request format"}`,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Check if username or email already exists
	collection := db.Database.Collection("users")
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"username": regReq.Username},
			{"email": regReq.Email},
		},
	}).Decode(&existingUser)

	if err == nil {
		return response.SimpleError(409, "Username or email already exists"), nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.SimpleError(500, "Failed to hash password"), nil
	}

	// Create new user
	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Username:     regReq.Username,
		Email:        regReq.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    regReq.FirstName,
		LastName:     regReq.LastName,
		Name:         regReq.FirstName + " " + regReq.LastName,
		Role:         models.RoleCreator, // Default role, admin can upgrade
		Status:       models.UserStatusActive,
		IsActive:     true, // Backward compatibility
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user to database
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return response.SimpleError(500, "Failed to create user"), nil
	}

	// Send welcome email (non-blocking, log error if it fails)
	go func() {
		if err := sendWelcomeEmailAsync(newUser.Email, newUser.Name); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", newUser.Email, err)
		}
	}()

	// Create notification for user registration
	if err := notification.CreateNotification(
		newUser.ID,
		notification.NotificationTypeUserRegistered,
		"Welcome to EgalDeutsch!",
		"Your account has been successfully created. Start learning German today!",
		"/stories",
		nil,
	); err != nil {
		log.Printf("Failed to create registration notification: %v", err)
	}

	// Return user response (without password)
	userResponse := UserResponse{
		ID:          newUser.ID.Hex(),
		Username:    newUser.Username,
		Email:       newUser.Email,
		FirstName:   newUser.FirstName,
		LastName:    newUser.LastName,
		Role:        string(newUser.Role),
		Status:      string(newUser.Status),
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		LastLoginAt: newUser.LastLoginAt,
	}

	return response.JSON(201, userResponse)
}

// loginUser handles user authentication
func loginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginReq UserLoginRequest
	if err := json.Unmarshal([]byte(request.Body), &loginReq); err != nil {
		return response.SimpleError(400, "Invalid request format"), nil
	}

	// Find user by username or email
	collection := db.Database.Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"username": loginReq.Username},
			{"email": loginReq.Username},
		},
	}).Decode(&user)

	if err != nil {
		return response.SimpleError(401, "Invalid credentials"), nil
	}

	// Check if user is active
	if user.Status != models.UserStatusActive {
		return response.SimpleError(403, "Account is not active"), nil
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		return response.SimpleError(401, "Invalid credentials"), nil
	}

	// Create session using auth package
	session, token, err := auth.CreateSession(user.ID, request)
	if err != nil {
		return response.SimpleError(500, "Failed to create session"), nil
	}

	// Update user's last login
	now := time.Now()
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		}},
	)
	if err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.ID.Hex(), err)
	}
	user.LastLoginAt = &now

	// Return authentication response
	userResponse := UserResponse{
		ID:          user.ID.Hex(),
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Role:        string(user.Role),
		Status:      string(user.Status),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	authResponse := AuthTokenResponse{
		Token:     token,
		ExpiresAt: session.ExpiresAt,
		User:      userResponse,
	}

	return response.JSON(200, authResponse)
}

// forgotPassword handles password reset request
func forgotPassword(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req ForgotPasswordRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.SimpleError(400, "Invalid request format"), nil
	}

	// Find user by email
	collection := db.Database.Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)

	// Always return success to prevent email enumeration
	if err != nil {
		return response.SuccessJSON(200, nil, "If the email exists, a password reset link has been sent")
	}

	// Check if user is active
	if user.Status != models.UserStatusActive {
		return response.SuccessJSON(200, nil, "If the email exists, a password reset link has been sent")
	}

	// Generate password reset token
	token, err := createPasswordResetToken(user.ID)
	if err != nil {
		log.Printf("Failed to create password reset token for user %s: %v", user.ID.Hex(), err)
		return response.SuccessJSON(200, nil, "If the email exists, a password reset link has been sent")
	}

	// Send password reset email (non-blocking)
	go func() {
		if err := sendPasswordResetEmailAsync(user.Email, user.Name, token); err != nil {
			log.Printf("Failed to send password reset email to %s: %v", user.Email, err)
		}
	}()

	return response.SuccessJSON(200, nil, "If the email exists, a password reset link has been sent")
}

// resetPassword handles password reset with token
func resetPassword(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req ResetPasswordRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.SimpleError(400, "Invalid request format"), nil
	}

	// Validate password length
	if len(req.NewPassword) < 8 {
		return response.SimpleError(400, "Password must be at least 8 characters long"), nil
	}

	// Validate reset token
	userID, err := validateResetToken(req.Token)
	if err != nil {
		return response.SimpleError(400, "Invalid or expired reset token"), nil
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return response.SimpleError(500, "Failed to hash password"), nil
	}

	// Update user password
	collection := db.Database.Collection("users")
	var user models.User
	err = collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"passwordHash": string(hashedPassword),
			"updatedAt":    time.Now(),
		}},
	).Decode(&user)

	if err != nil {
		return response.SimpleError(500, "Failed to update password"), nil
	}

	// Mark token as used
	if err := markTokenAsUsed(req.Token); err != nil {
		log.Printf("Failed to mark reset token as used: %v", err)
	}

	// Invalidate all existing sessions for this user
	sessionCollection := db.Database.Collection("sessions")
	_, err = sessionCollection.DeleteMany(context.TODO(), bson.M{"userId": userID})
	if err != nil {
		log.Printf("Failed to delete user sessions after password reset: %v", err)
	}

	// Send password changed email (we can skip this for now since it's not in our requirements)
	// TODO: Add password changed email template if needed

	// Create notification for password change
	if err := notification.CreateNotification(
		userID,
		notification.NotificationTypePasswordChanged,
		"Password Changed",
		"Your password has been successfully changed.",
		"",
		nil,
	); err != nil {
		log.Printf("Failed to create password change notification: %v", err)
	}

	return response.SuccessJSON(200, nil, "Password has been successfully reset")
}

// getUserProfile returns the current user's profile
func getUserProfile(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session using auth package
	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Return user profile
	userResponse := UserResponse{
		ID:          user.ID.Hex(),
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Role:        string(user.Role),
		Status:      string(user.Status),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	return response.JSON(200, userResponse)
}

func listUsers(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session and check admin/reviewer role using auth package
	_, err := auth.ValidateSessionWithRole(request, models.RoleAdmin, models.RoleReviewer)
	if err != nil {
		if err.Error() == "insufficient permissions" {
			return response.SimpleError(403, "Insufficient permissions"), nil
		}
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Get query parameters for filtering
	queryParams := request.QueryStringParameters
	filter := bson.M{}

	if role := queryParams["role"]; role != "" {
		filter["role"] = role
	}
	if status := queryParams["status"]; status != "" {
		filter["status"] = status
	}

	// Find users
	collection := db.Database.Collection("users")
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return response.SimpleError(500, "Failed to fetch users"), nil
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return response.SimpleError(500, "Failed to decode users"), nil
	}

	// Convert to response format (without passwords)
	var userResponses []UserResponse
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:          u.ID.Hex(),
			Username:    u.Username,
			Email:       u.Email,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Role:        string(u.Role),
			Status:      string(u.Status),
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
			LastLoginAt: u.LastLoginAt,
		})
	}

	return response.JSON(200, userResponses)
}

// UserUpdateRequest represents user profile update data
type UserUpdateRequest struct {
	FirstName   *string `json:"firstName,omitempty"`
	LastName    *string `json:"lastName,omitempty"`
	Email       *string `json:"email,omitempty"`
	OldPassword *string `json:"oldPassword,omitempty"`
	NewPassword *string `json:"newPassword,omitempty"`
}

func updateUserProfile(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session using auth package
	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	var updateReq UserUpdateRequest
	if err := json.Unmarshal([]byte(request.Body), &updateReq); err != nil {
		return response.SimpleError(400, "Invalid request format"), nil
	}

	collection := db.Database.Collection("users")
	updateFields := bson.M{"updatedAt": time.Now()}

	// Update basic profile fields
	if updateReq.FirstName != nil {
		updateFields["firstName"] = *updateReq.FirstName
	}
	if updateReq.LastName != nil {
		updateFields["lastName"] = *updateReq.LastName
	}

	// Update name field when first or last name changes
	if updateReq.FirstName != nil || updateReq.LastName != nil {
		firstName := user.FirstName
		lastName := user.LastName
		if updateReq.FirstName != nil {
			firstName = *updateReq.FirstName
		}
		if updateReq.LastName != nil {
			lastName = *updateReq.LastName
		}
		updateFields["name"] = firstName + " " + lastName
	}

	// Handle email update (check for uniqueness)
	if updateReq.Email != nil && *updateReq.Email != user.Email {
		err := collection.FindOne(context.TODO(), bson.M{"email": *updateReq.Email}).Err()
		if err == nil {
			return response.SimpleError(409, "Email already in use"), nil
		}
		updateFields["email"] = *updateReq.Email
	}

	// Handle password update
	if updateReq.OldPassword != nil && updateReq.NewPassword != nil {
		// Verify old password
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*updateReq.OldPassword))
		if err != nil {
			return response.SimpleError(400, "Current password is incorrect"), nil
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateReq.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return response.SimpleError(500, "Failed to hash new password"), nil
		}
		updateFields["passwordHash"] = string(hashedPassword)
	}

	// Update user in database
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return response.SimpleError(500, "Failed to update user"), nil
	}

	// Fetch updated user
	var updatedUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": user.ID}).Decode(&updatedUser)
	if err != nil {
		return response.SimpleError(500, "Failed to fetch updated user"), nil
	}

	// Return updated user profile
	userResponse := UserResponse{
		ID:          updatedUser.ID.Hex(),
		Username:    updatedUser.Username,
		Email:       updatedUser.Email,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		Role:        string(updatedUser.Role),
		Status:      string(updatedUser.Status),
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		LastLoginAt: updatedUser.LastLoginAt,
	}

	return response.JSON(200, userResponse)
}

func logoutUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract token using auth package
	token, err := auth.ExtractToken(request)
	if err != nil {
		return response.SimpleError(400, "No valid authorization token"), nil
	}

	// Delete session using auth package
	if err := auth.DeleteSession(token); err != nil {
		log.Printf("Failed to delete session: %v", err)
	}

	return response.SuccessJSON(200, nil, "Logged out successfully")
}

func deleteUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session and check admin role using auth package
	user, err := auth.ValidateSessionWithRole(request, models.RoleAdmin)
	if err != nil {
		if err.Error() == "insufficient permissions" {
			return response.SimpleError(403, "Insufficient permissions"), nil
		}
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Get user ID from path parameters
	userID := request.PathParameters["id"]
	if userID == "" {
		return response.SimpleError(400, "User ID is required"), nil
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return response.SimpleError(400, "Invalid user ID format"), nil
	}

	// Prevent admin from deleting themselves
	if objectID == user.ID {
		return response.SimpleError(400, "Cannot delete your own account"), nil
	}

	// Check if user exists
	collection := db.Database.Collection("users")
	var targetUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&targetUser)
	if err != nil {
		return response.SimpleError(404, "User not found"), nil
	}

	// Delete user's sessions first
	sessionCollection := db.Database.Collection("sessions")
	_, err = sessionCollection.DeleteMany(context.TODO(), bson.M{"userId": objectID})
	if err != nil {
		log.Printf("Failed to delete user sessions: %v", err)
	}

	// Instead of hard delete, we'll soft delete by setting status to suspended
	// This preserves content relationships and audit trail
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"status":    models.UserStatusSuspended,
			"isActive":  false,
			"updatedAt": time.Now(),
		}},
	)
	if err != nil {
		return response.SimpleError(500, "Failed to delete user"), nil
	}

	return response.SuccessJSON(200, nil, "User deleted successfully")
}

// Email helper functions

// sendWelcomeEmailAsync sends a welcome email via the email service
func sendWelcomeEmailAsync(email, userName string) error {
	emailServiceURL := getEmailServiceURL()

	payload := map[string]interface{}{
		"email":    email,
		"userName": userName,
	}

	return callEmailService(emailServiceURL+"/send-welcome", payload)
}

// sendPasswordResetEmailAsync sends a password reset email via the email service
func sendPasswordResetEmailAsync(email, userName, resetToken string) error {
	emailServiceURL := getEmailServiceURL()

	payload := map[string]interface{}{
		"email":      email,
		"userName":   userName,
		"resetToken": resetToken,
	}

	return callEmailService(emailServiceURL+"/send-password-reset", payload)
}

// callEmailService makes an HTTP call to the email service
func callEmailService(url string, payload map[string]interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to call email service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("email service returned status %d", resp.StatusCode)
	}

	return nil
}

// getEmailServiceURL returns the email service URL
func getEmailServiceURL() string {
	if baseURL := os.Getenv("NETLIFY_FUNCTIONS_URL"); baseURL != "" {
		return baseURL + "/email-service"
	}
	// Default for local development
	return "http://localhost:8888/.netlify/functions/email-service"
}

// Password reset token functions (simplified implementation)
type PasswordResetToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId"`
	Token     string             `bson:"token"`
	CreatedAt time.Time          `bson:"createdAt"`
	ExpiresAt time.Time          `bson:"expiresAt"`
	Used      bool               `bson:"used"`
}

// createPasswordResetToken creates a new password reset token
func createPasswordResetToken(userID primitive.ObjectID) (string, error) {
	collection := db.Database.Collection("password_reset_tokens")

	// Generate token (simple implementation)
	token := fmt.Sprintf("%d_%s", time.Now().Unix(), userID.Hex())

	resetToken := PasswordResetToken{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		Used:      false,
	}

	_, err := collection.InsertOne(context.TODO(), resetToken)
	if err != nil {
		return "", fmt.Errorf("failed to create password reset token: %w", err)
	}

	return token, nil
}

// validateResetToken validates and returns the user ID for a reset token
func validateResetToken(token string) (primitive.ObjectID, error) {
	collection := db.Database.Collection("password_reset_tokens")

	var resetToken PasswordResetToken
	err := collection.FindOne(context.TODO(), bson.M{
		"token":     token,
		"used":      false,
		"expiresAt": bson.M{"$gt": time.Now()},
	}).Decode(&resetToken)

	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid or expired token")
	}

	return resetToken.UserID, nil
}

// markTokenAsUsed marks a reset token as used
func markTokenAsUsed(token string) error {
	collection := db.Database.Collection("password_reset_tokens")

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"token": token},
		bson.M{"$set": bson.M{"used": true}},
	)

	return err
}

func main() {
	lambda.Start(handler)
}

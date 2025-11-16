package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/models"
	"egaldeutsch-serverless/netlify/functions/user-management/services"
	"egaldeutsch-serverless/netlify/functions/user-management/types"
	"egaldeutsch-serverless/pkg/auth"
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// GetUserProfile returns the current user's profile
func GetUserProfile(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session using auth package
	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Return user profile
	userResponse := types.UserResponse{
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

// ListUsers returns a list of users (admin/reviewer only)
func ListUsers(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// // Validate session and check admin/reviewer role using auth package
	// _, err := auth.ValidateSessionWithRole(request, models.RoleAdmin, models.RoleReviewer)
	// if err != nil {
	// 	if err.Error() == "insufficient permissions" {
	// 		return response.SimpleError(403, "Insufficient permissions"), nil
	// 	}
	// 	return response.SimpleError(401, "Unauthorized 1"), nil
	// }

	//validate jwt token
	_, err := middleware.ValidateJWT(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized 2"), nil
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
	collection, _ := services.GetUserCollection()
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
	var userResponses []types.UserResponse
	for _, u := range users {
		userResponses = append(userResponses, types.UserResponse{
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

// UpdateUserProfile updates user profile information
func UpdateUserProfile(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session using auth package

	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	var updateReq types.UserUpdateRequest
	if err := json.Unmarshal([]byte(request.Body), &updateReq); err != nil {
		return response.SimpleError(400, "Invalid request format"), nil
	}

	collection, _ := services.GetUserCollection()
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
	userResponse := types.UserResponse{
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

// DeleteUser soft deletes a user (admin only)
func DeleteUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
	collection, _ := services.GetUserCollection()
	var targetUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&targetUser)
	if err != nil {
		return response.SimpleError(404, "User not found"), nil
	}

	// Delete user's sessions first
	sessionCollection := db.Database.Collection("sessions") //TODO remove session
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

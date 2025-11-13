package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"egaldeutsch-serverless/db"
	"egaldeutsch-serverless/pkg/auth"
	"egaldeutsch-serverless/pkg/middleware"
	"egaldeutsch-serverless/pkg/notification"
	"egaldeutsch-serverless/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	case "GET":
		if strings.Contains(request.Path, "/unread-count") {
			return getUnreadCount(request)
		}
		return listNotifications(request)
	case "PUT":
		if strings.Contains(request.Path, "/read-all") {
			return markAllAsRead(request)
		}
		return markAsRead(request)
	case "DELETE":
		return deleteNotification(request)
	default:
		return response.SimpleError(405, "Method not allowed"), nil
	}
}

// listNotifications retrieves notifications for the authenticated user
func listNotifications(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session
	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Parse query parameters
	queryParams := request.QueryStringParameters
	page := 1
	limit := 20
	unreadOnly := false

	if p := queryParams["page"]; p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := queryParams["limit"]; l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if queryParams["unreadOnly"] == "true" {
		unreadOnly = true
	}

	// Get notifications
	notifications, total, err := notification.GetUserNotifications(user.ID, page, limit, unreadOnly)
	if err != nil {
		log.Printf("Failed to get notifications for user %s: %v", user.ID.Hex(), err)
		return response.SimpleError(500, "Failed to retrieve notifications"), nil
	}

	totalPages := (int(total) + limit - 1) / limit

	result := map[string]interface{}{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"limit":         limit,
		"totalPages":    totalPages,
	}

	return response.SuccessJSON(200, result, "Notifications retrieved successfully")
}

// getUnreadCount returns the count of unread notifications
func getUnreadCount(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session
	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Get unread count
	count, err := notification.GetUnreadCount(user.ID)
	if err != nil {
		log.Printf("Failed to get unread count for user %s: %v", user.ID.Hex(), err)
		return response.SimpleError(500, "Failed to get unread count"), nil
	}

	result := map[string]interface{}{
		"count": count,
	}

	return response.SuccessJSON(200, result, "Unread count retrieved successfully")
}

// markAsRead marks a single notification as read
func markAsRead(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session
	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Get notification ID from path or body
	notificationID := request.PathParameters["id"]
	if notificationID == "" {
		// Try to get from body
		var body map[string]string
		if err := json.Unmarshal([]byte(request.Body), &body); err == nil {
			notificationID = body["id"]
		}
	}

	if notificationID == "" {
		return response.SimpleError(400, "Notification ID is required"), nil
	}

	// Convert to ObjectID
	objectID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return response.SimpleError(400, "Invalid notification ID"), nil
	}

	// Mark as read
	if err := notification.MarkAsRead(objectID, user.ID); err != nil {
		log.Printf("Failed to mark notification as read: %v", err)
		return response.SimpleError(500, "Failed to mark notification as read"), nil
	}

	return response.SuccessJSON(200, nil, "Notification marked as read")
}

// markAllAsRead marks all notifications as read for the user
func markAllAsRead(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session
	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Mark all as read
	if err := notification.MarkAllAsRead(user.ID); err != nil {
		log.Printf("Failed to mark all notifications as read: %v", err)
		return response.SimpleError(500, "Failed to mark all notifications as read"), nil
	}

	return response.SuccessJSON(200, nil, "All notifications marked as read")
}

// deleteNotification deletes a notification
func deleteNotification(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate session
	user, err := auth.ValidateSession(request)
	if err != nil {
		return response.SimpleError(401, "Unauthorized"), nil
	}

	// Get notification ID
	notificationID := request.PathParameters["id"]
	if notificationID == "" {
		return response.SimpleError(400, "Notification ID is required"), nil
	}

	// Convert to ObjectID
	objectID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return response.SimpleError(400, "Invalid notification ID"), nil
	}

	// Delete notification
	if err := notification.DeleteNotification(objectID, user.ID); err != nil {
		log.Printf("Failed to delete notification: %v", err)
		return response.SimpleError(500, "Failed to delete notification"), nil
	}

	return response.SuccessJSON(200, nil, "Notification deleted successfully")
}

func main() {
	lambda.Start(handler)
}

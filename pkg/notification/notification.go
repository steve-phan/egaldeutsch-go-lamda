package notification

import (
	"context"
	"time"

	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeStoryPublished  NotificationType = "story_published"
	NotificationTypeStorySubmitted  NotificationType = "story_submitted"
	NotificationTypeStoryApproved   NotificationType = "story_approved"
	NotificationTypeStoryRejected   NotificationType = "story_rejected"
	NotificationTypeQuizPublished   NotificationType = "quiz_published"
	NotificationTypeQuizSubmitted   NotificationType = "quiz_submitted"
	NotificationTypeNewComment      NotificationType = "new_comment"
	NotificationTypeUserRegistered  NotificationType = "user_registered"
	NotificationTypePasswordChanged NotificationType = "password_changed"
)

// Notification represents a system notification
type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`                     // Recipient user ID
	Type      NotificationType   `bson:"type" json:"type"`                         // Type of notification
	Title     string             `bson:"title" json:"title"`                       // Notification title
	Message   string             `bson:"message" json:"message"`                   // Notification message
	Link      string             `bson:"link,omitempty" json:"link,omitempty"`     // Optional link to related content
	IsRead    bool               `bson:"isRead" json:"isRead"`                     // Whether notification has been read
	ReadAt    *time.Time         `bson:"readAt,omitempty" json:"readAt,omitempty"` // When notification was read
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`               // When notification was created
	Metadata  map[string]string  `bson:"metadata,omitempty" json:"metadata,omitempty"` // Additional metadata
}

// CreateNotification creates a new notification for a user
func CreateNotification(userID primitive.ObjectID, notifType NotificationType, title, message, link string, metadata map[string]string) error {
	notification := Notification{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Link:      link,
		IsRead:    false,
		CreatedAt: time.Now(),
		Metadata:  metadata,
	}

	collection := db.Database.Collection("notifications")
	_, err := collection.InsertOne(context.TODO(), notification)
	return err
}

// CreateNotificationForAllUsers creates a notification for all active users
func CreateNotificationForAllUsers(notifType NotificationType, title, message, link string, metadata map[string]string) error {
	// Get all active users
	usersCollection := db.Database.Collection("users")
	cursor, err := usersCollection.Find(context.TODO(), bson.M{"status": "active"})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	var userIDs []primitive.ObjectID
	for cursor.Next(context.TODO()) {
		var user struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		userIDs = append(userIDs, user.ID)
	}

	// Create notifications for all users
	collection := db.Database.Collection("notifications")
	var notifications []interface{}
	now := time.Now()

	for _, userID := range userIDs {
		notifications = append(notifications, Notification{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Type:      notifType,
			Title:     title,
			Message:   message,
			Link:      link,
			IsRead:    false,
			CreatedAt: now,
			Metadata:  metadata,
		})
	}

	if len(notifications) > 0 {
		_, err = collection.InsertMany(context.TODO(), notifications)
		return err
	}

	return nil
}

// GetUserNotifications retrieves notifications for a user with pagination
func GetUserNotifications(userID primitive.ObjectID, page, limit int, unreadOnly bool) ([]Notification, int64, error) {
	collection := db.Database.Collection("notifications")

	filter := bson.M{"userId": userID}
	if unreadOnly {
		filter["isRead"] = false
	}

	// Get total count
	total, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var notifications []Notification
	if err = cursor.All(context.TODO(), &notifications); err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetUnreadCount returns the count of unread notifications for a user
func GetUnreadCount(userID primitive.ObjectID) (int64, error) {
	collection := db.Database.Collection("notifications")
	return collection.CountDocuments(context.TODO(), bson.M{
		"userId": userID,
		"isRead": false,
	})
}

// MarkAsRead marks a notification as read
func MarkAsRead(notificationID primitive.ObjectID, userID primitive.ObjectID) error {
	collection := db.Database.Collection("notifications")
	now := time.Now()

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": notificationID, "userId": userID},
		bson.M{"$set": bson.M{
			"isRead": true,
			"readAt": now,
		}},
	)
	return err
}

// MarkAllAsRead marks all notifications as read for a user
func MarkAllAsRead(userID primitive.ObjectID) error {
	collection := db.Database.Collection("notifications")
	now := time.Now()

	_, err := collection.UpdateMany(
		context.TODO(),
		bson.M{"userId": userID, "isRead": false},
		bson.M{"$set": bson.M{
			"isRead": true,
			"readAt": now,
		}},
	)
	return err
}

// DeleteNotification deletes a notification
func DeleteNotification(notificationID primitive.ObjectID, userID primitive.ObjectID) error {
	collection := db.Database.Collection("notifications")
	_, err := collection.DeleteOne(
		context.TODO(),
		bson.M{"_id": notificationID, "userId": userID},
	)
	return err
}

// NotifyStoryPublished creates notifications when a story is published
func NotifyStoryPublished(storyID primitive.ObjectID, storyTitle string) error {
	return CreateNotificationForAllUsers(
		NotificationTypeStoryPublished,
		"New Story Available!",
		"A new story '"+storyTitle+"' has been published. Check it out!",
		"/stories?id="+storyID.Hex(),
		map[string]string{
			"storyId":    storyID.Hex(),
			"storyTitle": storyTitle,
		},
	)
}

// NotifyStorySubmitted notifies admins/reviewers when a story is submitted for review
func NotifyStorySubmitted(creatorID, storyID primitive.ObjectID, storyTitle string) error {
	// Get all admin and reviewer users
	usersCollection := db.Database.Collection("users")
	cursor, err := usersCollection.Find(context.TODO(), bson.M{
		"status": "active",
		"role":   bson.M{"$in": []string{"admin", "reviewer"}},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	collection := db.Database.Collection("notifications")
	var notifications []interface{}
	now := time.Now()

	for cursor.Next(context.TODO()) {
		var user struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		notifications = append(notifications, Notification{
			ID:      primitive.NewObjectID(),
			UserID:  user.ID,
			Type:    NotificationTypeStorySubmitted,
			Title:   "Story Submitted for Review",
			Message: "A new story '" + storyTitle + "' has been submitted for review.",
			Link:    "/admin/reviews?id=" + storyID.Hex(),
			IsRead:  false,
			CreatedAt: now,
			Metadata: map[string]string{
				"storyId":    storyID.Hex(),
				"storyTitle": storyTitle,
				"creatorId":  creatorID.Hex(),
			},
		})
	}

	if len(notifications) > 0 {
		_, err = collection.InsertMany(context.TODO(), notifications)
		return err
	}

	return nil
}

// NotifyStoryStatusChange notifies the creator about story status changes
func NotifyStoryStatusChange(creatorID, storyID primitive.ObjectID, storyTitle, newStatus, message string) error {
	var title string
	var notifType NotificationType

	switch newStatus {
	case "ready":
		title = "Story Approved!"
		notifType = NotificationTypeStoryApproved
	case "draft":
		title = "Story Needs Revision"
		notifType = NotificationTypeStoryRejected
	case "published":
		title = "Story Published!"
		notifType = NotificationTypeStoryPublished
	default:
		title = "Story Status Updated"
		notifType = NotificationTypeStorySubmitted
	}

	return CreateNotification(
		creatorID,
		notifType,
		title,
		message,
		"/stories/create?id="+storyID.Hex(),
		map[string]string{
			"storyId":    storyID.Hex(),
			"storyTitle": storyTitle,
			"newStatus":  newStatus,
		},
	)
}

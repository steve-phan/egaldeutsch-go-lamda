package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"egaldeutsch-serverless/models"
)

// NotificationService handles email notifications
type NotificationService struct {
	emailServiceURL string
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		emailServiceURL: getEmailServiceURL(),
	}
}

// SendNewStoryNotification sends email notification for new story
func (n *NotificationService) SendNewStoryNotification(storyID, storyTitle, storyLevel string) error {
	payload := map[string]interface{}{
		"type":    "new_story",
		"storyId": storyID,
		"title":   storyTitle,
		"level":   storyLevel,
	}

	return n.callEmailService(n.emailServiceURL, payload)
}

// callEmailService makes HTTP request to email service
func (n *NotificationService) callEmailService(url string, payload map[string]interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("email service returned status %d", resp.StatusCode)
	}

	return nil
}

// getEmailServiceURL returns the email service URL from environment
func getEmailServiceURL() string {
	if url := os.Getenv("EMAIL_SERVICE_URL"); url != "" {
		return url
	}
	// Default for local development
	return "http://localhost:8080/.netlify/functions/email-service"
}

// GetCommentType returns comment type based on status
func GetCommentType(status models.ContentStatus) string {
	switch status {
	case models.StatusPublished:
		return "approved"
	case models.StatusReady:
		return "ready"
	case models.StatusDraft:
		return "draft"
	case models.StatusPreview:
		return "preview"
	default:
		return "other"
	}
}

package email

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Service manages email operations with provider abstraction
type Service struct {
	provider EmailProvider
	config   *EmailConfig
}

// NewService creates a new email service with the specified provider
func NewService(config *EmailConfig) (*Service, error) {
	if config == nil {
		return nil, fmt.Errorf("email config is required")
	}

	var provider EmailProvider

	switch strings.ToLower(config.Provider) {
	case "brevo", "sendinblue":
		provider = NewBrevoProvider(config.APIKey, config.FromEmail, config.FromName)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", config.Provider)
	}

	if err := provider.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("provider validation failed: %w", err)
	}

	return &Service{
		provider: provider,
		config:   config,
	}, nil
}

// NewServiceFromEnv creates a new email service from environment variables
func NewServiceFromEnv() (*Service, error) {
	config := &EmailConfig{
		Provider:     getEnvOrDefault("EMAIL_PROVIDER", "brevo"),
		APIKey:       os.Getenv("EMAIL_API_KEY"),
		FromEmail:    os.Getenv("EMAIL_FROM"),
		FromName:     getEnvOrDefault("EMAIL_FROM_NAME", "EgalDeutsch"),
		ReplyToEmail: os.Getenv("EMAIL_REPLY_TO"),
	}

	// Parse max bulk recipients
	if maxBulkStr := os.Getenv("EMAIL_MAX_BULK_RECIPIENTS"); maxBulkStr != "" {
		if maxBulk, err := strconv.Atoi(maxBulkStr); err == nil {
			config.MaxBulkRecipients = maxBulk
		}
	}
	
	// Default to 10 for trial accounts
	if config.MaxBulkRecipients == 0 {
		config.MaxBulkRecipients = 10
	}

	return NewService(config)
}

// SendEmail sends a single email
func (s *Service) SendEmail(email *Email) error {
	// Set default from if not specified
	if email.From == "" {
		email.From = s.config.FromEmail
	}
	if email.FromName == "" {
		email.FromName = s.config.FromName
	}
	if email.ReplyTo == "" && s.config.ReplyToEmail != "" {
		email.ReplyTo = s.config.ReplyToEmail
	}

	return s.provider.SendEmail(email)
}

// SendBulkEmail sends emails to multiple recipients with bulk limits
func (s *Service) SendBulkEmail(emails []*Email) error {
	if len(emails) == 0 {
		return fmt.Errorf("no emails to send")
	}

	// Apply bulk recipient limit
	if len(emails) > s.config.MaxBulkRecipients {
		emails = emails[:s.config.MaxBulkRecipients]
	}

	// Set defaults for all emails
	for _, email := range emails {
		if email.From == "" {
			email.From = s.config.FromEmail
		}
		if email.FromName == "" {
			email.FromName = s.config.FromName
		}
		if email.ReplyTo == "" && s.config.ReplyToEmail != "" {
			email.ReplyTo = s.config.ReplyToEmail
		}
	}

	return s.provider.SendBulkEmail(emails)
}

// SendWelcomeEmail sends a welcome email to new users
func (s *Service) SendWelcomeEmail(userEmail, userName string) error {
	template := s.getWelcomeEmailTemplate(userName)
	
	email := &Email{
		To:          []string{userEmail},
		Subject:     template.Subject,
		HTMLContent: template.HTMLContent,
		TextContent: template.TextContent,
	}

	return s.SendEmail(email)
}

// SendPasswordResetEmail sends a password reset email
func (s *Service) SendPasswordResetEmail(userEmail, userName, resetToken string) error {
	template := s.getPasswordResetEmailTemplate(userName, resetToken)
	
	email := &Email{
		To:          []string{userEmail},
		Subject:     template.Subject,
		HTMLContent: template.HTMLContent,
		TextContent: template.TextContent,
	}

	return s.SendEmail(email)
}

// SendNewStoryNotification sends new story notification to subscribers
func (s *Service) SendNewStoryNotification(subscriberEmails []string, storyTitle, storyLevel string, storyID string) error {
	if len(subscriberEmails) == 0 {
		return fmt.Errorf("no subscribers to notify")
	}

	// Apply bulk limit
	if len(subscriberEmails) > s.config.MaxBulkRecipients {
		subscriberEmails = subscriberEmails[:s.config.MaxBulkRecipients]
	}

	template := s.getNewStoryEmailTemplate(storyTitle, storyLevel, storyID)
	
	var emails []*Email
	for _, subscriberEmail := range subscriberEmails {
		email := &Email{
			To:          []string{subscriberEmail},
			Subject:     template.Subject,
			HTMLContent: template.HTMLContent,
			TextContent: template.TextContent,
		}
		emails = append(emails, email)
	}

	return s.SendBulkEmail(emails)
}

// GetConfig returns the current email configuration (sensitive data masked)
func (s *Service) GetConfig() *EmailConfig {
	config := *s.config
	// Mask sensitive data
	if len(config.APIKey) > 8 {
		config.APIKey = config.APIKey[:4] + "****" + config.APIKey[len(config.APIKey)-4:]
	}
	return &config
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
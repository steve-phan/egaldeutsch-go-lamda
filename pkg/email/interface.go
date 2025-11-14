package email

import "time"

// EmailProvider defines the interface that all email providers must implement
type EmailProvider interface {
	SendEmail(email *Email) error
	SendBulkEmail(emails []*Email) error
	ValidateConfig() error
}

// Email represents an email message
type Email struct {
	To          []string          `json:"to"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	From        string            `json:"from"`
	FromName    string            `json:"fromName,omitempty"`
	Subject     string            `json:"subject"`
	HTMLContent string            `json:"htmlContent,omitempty"`
	TextContent string            `json:"textContent,omitempty"`
	ReplyTo     string            `json:"replyTo,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
}

// Attachment represents an email attachment
type Attachment struct {
	Name        string `json:"name"`
	Content     []byte `json:"content"`
	ContentType string `json:"contentType"`
}

// EmailTemplate represents predefined email templates
type EmailTemplate struct {
	Name        string            `json:"name"`
	Subject     string            `json:"subject"`
	HTMLContent string            `json:"htmlContent"`
	TextContent string            `json:"textContent"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Provider        string `json:"provider"`        // "brevo", "sendgrid", "ses", etc.
	APIKey          string `json:"apiKey"`
	FromEmail       string `json:"fromEmail"`
	FromName        string `json:"fromName"`
	ReplyToEmail    string `json:"replyToEmail"`
	MaxBulkRecipients int  `json:"maxBulkRecipients"` // Max recipients per bulk email (for trial accounts)
}

// EmailStats represents email sending statistics
type EmailStats struct {
	Sent      int       `json:"sent"`
	Failed    int       `json:"failed"`
	Timestamp time.Time `json:"timestamp"`
	Provider  string    `json:"provider"`
}
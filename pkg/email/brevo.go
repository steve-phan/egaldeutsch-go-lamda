package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// BrevoProvider implements EmailProvider interface for Brevo (formerly Sendinblue)
type BrevoProvider struct {
	APIKey    string
	BaseURL   string
	FromEmail string
	FromName  string
}

// NewBrevoProvider creates a new Brevo email provider
func NewBrevoProvider(apiKey, fromEmail, fromName string) *BrevoProvider {
	return &BrevoProvider{
		APIKey:    apiKey,
		BaseURL:   "https://api.brevo.com/v3",
		FromEmail: fromEmail,
		FromName:  fromName,
	}
}

// BrevoEmail represents Brevo API email structure
type BrevoEmail struct {
	Sender      BrevoContact      `json:"sender"`
	To          []BrevoContact    `json:"to"`
	CC          []BrevoContact    `json:"cc,omitempty"`
	BCC         []BrevoContact    `json:"bcc,omitempty"`
	Subject     string            `json:"subject"`
	HTMLContent string            `json:"htmlContent,omitempty"`
	TextContent string            `json:"textContent,omitempty"`
	ReplyTo     *BrevoContact     `json:"replyTo,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// BrevoContact represents a contact in Brevo API
type BrevoContact struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// BrevoBulkEmail represents Brevo bulk email structure
type BrevoBulkEmail struct {
	Sender      BrevoContact   `json:"sender"`
	To          []BrevoContact `json:"to"`
	Subject     string         `json:"subject"`
	HTMLContent string         `json:"htmlContent,omitempty"`
	TextContent string         `json:"textContent,omitempty"`
	ReplyTo     *BrevoContact  `json:"replyTo,omitempty"`
}

// BrevoResponse represents Brevo API response
type BrevoResponse struct {
	MessageID string `json:"messageId"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
}

// ValidateConfig validates the Brevo provider configuration
func (b *BrevoProvider) ValidateConfig() error {
	if b.APIKey == "" {
		return fmt.Errorf("brevo API key is required")
	}
	if b.FromEmail == "" {
		return fmt.Errorf("from email is required")
	}
	if !strings.Contains(b.FromEmail, "@") {
		return fmt.Errorf("invalid from email format")
	}
	return nil
}

// SendEmail sends a single email using Brevo API
func (b *BrevoProvider) SendEmail(email *Email) error {
	if err := b.ValidateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Convert to Brevo format
	brevoEmail := b.convertToBrevoEmail(email)

	// Make API request
	return b.makeAPIRequest("/smtp/email", brevoEmail)
}

// SendBulkEmail sends multiple emails using Brevo API
func (b *BrevoProvider) SendBulkEmail(emails []*Email) error {
	if err := b.ValidateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	if len(emails) == 0 {
		return fmt.Errorf("no emails to send")
	}

	// For bulk emails, we'll send them individually through Brevo
	// This gives us better error handling per email
	var errors []string
	for i, email := range emails {
		if err := b.SendEmail(email); err != nil {
			errors = append(errors, fmt.Sprintf("email %d: %v", i+1, err))
		}

		// Add small delay to avoid rate limiting
		if i < len(emails)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("bulk email errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// convertToBrevoEmail converts our Email struct to Brevo API format
func (b *BrevoProvider) convertToBrevoEmail(email *Email) *BrevoEmail {
	brevoEmail := &BrevoEmail{
		Sender: BrevoContact{
			Email: b.FromEmail,
			Name:  b.FromName,
		},
		Subject:     email.Subject,
		HTMLContent: email.HTMLContent,
		TextContent: email.TextContent,
		Headers:     email.Headers,
	}

	// Convert recipients
	for _, to := range email.To {
		brevoEmail.To = append(brevoEmail.To, BrevoContact{Email: to})
	}

	for _, cc := range email.CC {
		brevoEmail.CC = append(brevoEmail.CC, BrevoContact{Email: cc})
	}

	for _, bcc := range email.BCC {
		brevoEmail.BCC = append(brevoEmail.BCC, BrevoContact{Email: bcc})
	}

	// Set reply-to if provided
	if email.ReplyTo != "" {
		brevoEmail.ReplyTo = &BrevoContact{Email: email.ReplyTo}
	}

	return brevoEmail
}

// makeAPIRequest makes an HTTP request to Brevo API
func (b *BrevoProvider) makeAPIRequest(endpoint string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", b.BaseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", b.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var brevoResp BrevoResponse
	if err := json.NewDecoder(resp.Body).Decode(&brevoResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("brevo API error (status %d): %s - %s", resp.StatusCode, brevoResp.Code, brevoResp.Message)
	}

	return nil
}

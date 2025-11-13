package email

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"egaldeutsch-serverless/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmailService provides email sending functionality
type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expiresAt"`
	CreatedAt time.Time          `bson:"createdAt"`
	Used      bool               `bson:"used"`
	UsedAt    *time.Time         `bson:"usedAt,omitempty"`
}

// NewEmailService creates a new email service from environment variables
func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
		FromName:     os.Getenv("FROM_NAME"),
	}
}

// IsConfigured checks if email service is properly configured
func (s *EmailService) IsConfigured() bool {
	return s.SMTPHost != "" && s.SMTPPort != "" && s.FromEmail != ""
}

// SendEmail sends an email
func (s *EmailService) SendEmail(to, subject, body string) error {
	if !s.IsConfigured() {
		return fmt.Errorf("email service not configured")
	}

	// Prepare email message
	from := s.FromEmail
	if s.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.FromName, s.FromEmail)
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-version: 1.0;\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
			"\r\n"+
			"%s\r\n",
		from, to, subject, body,
	))

	// Setup authentication
	auth := smtp.PlainAuth("", s.SMTPUsername, s.SMTPPassword, s.SMTPHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)
	err := smtp.SendMail(addr, auth, s.FromEmail, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *EmailService) SendWelcomeEmail(to, firstName string) error {
	subject := "Welcome to EgalDeutsch - Learn German!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2563eb;">Welcome to EgalDeutsch!</h1>
        <p>Hello %s,</p>
        <p>Thank you for registering with EgalDeutsch - your platform for learning German through stories and interactive quizzes!</p>
        
        <h2 style="color: #2563eb;">What's Next?</h2>
        <ul>
            <li>Explore German stories at different difficulty levels (A1-C2)</li>
            <li>Test your comprehension with interactive quizzes</li>
            <li>Track your progress on the leaderboard</li>
        </ul>
        
        <p>Start your learning journey today!</p>
        
        <div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd;">
            <p style="font-size: 12px; color: #666;">
                If you have any questions, feel free to reach out to our support team.
            </p>
        </div>
    </div>
</body>
</html>
`, firstName)

	return s.SendEmail(to, subject, body)
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, firstName, resetToken string) error {
	// Construct reset URL (you should set this as an environment variable in production)
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", baseURL, resetToken)

	subject := "Reset Your EgalDeutsch Password"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2563eb;">Password Reset Request</h1>
        <p>Hello %s,</p>
        <p>We received a request to reset your password for your EgalDeutsch account.</p>
        
        <p>Click the button below to reset your password:</p>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #2563eb; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">
                Reset Password
            </a>
        </div>
        
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #2563eb;">%s</p>
        
        <p><strong>This link will expire in 1 hour.</strong></p>
        
        <p>If you didn't request a password reset, you can safely ignore this email. Your password will not be changed.</p>
        
        <div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd;">
            <p style="font-size: 12px; color: #666;">
                For security reasons, never share this link with anyone.
            </p>
        </div>
    </div>
</body>
</html>
`, firstName, resetURL, resetURL)

	return s.SendEmail(to, subject, body)
}

// GenerateResetToken generates a password reset token
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreatePasswordResetToken creates and stores a password reset token
func CreatePasswordResetToken(userID primitive.ObjectID) (string, error) {
	token, err := GenerateResetToken()
	if err != nil {
		return "", err
	}

	resetToken := PasswordResetToken{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token expires in 1 hour
		CreatedAt: time.Now(),
		Used:      false,
	}

	collection := db.Database.Collection("password_reset_tokens")
	_, err = collection.InsertOne(context.TODO(), resetToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateResetToken validates a password reset token and returns the user ID
func ValidateResetToken(token string) (primitive.ObjectID, error) {
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

// MarkTokenAsUsed marks a password reset token as used
func MarkTokenAsUsed(token string) error {
	collection := db.Database.Collection("password_reset_tokens")
	now := time.Now()

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"token": token},
		bson.M{"$set": bson.M{
			"used":   true,
			"usedAt": now,
		}},
	)
	return err
}

// SendPasswordChangedEmail sends a notification email when password is changed
func (s *EmailService) SendPasswordChangedEmail(to, firstName string) error {
	subject := "Your EgalDeutsch Password Has Been Changed"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2563eb;">Password Changed Successfully</h1>
        <p>Hello %s,</p>
        <p>This is a confirmation that your password for your EgalDeutsch account has been successfully changed.</p>
        
        <p><strong>If you made this change, no further action is needed.</strong></p>
        
        <p>If you did not make this change, please contact our support team immediately to secure your account.</p>
        
        <div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd;">
            <p style="font-size: 12px; color: #666;">
                This is an automated security notification. Please do not reply to this email.
            </p>
        </div>
    </div>
</body>
</html>
`, firstName)

	return s.SendEmail(to, subject, body)
}

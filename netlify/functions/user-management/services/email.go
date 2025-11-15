package services

import (
	"fmt"
	"log"

	"egaldeutsch-serverless/pkg/email"
)

// SendWelcomeEmail sends a welcome email directly using email service
func SendWelcomeEmail(userEmail, userName string) error {
	log.Printf("Sending welcome email to: %s", userEmail)

	// Initialize email service from environment
	emailService, err := email.NewServiceFromEnv()
	if err != nil {
		log.Printf("Failed to initialize email service: %v", err)
		return fmt.Errorf("failed to initialize email service: %w", err)
	}

	// Send welcome email
	if err := emailService.SendWelcomeEmail(userEmail, userName); err != nil {
		log.Printf("Failed to send welcome email: %v", err)
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	log.Printf("Welcome email sent successfully to: %s", userEmail)
	return nil
}

// SendPasswordResetEmail sends a password reset email directly using email service
func SendPasswordResetEmail(userEmail, userName, resetToken string) error {
	log.Printf("Sending password reset email to: %s", userEmail)

	// Initialize email service from environment
	emailService, err := email.NewServiceFromEnv()
	if err != nil {
		log.Printf("Failed to initialize email service: %v", err)
		return fmt.Errorf("failed to initialize email service: %w", err)
	}

	// Send password reset email
	if err := emailService.SendPasswordResetEmail(userEmail, userName, resetToken); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	log.Printf("Password reset email sent successfully to: %s", userEmail)
	return nil
}

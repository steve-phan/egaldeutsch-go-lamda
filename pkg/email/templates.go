package email

import (
	"fmt"
	"strings"
)

// getWelcomeEmailTemplate returns the welcome email template
func (s *Service) getWelcomeEmailTemplate(userName string) *EmailTemplate {
	baseURL := getEnvOrDefault("FRONTEND_URL", "http://localhost:8000")

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Welcome to EgalDeutsch!</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #f8fafc; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; margin: 10px 0; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎉 Welcome to EgalDeutsch!</h1>
    </div>
    <div class="content">
        <h2>Hello %s!</h2>
        <p>Thank you for joining EgalDeutsch, your platform for learning German through engaging stories and interactive quizzes.</p>
        <p>With EgalDeutsch, you can:</p>
        <ul>
            <li>📚 Read German stories at your CEFR level (A1-C2)</li>
            <li>🧩 Take interactive quizzes to test your comprehension</li>
            <li>📈 Track your progress and improve your German skills</li>
            <li>🎯 Focus on vocabulary, grammar, and reading comprehension</li>
        </ul>
        <p>Ready to start your German learning journey?</p>
        <a href="%s/stories" class="button">Explore Stories</a>
        <p>If you have any questions, feel free to reply to this email or contact our support team.</p>
        <p>Viel Erfolg beim Deutschlernen! (Good luck learning German!)</p>
        <p>Best regards,<br>The EgalDeutsch Team</p>
    </div>
    <div class="footer">
        <p>EgalDeutsch - Learn German Through Stories</p>
        <p>If you didn't create this account, please ignore this email.</p>
    </div>
</body>
</html>`, userName, baseURL)

	textContent := fmt.Sprintf(`
Welcome to EgalDeutsch!

Hello %s!

Thank you for joining EgalDeutsch, your platform for learning German through engaging stories and interactive quizzes.

With EgalDeutsch, you can:
- Read German stories at your CEFR level (A1-C2)
- Take interactive quizzes to test your comprehension
- Track your progress and improve your German skills
- Focus on vocabulary, grammar, and reading comprehension

Ready to start your German learning journey?
Visit: %s/stories

If you have any questions, feel free to reply to this email or contact our support team.

Viel Erfolg beim Deutschlernen! (Good luck learning German!)

Best regards,
The EgalDeutsch Team

---
EgalDeutsch - Learn German Through Stories
If you didn't create this account, please ignore this email.`, userName, baseURL)

	return &EmailTemplate{
		Name:        "welcome",
		Subject:     "🎉 Welcome to EgalDeutsch - Start Learning German Today!",
		HTMLContent: htmlContent,
		TextContent: textContent,
	}
}

// getPasswordResetEmailTemplate returns the password reset email template
func (s *Service) getPasswordResetEmailTemplate(userName, resetToken string) *EmailTemplate {
	baseURL := getEnvOrDefault("FRONTEND_URL", "http://localhost:8000")
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", baseURL, resetToken)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Reset Your Password - EgalDeutsch</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #dc2626; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #f8fafc; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background-color: #dc2626; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; margin: 10px 0; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
        .warning { background-color: #fef2f2; border-left: 4px solid #dc2626; padding: 15px; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔒 Password Reset Request</h1>
    </div>
    <div class="content">
        <h2>Hello %s!</h2>
        <p>We received a request to reset your password for your EgalDeutsch account.</p>
        <p>Click the button below to reset your password:</p>
        <a href="%s" class="button">Reset Password</a>
        <div class="warning">
            <p><strong>⚠️ Important:</strong></p>
            <ul>
                <li>This link will expire in 1 hour for security reasons</li>
                <li>If you didn't request this reset, please ignore this email</li>
                <li>Never share this link with anyone</li>
            </ul>
        </div>
        <p>If the button doesn't work, copy and paste this link into your browser:</p>
        <p style="word-break: break-all; background-color: #f1f5f9; padding: 10px; border-radius: 4px; font-family: monospace;">%s</p>
        <p>If you have any questions, feel free to contact our support team.</p>
        <p>Best regards,<br>The EgalDeutsch Team</p>
    </div>
    <div class="footer">
        <p>EgalDeutsch - Learn German Through Stories</p>
        <p>This email was sent because a password reset was requested for your account.</p>
    </div>
</body>
</html>`, userName, resetURL, resetURL)

	textContent := fmt.Sprintf(`
Password Reset Request - EgalDeutsch

Hello %s!

We received a request to reset your password for your EgalDeutsch account.

To reset your password, click on this link or copy it into your browser:
%s

IMPORTANT:
- This link will expire in 1 hour for security reasons
- If you didn't request this reset, please ignore this email
- Never share this link with anyone

If you have any questions, feel free to contact our support team.

Best regards,
The EgalDeutsch Team

---
EgalDeutsch - Learn German Through Stories
This email was sent because a password reset was requested for your account.`, userName, resetURL)

	return &EmailTemplate{
		Name:        "password_reset",
		Subject:     "🔒 Reset Your EgalDeutsch Password",
		HTMLContent: htmlContent,
		TextContent: textContent,
	}
}

// getNewStoryEmailTemplate returns the new story notification email template
func (s *Service) getNewStoryEmailTemplate(storyTitle, storyLevel, storyID string) *EmailTemplate {
	baseURL := getEnvOrDefault("FRONTEND_URL", "http://localhost:8000")
	storyURL := fmt.Sprintf("%s/story/%s", baseURL, storyID)

	levelEmoji := getLevelEmoji(storyLevel)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>New German Story Available - EgalDeutsch</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #059669; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #f0fdf4; padding: 30px; border-radius: 0 0 8px 8px; }
        .story-card { background-color: white; border: 2px solid #d1fae5; border-radius: 8px; padding: 20px; margin: 15px 0; }
        .level-badge { display: inline-block; background-color: #059669; color: white; padding: 4px 12px; border-radius: 20px; font-size: 14px; font-weight: bold; }
        .button { display: inline-block; background-color: #059669; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; margin: 10px 0; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>📚 New German Story Available!</h1>
    </div>
    <div class="content">
        <h2>A new story is waiting for you!</h2>
        <p>We're excited to share a new German story that's perfect for your learning journey.</p>
        
        <div class="story-card">
            <h3>%s %s</h3>
            <p><span class="level-badge">%s Level</span></p>
            <p>This story is designed to help you improve your German reading comprehension, vocabulary, and grammar skills at the %s level.</p>
            <a href="%s" class="button">Read Story Now</a>
        </div>
        
        <p>📖 <strong>What you'll find:</strong></p>
        <ul>
            <li>Engaging story content at your level</li>
            <li>Interactive vocabulary explanations</li>
            <li>Comprehension questions and quizzes</li>
            <li>Progress tracking for your learning journey</li>
        </ul>
        
        <p>Ready to dive into your next German learning adventure?</p>
        <a href="%s/stories" class="button">Browse All Stories</a>
        
        <p>Happy reading and learning!</p>
        <p>Viel Spaß beim Lesen! (Have fun reading!)</p>
        <p>Best regards,<br>The EgalDeutsch Team</p>
    </div>
    <div class="footer">
        <p>EgalDeutsch - Learn German Through Stories</p>
        <p>You're receiving this because you subscribed to new story notifications.</p>
    </div>
</body>
</html>`, levelEmoji, storyTitle, storyLevel, storyLevel, storyURL, baseURL)

	textContent := fmt.Sprintf(`
New German Story Available - EgalDeutsch

A new story is waiting for you!

We're excited to share a new German story that's perfect for your learning journey.

Story: %s %s
Level: %s

This story is designed to help you improve your German reading comprehension, vocabulary, and grammar skills at the %s level.

Read the story here: %s

What you'll find:
- Engaging story content at your level
- Interactive vocabulary explanations
- Comprehension questions and quizzes  
- Progress tracking for your learning journey

Browse all stories: %s/stories

Happy reading and learning!
Viel Spaß beim Lesen! (Have fun reading!)

Best regards,
The EgalDeutsch Team

---
EgalDeutsch - Learn German Through Stories
You're receiving this because you subscribed to new story notifications.`, levelEmoji, storyTitle, storyLevel, storyLevel, storyURL, baseURL)

	return &EmailTemplate{
		Name:        "new_story",
		Subject:     fmt.Sprintf("📚 New %s Story: %s", storyLevel, storyTitle),
		HTMLContent: htmlContent,
		TextContent: textContent,
	}
}

// getLevelEmoji returns appropriate emoji for CEFR level
func getLevelEmoji(level string) string {
	level = strings.ToUpper(level)
	switch level {
	case "A1":
		return "🌱"
	case "A2":
		return "🌿"
	case "B1":
		return "🌳"
	case "B2":
		return "🍃"
	case "C1":
		return "🎓"
	case "C2":
		return "👑"
	default:
		return "📚"
	}
}

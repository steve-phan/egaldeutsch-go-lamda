package types

import "time"

// UserResponse represents user data in responses (without sensitive info)
type UserResponse struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// AuthTokenResponse represents authentication response
type AuthTokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expiresAt"`
	User      UserResponse `json:"user"`
}
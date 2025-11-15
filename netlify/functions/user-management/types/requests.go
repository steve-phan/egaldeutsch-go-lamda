package types

// UserRegistrationRequest represents user registration data
type UserRegistrationRequest struct {
	Username      string `json:"username" validate:"required,min=3,max=50"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8,max=128"`
	FirstName     string `json:"firstName" validate:"required,min=2,max=50"`
	LastName      string `json:"lastName" validate:"required,min=2,max=50"`
	PreferredRole string `json:"preferredRole" validate:"required,oneof=creator reviewer"`
}

// UserLoginRequest represents login credentials
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// UserUpdateRequest represents user profile update data
type UserUpdateRequest struct {
	FirstName   *string `json:"firstName,omitempty"`
	LastName    *string `json:"lastName,omitempty"`
	Email       *string `json:"email,omitempty"`
	OldPassword *string `json:"oldPassword,omitempty"`
	NewPassword *string `json:"newPassword,omitempty"`
}

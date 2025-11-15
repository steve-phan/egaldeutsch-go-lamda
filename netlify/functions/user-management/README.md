# User Management Function - Code Structure

## Overview

The user-management Lambda function has been refactored from a single large file into a well-organized modular structure for better maintainability and separation of concerns.

## File Structure

```
netlify/functions/user-management/
├── main.go                     # Main entry point and routing
├── handlers/                   # HTTP request handlers
│   ├── registration.go         # User registration logic
│   ├── auth.go                 # User authentication (login)
│   ├── password.go             # Password reset functionality
│   ├── users.go                # User profile and management
│   └── session.go              # Session management (logout)
├── services/                   # Business logic services
│   ├── email.go                # Email sending services
│   └── tokens.go               # Password reset token management
└── types/                      # Common type definitions
    ├── requests.go             # Request DTOs
    └── responses.go            # Response DTOs
```

## File Responsibilities

### main.go

- Lambda entry point
- Environment variable initialization
- Database connection management
- HTTP routing logic
- CORS handling

### handlers/registration.go

- **RegisterUser**: User registration endpoint
- Handles user validation, password hashing, user creation
- Sends welcome emails and creates registration notifications

### handlers/auth.go

- **LoginUser**: User authentication endpoint
- Validates credentials, creates sessions, updates last login
- Returns authentication tokens

### handlers/password.go

- **ForgotPassword**: Password reset request endpoint
- **ResetPassword**: Password reset with token endpoint
- Handles token validation and password updates

### handlers/users.go

- **GetUserProfile**: Retrieves current user profile
- **ListUsers**: Lists users (admin/reviewer only)
- **UpdateUserProfile**: Updates user profile information
- **DeleteUser**: Soft deletes users (admin only)

### handlers/session.go

- **LogoutUser**: Invalidates user sessions

### services/email.go

- **SendWelcomeEmail**: Sends welcome emails to new users
- **SendPasswordResetEmail**: Sends password reset emails
- Centralized email service initialization

### services/tokens.go

- **CreatePasswordResetToken**: Creates password reset tokens
- **ValidateResetToken**: Validates and retrieves user ID from tokens
- **MarkTokenAsUsed**: Marks tokens as used
- **PasswordResetToken**: Token data structure

### types/requests.go

- **UserRegistrationRequest**: User registration data
- **UserLoginRequest**: Login credentials
- **ForgotPasswordRequest**: Password reset request
- **ResetPasswordRequest**: Password reset with token
- **UserUpdateRequest**: Profile update data

### types/responses.go

- **UserResponse**: User data for responses (without sensitive info)
- **AuthTokenResponse**: Authentication response with token

## Benefits of This Structure

1. **Separation of Concerns**: Each file has a single, well-defined responsibility
2. **Better Maintainability**: Changes to specific functionality are isolated
3. **Improved Readability**: Smaller, focused files are easier to understand
4. **Better Testing**: Individual components can be tested in isolation
5. **Code Reusability**: Services can be reused across different handlers
6. **Type Safety**: Centralized type definitions prevent inconsistencies

## Usage

The function still works exactly the same from the API perspective. All the same endpoints are available:

- POST `/register` - User registration
- POST `/login` - User login
- POST `/forgot-password` - Request password reset
- POST `/reset-password` - Reset password with token
- GET `/profile` - Get user profile
- GET `/` - List users (admin/reviewer)
- PUT `/` - Update user profile
- DELETE `/logout` - Logout user
- DELETE `/{id}` - Delete user (admin)

## Development

- Each handler can be developed and tested independently
- Services provide reusable business logic
- Types ensure consistency across the application
- Main.go provides clean routing and centralized configuration

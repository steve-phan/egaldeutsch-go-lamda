# Stories Function - Code Structure

## Overview

The stories Lambda function has been refactored from a single large file into a well-organized modular structure for better maintainability and separation of concerns.

## File Structure

```
netlify/functions/stories/
├── main.go                     # Main entry point and routing
├── handlers/                   # HTTP request handlers
│   └── stories.go              # Story CRUD operations
├── services/                   # Business logic services
│   ├── database.go             # Database connection management
│   └── cors.go                 # CORS handling utilities
└── types/                      # Common type definitions (empty for now)
```

## File Responsibilities

### main.go

- Lambda entry point
- HTTP routing logic
- CORS preflight handling
- Request delegation to appropriate handlers

### handlers/stories.go

- **GetAllStories**: Retrieves all published stories
- **GetStoryByID**: Retrieves a specific story by ID
- **CreateStory**: Creates a new story
- All handlers include proper error handling and CORS headers

### services/database.go

- **GetCollections**: Returns both stories and questions collections
- **GetStoriesCollection**: Returns only the stories collection
- **GetQuestionsCollection**: Returns only the questions collection
- Centralized database connection management

### services/cors.go

- **GetCORSHeaders**: Returns standard CORS headers
- **HandleCORSOptions**: Handles OPTIONS preflight requests
- Centralized CORS configuration

## API Endpoints

### GET `/` (no query parameters)

- Returns all published stories
- Response: `{"success": true, "data": [...], "message": "Found N stories"}`

### GET `/?id={storyId}`

- Returns a specific story by ID
- Only returns published stories
- Response: `{"success": true, "data": {...}, "message": "Story retrieved successfully"}`

### POST `/`

- Creates a new story
- Accepts story data in request body
- Automatically sets ID, createdAt, and updatedAt timestamps
- Response: Created story object

### OPTIONS `/`

- CORS preflight request handling
- Returns appropriate CORS headers

## Benefits of This Structure

1. **Separation of Concerns**: Each file has a single, well-defined responsibility
2. **Better Maintainability**: Changes to specific functionality are isolated
3. **Improved Readability**: Smaller, focused files are easier to understand
4. **Better Testing**: Individual components can be tested in isolation
5. **Code Reusability**: Services can be reused across different handlers
6. **Centralized Configuration**: CORS and database logic is centralized

## Key Features

- **Published Stories Only**: All GET operations filter for published stories
- **Proper Error Handling**: Comprehensive error responses with appropriate HTTP status codes
- **CORS Support**: Full CORS support for cross-origin requests
- **Database Connection Management**: Efficient database connection handling
- **Type Safety**: Leverages existing models.Story type for consistency

## Usage

The function works exactly the same from the API perspective. All the same endpoints are available with the same request/response formats.

## Development

- Each handler can be developed and tested independently
- Services provide reusable business logic
- Main.go provides clean routing and centralized request handling
- Easy to extend with additional story operations (update, delete, etc.)

## Future Enhancements

- Add update story functionality
- Add delete story functionality
- Add authentication/authorization
- Add story filtering and pagination
- Add story search functionality

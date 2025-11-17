package response

import (
	"encoding/json"
	"fmt"

	"egaldeutsch-serverless/pkg/middleware"

	"github.com/aws/aws-lambda-go/events"
)

// Standard error response structures
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SimpleErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response structure
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// JSON creates a JSON response with the given status code, body, and CORS config
func JSON(statusCode int, body interface{}, corsConfig middleware.CORSConfig) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return ErrorJSON(500, "Failed to encode response", err.Error(), corsConfig), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    middleware.GetCORSHeadersWithConfig(corsConfig),
		Body:       string(jsonBody),
	}, nil
}

// JSONWithDefault creates a JSON response with default authenticated CORS headers
func JSONWithDefault(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	return JSON(statusCode, body, middleware.AuthenticatedAPI)
}

// SuccessJSON creates a successful JSON response with specified CORS config
func SuccessJSON(statusCode int, data interface{}, message string, corsConfig middleware.CORSConfig) (events.APIGatewayProxyResponse, error) {
	return JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}, corsConfig)
}

// SuccessJSONWithDefault creates a successful JSON response with default CORS headers
func SuccessJSONWithDefault(statusCode int, data interface{}, message string) (events.APIGatewayProxyResponse, error) {
	return SuccessJSON(statusCode, data, message, middleware.AuthenticatedAPI)
}

// ErrorJSON creates an error JSON response with specified CORS config
func ErrorJSON(statusCode int, error string, message string, corsConfig middleware.CORSConfig) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(ErrorResponse{
		Error:   error,
		Message: message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    middleware.GetCORSHeadersWithConfig(corsConfig),
		Body:       string(body),
	}
}

// ErrorJSONWithDefault creates an error JSON response with default CORS headers
func ErrorJSONWithDefault(statusCode int, error string, message string) events.APIGatewayProxyResponse {
	return ErrorJSON(statusCode, error, message, middleware.AuthenticatedAPI)
}

// Error creates a standardized error response with specified CORS config
func Error(statusCode int, errorType string, message string, corsConfig middleware.CORSConfig) events.APIGatewayProxyResponse {
	response := ErrorResponse{
		Success: false,
		Error:   errorType,
		Message: message,
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    middleware.GetCORSHeadersWithConfig(corsConfig),
		Body:       string(body),
	}
}

// ErrorWithDefault creates an error response with default authenticated CORS
func ErrorWithDefault(statusCode int, errorType string, message string) events.APIGatewayProxyResponse {
	return Error(statusCode, errorType, message, middleware.AuthenticatedAPI)
}

// SimpleError creates a simple error response with specified CORS config
func SimpleError(statusCode int, message string, corsConfig middleware.CORSConfig) events.APIGatewayProxyResponse {
	return ErrorJSON(statusCode, message, "", corsConfig)
}

// SimpleErrorWithDefault creates a simple error response with default CORS headers
func SimpleErrorWithDefault(statusCode int, message string) events.APIGatewayProxyResponse {
	return SimpleError(statusCode, message, middleware.AuthenticatedAPI)
}

// Predefined common errors
func DatabaseError(message string) events.APIGatewayProxyResponse {
	return ErrorWithDefault(500, "database_error", message)
}

func ValidationError(message string) events.APIGatewayProxyResponse {
	return ErrorWithDefault(400, "validation_error", message)
}

func NotFoundError(resource string) events.APIGatewayProxyResponse {
	return ErrorWithDefault(404, "not_found", fmt.Sprintf("%s not found", resource))
}

func UnauthorizedError(message string) events.APIGatewayProxyResponse {
	return ErrorWithDefault(401, "unauthorized", message)
}

func ForbiddenError(message string) events.APIGatewayProxyResponse {
	return ErrorWithDefault(403, "forbidden", message)
}

func BadRequestError(message string) events.APIGatewayProxyResponse {
	return ErrorWithDefault(400, "bad_request", message)
}

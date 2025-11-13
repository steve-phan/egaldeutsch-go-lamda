package response

import (
	"encoding/json"

	"egaldeutsch-serverless/pkg/middleware"

	"github.com/aws/aws-lambda-go/events"
)

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response structure
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// JSON creates a JSON response with the given status code and body
func JSON(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return ErrorJSON(500, "Failed to encode response", err.Error()), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    middleware.GetCORSHeaders(),
		Body:       string(jsonBody),
	}, nil
}

// SuccessJSON creates a successful JSON response
func SuccessJSON(statusCode int, data interface{}, message string) (events.APIGatewayProxyResponse, error) {
	return JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorJSON creates an error JSON response
func ErrorJSON(statusCode int, error string, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(ErrorResponse{
		Error:   error,
		Message: message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    middleware.GetCORSHeaders(),
		Body:       string(body),
	}
}

// SimpleError creates a simple error response with just an error message
func SimpleError(statusCode int, message string) events.APIGatewayProxyResponse {
	return ErrorJSON(statusCode, message, "")
}

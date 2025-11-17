package middleware

import "github.com/aws/aws-lambda-go/events"

type CORSConfig struct {
	AllowOrigin      string
	AllowHeaders     string
	AllowMethods     string
	AllowCredentials string
	ContentType      string
}

var (
	PublicAPI = CORSConfig{
		AllowOrigin:      "*",
		AllowHeaders:     "Content-Type",
		AllowMethods:     "GET, POST, OPTIONS",
		AllowCredentials: "true",
		ContentType:      "application/json",
	}

	AuthenticatedAPI = CORSConfig{
		AllowOrigin:      "*",
		AllowHeaders:     "Content-Type, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: "true",
		ContentType:      "application/json",
	}

	ReadOnlyAPI = CORSConfig{
		AllowOrigin:      "*",
		AllowHeaders:     "Content-Type",
		AllowMethods:     "GET, OPTIONS",
		AllowCredentials: "true",
		ContentType:      "application/json",
	}
)

// GetCORSHeadersWithConfig returns CORS headers for a specific configuration
func GetCORSHeadersWithConfig(config CORSConfig) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      config.AllowOrigin,
		"Access-Control-Allow-Headers":     config.AllowHeaders,
		"Access-Control-Allow-Methods":     config.AllowMethods,
		"Access-Control-Allow-Credentials": config.AllowCredentials,
		"Content-Type":                     config.ContentType,
	}
}

// GetPublicCORSHeaders returns CORS headers for public APIs
func GetPublicCORSHeaders() map[string]string {
	return GetCORSHeadersWithConfig(PublicAPI)
}

// GetAuthenticatedCORSHeaders returns CORS headers for authenticated APIs
func GetAuthenticatedCORSHeaders() map[string]string {
	return GetCORSHeadersWithConfig(AuthenticatedAPI)
}

// GetReadOnlyCORSHeaders returns CORS headers for read-only APIs
func GetReadOnlyCORSHeaders() map[string]string {
	return GetCORSHeadersWithConfig(ReadOnlyAPI)
}

// HandleCORS handles OPTIONS preflight requests
func HandleCORSWithConfig(request events.APIGatewayProxyRequest, config CORSConfig) (events.APIGatewayProxyResponse, bool) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    GetCORSHeadersWithConfig(config),
		}, true
	}
	return events.APIGatewayProxyResponse{}, false
}

// HandlePublicCORS handles OPTIONS for public APIs
func HandlePublicCORS(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, bool) {
	return HandleCORSWithConfig(request, PublicAPI)
}

// HandleAuthenticatedCORS handles OPTIONS for authenticated APIs
func HandleAuthenticatedCORS(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, bool) {
	return HandleCORSWithConfig(request, AuthenticatedAPI)
}

// HandleReadOnlyCORS handles OPTIONS for read-only APIs
func HandleReadOnlyCORS(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, bool) {
	return HandleCORSWithConfig(request, ReadOnlyAPI)
}

// HandleCORSOptions is a convenience function that returns an OPTIONS response
func HandleCORSOptions(config CORSConfig) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    GetCORSHeadersWithConfig(config),
	}
}

// GetCORSHeaders returns standard CORS headers for API responses
func GetCORSHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		"Access-Control-Allow-Credentials": "true",
		"Content-Type":                     "application/json",
	}
}

// HandleCORS handles OPTIONS preflight requests
func HandleCORS(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, bool) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    GetCORSHeaders(),
		}, true
	}
	return events.APIGatewayProxyResponse{}, false
}

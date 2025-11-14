package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvInDevelopment loads .env file in development environment
// This should be called at the beginning of each function handler
func LoadEnvInDevelopment() {
	// Only load .env in development
	if os.Getenv("NODE_ENV") != "production" {
		// Try multiple paths for .env file
		paths := []string{
			"../../../.env", // For Netlify functions (3 levels up)
			"../../.env",    // Alternative path
			".env",          // Current directory
		}

		var loaded bool
		for _, path := range paths {
			if err := godotenv.Load(path); err == nil {
				log.Printf("Loaded environment variables from %s", path)
				loaded = true
				break
			}
		}

		if !loaded {
			log.Printf("No .env file found (this is normal in production)")
		}
	}
}

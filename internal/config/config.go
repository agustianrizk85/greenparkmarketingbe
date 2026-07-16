// Package config loads runtime configuration from environment variables with
// sensible defaults, so the service runs out of the box for local development.
package config

import (
	"os"
	"time"
)

// Config holds the server runtime configuration.
type Config struct {
	Port        string        // HTTP port to listen on
	AllowOrigin string        // CORS allowed origin
	DataPath    string        // JSON file the master data is persisted to (when no DB)
	DatabaseURL string        // PostgreSQL DSN; when set, used instead of the JSON file
	SessionTTL  time.Duration // bearer-token session lifetime (legacy; unused under SSO)

	// Master auth (SSO): verify access tokens via the central auth public key.
	AuthJWKSURL string
	AuthIssuer  string
}

// Load reads configuration from the environment, applying defaults.
func Load() Config {
	return Config{
		Port:        getenv("MARKETING_PORT", "8086"),
		AllowOrigin: getenv("MARKETING_ALLOW_ORIGIN", "*"),
		DataPath:    getenv("MARKETING_DATA_PATH", "data/marketing-data.json"),
		DatabaseURL: getenv("MARKETING_DATABASE_URL", ""),
		SessionTTL:  12 * time.Hour,
		AuthJWKSURL: getenv("AUTH_JWKS_URL", "http://localhost:8090/.well-known/jwks.json"),
		AuthIssuer:  getenv("AUTH_ISSUER", "greenpark-auth"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

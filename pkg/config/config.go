package config

import (
	"os"
)

type Config struct {
	ApiPort    string
	AuthConfig *AuthConfig
}

type AuthConfig struct {
	URL string
}

func LoadConfig() *Config {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "0.0.0.0:3011"
	}

	return &Config{
		ApiPort: appPort,
		AuthConfig: &AuthConfig{
			URL: authServiceURL,
		},
	}
}

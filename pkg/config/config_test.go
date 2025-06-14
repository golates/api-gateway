package config

import (
	"github.com/joho/godotenv"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Load config with default values", func(t *testing.T) {
		testConfig := &Config{
			"3000",
			&AuthConfig{"0.0.0.0:3011"},
		}
		defaultConfig := LoadConfig()

		if defaultConfig.ApiPort != testConfig.ApiPort {
			t.Errorf("API Port %v not equal to %v", defaultConfig.ApiPort, testConfig.ApiPort)
		}

		if defaultConfig.AuthConfig.URL != testConfig.AuthConfig.URL {
			t.Errorf("Auth Service %v not equal to %v", defaultConfig.AuthConfig.URL, testConfig.AuthConfig.URL)
		}
	})
}

func TestLoadConfigTableDriven(t *testing.T) {
	var tests = []struct {
		filename string
		config   *Config
	}{
		{".sample_env", &Config{"3001", &AuthConfig{"0.0.0.0:3012"}}},
	}

	for _, test := range tests {
		err := godotenv.Load(test.filename)
		if err != nil {
			t.Errorf(err.Error())
		}

		config := LoadConfig()
		if config.ApiPort != test.config.ApiPort {
			t.Errorf("API Port %v not equal to %v", config.ApiPort, test.config.ApiPort)
		}

		if config.AuthConfig.URL != test.config.AuthConfig.URL {
			t.Errorf("Auth Service %v not equal to %v", config.AuthConfig.URL, test.config.AuthConfig.URL)
		}
	}
}

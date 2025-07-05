package config_test

import (
	"testing"

	"mcp-bridge/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate_BasicAuth_Success(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Auth: &config.AuthConfig{
					Type: "basic",
					Basic: &config.BasicAuthConfig{
						Username: "testuser",
						Password: "testpass",
					},
				},
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_Validate_BasicAuth_MissingType(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Auth: &config.AuthConfig{
					Type: "",
					Basic: &config.BasicAuthConfig{
						Username: "testuser",
						Password: "testpass",
					},
				},
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "auth type is required when auth is configured")
}

func TestConfig_Validate_BasicAuth_MissingConfig(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Auth: &config.AuthConfig{
					Type:  "basic",
					Basic: nil,
				},
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "basic auth configuration is required when type is 'basic'")
}

func TestConfig_Validate_BasicAuth_MissingUsername(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Auth: &config.AuthConfig{
					Type: "basic",
					Basic: &config.BasicAuthConfig{
						Username: "",
						Password: "testpass",
					},
				},
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "basic auth username is required")
}

func TestConfig_Validate_BasicAuth_MissingPassword(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Auth: &config.AuthConfig{
					Type: "basic",
					Basic: &config.BasicAuthConfig{
						Username: "testuser",
						Password: "",
					},
				},
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "basic auth password is required")
}

func TestConfig_Validate_UnsupportedAuthType(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Auth: &config.AuthConfig{
					Type: "oauth2",
				},
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported auth type 'oauth2'")
}

func TestConfig_Validate_NoAuth(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Auth:    nil, // No authentication
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_Validate_WithHeaders(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Headers: map[string]string{
					"X-API-Key":      "test-key",
					"Authorization":  "Bearer token",
					"X-Custom":       "custom-value",
				},
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "test-key", cfg.APIs[0].Headers["X-API-Key"])
	assert.Equal(t, "Bearer token", cfg.APIs[0].Headers["Authorization"])
	assert.Equal(t, "custom-value", cfg.APIs[0].Headers["X-Custom"])
}

func TestConfig_Validate_HeadersWithAuth(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Headers: map[string]string{
					"X-API-Key": "test-key",
				},
				Auth: &config.AuthConfig{
					Type: "basic",
					Basic: &config.BasicAuthConfig{
						Username: "testuser",
						Password: "testpass",
					},
				},
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "test-key", cfg.APIs[0].Headers["X-API-Key"])
}
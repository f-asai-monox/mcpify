package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"mcp-bridge/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_DefaultPath(t *testing.T) {
	cfg, err := config.LoadConfig("")
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "mcp-bridge", cfg.Server.Name)
	assert.Equal(t, "1.0.0", cfg.Server.Version)
	assert.Len(t, cfg.APIs, 1)
	assert.Equal(t, "default-api", cfg.APIs[0].Name)
}

func TestLoadConfig_ValidFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testConfig := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:3000",
				Timeout: 60,
			},
		},
		Server: config.ServerConfig{
			Name:        "test-server",
			Version:     "2.0.0",
			Description: "Test server",
		},
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	require.NoError(t, err)

	err = os.WriteFile(configPath, data, 0644)
	require.NoError(t, err)

	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, "test-api", cfg.APIs[0].Name)
	assert.Equal(t, "http://localhost:3000", cfg.APIs[0].BaseURL)
	assert.Equal(t, 60, cfg.APIs[0].Timeout)
	assert.Equal(t, "test-server", cfg.Server.Name)
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	err := os.WriteFile(configPath, []byte("invalid json"), 0644)
	require.NoError(t, err)

	cfg, err := config.LoadConfig(configPath)
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "error parsing config file")
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	cfg, err := config.LoadConfig("/nonexistent/path/config.json")
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "mcp-bridge", cfg.Server.Name)
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testConfig := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "save-test",
				BaseURL: "http://localhost:4000",
				Timeout: 45,
			},
		},
		Server: config.ServerConfig{
			Name:        "save-server",
			Version:     "3.0.0",
			Description: "Save test server",
		},
	}

	err := config.SaveConfig(testConfig, configPath)
	require.NoError(t, err)

	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var loadedConfig config.Config
	err = json.Unmarshal(data, &loadedConfig)
	require.NoError(t, err)

	assert.Equal(t, testConfig.APIs[0].Name, loadedConfig.APIs[0].Name)
	assert.Equal(t, testConfig.APIs[0].BaseURL, loadedConfig.APIs[0].BaseURL)
	assert.Equal(t, testConfig.Server.Name, loadedConfig.Server.Name)
}

func TestSaveConfig_CreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "subdir", "config.json")

	testConfig := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "default-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
			},
		},
		Server: config.ServerConfig{
			Name:        "mcp-bridge",
			Version:     "1.0.0",
			Description: "REST API to MCP Bridge Server",
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	err := config.SaveConfig(testConfig, configPath)
	require.NoError(t, err)

	assert.FileExists(t, configPath)
}

func TestValidate_Success(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "valid-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Endpoints: []config.CustomEndpoint{
					{
						Name:        "get-users",
						Description: "Get all users",
						Method:      "GET",
						Path:        "/users",
						Parameters: []config.CustomParameter{
							{
								Name:        "limit",
								Type:        "integer",
								Required:    false,
								Description: "Limit results",
								In:          "query",
							},
						},
					},
				},
			},
		},
		Server: config.ServerConfig{
			Name:        "test-server",
			Version:     "1.0.0",
			Description: "Test server",
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestValidate_NoAPIs(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one API configuration is required")
}

func TestValidate_MissingAPIName(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "",
				BaseURL: "http://localhost:8080",
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestValidate_MissingBaseURL(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "",
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "base URL is required")
}

func TestValidate_DefaultTimeout(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 0,
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
	assert.Equal(t, 30, cfg.APIs[0].Timeout)
}

func TestValidate_EndpointValidation(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Endpoints: []config.CustomEndpoint{
					{
						Name:        "",
						Description: "Test endpoint",
						Method:      "GET",
						Path:        "/test",
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
	assert.Contains(t, err.Error(), "endpoint 0: name is required")
}

func TestValidate_ParameterDefaults(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
				Endpoints: []config.CustomEndpoint{
					{
						Name:        "test-endpoint",
						Description: "Test endpoint",
						Method:      "GET",
						Path:        "/test",
						Parameters: []config.CustomParameter{
							{
								Name:        "param1",
								Required:    true,
								Description: "Test parameter",
							},
						},
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
	assert.Equal(t, "query", cfg.APIs[0].Endpoints[0].Parameters[0].In)
	assert.Equal(t, "string", cfg.APIs[0].Endpoints[0].Parameters[0].Type)
}

func TestValidate_ServerDefaults(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
			},
		},
		Server: config.ServerConfig{
			Name:    "",
			Version: "",
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "mcp-bridge", cfg.Server.Name)
	assert.Equal(t, "1.0.0", cfg.Server.Version)
}

func TestValidate_WithTransportConfig(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
		Transport: config.TransportConfig{
			Type: "http",
			HTTP: &config.HTTPTransportConfig{
				Host: "localhost",
				Port: 8080,
				CORS: true,
			},
		},
	}

	// Transport config is ignored during validation
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestValidate_WithoutTransportConfig(t *testing.T) {
	cfg := &config.Config{
		APIs: []config.APIConfig{
			{
				Name:    "test-api",
				BaseURL: "http://localhost:8080",
				Timeout: 30,
			},
		},
		Server: config.ServerConfig{
			Name:    "test-server",
			Version: "1.0.0",
		},
		// No transport config - should be fine
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}
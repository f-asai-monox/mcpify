package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	API      APIConfig               `json:"api"`
	Server   ServerConfig            `json:"server"`
	Headers  map[string]string       `json:"headers,omitempty"`
	Endpoints []CustomEndpoint       `json:"endpoints,omitempty"`
}

type APIConfig struct {
	BaseURL string `json:"baseUrl"`
	Timeout int    `json:"timeout"`
}

type ServerConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type CustomEndpoint struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Parameters  []CustomParameter      `json:"parameters"`
	Headers     map[string]string      `json:"headers,omitempty"`
}

type CustomParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	In          string      `json:"in"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return getDefaultConfig(), nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}
	
	return &config, nil
}

func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}
	
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	
	return nil
}

func getDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./mcp-bridge-config.json"
	}
	return filepath.Join(home, ".config", "mcp-bridge", "config.json")
}

func getDefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			BaseURL: "http://localhost:8080",
			Timeout: 30,
		},
		Server: ServerConfig{
			Name:        "mcp-bridge",
			Version:     "1.0.0",
			Description: "REST API to MCP Bridge Server",
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Endpoints: []CustomEndpoint{},
	}
}

func (c *Config) Validate() error {
	if c.API.BaseURL == "" {
		return fmt.Errorf("API base URL is required")
	}
	
	if c.API.Timeout <= 0 {
		c.API.Timeout = 30
	}
	
	if c.Server.Name == "" {
		c.Server.Name = "mcp-bridge"
	}
	
	if c.Server.Version == "" {
		c.Server.Version = "1.0.0"
	}
	
	for i, endpoint := range c.Endpoints {
		if endpoint.Name == "" {
			return fmt.Errorf("endpoint %d: name is required", i)
		}
		
		if endpoint.Method == "" {
			return fmt.Errorf("endpoint %s: method is required", endpoint.Name)
		}
		
		if endpoint.Path == "" {
			return fmt.Errorf("endpoint %s: path is required", endpoint.Name)
		}
		
		for j, param := range endpoint.Parameters {
			if param.Name == "" {
				return fmt.Errorf("endpoint %s, parameter %d: name is required", endpoint.Name, j)
			}
			
			if param.In == "" {
				param.In = "query"
			}
			
			if param.Type == "" {
				param.Type = "string"
			}
		}
	}
	
	return nil
}
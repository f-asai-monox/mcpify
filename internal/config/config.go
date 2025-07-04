package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIs     []APIConfig             `json:"apis"`
	Server   ServerConfig            `json:"server"`
	Headers  map[string]string       `json:"headers,omitempty"`
}

type APIConfig struct {
	Name      string            `json:"name"`
	BaseURL   string            `json:"baseUrl"`
	Timeout   int               `json:"timeout"`
	Endpoints []CustomEndpoint  `json:"endpoints,omitempty"`
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
		APIs: []APIConfig{
			{
				Name:      "default-api",
				BaseURL:   "http://localhost:8080",
				Timeout:   30,
				Endpoints: []CustomEndpoint{},
			},
		},
		Server: ServerConfig{
			Name:        "mcp-bridge",
			Version:     "1.0.0",
			Description: "REST API to MCP Bridge Server",
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func (c *Config) Validate() error {
	if len(c.APIs) == 0 {
		return fmt.Errorf("at least one API configuration is required")
	}
	
	for i, api := range c.APIs {
		if api.Name == "" {
			return fmt.Errorf("API %d: name is required", i)
		}
		
		if api.BaseURL == "" {
			return fmt.Errorf("API %s: base URL is required", api.Name)
		}
		
		if api.Timeout <= 0 {
			c.APIs[i].Timeout = 30
		}
		
		for j, endpoint := range api.Endpoints {
			if endpoint.Name == "" {
				return fmt.Errorf("API %s, endpoint %d: name is required", api.Name, j)
			}
			
			if endpoint.Method == "" {
				return fmt.Errorf("API %s, endpoint %s: method is required", api.Name, endpoint.Name)
			}
			
			if endpoint.Path == "" {
				return fmt.Errorf("API %s, endpoint %s: path is required", api.Name, endpoint.Name)
			}
			
			for k, param := range endpoint.Parameters {
				if param.Name == "" {
					return fmt.Errorf("API %s, endpoint %s, parameter %d: name is required", api.Name, endpoint.Name, k)
				}
				
				if param.In == "" {
					api.Endpoints[j].Parameters[k].In = "query"
				}
				
				if param.Type == "" {
					api.Endpoints[j].Parameters[k].Type = "string"
				}
			}
		}
	}
	
	if c.Server.Name == "" {
		c.Server.Name = "mcp-bridge"
	}
	
	if c.Server.Version == "" {
		c.Server.Version = "1.0.0"
	}
	
	return nil
}
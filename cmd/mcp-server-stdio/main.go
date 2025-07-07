package main

import (
	"flag"
	"log"
	"os"

	"mcp-bridge/internal/bridge"
	"mcp-bridge/internal/config"
	"mcp-bridge/internal/transport"
)

func main() {
	var (
		configPath = flag.String("config", "", "Path to configuration file")
		apiURL     = flag.String("api-url", "", "REST API base URL (overrides config)")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Always set log output to stderr to avoid interfering with MCP protocol
	log.SetOutput(os.Stderr)
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if *apiURL != "" && len(cfg.APIs) > 0 {
		cfg.APIs[0].BaseURL = *apiURL
	}

	// Stdio server always uses stdio transport, ignoring config file transport settings

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create stdio transport
	mcpTransport := transport.NewStdioTransport()
	mcpBridge := bridge.NewMCPBridge(mcpTransport)

	for key, value := range cfg.Headers {
		mcpBridge.SetAPIHeader(key, value)
	}

	for _, api := range cfg.APIs {
		for _, endpoint := range api.Endpoints {
			// Merge API-level headers with endpoint-level headers
			mergedHeaders := make(map[string]string)
			// Add API-level headers first
			for key, value := range api.Headers {
				mergedHeaders[key] = value
			}
			// Override with endpoint-level headers
			for key, value := range endpoint.Headers {
				mergedHeaders[key] = value
			}

			apiEndpoint := bridge.APIEndpoint{
				Name:        api.Name + "__" + endpoint.Name,
				Description: endpoint.Description,
				Method:      endpoint.Method,
				Path:        endpoint.Path,
				Headers:     mergedHeaders,
				Parameters:  make([]bridge.APIParameter, len(endpoint.Parameters)),
				APIName:     api.Name,
				BaseURL:     api.BaseURL,
				Auth:        api.Auth,
			}

			for i, param := range endpoint.Parameters {
				apiEndpoint.Parameters[i] = bridge.APIParameter{
					Name:        param.Name,
					Type:        param.Type,
					Required:    param.Required,
					Description: param.Description,
					Default:     param.Default,
					In:          param.In,
				}
			}

			mcpBridge.AddCustomEndpoint(apiEndpoint)
		}
	}

	if err := mcpBridge.Start(); err != nil {
		log.Fatalf("Error starting MCP bridge: %v", err)
	}
}
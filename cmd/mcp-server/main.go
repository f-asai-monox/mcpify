package main

import (
	"flag"
	"log"
	"os"

	"mcp-bridge/internal/bridge"
	"mcp-bridge/internal/config"
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

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	mcpBridge := bridge.NewMCPBridge()

	for key, value := range cfg.Headers {
		mcpBridge.SetAPIHeader(key, value)
	}

	for _, api := range cfg.APIs {
		for _, endpoint := range api.Endpoints {
			apiEndpoint := bridge.APIEndpoint{
				Name:        api.Name + "__" + endpoint.Name,
				Description: endpoint.Description,
				Method:      endpoint.Method,
				Path:        endpoint.Path,
				Headers:     endpoint.Headers,
				Parameters:  make([]bridge.APIParameter, len(endpoint.Parameters)),
				APIName:     api.Name,
				BaseURL:     api.BaseURL,
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

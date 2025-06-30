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
	
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(os.Stderr)
	}
	
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	
	if *apiURL != "" {
		cfg.API.BaseURL = *apiURL
	}
	
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}
	
	log.Printf("Starting MCP Bridge Server...")
	log.Printf("API Base URL: %s", cfg.API.BaseURL)
	log.Printf("Server: %s v%s", cfg.Server.Name, cfg.Server.Version)
	
	mcpBridge := bridge.NewMCPBridge(cfg.API.BaseURL)
	
	for key, value := range cfg.Headers {
		mcpBridge.SetAPIHeader(key, value)
	}
	
	for _, endpoint := range cfg.Endpoints {
		apiEndpoint := bridge.APIEndpoint{
			Name:        endpoint.Name,
			Description: endpoint.Description,
			Method:      endpoint.Method,
			Path:        endpoint.Path,
			Headers:     endpoint.Headers,
			Parameters:  make([]bridge.APIParameter, len(endpoint.Parameters)),
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
	
	log.Printf("MCP Bridge Server ready - listening on stdin/stdout")
	
	if err := mcpBridge.Start(); err != nil {
		log.Fatalf("Error starting MCP bridge: %v", err)
	}
}
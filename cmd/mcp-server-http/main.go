package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mcp-bridge/internal/bridge"
	"mcp-bridge/internal/config"
	"mcp-bridge/internal/transport"
)

func main() {
	var (
		configPath = flag.String("config", "", "Path to configuration file")
		apiURL     = flag.String("api-url", "", "REST API base URL (overrides config)")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
		httpHost   = flag.String("host", "localhost", "HTTP host")
		httpPort   = flag.Int("port", 8080, "HTTP port")
		httpCORS   = flag.Bool("cors", true, "Enable CORS for HTTP transport")
	)
	flag.Parse()

	// Set log output to stderr
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

	// HTTP server uses command line args, ignoring config file transport settings

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create HTTP transport with command line settings
	httpConfig := &transport.HTTPConfig{
		Host: *httpHost,
		Port: *httpPort,
		CORS: *httpCORS,
	}
	mcpTransport := transport.NewHTTPTransport(httpConfig)
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

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting HTTP MCP server on %s:%d", *httpHost, *httpPort)
		if err := mcpBridge.Start(); err != nil {
			log.Fatalf("Error starting MCP bridge: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down HTTP MCP server...")
	
	// Close the transport
	if err := mcpTransport.Close(); err != nil {
		log.Printf("Error closing transport: %v", err)
	}
}
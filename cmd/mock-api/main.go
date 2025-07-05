package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type MockConfig struct {
	Server    ServerConfig     `json:"server"`
	Auth      AuthConfig       `json:"auth"`
	Resources []ResourceConfig `json:"resources"`
	Endpoints []EndpointConfig `json:"endpoints"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

type AuthConfig struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResourceConfig struct {
	Name       string                   `json:"name"`
	Path       string                   `json:"path"`
	Enabled    bool                     `json:"enabled"`
	Data       []map[string]interface{} `json:"data"`
	Methods    []string                 `json:"methods"`
	SupportsID bool                     `json:"supportsId"`
}

type EndpointConfig struct {
	Path     string                 `json:"path"`
	Method   string                 `json:"method"`
	Enabled  bool                   `json:"enabled"`
	Response map[string]interface{} `json:"response"`
}

var (
	config       *MockConfig
	authEnabled  bool
	authUsername string
	authPassword string
)

func loadConfig(configPath string) (*MockConfig, error) {
	if configPath == "" {
		configPath = "configs/mock/users.json"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg MockConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &cfg, nil
}

func basicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authEnabled {
			next(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Mock API"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth, "Basic ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		encoded := strings.TrimPrefix(auth, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusUnauthorized)
			return
		}

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			http.Error(w, "Invalid credentials format", http.StatusUnauthorized)
			return
		}

		username, password := creds[0], creds[1]
		if username != authUsername || password != authPassword {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func corsHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func createResourceHandler(resource *ResourceConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if !contains(resource.Methods, r.Method) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		switch r.Method {
		case "GET":
			if err := json.NewEncoder(w).Encode(resource.Data); err != nil {
				http.Error(w, "Error encoding data", http.StatusInternalServerError)
			}
		case "POST":
			var newItem map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			nextID := len(resource.Data) + 1
			newItem["id"] = nextID
			newItem["created"] = time.Now().Format(time.RFC3339)
			resource.Data = append(resource.Data, newItem)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newItem)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func createResourceIDHandler(resource *ResourceConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if !resource.SupportsID {
			http.Error(w, "ID operations not supported", http.StatusNotFound)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, resource.Path+"/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var item map[string]interface{}
		var index int = -1
		for i, data := range resource.Data {
			if itemID, ok := data["id"].(float64); ok && int(itemID) == id {
				item = data
				index = i
				break
			} else if itemID, ok := data["id"].(int); ok && itemID == id {
				item = data
				index = i
				break
			}
		}

		if item == nil {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			json.NewEncoder(w).Encode(item)
		case "PUT":
			if !contains(resource.Methods, "PUT") {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var updatedItem map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			updatedItem["id"] = id
			if created, exists := item["created"]; exists {
				updatedItem["created"] = created
			}

			resource.Data[index] = updatedItem
			json.NewEncoder(w).Encode(updatedItem)
		case "DELETE":
			if !contains(resource.Methods, "DELETE") {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			resource.Data = append(resource.Data[:index], resource.Data[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func createMultiMethodEndpointHandler(endpoints []*EndpointConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, endpoint := range endpoints {
			if r.Method == endpoint.Method {
				createEndpointHandler(endpoint)(w, r)
				return
			}
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createEndpointHandler(endpoint *EndpointConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != endpoint.Method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response := make(map[string]interface{})
		for key, value := range endpoint.Response {
			if str, ok := value.(string); ok && str == "{{timestamp}}" {
				response[key] = time.Now().Format(time.RFC3339)
			} else {
				response[key] = value
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	// Load configuration
	configPath := os.Getenv("MOCK_CONFIG")
	var err error
	config, err = loadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Override port from environment variable
	port := os.Getenv("PORT")
	if port != "" {
		config.Server.Port = port
	}

	// Override auth from environment variables
	if os.Getenv("AUTH_ENABLED") == "true" {
		config.Auth.Enabled = true
		if username := os.Getenv("AUTH_USERNAME"); username != "" {
			config.Auth.Username = username
		}
		if password := os.Getenv("AUTH_PASSWORD"); password != "" {
			config.Auth.Password = password
		}
	}

	// Set auth globals
	authEnabled = config.Auth.Enabled
	authUsername = config.Auth.Username
	authPassword = config.Auth.Password

	// Register endpoint handlers grouped by path
	endpointsByPath := make(map[string][]*EndpointConfig)
	for i := range config.Endpoints {
		endpoint := &config.Endpoints[i]
		if endpoint.Enabled {
			endpointsByPath[endpoint.Path] = append(endpointsByPath[endpoint.Path], endpoint)
		}
	}
	
	for path, endpoints := range endpointsByPath {
		http.HandleFunc(path, corsHandler(basicAuthMiddleware(createMultiMethodEndpointHandler(endpoints))))
	}

	// Register resource handlers
	for i := range config.Resources {
		resource := &config.Resources[i]
		if resource.Enabled {
			http.HandleFunc(resource.Path, corsHandler(basicAuthMiddleware(createResourceHandler(resource))))
			if resource.SupportsID {
				http.HandleFunc(resource.Path+"/", corsHandler(basicAuthMiddleware(createResourceIDHandler(resource))))
			}
		}
	}

	// Display server information
	fmt.Printf("%s starting on port %s...\n", config.Server.Name, config.Server.Port)
	if configPath := os.Getenv("MOCK_CONFIG"); configPath != "" {
		fmt.Printf("Using config file: %s\n", configPath)
	} else {
		fmt.Printf("Using default config file: configs/mock/users.json\n")
	}
	if authEnabled {
		fmt.Printf("Basic Authentication: ENABLED (username: %s)\n", authUsername)
	} else {
		fmt.Println("Basic Authentication: DISABLED")
	}

	fmt.Println("Available endpoints:")
	for _, endpoint := range config.Endpoints {
		if endpoint.Enabled {
			fmt.Printf("  %s    %s\n", endpoint.Method, endpoint.Path)
		}
	}

	for _, resource := range config.Resources {
		if resource.Enabled {
			for _, method := range resource.Methods {
				fmt.Printf("  %s    %s\n", method, resource.Path)
				if resource.SupportsID && (method == "GET" || method == "PUT" || method == "DELETE") {
					fmt.Printf("  %s    %s/{id}\n", method, resource.Path)
				}
			}
		}
	}

	if authEnabled {
		fmt.Println("\nTo test with authentication:")
		fmt.Printf("  curl -u %s:%s http://localhost:%s/users\n", authUsername, authPassword, config.Server.Port)
	}

	if err := http.ListenAndServe(":"+config.Server.Port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

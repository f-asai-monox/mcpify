package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"mcp-bridge/pkg/types"
)

// HTTPTransport implements the Transport interface for HTTP communication
type HTTPTransport struct {
	config     *HTTPConfig
	server     *http.Server
	messageCh  chan *types.JSONRPCMessage
	responseCh chan *types.JSONRPCMessage
	closed     bool
	mu         sync.RWMutex
	wg         sync.WaitGroup
}

// NewHTTPTransport creates a new HTTP transport
func NewHTTPTransport(config *HTTPConfig) *HTTPTransport {
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 8080
	}

	return &HTTPTransport{
		config:     config,
		messageCh:  make(chan *types.JSONRPCMessage, 100),
		responseCh: make(chan *types.JSONRPCMessage, 100),
		closed:     false,
	}
}

// Start begins listening for HTTP requests
func (t *HTTPTransport) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", t.handleMCPRequest)
	if t.config.CORS {
		mux.HandleFunc("/", t.handleCORS)
	}

	t.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", t.config.Host, t.config.Port),
		Handler: mux,
	}

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	return nil
}

// ReadMessage reads a single JSON-RPC message from the HTTP transport
func (t *HTTPTransport) ReadMessage() (*types.JSONRPCMessage, error) {
	t.mu.RLock()
	closed := t.closed
	t.mu.RUnlock()

	if closed {
		return nil, fmt.Errorf("transport is closed")
	}

	select {
	case msg := <-t.messageCh:
		return msg, nil
	case <-time.After(100 * time.Millisecond):
		return nil, nil // No message available, return nil to continue polling
	}
}

// WriteMessage writes a JSON-RPC message to the HTTP transport
func (t *HTTPTransport) WriteMessage(msg *types.JSONRPCMessage) error {
	t.mu.RLock()
	closed := t.closed
	t.mu.RUnlock()

	if closed {
		return fmt.Errorf("transport is closed")
	}

	select {
	case t.responseCh <- msg:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout writing response")
	}
}

// Close closes the HTTP transport
func (t *HTTPTransport) Close() error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true
	t.mu.Unlock()

	close(t.messageCh)
	close(t.responseCh)

	if t.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := t.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("error shutting down HTTP server: %w", err)
		}
	}

	t.wg.Wait()
	return nil
}

// handleMCPRequest handles incoming MCP JSON-RPC requests
func (t *HTTPTransport) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	if t.config.CORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var msg types.JSONRPCMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Invalid JSON-RPC message", http.StatusBadRequest)
		return
	}

	// Send the message to the message channel
	select {
	case t.messageCh <- &msg:
	case <-time.After(5 * time.Second):
		http.Error(w, "Timeout processing request", http.StatusRequestTimeout)
		return
	}

	// Wait for response
	select {
	case response := <-t.responseCh:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	case <-time.After(30 * time.Second):
		http.Error(w, "Timeout waiting for response", http.StatusRequestTimeout)
	}
}

// handleCORS handles CORS preflight requests
func (t *HTTPTransport) handleCORS(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.NotFound(w, r)
}
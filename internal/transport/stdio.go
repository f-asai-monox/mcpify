package transport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"mcp-bridge/pkg/types"
)

// StdioTransport implements the Transport interface for stdin/stdout communication
type StdioTransport struct {
	reader *bufio.Scanner
	writer io.Writer
	closed bool
}

// NewStdioTransport creates a new stdio transport
func NewStdioTransport() *StdioTransport {
	return &StdioTransport{
		reader: bufio.NewScanner(os.Stdin),
		writer: os.Stdout,
		closed: false,
	}
}

// Start begins listening for messages from stdin
func (t *StdioTransport) Start() error {
	// For stdio transport, Start() doesn't need to do anything special
	// The actual message reading happens in ReadMessage()
	return nil
}

// ReadMessage reads a single JSON-RPC message from stdin
func (t *StdioTransport) ReadMessage() (*types.JSONRPCMessage, error) {
	if t.closed {
		return nil, fmt.Errorf("transport is closed")
	}

	for t.reader.Scan() {
		line := strings.TrimSpace(t.reader.Text())
		if line == "" {
			continue
		}

		var msg types.JSONRPCMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}

		return &msg, nil
	}

	if err := t.reader.Err(); err != nil {
		return nil, fmt.Errorf("error reading from stdin: %w", err)
	}

	// EOF reached
	return nil, io.EOF
}

// WriteMessage writes a JSON-RPC message to stdout
func (t *StdioTransport) WriteMessage(msg *types.JSONRPCMessage) error {
	if t.closed {
		return fmt.Errorf("transport is closed")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	_, err = fmt.Fprintf(t.writer, "%s\n", string(data))
	return err
}

// Close closes the transport
func (t *StdioTransport) Close() error {
	t.closed = true
	return nil
}
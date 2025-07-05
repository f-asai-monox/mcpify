package bridge_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mcp-bridge/internal/bridge"
	"mcp-bridge/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestClient_BasicAuth_Success(t *testing.T) {
	expectedUsername := "testuser"
	expectedPassword := "testpass"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Authorization header
		auth := r.Header.Get("Authorization")
		require.NotEmpty(t, auth, "Authorization header should be present")
		require.True(t, len(auth) > 6, "Authorization header should start with 'Basic '")
		
		// Decode and verify credentials
		encoded := auth[6:] // Remove "Basic " prefix
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		require.NoError(t, err, "Should be valid base64")
		
		credentials := string(decoded)
		expected := expectedUsername + ":" + expectedPassword
		assert.Equal(t, expected, credentials, "Credentials should match")
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "authenticated successfully",
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "test-auth",
		Description: "Test endpoint with Basic Auth",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     server.URL,
		Auth: &config.AuthConfig{
			Type: "basic",
			Basic: &config.BasicAuthConfig{
				Username: expectedUsername,
				Password: expectedPassword,
			},
		},
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Data)
}

func TestRestClient_BasicAuth_MissingCredentials(t *testing.T) {
	client := bridge.NewRestClient()
	
	// Test with nil Basic config
	endpoint := bridge.APIEndpoint{
		Name:        "test-auth",
		Description: "Test endpoint with invalid Basic Auth",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     "http://localhost:8080",
		Auth: &config.AuthConfig{
			Type:  "basic",
			Basic: nil,
		},
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "basic auth configuration is nil")
}

func TestRestClient_UnsupportedAuthType(t *testing.T) {
	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "test-auth",
		Description: "Test endpoint with unsupported auth",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     "http://localhost:8080",
		Auth: &config.AuthConfig{
			Type: "unsupported",
		},
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unsupported authentication type: unsupported")
}

func TestRestClient_NoAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify no Authorization header
		auth := r.Header.Get("Authorization")
		assert.Empty(t, auth, "Authorization header should not be present")
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "no auth required",
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "test-no-auth",
		Description: "Test endpoint without auth",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     server.URL,
		Auth:        nil, // No authentication
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Data)
}

func TestRestClient_BasicAuth_WithOtherHeaders(t *testing.T) {
	expectedUsername := "user"
	expectedPassword := "pass"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Authorization header
		auth := r.Header.Get("Authorization")
		require.NotEmpty(t, auth)
		
		// Verify other headers are preserved
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom-Header"))
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	client.SetHeader("Content-Type", "application/json")
	
	endpoint := bridge.APIEndpoint{
		Name:        "test-auth-headers",
		Description: "Test endpoint with auth and custom headers",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     server.URL,
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
		},
		Auth: &config.AuthConfig{
			Type: "basic",
			Basic: &config.BasicAuthConfig{
				Username: expectedUsername,
				Password: expectedPassword,
			},
		},
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
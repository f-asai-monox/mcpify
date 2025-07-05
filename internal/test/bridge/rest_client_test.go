package bridge_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mcp-bridge/internal/bridge"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRestClient(t *testing.T) {
	client := bridge.NewRestClient()
	assert.NotNil(t, client)
}

func TestRestClient_SetHeader(t *testing.T) {
	client := bridge.NewRestClient()
	client.SetHeader("Authorization", "Bearer token123")
}

func TestRestClient_MakeRequest_GET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/users", r.URL.Path)
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": []string{"user1", "user2"},
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "get-users",
		Description: "Get users",
		Method:      "GET",
		Path:        "/users",
		BaseURL:     server.URL,
		Parameters: []bridge.APIParameter{
			{
				Name:        "limit",
				Type:        "integer",
				Required:    false,
				Description: "Limit results",
				In:          "query",
			},
		},
	}

	args := map[string]interface{}{
		"limit": "10",
	}

	resp, err := client.MakeRequest(endpoint, args)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Headers["Content-Type"])
	assert.NotNil(t, resp.Data)
	assert.Empty(t, resp.Error)
}

func TestRestClient_MakeRequest_POST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/users", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "John Doe", body["name"])
		assert.Equal(t, "john@example.com", body["email"])
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   1,
			"name": "John Doe",
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "create-user",
		Description: "Create user",
		Method:      "POST",
		Path:        "/users",
		BaseURL:     server.URL,
		Parameters: []bridge.APIParameter{
			{
				Name:        "name",
				Type:        "string",
				Required:    true,
				Description: "User name",
				In:          "body",
			},
			{
				Name:        "email",
				Type:        "string",
				Required:    true,
				Description: "User email",
				In:          "body",
			},
		},
	}

	args := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	resp, err := client.MakeRequest(endpoint, args)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, resp.Data)
}

func TestRestClient_MakeRequest_PathParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/users/123", r.URL.Path)
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   123,
			"name": "John Doe",
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "get-user",
		Description: "Get user by ID",
		Method:      "GET",
		Path:        "/users/{id}",
		BaseURL:     server.URL,
		Parameters: []bridge.APIParameter{
			{
				Name:        "id",
				Type:        "integer",
				Required:    true,
				Description: "User ID",
				In:          "path",
			},
		},
	}

	args := map[string]interface{}{
		"id": "123",
	}

	resp, err := client.MakeRequest(endpoint, args)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Data)
}

func TestRestClient_MakeRequest_HeaderParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "success",
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "protected-endpoint",
		Description: "Protected endpoint",
		Method:      "GET",
		Path:        "/protected",
		BaseURL:     server.URL,
		Parameters: []bridge.APIParameter{
			{
				Name:        "Authorization",
				Type:        "string",
				Required:    true,
				Description: "Auth token",
				In:          "header",
			},
			{
				Name:        "Accept",
				Type:        "string",
				Required:    false,
				Description: "Accept header",
				In:          "header",
			},
		},
	}

	args := map[string]interface{}{
		"Authorization": "Bearer token123",
		"Accept":        "application/json",
	}

	resp, err := client.MakeRequest(endpoint, args)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_MakeRequest_RequiredParameterMissing(t *testing.T) {
	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "get-user",
		Description: "Get user by ID",
		Method:      "GET",
		Path:        "/users/{id}",
		BaseURL:     "http://localhost:8080",
		Parameters: []bridge.APIParameter{
			{
				Name:        "id",
				Type:        "integer",
				Required:    true,
				Description: "User ID",
				In:          "path",
			},
		},
	}

	args := map[string]interface{}{}

	resp, err := client.MakeRequest(endpoint, args)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "required parameter 'id' is missing")
}

func TestRestClient_MakeRequest_DefaultParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": []string{"user1"},
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "get-users",
		Description: "Get users",
		Method:      "GET",
		Path:        "/users",
		BaseURL:     server.URL,
		Parameters: []bridge.APIParameter{
			{
				Name:        "limit",
				Type:        "integer",
				Required:    false,
				Description: "Limit results",
				In:          "query",
				Default:     10,
			},
		},
	}

	args := map[string]interface{}{}

	resp, err := client.MakeRequest(endpoint, args)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_MakeRequest_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "get-user",
		Description: "Get user by ID",
		Method:      "GET",
		Path:        "/users/999",
		BaseURL:     server.URL,
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Contains(t, resp.Error, "HTTP 404")
	assert.Contains(t, resp.Error, "User not found")
}

func TestRestClient_MakeRequest_CustomBaseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "success",
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "test-endpoint",
		Description: "Test endpoint",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     server.URL,
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_MakeRequest_MissingBaseURL(t *testing.T) {
	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "test-endpoint",
		Description: "Test endpoint",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     "",
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "endpoint BaseURL is required")
}

func TestRestClient_MakeRequest_WithAPIHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom-Header"))
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "success",
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "api-with-headers",
		Description: "API with custom headers",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     server.URL,
		Headers: map[string]string{
			"X-API-Key":      "test-api-key",
			"Authorization":  "Bearer test-token",
			"X-Custom-Header": "custom-value",
		},
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_MakeRequest_HeaderPrecedence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Endpoint-level header should override API-level header
		assert.Equal(t, "endpoint-value", r.Header.Get("X-Override"))
		assert.Equal(t, "api-value", r.Header.Get("X-API-Only"))
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "success",
		})
	}))
	defer server.Close()

	client := bridge.NewRestClient()
	endpoint := bridge.APIEndpoint{
		Name:        "header-precedence-test",
		Description: "Test header precedence",
		Method:      "GET",
		Path:        "/test",
		BaseURL:     server.URL,
		Headers: map[string]string{
			"X-Override":  "endpoint-value", // This should override API-level
			"X-API-Only":  "api-value",     // This should remain
		},
	}

	resp, err := client.MakeRequest(endpoint, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
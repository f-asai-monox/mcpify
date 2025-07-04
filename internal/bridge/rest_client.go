package bridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RestClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

type APIEndpoint struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Parameters  []APIParameter    `json:"parameters"`
	Headers     map[string]string `json:"headers"`
	APIName     string            `json:"apiName"`
	BaseURL     string            `json:"baseUrl"`
}

type APIParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	In          string      `json:"in"`
}

type APIResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Data       interface{}       `json:"data,omitempty"`
	Error      string            `json:"error,omitempty"`
}

func NewRestClient(baseURL string) *RestClient {
	return &RestClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

func (c *RestClient) SetHeader(key, value string) {
	c.headers[key] = value
}

func (c *RestClient) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

func (c *RestClient) MakeRequest(endpoint APIEndpoint, args map[string]interface{}) (*APIResponse, error) {
	baseURL := c.baseURL
	if endpoint.BaseURL != "" {
		baseURL = endpoint.BaseURL
	}

	fullURL, err := c.buildURLWithBase(endpoint, args, baseURL)
	if err != nil {
		return nil, fmt.Errorf("error building URL: %w", err)
	}

	var reqBody io.Reader
	if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "PATCH" {
		bodyData := c.extractBodyData(endpoint, args)
		if bodyData != nil {
			jsonData, err := json.Marshal(bodyData)
			if err != nil {
				return nil, fmt.Errorf("error marshaling request body: %w", err)
			}
			reqBody = bytes.NewBuffer(jsonData)
		}
	}

	req, err := http.NewRequest(endpoint.Method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, param := range endpoint.Parameters {
		if param.In == "header" {
			if value, exists := args[param.Name]; exists {
				req.Header.Set(param.Name, fmt.Sprintf("%v", value))
			}
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	apiResp := &APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    responseHeaders,
		Body:       string(body),
	}

	if resp.StatusCode >= 400 {
		apiResp.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))
	} else {
		var data interface{}
		if err := json.Unmarshal(body, &data); err == nil {
			apiResp.Data = data
		}
	}

	return apiResp, nil
}

func (c *RestClient) buildURLWithBase(endpoint APIEndpoint, args map[string]interface{}, baseURL string) (string, error) {
	path := endpoint.Path
	queryParams := url.Values{}

	for _, param := range endpoint.Parameters {
		value, exists := args[param.Name]
		if !exists {
			if param.Required && param.Default == nil {
				return "", fmt.Errorf("required parameter '%s' is missing", param.Name)
			}
			if param.Default != nil {
				value = param.Default
			}
		}

		if value == nil {
			continue
		}

		switch param.In {
		case "path":
			placeholder := "{" + param.Name + "}"
			path = strings.ReplaceAll(path, placeholder, fmt.Sprintf("%v", value))
		case "query":
			queryParams.Add(param.Name, fmt.Sprintf("%v", value))
		}
	}

	fullURL := strings.TrimSuffix(baseURL, "/") + path
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	return fullURL, nil
}

func (c *RestClient) extractBodyData(endpoint APIEndpoint, args map[string]interface{}) map[string]interface{} {
	bodyData := make(map[string]interface{})

	for _, param := range endpoint.Parameters {
		if param.In == "body" {
			value, exists := args[param.Name]
			if exists {
				bodyData[param.Name] = value
			} else if param.Required && param.Default != nil {
				bodyData[param.Name] = param.Default
			}
		}
	}

	if len(bodyData) == 0 {
		for key, value := range args {
			isPathOrQuery := false
			for _, param := range endpoint.Parameters {
				if param.Name == key && (param.In == "path" || param.In == "query" || param.In == "header") {
					isPathOrQuery = true
					break
				}
			}
			if !isPathOrQuery {
				bodyData[key] = value
			}
		}
	}

	if len(bodyData) == 0 {
		return nil
	}

	return bodyData
}

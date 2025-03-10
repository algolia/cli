package genai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	// DefaultBaseURL is the default base URL for the Algolia GenAI API
	DefaultBaseURL = "https://generative-ai.algolia.com"
)

// Client provides methods to interact with the Algolia GenAI API
type Client struct {
	AppID  string
	APIKey string
	client *http.Client
}

// NewClient returns a new GenAI API client
func NewClient(appID, apiKey string) *Client {
	return &Client{
		AppID:  appID,
		APIKey: apiKey,
		client: http.DefaultClient,
	}
}

// NewClientWithHTTPClient returns a new GenAI API client with a custom HTTP client
func NewClientWithHTTPClient(appID, apiKey string, client *http.Client) *Client {
	return &Client{
		AppID:  appID,
		APIKey: apiKey,
		client: client,
	}
}

func (c *Client) request(res interface{}, method string, path string, body interface{}, urlParams map[string]string) error {
	r, err := c.buildRequest(method, path, body, urlParams)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return &APIError{
				StatusCode: resp.StatusCode,
				Method:     method,
				Path:       path,
				Response:   errResp,
			}
		}

		// Fallback if error response couldn't be parsed
		return &APIError{
			StatusCode: resp.StatusCode,
			Method:     method,
			Path:       path,
			Response: ErrorResponse{
				Message: "Error accessing API",
			},
		}
	}

	if res != nil {
		if err := unmarshalTo(resp, res); err != nil {
			return err
		}
	}

	return nil
}

// buildRequest builds an HTTP request
func (c *Client) buildRequest(method, path string, body interface{}, urlParams map[string]string) (*http.Request, error) {
	path = strings.TrimSuffix(path, "/")
	url := DefaultBaseURL + path

	var req *http.Request
	var err error

	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Algolia-Application-Id", c.AppID)
	req.Header.Set("X-Algolia-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Add URL params
	values := req.URL.Query()
	for k, v := range urlParams {
		values.Set(k, v)
	}
	req.URL.RawQuery = values.Encode()

	return req, nil
}

// unmarshalTo unmarshals an HTTP response body to a given interface
func unmarshalTo(r *http.Response, v interface{}) error {
	// Don't close the body here as it's now closed in the request method
	return json.NewDecoder(r.Body).Decode(v)
}

// CreateDataSource creates a new data source
func (c *Client) CreateDataSource(input CreateDataSourceInput) (*DataSourceResponse, error) {
	var res DataSourceResponse
	err := c.request(&res, http.MethodPost, "/create/data_source", input, nil)
	return &res, err
}

// GetDataSource gets a data source by ID
func (c *Client) GetDataSource(objectID string) (*DataSourceDetails, error) {
	var res DataSourceDetails
	err := c.request(&res, http.MethodGet, fmt.Sprintf("/get/data_source/%s", objectID), nil, nil)
	return &res, err
}

// ListDataSources lists all data sources
// Note: Not supported by the API yet
func (c *Client) ListDataSources() (*ListDataSourcesResponse, error) {
	// The API doesn't have a list endpoint (yet)
	return nil, fmt.Errorf("listing data sources is not currently supported by the Algolia GenAI API")
}

// UpdateDataSource updates an existing data source
func (c *Client) UpdateDataSource(input UpdateDataSourceInput) (*DataSourceResponse, error) {
	var res DataSourceResponse
	err := c.request(&res, http.MethodPost, "/update/data_source", input, nil)
	return &res, err
}

// DeleteDataSources deletes one or more data sources
func (c *Client) DeleteDataSources(input DeleteDataSourcesInput) (*DeleteResponse, error) {
	var res DeleteResponse
	err := c.request(&res, http.MethodPost, "/delete/data_sources", input, nil)
	return &res, err
}

// CreatePrompt creates a new prompt
func (c *Client) CreatePrompt(input CreatePromptInput) (*PromptResponse, error) {
	var res PromptResponse
	err := c.request(&res, http.MethodPost, "/create/prompt", input, nil)
	return &res, err
}

// GetPrompt gets a prompt by ID
func (c *Client) GetPrompt(objectID string) (*PromptDetails, error) {
	var res PromptDetails
	err := c.request(&res, http.MethodGet, fmt.Sprintf("/get/prompt/%s", objectID), nil, nil)
	return &res, err
}

// ListPrompts lists all prompts
// Note: Not supported by the API yet
func (c *Client) ListPrompts() (*ListPromptsResponse, error) {
	// The API doesn't seem to have a list endpoint
	return nil, fmt.Errorf("listing prompts is not currently supported by the Algolia GenAI API")
}

// UpdatePrompt updates an existing prompt
func (c *Client) UpdatePrompt(input UpdatePromptInput) (*PromptResponse, error) {
	var res PromptResponse
	err := c.request(&res, http.MethodPost, "/update/prompt", input, nil)
	return &res, err
}

// DeletePrompts deletes one or more prompts
func (c *Client) DeletePrompts(input DeletePromptsInput) (*DeleteResponse, error) {
	var res DeleteResponse
	err := c.request(&res, http.MethodPost, "/delete/prompts", input, nil)
	return &res, err
}

// GenerateResponse generates a response using a prompt
func (c *Client) GenerateResponse(input GenerateResponseInput) (*GenerateResponseOutput, error) {
	var res GenerateResponseOutput
	err := c.request(&res, http.MethodPost, "/generate/response", input, nil)
	return &res, err
}

// ListResponses lists all responses
// Note: Not supported by the API yet
func (c *Client) ListResponses() (*ListResponsesResponse, error) {
	// The API doesn't seem to have a list endpoint
	return nil, fmt.Errorf("listing responses is not currently supported by the Algolia GenAI API")
}

// GetResponse retrieves a response by ID
func (c *Client) GetResponse(objectID string) (*ResponseDetails, error) {
	var res ResponseDetails
	err := c.request(&res, http.MethodGet, fmt.Sprintf("/get/response/%s", objectID), nil, nil)
	return &res, err
}

// DeleteResponses deletes one or more responses
func (c *Client) DeleteResponses(input DeleteResponsesInput) (*DeleteResponse, error) {
	var res DeleteResponse
	err := c.request(&res, http.MethodPost, "/delete/responses", input, nil)
	return &res, err
}

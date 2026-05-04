package agentstudio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Client provides methods to interact with the Algolia Agent Studio (RAG) API.
type Client struct {
	AppID   string
	APIKey  string
	BaseURL string

	client *http.Client
}

// DefaultBaseURL builds the standard per-app server URL declared by the spec.
func DefaultBaseURL(appID string) string {
	return fmt.Sprintf("https://%s.algolia.net/agent-studio/", appID)
}

// NewClient returns a new Agent Studio API client using the default per-app
// server URL.
func NewClient(appID, apiKey string) *Client {
	return &Client{
		AppID:   appID,
		APIKey:  apiKey,
		BaseURL: DefaultBaseURL(appID),
		client:  http.DefaultClient,
	}
}

// NewClientWithHTTPClient returns a new Agent Studio API client with a custom
// HTTP client. Tests use this to inject an httptest server.
func NewClientWithHTTPClient(appID, apiKey string, hc *http.Client) *Client {
	return &Client{
		AppID:   appID,
		APIKey:  apiKey,
		BaseURL: DefaultBaseURL(appID),
		client:  hc,
	}
}

// request sends an HTTP request and unmarshals the response body to res when
// non-nil. Status >= 400 is converted to a formatted error.
func (c *Client) request(
	res interface{},
	method, path string,
	body interface{},
	urlParams map[string]string,
) error {
	r, err := c.buildRequest(method, path, body, urlParams)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errResp ErrResponse
		_ = unmarshalTo(resp, &errResp)
		return fmt.Errorf("agentstudio: %s %s -> %d %s", method, path, resp.StatusCode, formatDetail(errResp.Detail))
	}

	if res != nil {
		if err := unmarshalTo(resp, res); err != nil {
			return err
		}
	} else {
		_ = resp.Body.Close()
	}

	return nil
}

func (c *Client) buildRequest(
	method, path string,
	body interface{},
	urlParams map[string]string,
) (*http.Request, error) {
	url := strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")

	var reader io.ReadCloser
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = io.NopCloser(bytes.NewReader(b))
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Algolia-Application-Id", c.AppID)
	req.Header.Set("X-Algolia-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if len(urlParams) > 0 {
		q := req.URL.Query()
		for k, v := range urlParams {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

func unmarshalTo(r *http.Response, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// formatDetail renders a FastAPI `detail` field for human consumption.
// If detail is a JSON string, return it; if it's an array, join entries; else
// return the raw JSON.
func formatDetail(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var entries []map[string]any
	if err := json.Unmarshal(raw, &entries); err == nil {
		var parts []string
		for _, e := range entries {
			parts = append(parts, fmt.Sprintf("%v", e))
		}
		return strings.Join(parts, "; ")
	}
	return string(raw)
}

func paginationParams(page, limit int) map[string]string {
	params := map[string]string{}
	if page > 0 {
		params["page"] = strconv.Itoa(page)
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	return params
}

// Agents -----------------------------------------------------------------

func (c *Client) ListAgents(page, limit int) (*PaginatedAgentsResponse, error) {
	var res PaginatedAgentsResponse
	if err := c.request(&res, http.MethodGet, "1/agents", nil, paginationParams(page, limit)); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetAgent(id string) (*Agent, error) {
	var res Agent
	if err := c.request(&res, http.MethodGet, fmt.Sprintf("1/agents/%s", id), nil, nil); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) CreateAgent(req AgentConfigCreate) (*Agent, error) {
	var res Agent
	if err := c.request(&res, http.MethodPost, "1/agents", req, nil); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) UpdateAgent(id string, req AgentConfigUpdate) (*Agent, error) {
	var res Agent
	if err := c.request(&res, http.MethodPatch, fmt.Sprintf("1/agents/%s", id), req, nil); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) DeleteAgent(id string) error {
	return c.request(nil, http.MethodDelete, fmt.Sprintf("1/agents/%s", id), nil, nil)
}

func (c *Client) DuplicateAgent(id string) (*Agent, error) {
	var res Agent
	if err := c.request(&res, http.MethodPost, fmt.Sprintf("1/agents/%s/duplicate", id), nil, nil); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) PublishAgent(id string) (*Agent, error) {
	var res Agent
	if err := c.request(&res, http.MethodPost, fmt.Sprintf("1/agents/%s/publish", id), nil, nil); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) UnpublishAgent(id string) (*Agent, error) {
	var res Agent
	if err := c.request(&res, http.MethodPost, fmt.Sprintf("1/agents/%s/unpublish", id), nil, nil); err != nil {
		return nil, err
	}
	return &res, nil
}

// Completions ------------------------------------------------------------

// CompletionParams carries the query-string options for CreateCompletion.
// CompatibilityMode is required by the API; the rest are optional and default
// to the API's defaults when zero-valued.
type CompletionParams struct {
	CompatibilityMode string // "ai-sdk-4" | "ai-sdk-5" (required)
	Stream            *bool  // default true on the server; set false for non-streaming JSON
	Cache             *bool  // default true on the server
}

// CreateCompletion invokes an agent. When stream=false the response is a JSON
// document; when stream=true (server default) the body is SSE bytes. Either
// way the raw bytes are returned for the caller to print.
func (c *Client) CreateCompletion(agentID string, req AgentCompletionRequest, params CompletionParams) ([]byte, error) {
	q := map[string]string{
		"compatibilityMode": params.CompatibilityMode,
	}
	if params.Stream != nil {
		q["stream"] = strconv.FormatBool(*params.Stream)
	}
	if params.Cache != nil {
		q["cache"] = strconv.FormatBool(*params.Cache)
	}

	r, err := c.buildRequest(http.MethodPost, fmt.Sprintf("1/agents/%s/completions", agentID), req, q)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		var errResp ErrResponse
		_ = json.Unmarshal(body, &errResp)
		return nil, fmt.Errorf("agentstudio: POST 1/agents/%s/completions -> %d %s", agentID, resp.StatusCode, formatDetail(errResp.Detail))
	}
	return body, nil
}

// Conversations ----------------------------------------------------------

func (c *Client) ListConversations(agentID string, page, limit int) (*PaginatedConversationsResponse, error) {
	var res PaginatedConversationsResponse
	if err := c.request(&res, http.MethodGet, fmt.Sprintf("1/agents/%s/conversations", agentID), nil, paginationParams(page, limit)); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetConversation(agentID, convID string) (*Conversation, error) {
	var res Conversation
	if err := c.request(&res, http.MethodGet, fmt.Sprintf("1/agents/%s/conversations/%s", agentID, convID), nil, nil); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) DeleteConversation(agentID, convID string) error {
	return c.request(nil, http.MethodDelete, fmt.Sprintf("1/agents/%s/conversations/%s", agentID, convID), nil, nil)
}

func (c *Client) DeleteAllConversations(agentID string) error {
	return c.request(nil, http.MethodDelete, fmt.Sprintf("1/agents/%s/conversations", agentID), nil, nil)
}

// ExportConversations returns the raw response body. The spec does not
// constrain the type so we leave it to the caller to write the bytes (the
// dashboard call returns an attachment).
func (c *Client) ExportConversations(agentID string) ([]byte, error) {
	r, err := c.buildRequest(http.MethodGet, fmt.Sprintf("1/agents/%s/conversations/export", agentID), nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp ErrResponse
		_ = json.Unmarshal(body, &errResp)
		return nil, fmt.Errorf("agentstudio: GET conversations/export -> %d %s", resp.StatusCode, formatDetail(errResp.Detail))
	}

	return body, nil
}

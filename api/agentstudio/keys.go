package agentstudio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) ListSecretKeys(ctx context.Context, p ListSecretKeysParams) (*PaginatedSecretKeysResponse, error) {
	q := url.Values{}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	endpoint := c.cfg.BaseURL + "/1/secret-keys"
	if encoded := q.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: list secret keys: %w", err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out PaginatedSecretKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode list secret keys response: %w", err)
	}
	return &out, nil
}

func (c *Client) GetSecretKey(ctx context.Context, id string) (*SecretKey, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent studio: get secret key: id is required")
	}
	return c.doSecretKey(ctx, http.MethodGet, "/1/secret-keys/"+url.PathEscape(id), nil, "get secret key")
}

func (c *Client) CreateSecretKey(ctx context.Context, body SecretKeyCreate) (*SecretKey, error) {
	if strings.TrimSpace(body.Name) == "" {
		return nil, fmt.Errorf("agent studio: create secret key: name is required")
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.doSecretKey(ctx, http.MethodPost, "/1/secret-keys", raw, "create secret key")
}

func (c *Client) UpdateSecretKey(ctx context.Context, id string, body SecretKeyPatch) (*SecretKey, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent studio: update secret key: id is required")
	}
	if body.Name == nil && body.AgentIDs == nil {
		return nil, fmt.Errorf("agent studio: update secret key: nothing to update")
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.doSecretKey(ctx, http.MethodPatch, "/1/secret-keys/"+url.PathEscape(id), raw, "update secret key")
}

func (c *Client) DeleteSecretKey(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("agent studio: delete secret key: id is required")
	}
	return c.doDeleteNoBody(ctx, c.cfg.BaseURL+"/1/secret-keys/"+url.PathEscape(id), "delete secret key")
}

func (c *Client) doSecretKey(ctx context.Context, method, path string, body []byte, label string) (*SecretKey, error) {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	var req *http.Request
	var err error
	if bodyReader != nil {
		req, err = http.NewRequestWithContext(ctx, method, c.cfg.BaseURL+path, bodyReader)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, c.cfg.BaseURL+path, nil)
	}
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: %s: %w", label, err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out SecretKey
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode %s response: %w", label, err)
	}
	return &out, nil
}

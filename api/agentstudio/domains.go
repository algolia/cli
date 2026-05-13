package agentstudio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListAllowedDomains(ctx context.Context, agentID string) (*AllowedDomainListResponse, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: list allowed domains: agent id is required")
	}
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) + "/allowed-domains"
	return c.getDomainList(ctx, endpoint, "list allowed domains")
}

func (c *Client) GetAllowedDomain(ctx context.Context, agentID, domainID string) (*AllowedDomain, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: get allowed domain: agent id is required")
	}
	if strings.TrimSpace(domainID) == "" {
		return nil, fmt.Errorf("agent studio: get allowed domain: domain id is required")
	}
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) +
		"/allowed-domains/" + url.PathEscape(domainID)
	return c.getDomain(ctx, endpoint, "get allowed domain")
}

func (c *Client) CreateAllowedDomain(ctx context.Context, agentID, domain string) (*AllowedDomain, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: create allowed domain: agent id is required")
	}
	if strings.TrimSpace(domain) == "" {
		return nil, fmt.Errorf("agent studio: create allowed domain: domain is required")
	}
	body, _ := json.Marshal(map[string]string{"domain": domain})
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) + "/allowed-domains"
	return c.postDomain(ctx, endpoint, body, "create allowed domain")
}

func (c *Client) DeleteAllowedDomain(ctx context.Context, agentID, domainID string) error {
	if strings.TrimSpace(agentID) == "" {
		return fmt.Errorf("agent studio: delete allowed domain: agent id is required")
	}
	if strings.TrimSpace(domainID) == "" {
		return fmt.Errorf("agent studio: delete allowed domain: domain id is required")
	}
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) +
		"/allowed-domains/" + url.PathEscape(domainID)
	return c.doDeleteNoBody(ctx, endpoint, "delete allowed domain")
}

func (c *Client) BulkInsertAllowedDomains(
	ctx context.Context,
	agentID string,
	domains []string,
) (*AllowedDomainListResponse, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agent studio: bulk insert: agent id is required")
	}
	if len(domains) == 0 {
		return nil, fmt.Errorf("agent studio: bulk insert: at least one domain is required")
	}
	body, _ := json.Marshal(map[string][]string{"domains": domains})
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) + "/allowed-domains/bulk"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: bulk insert: %w", err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out AllowedDomainListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode bulk insert response: %w", err)
	}
	return &out, nil
}

func (c *Client) BulkDeleteAllowedDomains(ctx context.Context, agentID string, domainIDs []string) error {
	if strings.TrimSpace(agentID) == "" {
		return fmt.Errorf("agent studio: bulk delete: agent id is required")
	}
	if len(domainIDs) == 0 {
		return fmt.Errorf("agent studio: bulk delete: at least one domain id is required")
	}
	body, _ := json.Marshal(map[string][]string{"domainIds": domainIDs})
	endpoint := c.cfg.BaseURL + "/1/agents/" + url.PathEscape(agentID) + "/allowed-domains/bulk"

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("agent studio: bulk delete: %w", err)
	}
	defer resp.Body.Close()
	return checkResponse(resp)
}

func (c *Client) getDomainList(ctx context.Context, endpoint, label string) (*AllowedDomainListResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: %s: %w", label, err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out AllowedDomainListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode %s response: %w", label, err)
	}
	return &out, nil
}

func (c *Client) getDomain(ctx context.Context, endpoint, label string) (*AllowedDomain, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: %s: %w", label, err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out AllowedDomain
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode %s response: %w", label, err)
	}
	return &out, nil
}

func (c *Client) postDomain(ctx context.Context, endpoint string, body []byte, label string) (*AllowedDomain, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: %s: %w", label, err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out AllowedDomain
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode %s response: %w", label, err)
	}
	return &out, nil
}

func (c *Client) doDeleteNoBody(ctx context.Context, endpoint, label string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("agent studio: %s: %w", label, err)
	}
	defer resp.Body.Close()
	return checkResponse(resp)
}

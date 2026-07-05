package agentstudio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// UpdateProvider calls PATCH /1/providers/{id}.
//
// Provider CRUD otherwise lives in the official SDK. Update stays here because
// the SDK's ProviderAuthenticationPatch models "input" as a discriminated
// oneOf union that can't round-trip a partial pass-through patch (e.g. an
// apiKey rotation) the way this verbatim PATCH does.
func (c *Client) UpdateProvider(ctx context.Context, id string, body json.RawMessage) (*Provider, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent studio: provider id is required")
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("agent studio: update provider: body is required")
	}
	endpoint := c.cfg.BaseURL + "/1/providers/" + url.PathEscape(id)
	return c.doProviderMutation(ctx, http.MethodPatch, endpoint, body, "update provider")
}

// ListProviderModels calls GET /1/providers/models — the static catalog
// of supported models per provider type. Useful before creating a provider.
func (c *Client) ListProviderModels(ctx context.Context) (map[string][]string, error) {
	endpoint := c.cfg.BaseURL + "/1/providers/models"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: list provider models: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	out := map[string][]string{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode list provider models response: %w", err)
	}
	return out, nil
}

// ListModelsForProvider calls GET /1/providers/{id}/models — the
// per-account catalog (includes fine-tunes, Azure deployments, etc).
// Returns raw JSON because the spec leaves the response shape unpinned.
func (c *Client) ListModelsForProvider(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent studio: provider id is required")
	}
	endpoint := c.cfg.BaseURL + "/1/providers/" + url.PathEscape(id) + "/models"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: list models for provider: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	body, err := readAllRawJSON(resp)
	if err != nil {
		return nil, fmt.Errorf("agent studio: read list models for provider response: %w", err)
	}
	return body, nil
}

// readAllRawJSON decodes resp.Body as a single JSON value, validating
// well-formedness so json.Encoder downstream can't refuse it.
func readAllRawJSON(resp *http.Response) (json.RawMessage, error) {
	dec := json.NewDecoder(resp.Body)
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *Client) doProviderMutation(
	ctx context.Context,
	method, endpoint string,
	body json.RawMessage,
	errLabel string,
) (*Provider, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: %s: %w", errLabel, err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var out Provider
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode %s response: %w", errLabel, err)
	}
	return &out, nil
}

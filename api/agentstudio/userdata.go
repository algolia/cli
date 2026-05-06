package agentstudio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetUserData(ctx context.Context, userToken string) (*UserDataResponse, error) {
	if strings.TrimSpace(userToken) == "" {
		return nil, fmt.Errorf("agent studio: get user data: user token is required")
	}
	endpoint := c.cfg.BaseURL + "/1/user-data/" + url.PathEscape(userToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent studio: get user data: %w", err)
	}
	defer resp.Body.Close()
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	var out UserDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("agent studio: decode get user data response: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteUserData(ctx context.Context, userToken string) error {
	if strings.TrimSpace(userToken) == "" {
		return fmt.Errorf("agent studio: delete user data: user token is required")
	}
	return c.doDeleteNoBody(ctx,
		c.cfg.BaseURL+"/1/user-data/"+url.PathEscape(userToken),
		"delete user data")
}

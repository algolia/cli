package dashboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// DefaultDashboardURL and DefaultAPIURL are empty by default and must be
// injected at build time via ldflags, e.g.:
//
//	go build -ldflags "-X github.com/algolia/cli/api/dashboard.DefaultDashboardURL=https://..."
//
// They can also be overridden at runtime with ALGOLIA_DASHBOARD_URL / ALGOLIA_API_URL / ALGOLIA_OAUTH_SCOPE
// environment variables.
var (
	DefaultDashboardURL = ""
	DefaultAPIURL       = ""
	DefaultOAuthScope   = ""
)

// Client interacts with the Algolia Dashboard OAuth endpoint and the Public API.
type Client struct {
	DashboardURL string
	APIURL       string
	OAuthScope   string
	ClientID     string
	client       *http.Client
}

// NewClient creates a new dashboard client with the given OAuth client ID.
// Respects ALGOLIA_DASHBOARD_URL, ALGOLIA_API_URL, and ALGOLIA_OAUTH_SCOPE
// environment variables, falling back to the compiled-in defaults (set via ldflags).
func NewClient(clientID string) *Client {
	dashboardURL := DefaultDashboardURL
	if v := os.Getenv("ALGOLIA_DASHBOARD_URL"); v != "" {
		dashboardURL = strings.TrimRight(v, "/")
	}
	if dashboardURL == "" {
		fmt.Fprintln(os.Stderr, "fatal: ALGOLIA_DASHBOARD_URL is not set and no default was compiled in")
		os.Exit(1)
	}

	apiURL := DefaultAPIURL
	if v := os.Getenv("ALGOLIA_API_URL"); v != "" {
		apiURL = strings.TrimRight(v, "/")
	}
	if apiURL == "" {
		fmt.Fprintln(os.Stderr, "fatal: ALGOLIA_API_URL is not set and no default was compiled in")
		os.Exit(1)
	}

	oauthScope := DefaultOAuthScope
	if v := os.Getenv("ALGOLIA_OAUTH_SCOPE"); v != "" {
		oauthScope = v
	}
	if oauthScope == "" {
		fmt.Fprintln(os.Stderr, "fatal: ALGOLIA_OAUTH_SCOPE is not set and no default was compiled in")
		os.Exit(1)
	}

	return &Client{
		DashboardURL: dashboardURL,
		APIURL:       apiURL,
		OAuthScope:   oauthScope,
		ClientID:     clientID,
		client:       http.DefaultClient,
	}
}

// NewClientWithHTTPClient creates a new dashboard client with a custom HTTP client.
// Used primarily in tests; callers must set DashboardURL and APIURL explicitly.
func NewClientWithHTTPClient(clientID string, httpClient *http.Client) *Client {
	return &Client{
		ClientID: clientID,
		client:   httpClient,
	}
}

// AuthorizeURL builds the OAuth 2.0 authorization URL for the browser-based sign-in flow.
func (c *Client) AuthorizeURL(codeChallenge, redirectURI string) string {
	return c.buildAuthorizeURL(codeChallenge, redirectURI, nil)
}

// SignupAuthorizeURL builds an OAuth 2.0 authorization URL that opens the
// sign-up page instead of the default sign-in page.
func (c *Client) SignupAuthorizeURL(codeChallenge, redirectURI string) string {
	return c.buildAuthorizeURL(codeChallenge, redirectURI, map[string]string{"screen": "signup"})
}

func (c *Client) buildAuthorizeURL(codeChallenge, redirectURI string, extra map[string]string) string {
	params := url.Values{
		"client_id":             {c.ClientID},
		"response_type":         {"code"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"scope":                 {c.OAuthScope},
		"redirect_uri":          {redirectURI},
	}
	for k, v := range extra {
		params.Set(k, v)
	}
	return c.DashboardURL + "/2/oauth/authorize?" + params.Encode()
}

// AuthorizationCodeGrant exchanges an authorization code + PKCE code_verifier
// for an access token. The redirectURI must match the one used in the authorize URL.
func (c *Client) AuthorizationCodeGrant(code, codeVerifier, redirectURI string) (*OAuthTokenResponse, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {c.ClientID},
		"code":          {code},
		"code_verifier": {codeVerifier},
		"redirect_uri":  {redirectURI},
	}

	req, err := http.NewRequest(http.MethodPost, c.DashboardURL+"/2/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseOAuthError(resp, "authorization code exchange")
	}

	var tokenResp OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshToken uses a refresh token to obtain a new access token.
func (c *Client) RefreshToken(refreshToken string) (*OAuthTokenResponse, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {c.ClientID},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequest(http.MethodPost, c.DashboardURL+"/2/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseOAuthError(resp, "token refresh")
	}

	var tokenResp OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh token response: %w", err)
	}

	return &tokenResp, nil
}

// RevokeToken revokes an OAuth access or refresh token via POST /2/oauth/revoke.
func (c *Client) RevokeToken(token string) error {
	form := url.Values{
		"client_id": {c.ClientID},
		"token":     {token},
	}

	req, err := http.NewRequest(http.MethodPost, c.DashboardURL+"/2/oauth/revoke", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("token revocation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseOAuthError(resp, "token revocation")
	}

	return nil
}

// ListApplications returns all applications for the authenticated user,
// following pagination to fetch every page.
func (c *Client) ListApplications(accessToken string) ([]Application, error) {
	var allApps []Application
	page := 1

	for {
		endpoint := fmt.Sprintf("%s/1/applications?page=%d", c.APIURL, page)
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		c.setAPIHeaders(req, accessToken)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("list applications request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrSessionExpired
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("list applications failed with status: %d", resp.StatusCode)
		}

		var appsResp ApplicationsResponse
		if err := json.NewDecoder(resp.Body).Decode(&appsResp); err != nil {
			return nil, fmt.Errorf("failed to parse applications response: %w", err)
		}

		for i := range appsResp.Data {
			allApps = append(allApps, appsResp.Data[i].toApplication())
		}

		if appsResp.Meta.CurrentPage >= appsResp.Meta.TotalPages {
			break
		}
		page++
	}

	return allApps, nil
}

// GetApplication returns a single application by its ID.
func (c *Client) GetApplication(accessToken, appID string) (*Application, error) {
	req, err := http.NewRequest(http.MethodGet, c.APIURL+"/1/application/"+url.PathEscape(appID), nil)
	if err != nil {
		return nil, err
	}
	c.setAPIHeaders(req, accessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get application request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrSessionExpired
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get application failed with status: %d", resp.StatusCode)
	}

	var singleResp SingleApplicationResponse
	if err := json.NewDecoder(resp.Body).Decode(&singleResp); err != nil {
		return nil, fmt.Errorf("failed to parse application response: %w", err)
	}

	app := singleResp.Data.toApplication()
	return &app, nil
}

// CreateApplication creates a new application for the authenticated user.
func (c *Client) CreateApplication(accessToken, region, name string) (*Application, error) {
	payload := CreateApplicationRequest{RegionCode: region, Name: name}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.APIURL+"/1/applications", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	c.setAPIHeaders(req, accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create application request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		respStr := string(respBody)

		if strings.Contains(strings.ToLower(respStr), "cluster") && strings.Contains(strings.ToLower(respStr), "not available") ||
			strings.Contains(strings.ToLower(respStr), "no cluster") {
			return nil, &ErrClusterUnavailable{Region: region, Message: fmt.Sprintf("no cluster available in region %q", region)}
		}

		return nil, fmt.Errorf("create application failed with status %d: %s", resp.StatusCode, respStr)
	}

	var singleResp SingleApplicationResponse
	if err := json.NewDecoder(resp.Body).Decode(&singleResp); err != nil {
		return nil, fmt.Errorf("failed to parse application response: %w", err)
	}

	app := singleResp.Data.toApplication()
	return &app, nil
}

// ListRegions returns the allowed hosting regions for application creation.
func (c *Client) ListRegions(accessToken string) ([]Region, error) {
	req, err := http.NewRequest(http.MethodGet, c.APIURL+"/1/hosting/regions", nil)
	if err != nil {
		return nil, err
	}
	c.setAPIHeaders(req, accessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list regions request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list regions failed with status: %d", resp.StatusCode)
	}

	var regionsResp RegionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&regionsResp); err != nil {
		return nil, fmt.Errorf("failed to parse regions response: %w", err)
	}

	return regionsResp.RegionCodes, nil
}

// WriteACL is the set of permissions for API keys created by the CLI.
var WriteACL = []string{
	"search", "browse", "seeUnretrievableAttributes", "listIndexes",
	"analytics", "logs", "addObject", "deleteObject", "deleteIndex",
	"settings", "editSettings", "recommendation",
}

// CreateAPIKey creates a new API key with the given ACL for the specified application.
func (c *Client) CreateAPIKey(accessToken, appID string, acl []string, description string) (string, error) {
	payload := CreateAPIKeyRequest{ACL: acl, Description: description}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/1/applications/%s/api-keys", c.APIURL, url.PathEscape(appID))
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	c.setAPIHeaders(req, accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("create API key request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read API key response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create API key failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var keyResp CreateAPIKeyResponse
	if err := json.Unmarshal(respBody, &keyResp); err != nil {
		return "", fmt.Errorf("failed to parse API key response: %w (body: %s)", err, string(respBody))
	}

	key := keyResp.Data.Attributes.Value
	if key == "" {
		return "", fmt.Errorf("API key creation succeeded but no key was returned in the response: %s", string(respBody))
	}

	return key, nil
}

func (c *Client) setAPIHeaders(req *http.Request, accessToken string) {
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.api+json")
}

func parseOAuthError(resp *http.Response, context string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s failed with status %d", context, resp.StatusCode)
	}

	var oauthErr OAuthErrorResponse
	if json.Unmarshal(body, &oauthErr) == nil && oauthErr.ErrorDescription != "" {
		return fmt.Errorf("%s: %s", context, oauthErr.ErrorDescription)
	}

	return fmt.Errorf("%s failed with status %d: %s", context, resp.StatusCode, string(body))
}

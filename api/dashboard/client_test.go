package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(handler http.Handler) (*httptest.Server, *Client) {
	ts := httptest.NewServer(handler)
	client := NewClientWithHTTPClient("test-client-id", ts.Client())
	client.DashboardURL = ts.URL
	client.APIURL = ts.URL
	client.OAuthScope = "scope:test"
	return ts, client
}

func TestAuthorizeURL(t *testing.T) {
	client := &Client{
		DashboardURL: "https://dashboard.example.com",
		ClientID:     "my-client-id",
		OAuthScope:   "scope:test",
	}

	url := client.AuthorizeURL("test-challenge", "http://localhost:12345")
	assert.Contains(t, url, "https://dashboard.example.com/2/oauth/authorize?")
	assert.Contains(t, url, "client_id=my-client-id")
	assert.Contains(t, url, "response_type=code")
	assert.Contains(t, url, "code_challenge=test-challenge")
	assert.Contains(t, url, "code_challenge_method=S256")
	assert.Contains(t, url, "redirect_uri=http")
}

func TestAuthorizationCodeGrant_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/2/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		require.NoError(t, r.ParseForm())
		assert.Equal(t, "authorization_code", r.FormValue("grant_type"))
		assert.Equal(t, "test-client-id", r.FormValue("client_id"))
		assert.Equal(t, "auth-code-123", r.FormValue("code"))
		assert.Equal(t, "verifier-xyz", r.FormValue("code_verifier"))

		require.NoError(t, json.NewEncoder(w).Encode(OAuthTokenResponse{
			AccessToken:  "access-token-123",
			RefreshToken: "refresh-token-456",
			TokenType:    "Bearer",
			ExpiresIn:    7200,
			Scope:        "scope:test",
			User: &User{
				ID:    1,
				Email: "user@test.com",
				Name:  "Test User",
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	resp, err := client.AuthorizationCodeGrant("auth-code-123", "verifier-xyz", "http://localhost:12345")
	require.NoError(t, err)
	assert.Equal(t, "access-token-123", resp.AccessToken)
	assert.Equal(t, "refresh-token-456", resp.RefreshToken)
	assert.Equal(t, "user@test.com", resp.User.Email)
}

func TestAuthorizationCodeGrant_InvalidGrant(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/2/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		require.NoError(t, json.NewEncoder(w).Encode(OAuthErrorResponse{
			Error:            "invalid_grant",
			ErrorDescription: "Authorization code has expired.",
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.AuthorizationCodeGrant("expired-code", "verifier", "http://localhost:12345")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Authorization code has expired")
}

func TestRefreshToken_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/2/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "refresh_token", r.FormValue("grant_type"))
		assert.Equal(t, "old-refresh-token", r.FormValue("refresh_token"))

		require.NoError(t, json.NewEncoder(w).Encode(OAuthTokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
			ExpiresIn:    7200,
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	resp, err := client.RefreshToken("old-refresh-token")
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", resp.AccessToken)
	assert.Equal(t, "new-refresh-token", resp.RefreshToken)
}

func TestListApplications_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		require.NoError(t, json.NewEncoder(w).Encode(ApplicationsResponse{
			Data: []ApplicationResource{
				{ID: "APP1", Type: "application", Attributes: ApplicationAttributes{ApplicationID: "APP1", Name: "My App", APIKey: "key1"}},
				{ID: "APP2", Type: "application", Attributes: ApplicationAttributes{ApplicationID: "APP2", Name: "Other App", APIKey: "key2"}},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	apps, err := client.ListApplications("test-token")
	require.NoError(t, err)
	assert.Len(t, apps, 2)
	assert.Equal(t, "APP1", apps[0].ID)
}

func TestListApplications_Unauthorized(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.ListApplications("expired-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session expired")
}

func TestGetApplication_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/application/APP1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		require.NoError(t, json.NewEncoder(w).Encode(SingleApplicationResponse{
			Data: ApplicationResource{
				ID: "APP1", Type: "application",
				Attributes: ApplicationAttributes{ApplicationID: "APP1", Name: "My App", APIKey: "api-key-123"},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	app, err := client.GetApplication("test-token", "APP1")
	require.NoError(t, err)
	assert.Equal(t, "APP1", app.ID)
	assert.Equal(t, "api-key-123", app.APIKey)
}

func TestCreateApplication_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload CreateApplicationRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Equal(t, "us", payload.RegionCode)
		assert.Equal(t, "My App", payload.Name)

		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(SingleApplicationResponse{
			Data: ApplicationResource{
				ID: "NEW_APP", Type: "application",
				Attributes: ApplicationAttributes{ApplicationID: "NEW_APP", Name: "My App"},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	app, err := client.CreateApplication("test-token", "us", "My App")
	require.NoError(t, err)
	assert.Equal(t, "NEW_APP", app.ID)
	assert.Equal(t, "My App", app.Name)
}

func TestGetCrawlerUser_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/crawler/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		require.NoError(t, json.NewEncoder(w).Encode(CrawlerUserResponse{
			Data: CrawlerUserData{
				ID:     "crawler-user-id",
				Email:  "crawler@example.com",
				Name:   "Crawler User",
				APIKey: "crawler-api-key",
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	user, err := client.GetCrawlerUser("test-token")
	require.NoError(t, err)
	assert.Equal(t, "crawler-user-id", user.ID)
	assert.Equal(t, "crawler@example.com", user.Email)
	assert.Equal(t, "Crawler User", user.Name)
	assert.Equal(t, "crawler-api-key", user.APIKey)
}

func TestGetCrawlerUser_HTTPError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/crawler/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.GetCrawlerUser("test-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "crawler user failed with status: 403")
}

func TestGetCrawlerUser_InvalidJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/crawler/user", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"data":`))
		require.NoError(t, err)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.GetCrawlerUser("test-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse crawler response")
}

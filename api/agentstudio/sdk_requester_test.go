package agentstudio

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func doSDKRequest(t *testing.T, r *SDKRequester, ts *httptest.Server, path string) (*http.Response, string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL+path, strings.NewReader(`{}`))
	require.NoError(t, err)
	req.Header.Set("X-Algolia-API-Key", "secret-key")
	resp, err := r.Request(req, 0, 0)
	require.NoError(t, err)
	t.Cleanup(func() { _ = resp.Body.Close() })
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, string(body)
}

func TestSDKRequester_RewritesPath(t *testing.T) {
	cases := []struct {
		name     string
		baseURL  string
		wantPath string
	}{
		{"cluster proxy keeps the prefix", "https://APP.algolia.net/agent-studio", "/agent-studio/1/agents"},
		{"override without a path drops the prefix", "https://staging.example.com", "/1/agents"},
		{"override with a path is prefixed", "https://example.com/custom/", "/custom/1/agents"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.URL.Path
				w.WriteHeader(http.StatusOK)
			}))
			t.Cleanup(ts.Close)
			doSDKRequest(t, NewSDKRequester(tc.baseURL), ts, "/agent-studio/1/agents")
			assert.Equal(t, tc.wantPath, got)
		})
	}
}

func TestSDKRequester_DoesNotFollowRedirects(t *testing.T) {
	var otherHostHits int
	other := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		otherHostHits++
	}))
	t.Cleanup(other.Close)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL, http.StatusTemporaryRedirect)
	}))
	t.Cleanup(ts.Close)

	resp, _ := doSDKRequest(t, NewSDKRequester(""), ts, "/agent-studio/1/agents")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Zero(t, otherHostHits, "the API key must never reach the redirect target")
}

func TestSDKRequester_SanitizesErrorBodies(t *testing.T) {
	const secret = "sk-SUPERSECRET"
	cases := []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{
			"422 detail list never echoes input",
			http.StatusUnprocessableEntity,
			`{"detail":[{"type":"missing","loc":["body","name"],"msg":"Field required","input":{"input":{"apiKey":"` + secret + `"}}}]}`,
			"Field required",
		},
		{"401 gets the auth hint", http.StatusUnauthorized, `{"message":"Invalid API key"}`, "check your application id and api key"},
		{"403 gets the ACL hint", http.StatusForbidden, `{"detail":"nope"}`, "missing a required ACL"},
		{"empty body falls back to the status text", http.StatusBadRequest, ``, "Bad Request"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			t.Cleanup(ts.Close)
			resp, body := doSDKRequest(t, NewSDKRequester(""), ts, "/agent-studio/1/providers")
			assert.Equal(t, tc.status, resp.StatusCode)
			assert.Contains(t, body, tc.want)
			assert.NotContains(t, body, secret)
		})
	}
}

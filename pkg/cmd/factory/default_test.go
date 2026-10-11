package factory

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/config"
)

func Test_resolveAgentStudioBaseURL(t *testing.T) {
	tests := []struct {
		name            string
		profileOverride string
		buildDefault    string
		appID           string
		want            string
		wantErr         bool
	}{
		{
			name:            "profile override wins over build default",
			profileOverride: "https://debug.example.com",
			buildDefault:    "https://agent-studio.staging.eu.algolia.com",
			appID:           "betaXYZ",
			want:            "https://debug.example.com",
		},
		{
			name:            "build default wins over cluster-proxy fallback when profile is empty",
			profileOverride: "",
			buildDefault:    "https://agent-studio.staging.eu.algolia.com",
			appID:           "betaXYZ",
			want:            "https://agent-studio.staging.eu.algolia.com",
		},
		{
			name:            "cluster-proxy fallback when both overrides are empty",
			profileOverride: "",
			buildDefault:    "",
			appID:           "APP123",
			want:            "https://APP123.algolia.net/agent-studio",
		},
		{
			name:            "trailing slash on profile override is trimmed",
			profileOverride: "https://debug.example.com/",
			buildDefault:    "",
			appID:           "APP123",
			want:            "https://debug.example.com",
		},
		{
			name:            "missing appID with no overrides errors out",
			profileOverride: "",
			buildDefault:    "",
			appID:           "",
			wantErr:         true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveAgentStudioBaseURL(tc.profileOverride, tc.buildDefault, tc.appID)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_crawlerClient_UsesConfiguredUserID(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("ALGOLIA_CRAWLER_USER_ID", "configured-user")
	t.Setenv("ALGOLIA_CRAWLER_API_KEY", "crawler-key")

	called := false
	old := fetchCrawlerUserID
	fetchCrawlerUserID = func() (string, error) {
		called = true
		return "lazy-user", nil
	}
	t.Cleanup(func() { fetchCrawlerUserID = old })

	f := New("1.0.0", &config.Config{})
	_, err := f.CrawlerClient()
	require.NoError(t, err)
	assert.False(t, called)
}

func Test_crawlerClient_LazyFetchesUserID(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("ALGOLIA_CRAWLER_USER_ID", "")
	t.Setenv("ALGOLIA_CRAWLER_API_KEY", "crawler-key")

	called := false
	old := fetchCrawlerUserID
	fetchCrawlerUserID = func() (string, error) {
		called = true
		return "lazy-user", nil
	}
	t.Cleanup(func() { fetchCrawlerUserID = old })

	f := New("1.0.0", &config.Config{})
	_, err := f.CrawlerClient()
	require.NoError(t, err)
	assert.True(t, called)
}

func Test_crawlerClient_ErrorsWhenUserIDUnavailable(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("ALGOLIA_CRAWLER_USER_ID", "")
	t.Setenv("ALGOLIA_CRAWLER_API_KEY", "crawler-key")

	old := fetchCrawlerUserID
	fetchCrawlerUserID = func() (string, error) {
		return "", errors.New("no stored token")
	}
	t.Cleanup(func() { fetchCrawlerUserID = old })

	f := New("1.0.0", &config.Config{})
	_, err := f.CrawlerClient()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "auth login")
}

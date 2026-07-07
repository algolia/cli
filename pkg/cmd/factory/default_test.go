package factory

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/config"
)

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

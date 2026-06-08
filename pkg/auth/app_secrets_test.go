package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestAppSecrets_SaveAndLoadRoundTrip(t *testing.T) {
	keyring.MockInit()

	require.NoError(t, SaveAppSecrets("APP1", AppSecrets{
		APIKey:        "key-1",
		CrawlerAPIKey: "crawler-1",
	}))

	loaded, err := LoadAppSecrets("APP1")
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "key-1", loaded.APIKey)
	assert.Equal(t, "crawler-1", loaded.CrawlerAPIKey)
}

func TestAppSecrets_LoadMissingReturnsNil(t *testing.T) {
	keyring.MockInit()

	loaded, err := LoadAppSecrets("UNKNOWN")
	require.NoError(t, err)
	assert.Nil(t, loaded)
}

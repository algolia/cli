package keychain

import (
	"errors"
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

func TestAppSecrets_PerAppIsolationAndOptionalCrawlerKey(t *testing.T) {
	keyring.MockInit()

	require.NoError(t, SaveAppSecrets("APP1", AppSecrets{APIKey: "key-1"}))
	require.NoError(
		t,
		SaveAppSecrets("APP2", AppSecrets{APIKey: "key-2", CrawlerAPIKey: "crawler-2"}),
	)

	app1, err := LoadAppSecrets("APP1")
	require.NoError(t, err)
	require.NotNil(t, app1)
	assert.Equal(t, "key-1", app1.APIKey)
	assert.Empty(t, app1.CrawlerAPIKey) // never set → stays empty

	app2, err := LoadAppSecrets("APP2")
	require.NoError(t, err)
	require.NotNil(t, app2)
	assert.Equal(t, "key-2", app2.APIKey)
	assert.Equal(t, "crawler-2", app2.CrawlerAPIKey)
}

func TestAppSecrets_EmptyAppIDIsRejected(t *testing.T) {
	keyring.MockInit()

	require.Error(t, SaveAppSecrets("", AppSecrets{APIKey: "key-1"}))

	_, err := LoadAppSecrets("")
	require.Error(t, err)
}

func TestAppSecrets_LoadKeychainErrorPropagates(t *testing.T) {
	keyring.MockInitWithError(errors.New("keychain unavailable"))

	loaded, err := LoadAppSecrets("APP1")
	require.Error(t, err)
	assert.Nil(t, loaded)
}

func TestAppSecrets_LoadMalformedJSONReturnsError(t *testing.T) {
	keyring.MockInit()
	require.NoError(t, keyring.Set(service, appSecretsUser("BAD"), "not-json"))

	loaded, err := LoadAppSecrets("BAD")
	require.Error(t, err)
	assert.Nil(t, loaded)
}

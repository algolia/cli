package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestSecrets_SetGetDelete(t *testing.T) {
	keyring.MockInit()

	require.NoError(t, SetSecret("APP1", SecretAPIKey, "secret-key"))

	got, err := GetSecret("APP1", SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "secret-key", got)

	require.NoError(t, DeleteSecret("APP1", SecretAPIKey))

	got, err = GetSecret("APP1", SecretAPIKey)
	require.NoError(t, err)
	assert.Empty(t, got, "deleted secret reads back empty")
}

func TestGetSecret_MissingReturnsEmpty(t *testing.T) {
	keyring.MockInit()

	got, err := GetSecret("UNKNOWN", SecretAPIKey)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestDeleteSecret_MissingIsNoError(t *testing.T) {
	keyring.MockInit()
	require.NoError(t, DeleteSecret("UNKNOWN", SecretCrawlerAPIKey))
}

func TestSecrets_NamespacedPerAppAndKind(t *testing.T) {
	keyring.MockInit()

	require.NoError(t, SetSecret("APP1", SecretAPIKey, "app1-api"))
	require.NoError(t, SetSecret("APP1", SecretCrawlerAPIKey, "app1-crawler"))
	require.NoError(t, SetSecret("APP2", SecretAPIKey, "app2-api"))

	api1, err := GetSecret("APP1", SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "app1-api", api1)

	crawler1, err := GetSecret("APP1", SecretCrawlerAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "app1-crawler", crawler1)

	api2, err := GetSecret("APP2", SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "app2-api", api2)

	// Deleting one secret must not affect the others.
	require.NoError(t, DeleteSecret("APP1", SecretAPIKey))

	crawler1, err = GetSecret("APP1", SecretCrawlerAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "app1-crawler", crawler1)

	api2, err = GetSecret("APP2", SecretAPIKey)
	require.NoError(t, err)
	assert.Equal(t, "app2-api", api2)
}

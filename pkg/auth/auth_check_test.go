package auth

import (
	"os"
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestRemediation(t *testing.T) {
	t.Run("signed out", func(t *testing.T) {
		keyring.MockInit()

		assert.Contains(
			t,
			Remediation(config.ErrApplicationIDNotConfigured),
			"algolia auth login",
		)
	})

	t.Run("signed in", func(t *testing.T) {
		keyring.MockInit()
		require.NoError(t, SaveToken(&dashboard.OAuthTokenResponse{
			AccessToken: "tok-1",
			ExpiresIn:   3600,
			CreatedAt:   time.Now().Unix(),
		}))

		assert.Contains(
			t,
			Remediation(config.ErrAPIKeyNotConfigured),
			"algolia application select",
		)
	})

	t.Run("api key missing from the keychain", func(t *testing.T) {
		keyring.MockInit()

		assert.Empty(t, Remediation(config.ErrAPIKeyMissingFromKeychain))
	})
}

func TestWithRemediation(t *testing.T) {
	keyring.MockInit()

	wrapped := WithRemediation(config.ErrAPIKeyNotConfigured)
	require.ErrorIs(t, wrapped, config.ErrAPIKeyNotConfigured)
	assert.Contains(t, wrapped.Error(), "algolia auth login")

	other := assert.AnError
	assert.Same(t, other, WithRemediation(other))
}

func Test_CheckACLs(t *testing.T) {
	// Remove these environment variables before the tests
	os.Unsetenv("ALGOLIA_APPLICATION_ID")
	os.Unsetenv("ALGOLIA_API_KEY")

	tests := []struct {
		name           string
		cmd            *cobra.Command
		adminKey       bool
		ACLs           []search.Acl
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "need no acls",
			cmd: &cobra.Command{
				Annotations: map[string]string{},
			},
			adminKey: false,
			ACLs:     []search.Acl{},
			wantErr:  false,
		},
		{
			name: "need admin key, not admin key",
			cmd: &cobra.Command{
				Annotations: map[string]string{
					"acls": "admin",
				},
			},
			adminKey:       false,
			ACLs:           []search.Acl{},
			wantErr:        true,
			wantErrMessage: "this command requires an admin API key. Use the `--api-key` flag with a valid admin API key",
		},
		{
			name: "need admin key, admin key",
			cmd: &cobra.Command{
				Annotations: map[string]string{
					"acls": "admin",
				},
			},
			adminKey:       true,
			ACLs:           []search.Acl{},
			wantErr:        false,
			wantErrMessage: "",
		},
		{
			name: "need ACLs, missing ACLs",
			cmd: &cobra.Command{
				Annotations: map[string]string{
					"acls": "search",
				},
			},
			adminKey: false,
			ACLs:     []search.Acl{},
			wantErr:  true,
			wantErrMessage: `Missing API key ACL(s): search
Edit your profile or use the ` + "`--api-key`" + ` flag to provide an API key with the missing ACLs.
See https://www.algolia.com/doc/guides/security/api-keys/#rights-and-restrictions for more information`,
		},
		{
			name: "need ACLs, has ACLs",
			cmd: &cobra.Command{
				Annotations: map[string]string{
					"acls": "search",
				},
			},
			adminKey: false,
			ACLs:     []search.Acl{search.ACL_SEARCH},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.adminKey {
				r.Register(
					httpmock.REST("GET", "1/keys"),
					httpmock.JSONResponse(search.ListApiKeysResponse{}),
				)
			} else {
				r.Register(
					httpmock.REST("GET", "1/keys"),
					httpmock.ErrorResponse(),
				)
			}

			if tt.ACLs != nil && !tt.adminKey {
				r.Register(
					httpmock.REST("GET", "1/keys/test"),
					httpmock.JSONResponse(search.ApiKey{Acl: tt.ACLs}),
				)
			}

			f, _ := test.NewFactory(false, &r, nil, "")
			f.Config.Profile().APIKey = "test"

			err := CheckACLs(tt.cmd, f)
			if tt.wantErr {
				assert.EqualError(t, err, tt.wantErrMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

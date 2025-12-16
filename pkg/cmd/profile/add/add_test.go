package add

import (
	"errors"
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewAddCmd(t *testing.T) {
	cfg := test.NewDefaultConfigStub()
	tests := []struct {
		name      string
		tty       bool
		cli       string
		cfg       config.IConfig
		wantsErr  bool
		wantsOpts AddOptions
	}{
		{
			name:     "not interactive, missing flags",
			cli:      "",
			cfg:      cfg,
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "not interactive, all flags",
			cli:      "--name my-app --app-id my-app-id --api-key my-admin-api-key",
			cfg:      cfg,
			tty:      false,
			wantsErr: false,
			wantsOpts: AddOptions{
				Profile: config.Profile{
					Name:          "my-app",
					ApplicationID: "my-app-id",
					APIKey:        "my-admin-api-key",
				},
			},
		},
		{
			name:     "not interactive, all flags, existing profile",
			cli:      "--name default --app-id my-app-id --api-key my-admin-api-key",
			cfg:      cfg,
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "not interactive, all flags, existing app ID",
			cli:      "--name my-app --app-id default --api-key my-admin-api-key",
			cfg:      cfg,
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "interactive, no flags",
			cli:      "",
			cfg:      cfg,
			tty:      true,
			wantsErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			io.SetStdinTTY(tt.tty)
			io.SetStdoutTTY(tt.tty)

			f := &cmdutil.Factory{
				IOStreams: io,
				Config:    tt.cfg,
			}

			var opts *AddOptions
			cmd := NewAddCmd(f, func(o *AddOptions) error {
				opts = o
				return nil
			})

			args, err := shlex.Split(tt.cli)
			require.NoError(t, err)
			cmd.SetArgs(args)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantsOpts.Profile.Name, opts.Profile.Name)
			assert.Equal(t, tt.wantsOpts.Profile.ApplicationID, opts.Profile.ApplicationID)
			assert.Equal(t, tt.wantsOpts.Profile.APIKey, opts.Profile.APIKey)
			assert.Equal(t, tt.wantsOpts.Profile.Default, opts.Profile.Default)
		})
	}
}

type stubAPIKeyInspector struct {
	listErr error

	getResp   *search.GetApiKeyResponse
	getErr    error
	getCalled bool
}

func (s *stubAPIKeyInspector) ListApiKeys(opts ...search.RequestOption) (*search.ListApiKeysResponse, error) {
	return &search.ListApiKeysResponse{}, s.listErr
}

func (s *stubAPIKeyInspector) GetApiKey(r search.ApiGetApiKeyRequest, opts ...search.RequestOption) (*search.GetApiKeyResponse, error) {
	s.getCalled = true
	return s.getResp, s.getErr
}

func (s *stubAPIKeyInspector) NewApiGetApiKeyRequest(key string) search.ApiGetApiKeyRequest {
	return search.ApiGetApiKeyRequest{}
}

func TestInspectAPIKey_AdminKeySkipsGetApiKey(t *testing.T) {
	stub := &stubAPIKeyInspector{
		listErr: nil, // admin keys can list API keys
		getErr:  errors.New("should not be called"),
	}

	isAdmin, acls, err := inspectAPIKey(stub, "my-admin-key")
	require.NoError(t, err)
	assert.True(t, isAdmin)
	assert.Nil(t, acls)
	assert.False(t, stub.getCalled)
}

func TestInspectAPIKey_NonAdminKeyReturnsACLs(t *testing.T) {
	stub := &stubAPIKeyInspector{
		listErr: errors.New("API error [403] forbidden"), // non-admin keys cannot list API keys
		getResp: &search.GetApiKeyResponse{
			Acl: []search.Acl{search.ACL_SEARCH, search.ACL_ADD_OBJECT},
		},
	}

	isAdmin, acls, err := inspectAPIKey(stub, "my-write-key")
	require.NoError(t, err)
	assert.False(t, isAdmin)
	assert.Equal(t, []string{"search", "addObject"}, acls)
	assert.True(t, stub.getCalled)
}

func TestInspectAPIKey_InvalidCredentials(t *testing.T) {
	stub := &stubAPIKeyInspector{
		listErr: errors.New("API error [403] invalid"), // fall back to GetApiKey
		getErr:  errors.New("API error [403] invalid"),
	}

	_, _, err := inspectAPIKey(stub, "bad-key")
	require.Error(t, err)
	assert.Equal(t, "invalid application credentials", err.Error())
}

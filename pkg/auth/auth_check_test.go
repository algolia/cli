package auth

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_CheckACLs(t *testing.T) {
	tests := []struct {
		name           string
		cmd            *cobra.Command
		adminKey       bool
		ACLs           []string
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "need no acls",
			cmd: &cobra.Command{
				Annotations: map[string]string{},
			},
			adminKey: false,
			ACLs:     []string{},
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
			ACLs:           []string{},
			wantErr:        true,
			wantErrMessage: "This command requires an admin API Key. Please use the `--api-key` flag to provide a valid admin API Key.\n",
		},
		{
			name: "need admin key, admin key",
			cmd: &cobra.Command{
				Annotations: map[string]string{
					"acls": "admin",
				},
			},
			adminKey:       true,
			ACLs:           []string{},
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
			ACLs:     []string{},
			wantErr:  true,
			wantErrMessage: `Missing API Key ACL(s): search
Either edit your profile or use the ` + "`--api-key`" + ` flag to provide an API Key with the missing ACLs.
See https://www.algolia.com/doc/guides/security/api-keys/#rights-and-restrictions for more information.
`,
		},
		{
			name: "need ACLs, has ACLs",
			cmd: &cobra.Command{
				Annotations: map[string]string{
					"acls": "search",
				},
			},
			adminKey: false,
			ACLs:     []string{"search"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.adminKey {
				r.Register(
					httpmock.REST("GET", "1/keys"),
					httpmock.JSONResponse(search.ListAPIKeysRes{}),
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
					httpmock.JSONResponse(search.Key{ACL: tt.ACLs}),
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

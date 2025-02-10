package list

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock/v4"
	"github.com/algolia/cli/test/v4"
)

func Test_runListCmd(t *testing.T) {
	tests := []struct {
		name    string
		isTTY   bool
		wantOut string
	}{
		{
			name:    "list",
			isTTY:   false,
			wantOut: "foo\ttest\t[search]\t[]\tNever expire\t0\t0\t[]\t5 years ago\n",
		},
		{
			name:    "list_tty",
			isTTY:   true,
			wantOut: "KEY  DESCRIPTION  ACL      INDICES  VALI...  MAX ...  MAX ...  REFE...  CREA...\nfoo  test         [sea...  []       Neve...  0        0        []       5 ye...\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := "test"
			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "1/keys"),
				httpmock.JSONResponse(search.ListApiKeysResponse{
					Keys: []search.GetApiKeyResponse{
						{
							Value:       "foo",
							Description: &name,
							Acl:         []search.Acl{search.ACL_SEARCH},
							CreatedAt:   1577836800,
						},
					},
				}),
			)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewListCmd(f, nil)
			out, err := test.Execute(cmd, "", out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}

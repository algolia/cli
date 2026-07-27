package list

import (
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const unroutableAPIURL = "http://127.0.0.1:1"

func withoutSession(t *testing.T) {
	t.Helper()
	t.Setenv("ALGOLIA_API_KEY", "")
	t.Setenv("ALGOLIA_ADMIN_API_KEY", "")
	t.Setenv("ALGOLIA_APPLICATION_ID", "")
	t.Setenv("ALGOLIA_API_URL", unroutableAPIURL)
	keyring.MockInit()
}

func explicitKeyConfig() *test.ConfigStub {
	return &test.ConfigStub{
		CurrentProfile: config.Profile{ApplicationID: "APP1", APIKey: "adm"},
	}
}

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
			withoutSession(t)

			oldNowFn := nowFn
			nowFn = func() time.Time { return time.Unix(1735689600, 0) } // 2025-01-01T00:00:00Z
			t.Cleanup(func() { nowFn = oldNowFn })

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

			f, out := test.NewFactory(tt.isTTY, &r, explicitKeyConfig(), "")
			cmd := NewListCmd(f, nil)
			out, err := test.Execute(cmd, "", out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}

func Test_runListCmd_outputJSON(t *testing.T) {
	withoutSession(t)

	oldNowFn := nowFn
	nowFn = func() time.Time { return time.Unix(1735689600, 0) } // 2025-01-01T00:00:00Z
	t.Cleanup(func() { nowFn = oldNowFn })

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

	f, out := test.NewFactory(false, &r, explicitKeyConfig(), "")
	cmd := NewListCmd(f, nil)
	out, err := test.Execute(cmd, "--output json", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, out.String(), `"keys":[`)
	assert.Contains(t, out.String(), `"value":"foo"`)
}

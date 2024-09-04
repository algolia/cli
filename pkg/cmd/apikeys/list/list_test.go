package list

import (
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
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
			wantOut: "\ttest\t[*]\t[]\tNever expire\t0\t0\t[]\ta long while ago\n",
		},
		{
			name:  "list_tty",
			isTTY: true,
			wantOut: `KEY  DESCRIPTION  ACL  INDICES  VALIDITY  MAX H...  MAX Q...  REFERERS  CREAT...
     test         [*]  []       Never...  0         0         []        a lon...
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "1/keys"),
				httpmock.JSONResponse(search.ListApiKeysResponse{
					Keys: []search.GetApiKeyResponse{
						{
							Value:                  test.Pointer("foo"),
							Description:            test.Pointer("test"),
							Acl:                    []search.Acl{"*"},
							Validity:               test.Pointer(int32(0)),
							MaxHitsPerQuery:        test.Pointer(int32(0)),
							MaxQueriesPerIPPerHour: test.Pointer(int32(0)),
							Referers:               []string{},
							CreatedAt:              time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
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

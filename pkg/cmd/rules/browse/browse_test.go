package browse

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runBrowseCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		hits    []search.Rule
		wantOut string
	}{
		{
			name:    "single rule",
			cli:     "foo",
			hits:    []search.Rule{{ObjectID: "foo"}},
			wantOut: "{\"consequence\":{},\"objectID\":\"foo\"}\n",
		},
		{
			name:    "multiple rules",
			cli:     "foo",
			hits:    []search.Rule{{ObjectID: "foo"}, {ObjectID: "bar"}},
			wantOut: "{\"consequence\":{},\"objectID\":\"foo\"}\n{\"consequence\":{},\"objectID\":\"bar\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "1/indexes/foo/settings"),
				httpmock.JSONResponse(search.SettingsResponse{}),
			)
			r.Register(
				httpmock.REST("POST", "1/indexes/foo/rules/search"),
				httpmock.JSONResponse(search.SearchRulesResponse{
					Hits: tt.hits,
				}),
			)
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, "")
			cmd := NewBrowseCmd(f)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}

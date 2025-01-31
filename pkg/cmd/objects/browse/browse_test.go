package browse

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runBrowseCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		hits    []map[string]interface{}
		wantOut string
	}{
		{
			name:    "single object",
			cli:     "foo",
			hits:    []map[string]interface{}{{"objectID": "foo"}},
			wantOut: "{\"objectID\":\"foo\"}\n",
		},
		{
			name:    "multiple objects",
			cli:     "foo",
			hits:    []map[string]interface{}{{"objectID": "foo"}, {"objectID": "bar"}},
			wantOut: "{\"objectID\":\"foo\"}\n{\"objectID\":\"bar\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("POST", "1/indexes/foo/browse"),
				httpmock.JSONResponse(search.QueryRes{
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

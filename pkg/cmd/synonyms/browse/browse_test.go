package browse

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock/v4"
	"github.com/algolia/cli/test/v4"
)

func Test_runBrowseCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		hits    []search.SynonymHit
		wantOut string
	}{
		{
			name:    "single synonym",
			cli:     "foo",
			hits:    []search.SynonymHit{{ObjectID: "foo", Type: "synonym"}},
			wantOut: "{\"objectID\":\"foo\",\"type\":\"synonym\"}\n",
		},
		{
			name: "multiple synonyms",
			cli:  "foo",
			hits: []search.SynonymHit{
				{ObjectID: "foo", Type: "synonym"},
				{ObjectID: "bar", Type: "synonym"},
			},
			wantOut: "{\"objectID\":\"foo\",\"type\":\"synonym\"}\n{\"objectID\":\"bar\",\"type\":\"synonym\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			// Check if index exists
			r.Register(
				httpmock.REST("GET", "1/indexes/foo/settings"),
				httpmock.JSONResponse(search.SettingsResponse{}),
			)
			r.Register(
				httpmock.REST("POST", "1/indexes/foo/synonyms/search"),
				httpmock.JSONResponse(search.SearchSynonymsResponse{
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

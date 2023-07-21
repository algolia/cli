package browse

import (
	"fmt"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runBrowseCmd(t *testing.T) {
	tests := []struct {
		name         string
		cli          string
		dictionaries []search.DictionaryName
		entries      bool
		isTTY        bool
		wantOut      string
	}{
		{
			name: "one dictionary",
			cli:  "plurals",
			dictionaries: []search.DictionaryName{
				search.Plurals,
			},
			entries: true,
			isTTY:   false,
			wantOut: "{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n",
		},
		{
			name: "multiple dictionaries",
			cli:  "plurals compounds",
			dictionaries: []search.DictionaryName{
				search.Plurals,
				search.Compounds,
			},
			entries: true,
			isTTY:   false,
			wantOut: "{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n",
		},
		{
			name: "all dictionaries",
			cli:  "--all",
			dictionaries: []search.DictionaryName{
				search.Stopwords,
				search.Plurals,
				search.Compounds,
			},
			entries: true,
			isTTY:   false,
			wantOut: "{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n",
		},
		{
			name: "one dictionary with default stopwords",
			cli:  "--all --include-defaults",
			dictionaries: []search.DictionaryName{
				search.Stopwords,
				search.Plurals,
				search.Compounds,
			},
			entries: true,
			isTTY:   false,
			wantOut: "{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}\n",
		},
		{
			name: "no entries",
			cli:  "plurals",
			dictionaries: []search.DictionaryName{
				search.Plurals,
			},
			entries: false,
			isTTY:   false,
			wantOut: "! No entries found.\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			for _, d := range tt.dictionaries {
				var entries []DictionaryEntry
				if tt.entries {
					entries = append(entries, DictionaryEntry{Type: "custom"})
				}
				r.Register(httpmock.REST("POST", fmt.Sprintf("1/dictionaries/%s/search", d)), httpmock.JSONResponse(search.SearchDictionariesRes{
					Hits: entries,
				}))
				r.Register(httpmock.REST("POST", fmt.Sprintf("1/dictionaries/%s/batch", d)), httpmock.JSONResponse(search.TaskStatusRes{}))
			}

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewBrowseCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}

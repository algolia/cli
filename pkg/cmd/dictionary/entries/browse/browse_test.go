package browse

import (
	"fmt"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runBrowseCmd(t *testing.T) {
	tests := []struct {
		name         string
		cli          string
		dictionaries []search.DictionaryType
		entries      bool
		isTTY        bool
		wantOut      string
	}{
		{
			name: "one dictionary",
			cli:  "plurals",
			dictionaries: []search.DictionaryType{
				search.DICTIONARY_TYPE_PLURALS,
			},
			entries: true,
			isTTY:   false,
			wantOut: "{\"objectID\":\"\",\"type\":\"custom\"}\n",
		},
		{
			name: "multiple dictionaries",
			cli:  "plurals compounds",
			dictionaries: []search.DictionaryType{
				search.DICTIONARY_TYPE_PLURALS,
				search.DICTIONARY_TYPE_COMPOUNDS,
			},
			entries: true,
			isTTY:   false,
			wantOut: "{\"objectID\":\"\",\"type\":\"custom\"}\n{\"objectID\":\"\",\"type\":\"custom\"}\n",
		},
		{
			name: "all dictionaries",
			cli:  "--all",
			dictionaries: []search.DictionaryType{
				search.DICTIONARY_TYPE_STOPWORDS,
				search.DICTIONARY_TYPE_PLURALS,
				search.DICTIONARY_TYPE_COMPOUNDS,
			},
			entries: true,
			isTTY:   false,
			wantOut: "{\"objectID\":\"\",\"type\":\"custom\"}\n{\"objectID\":\"\",\"type\":\"custom\"}\n{\"objectID\":\"\",\"type\":\"custom\"}\n",
		},
		{
			name: "one dictionary with default stopwords",
			cli:  "--all --include-defaults",
			dictionaries: []search.DictionaryType{
				search.DICTIONARY_TYPE_STOPWORDS,
				search.DICTIONARY_TYPE_PLURALS,
				search.DICTIONARY_TYPE_COMPOUNDS,
			},
			entries: true,
			isTTY:   false,
			wantOut: "{\"objectID\":\"\",\"type\":\"custom\"}\n{\"objectID\":\"\",\"type\":\"custom\"}\n{\"objectID\":\"\",\"type\":\"custom\"}\n",
		},
		{
			name: "no entries",
			cli:  "plurals",
			dictionaries: []search.DictionaryType{
				search.DICTIONARY_TYPE_PLURALS,
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
				var entries []search.DictionaryEntry
				if tt.entries {
					entries = append(
						entries,
						search.DictionaryEntry{
							Type: search.DICTIONARY_ENTRY_TYPE_CUSTOM.Ptr(),
						},
					)
				}
				r.Register(
					httpmock.REST("POST", fmt.Sprintf("1/dictionaries/%s/search", d)),
					httpmock.JSONResponse(search.SearchDictionaryEntriesResponse{
						Hits: entries,
					}),
				)
				r.Register(
					httpmock.REST("POST", fmt.Sprintf("1/dictionaries/%s/batch", d)),
					httpmock.JSONResponse(search.UpdatedAtResponse{}),
				)
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

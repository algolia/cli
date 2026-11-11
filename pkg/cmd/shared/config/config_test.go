package config

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestGetIndexConfig(t *testing.T) {
	tests := []struct {
		name         string
		scope        []string
		wantSettings bool
		wantRules    []string
		wantSynonyms []string
	}{
		{
			name:         "all scopes",
			scope:        []string{"settings", "synonyms", "rules"},
			wantSettings: true,
			wantRules:    []string{"rule-1", "rule-2"},
			wantSynonyms: []string{"syn-1"},
		},
		{
			name:      "rules only",
			scope:     []string{"rules"},
			wantRules: []string{"rule-1", "rule-2"},
		},
		{
			name:         "synonyms only",
			scope:        []string{"synonyms"},
			wantSynonyms: []string{"syn-1"},
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
					Hits: []search.Rule{{ObjectID: "rule-1"}, {ObjectID: "rule-2"}},
				}),
			)
			r.Register(
				httpmock.REST("POST", "1/indexes/foo/synonyms/search"),
				httpmock.JSONResponse(search.SearchSynonymsResponse{
					Hits: []search.SynonymHit{{ObjectID: "syn-1", Type: search.SYNONYM_TYPE_SYNONYM}},
				}),
			)

			f, _ := test.NewFactory(true, &r, nil, "")
			client, err := f.SearchClient()
			require.NoError(t, err)

			got, err := GetIndexConfig(client, "foo", tt.scope, f.IOStreams.ColorScheme())
			require.NoError(t, err)

			assert.Equal(t, tt.wantSettings, got.Settings != nil)

			var rules []string
			for _, rule := range got.Rules {
				rules = append(rules, rule.ObjectID)
			}
			assert.Equal(t, tt.wantRules, rules)

			var synonyms []string
			for _, synonym := range got.Synonyms {
				synonyms = append(synonyms, synonym.ObjectID)
			}
			assert.Equal(t, tt.wantSynonyms, synonyms)
		})
	}
}

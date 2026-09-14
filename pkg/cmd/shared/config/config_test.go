package config

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

func newTestClient(t *testing.T, r *httpmock.Registry) *search.APIClient {
	t.Helper()

	cfg := search.SearchConfiguration{
		Configuration: transport.Configuration{
			AppID:     "default",
			ApiKey:    "default",
			Requester: r,
		},
	}
	client, err := search.NewClientWithConfig(cfg)
	require.NoError(t, err)

	return client
}

func Test_GetSynonyms(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("POST", "1/indexes/foo/synonyms/search"),
		httpmock.JSONResponse(search.SearchSynonymsResponse{
			Hits: []search.SynonymHit{
				{ObjectID: "foo", Type: "synonym"},
				{ObjectID: "bar", Type: "synonym"},
			},
		}),
	)
	defer r.Verify(t)

	client := newTestClient(t, &r)

	synonyms, err := GetSynonyms(client, "foo")
	require.NoError(t, err)
	assert.Equal(t, []search.SynonymHit{
		{ObjectID: "foo", Type: "synonym"},
		{ObjectID: "bar", Type: "synonym"},
	}, synonyms)
}

func Test_GetRules(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("POST", "1/indexes/foo/rules/search"),
		httpmock.JSONResponse(search.SearchRulesResponse{
			Hits: []search.Rule{
				{ObjectID: "rule-1"},
				{ObjectID: "rule-2"},
			},
		}),
	)
	defer r.Verify(t)

	client := newTestClient(t, &r)

	rules, err := GetRules(client, "foo")
	require.NoError(t, err)
	assert.Equal(t, []search.Rule{
		{ObjectID: "rule-1"},
		{ObjectID: "rule-2"},
	}, rules)
}

func Test_GetIndexConfig(t *testing.T) {
	cs := iostreams.NewColorScheme(false, false, false)

	tests := []struct {
		name        string
		scope       []string
		synonymHits []search.SynonymHit
		ruleHits    []search.Rule
		settings    search.SettingsResponse
		wantErr     bool
		assertFn    func(t *testing.T, cfg *ExportConfigJSON)
	}{
		{
			name:  "exports synonyms",
			scope: []string{"synonyms"},
			synonymHits: []search.SynonymHit{
				{ObjectID: "foo", Type: "synonym"},
			},
			assertFn: func(t *testing.T, cfg *ExportConfigJSON) {
				assert.Len(t, cfg.Synonyms, 1)
				assert.Equal(t, "foo", cfg.Synonyms[0].ObjectID)
			},
		},
		{
			name:  "exports rules",
			scope: []string{"rules"},
			ruleHits: []search.Rule{
				{ObjectID: "rule-1"},
			},
			assertFn: func(t *testing.T, cfg *ExportConfigJSON) {
				assert.Len(t, cfg.Rules, 1)
				assert.Equal(t, "rule-1", cfg.Rules[0].ObjectID)
			},
		},
		{
			name:    "no config to export returns error",
			scope:   []string{"rules", "synonyms"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if utils.Contains(tt.scope, "synonyms") {
				r.Register(
					httpmock.REST("POST", "1/indexes/foo/synonyms/search"),
					httpmock.JSONResponse(search.SearchSynonymsResponse{Hits: tt.synonymHits}),
				)
			}
			if utils.Contains(tt.scope, "rules") {
				r.Register(
					httpmock.REST("POST", "1/indexes/foo/rules/search"),
					httpmock.JSONResponse(search.SearchRulesResponse{Hits: tt.ruleHits}),
				)
			}
			if utils.Contains(tt.scope, "settings") {
				r.Register(
					httpmock.REST("GET", "1/indexes/foo/settings"),
					httpmock.JSONResponse(tt.settings),
				)
			}
			defer r.Verify(t)

			client := newTestClient(t, &r)

			cfg, err := GetIndexConfig(client, "foo", tt.scope, cs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.assertFn(t, cfg)
		})
	}
}

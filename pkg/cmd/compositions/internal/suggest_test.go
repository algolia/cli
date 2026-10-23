package internal_test

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/httpmock"
)

func searchClientFor(r *httpmock.Registry) func() (*search.APIClient, error) {
	return func() (*search.APIClient, error) {
		return search.NewClientWithConfig(search.SearchConfiguration{
			Configuration: transport.Configuration{AppID: "default", ApiKey: "default", Requester: r},
		})
	}
}

func compositionClientFor(r *httpmock.Registry) func() (*composition.APIClient, error) {
	return func() (*composition.APIClient, error) {
		return composition.NewClientWithConfig(composition.CompositionConfiguration{
			Configuration: transport.Configuration{AppID: "default", ApiKey: "default", Requester: r},
		})
	}
}

func TestCompositionSuggester_Indices(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", "1/indexes"),
		httpmock.JSONResponse(search.ListIndicesResponse{Items: []search.FetchedIndex{{Name: "alpha"}, {Name: "beta"}}}))

	suggest := compinternal.CompositionSuggester(compinternal.SuggestDeps{Search: searchClientFor(r)})
	require.NotNil(t, suggest)

	assert.Equal(t, []string{"alpha", "beta"}, suggest("Composition.index", "index", nil))
	assert.Equal(t, []string{"alpha", "beta"}, suggest("MainRecommend.indexName", "indexName", nil))
	assert.Equal(t, []string{"alpha", "beta"}, suggest("Composition.sortingStrategy", "sortingStrategy", nil))
	assert.Nil(t, suggest("Composition.name", "name", nil))
}

func TestCompositionSuggester_Attributes(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", "1/indexes/alpha/settings"),
		httpmock.JSONResponse(search.SettingsResponse{
			SearchableAttributes:  []string{"title", "unordered(brand)"},
			AttributesForFaceting: []string{"filterOnly(category)"},
		}))

	suggest := compinternal.CompositionSuggester(compinternal.SuggestDeps{Search: searchClientFor(r)})

	got := suggest("params.attributesToRetrieve", "attributesToRetrieve", map[string]string{"index": "alpha"})
	assert.Equal(t, []string{"title", "brand", "category"}, got)
	assert.Nil(t, suggest("params.attributesToRetrieve", "attributesToRetrieve", map[string]string{}))
}

func TestCompositionSuggester_Filters(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", "1/indexes/alpha/settings"),
		httpmock.JSONResponse(search.SettingsResponse{
			SearchableAttributes:  []string{"title"},
			AttributesForFaceting: []string{"brand"},
		}))

	suggest := compinternal.CompositionSuggester(compinternal.SuggestDeps{Search: searchClientFor(r)})

	// filters offer "attribute:" templates from the in-scope index.
	got := suggest("params.filters", "filters", map[string]string{"index": "alpha"})
	assert.Equal(t, []string{"title:", "brand:"}, got)
	assert.Nil(t, suggest("params.filters", "filters", nil))
}

func TestCompositionSuggester_SortBy(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", "1/compositions/my-comp"),
		httpmock.StringResponse(`{"objectID":"my-comp","sortingStrategy":{"Price (asc)":"idx_asc"}}`))

	suggest := compinternal.CompositionSuggester(compinternal.SuggestDeps{
		Composition:   compositionClientFor(r),
		CompositionID: "my-comp",
	})

	assert.Equal(t, []string{"Price (asc)"}, suggest("params.sortBy", "sortBy", nil))
}

func TestCompositionSuggester_InjectedItemKey(t *testing.T) {
	suggest := compinternal.CompositionSuggester(compinternal.SuggestDeps{})
	assert.Equal(t, []string{"my-injected-group-1"}, suggest("injectedItems[0].key", "key", nil))
}

func TestCompositionSuggester_EmptyDeps(t *testing.T) {
	suggest := compinternal.CompositionSuggester(compinternal.SuggestDeps{})
	require.NotNil(t, suggest)
	assert.Nil(t, suggest("Composition.index", "index", nil))
}

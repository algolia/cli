package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mergeAlgoliaSearchIndexTool_emptyTools(t *testing.T) {
	got, err := mergeAlgoliaSearchIndexTool(nil, "catalog_search", "products", "catalog",
		nil)
	require.NoError(t, err)
	assert.JSONEq(t,
		`[{"name":"catalog_search","type":"algolia_search_index","indices":[{"index":"products","description":"catalog"}]}]`,
		string(got))
}

func Test_mergeAlgoliaSearchIndexTool_appendsIndexToExistingTool(t *testing.T) {
	existing := `[{"name":"s","type":"algolia_search_index","indices":[{"index":"a","description":"A"}]}]`
	got, err := mergeAlgoliaSearchIndexTool([]byte(existing), "catalog_search", "b", "B", nil)
	require.NoError(t, err)
	assert.JSONEq(t,
		`[{"name":"s","type":"algolia_search_index","indices":[{"index":"a","description":"A"},{"index":"b","description":"B"}]}]`,
		string(got))
}

func Test_mergeAlgoliaSearchIndexTool_fillsMissingToolName(t *testing.T) {
	existing := `[{"type":"algolia_search_index","indices":[{"index":"a","description":"A"}]}]`
	got, err := mergeAlgoliaSearchIndexTool([]byte(existing), "fill_name", "b", "B", nil)
	require.NoError(t, err)
	assert.JSONEq(t,
		`[{"name":"fill_name","type":"algolia_search_index","indices":[{"index":"a","description":"A"},{"index":"b","description":"B"}]}]`,
		string(got))
}

func Test_mergeAlgoliaSearchIndexTool_preservesOtherTools(t *testing.T) {
	existing := `[{"type":"other"},{"name":"s","type":"algolia_search_index","indices":[{"index":"a","description":"A"}]}]`
	got, err := mergeAlgoliaSearchIndexTool([]byte(existing), "catalog_search", "b", "B", nil)
	require.NoError(t, err)
	assert.JSONEq(t,
		`[{"type":"other"},{"name":"s","type":"algolia_search_index","indices":[{"index":"a","description":"A"},{"index":"b","description":"B"}]}]`,
		string(got))
}

func Test_mergeAlgoliaSearchIndexTool_duplicateIndex(t *testing.T) {
	existing := `[{"name":"s","type":"algolia_search_index","indices":[{"index":"a","description":"A"}]}]`
	_, err := mergeAlgoliaSearchIndexTool([]byte(existing), "catalog_search", "a", "dup", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `already present`)
}

func Test_mergeAlgoliaSearchIndexTool_rejectsNonArray(t *testing.T) {
	_, err := mergeAlgoliaSearchIndexTool([]byte(`{}`), "catalog_search", "x", "d", nil)
	require.Error(t, err)
}

func Test_mergeAlgoliaSearchIndexTool_searchParameters(t *testing.T) {
	got, err := mergeAlgoliaSearchIndexTool(nil, "catalog_search", "products", "cat",
		[]byte(`{"filters":"inStock:true"}`))
	require.NoError(t, err)
	assert.JSONEq(t,
		`[{"name":"catalog_search","type":"algolia_search_index","indices":[{"index":"products","description":"cat","searchParameters":{"filters":"inStock:true"}}]}]`,
		string(got))
}

func Test_mergeAlgoliaSearchIndexTool_rejectsBadSearchParameters(t *testing.T) {
	_, err := mergeAlgoliaSearchIndexTool(nil, "catalog_search", "p", "d", []byte(`not json`))
	require.Error(t, err)
}

func Test_mergeAlgoliaSearchIndexTool_rejectsBadToolName(t *testing.T) {
	_, err := mergeAlgoliaSearchIndexTool(nil, "ab", "p", "d", nil)
	require.Error(t, err)
}

func Test_validateToolName(t *testing.T) {
	assert.NoError(t, validateToolName("abc"))
	assert.NoError(t, validateToolName("abcdefghijklmnopqrstuvwxyzabcdef")) // 32 runes
	assert.Error(t, validateToolName("ab"))
	assert.Error(t, validateToolName("abcdefghijklmnopqrstuvwxyzabcdefg")) // 33
}

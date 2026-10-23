package internal

import (
	"strings"
	"sync"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"

	"github.com/algolia/cli/pkg/interactive"
)

// indexFields hold an index name as their value (sortingStrategy maps sort
// labels to index names).
var indexFields = map[string]bool{"index": true, "indexName": true, "sortingStrategy": true}

// attributeFields hold a single attribute name per entry, completed from the
// in-scope index's settings.
var attributeFields = map[string]bool{
	"attributesToRetrieve":             true,
	"restrictSearchableAttributes":     true,
	"attributesToHighlight":            true,
	"attributesToSnippet":              true,
	"disableExactOnAttributes":         true,
	"disableTypoToleranceOnAttributes": true,
	"facets":                           true,
}

// SuggestDeps carries what the composition suggester needs to reach live data.
// Any nil client disables the suggestions that depend on it.
type SuggestDeps struct {
	Search        func() (*search.APIClient, error)
	Composition   func() (*composition.APIClient, error)
	CompositionID string
}

// CompositionSuggester completes composition text fields from live data:
//   - index fields with the application's index names,
//   - attribute fields with the in-scope index's attributes (scope["index"]),
//   - filters with "attribute:" templates from the in-scope index,
//   - sortBy with the composition's sorting-strategy labels,
//   - the injected item key with a static format hint.
//
// Results are cached for the session and any error yields no suggestions.
func CompositionSuggester(deps SuggestDeps) interactive.SuggestFunc {
	indices := onceIndices(deps.Search)
	attributes := indexAttributes(deps.Search)
	sortLabels := onceSortLabels(deps.Composition, deps.CompositionID)

	return func(_, jsonName string, scope map[string]string) []string {
		switch {
		case jsonName == "key":
			return []string{"my-injected-group-1"}
		case jsonName == "sortBy":
			return sortLabels()
		case jsonName == "filters":
			return filterTemplates(attributes, scopeIndex(scope))
		case indexFields[jsonName]:
			return indices()
		case attributeFields[jsonName]:
			return attributes(scopeIndex(scope))
		}
		return nil
	}
}

func scopeIndex(scope map[string]string) string {
	if idx := scope["index"]; idx != "" {
		return idx
	}
	return scope["indexName"]
}

// filterTemplates turns an index's attributes into "attribute:" starting points
// so a filter expression can be built from a real attribute.
func filterTemplates(attributes func(string) []string, index string) []string {
	if index == "" {
		return nil
	}
	attrs := attributes(index)
	out := make([]string, 0, len(attrs))
	for _, a := range attrs {
		out = append(out, a+":")
	}
	return out
}

func onceIndices(searchClient func() (*search.APIClient, error)) func() []string {
	if searchClient == nil {
		return func() []string { return nil }
	}
	load := sync.OnceValues(func() ([]string, error) {
		client, err := searchClient()
		if err != nil {
			return nil, err
		}
		res, err := client.ListIndices(client.NewApiListIndicesRequest())
		if err != nil {
			return nil, err
		}
		names := make([]string, 0, len(res.Items))
		for _, idx := range res.Items {
			names = append(names, idx.Name)
		}
		return names, nil
	})
	return func() []string { names, _ := load(); return names }
}

func indexAttributes(searchClient func() (*search.APIClient, error)) func(string) []string {
	if searchClient == nil {
		return func(string) []string { return nil }
	}
	var mu sync.Mutex
	cache := map[string][]string{}
	return func(index string) []string {
		if index == "" {
			return nil
		}
		mu.Lock()
		defer mu.Unlock()
		if cached, ok := cache[index]; ok {
			return cached
		}
		attrs := fetchAttributes(searchClient, index)
		cache[index] = attrs
		return attrs
	}
}

func onceSortLabels(compositionClient func() (*composition.APIClient, error), compositionID string) func() []string {
	if compositionClient == nil || compositionID == "" {
		return func() []string { return nil }
	}
	load := sync.OnceValues(func() ([]string, error) {
		client, err := compositionClient()
		if err != nil {
			return nil, err
		}
		comp, err := client.GetComposition(client.NewApiGetCompositionRequest(compositionID))
		if err != nil || comp.SortingStrategy == nil {
			return nil, err
		}
		labels := make([]string, 0, len(*comp.SortingStrategy))
		for label := range *comp.SortingStrategy {
			labels = append(labels, label)
		}
		return labels, nil
	})
	return func() []string { labels, _ := load(); return labels }
}

// fetchAttributes returns an index's searchable and facet attribute names, with
// any modifiers (like unordered(name)) stripped. Returns nil on any error.
func fetchAttributes(searchClient func() (*search.APIClient, error), index string) []string {
	client, err := searchClient()
	if err != nil {
		return nil
	}
	settings, err := client.GetSettings(client.NewApiGetSettingsRequest(index))
	if err != nil {
		return nil
	}

	seen := map[string]bool{}
	var attrs []string
	for _, list := range [][]string{settings.SearchableAttributes, settings.AttributesForFaceting} {
		for _, a := range list {
			name := bareAttribute(a)
			if name != "" && !seen[name] {
				seen[name] = true
				attrs = append(attrs, name)
			}
		}
	}
	return attrs
}

// bareAttribute strips Algolia attribute modifiers, so "unordered(name)" and
// "filterOnly(brand)" become "name" and "brand".
func bareAttribute(a string) string {
	if open := strings.IndexByte(a, '('); open != -1 {
		if close := strings.IndexByte(a[open:], ')'); close != -1 {
			return strings.TrimSpace(a[open+1 : open+close])
		}
	}
	return strings.TrimSpace(a)
}

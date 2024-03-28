// This file is generated; DO NOT EDIT.

package cmdutil

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var BrowseParamsObject = []string{
	"advancedSyntax",
	"advancedSyntaxFeatures",
	"allowTyposOnNumericTokens",
	"alternativesAsExact",
	"analytics",
	"analyticsTags",
	"aroundLatLng",
	"aroundLatLngViaIP",
	"aroundPrecision",
	"aroundRadius",
	"attributeCriteriaComputedByMinProximity",
	"attributesToHighlight",
	"attributesToRetrieve",
	"attributesToSnippet",
	"clickAnalytics",
	"cursor",
	"customRanking",
	"decompoundQuery",
	"disableExactOnAttributes",
	"disableTypoToleranceOnAttributes",
	"distinct",
	"enableABTest",
	"enablePersonalization",
	"enableReRanking",
	"enableRules",
	"exactOnSingleWordQuery",
	"facetFilters",
	"facetingAfterDistinct",
	"facets",
	"filters",
	"getRankingInfo",
	"highlightPostTag",
	"highlightPreTag",
	"hitsPerPage",
	"ignorePlurals",
	"insideBoundingBox",
	"insidePolygon",
	"keepDiacriticsOnCharacters",
	"length",
	"maxFacetHits",
	"maxValuesPerFacet",
	"minProximity",
	"minWordSizefor1Typo",
	"minWordSizefor2Typos",
	"minimumAroundRadius",
	"mode",
	"naturalLanguages",
	"numericFilters",
	"offset",
	"optionalFilters",
	"optionalWords",
	"page",
	"percentileComputation",
	"personalizationImpact",
	"query",
	"queryLanguages",
	"queryType",
	"ranking",
	"reRankingApplyFilter",
	"relevancyStrictness",
	"removeStopWords",
	"removeWordsIfNoResults",
	"renderingContent",
	"replaceSynonymsInHighlight",
	"responseFields",
	"restrictHighlightAndSnippetArrays",
	"restrictSearchableAttributes",
	"ruleContexts",
	"semanticSearch",
	"similarQuery",
	"snippetEllipsisText",
	"sortFacetValuesBy",
	"sumOrFiltersScores",
	"synonyms",
	"tagFilters",
	"typoTolerance",
	"userToken",
}

var DeleteByParams = []string{
	"aroundLatLng",
	"aroundRadius",
	"facetFilters",
	"filters",
	"insideBoundingBox",
	"insidePolygon",
	"numericFilters",
	"tagFilters",
}

var IndexSettings = []string{
	"advancedSyntax",
	"advancedSyntaxFeatures",
	"allowCompressionOfIntegerArray",
	"allowTyposOnNumericTokens",
	"alternativesAsExact",
	"attributeCriteriaComputedByMinProximity",
	"attributeForDistinct",
	"attributesForFaceting",
	"attributesToHighlight",
	"attributesToRetrieve",
	"attributesToSnippet",
	"attributesToTransliterate",
	"camelCaseAttributes",
	"customNormalization",
	"customRanking",
	"decompoundQuery",
	"decompoundedAttributes",
	"disableExactOnAttributes",
	"disablePrefixOnAttributes",
	"disableTypoToleranceOnAttributes",
	"disableTypoToleranceOnWords",
	"distinct",
	"enablePersonalization",
	"enableReRanking",
	"enableRules",
	"exactOnSingleWordQuery",
	"highlightPostTag",
	"highlightPreTag",
	"hitsPerPage",
	"ignorePlurals",
	"indexLanguages",
	"keepDiacriticsOnCharacters",
	"maxFacetHits",
	"maxValuesPerFacet",
	"minProximity",
	"minWordSizefor1Typo",
	"minWordSizefor2Typos",
	"mode",
	"numericAttributesForFiltering",
	"optionalWords",
	"paginationLimitedTo",
	"queryLanguages",
	"queryType",
	"ranking",
	"reRankingApplyFilter",
	"relevancyStrictness",
	"removeStopWords",
	"removeWordsIfNoResults",
	"renderingContent",
	"replaceSynonymsInHighlight",
	"replicas",
	"responseFields",
	"restrictHighlightAndSnippetArrays",
	"searchableAttributes",
	"semanticSearch",
	"separatorsToIndex",
	"snippetEllipsisText",
	"sortFacetValuesBy",
	"typoTolerance",
	"unretrievableAttributes",
	"userData",
}

var SearchParamsObject = []string{
	"advancedSyntax",
	"advancedSyntaxFeatures",
	"allowTyposOnNumericTokens",
	"alternativesAsExact",
	"analytics",
	"analyticsTags",
	"aroundLatLng",
	"aroundLatLngViaIP",
	"aroundPrecision",
	"aroundRadius",
	"attributeCriteriaComputedByMinProximity",
	"attributesToHighlight",
	"attributesToRetrieve",
	"attributesToSnippet",
	"clickAnalytics",
	"customRanking",
	"decompoundQuery",
	"disableExactOnAttributes",
	"disableTypoToleranceOnAttributes",
	"distinct",
	"enableABTest",
	"enablePersonalization",
	"enableReRanking",
	"enableRules",
	"exactOnSingleWordQuery",
	"facetFilters",
	"facetingAfterDistinct",
	"facets",
	"filters",
	"getRankingInfo",
	"highlightPostTag",
	"highlightPreTag",
	"hitsPerPage",
	"ignorePlurals",
	"insideBoundingBox",
	"insidePolygon",
	"keepDiacriticsOnCharacters",
	"length",
	"maxFacetHits",
	"maxValuesPerFacet",
	"minProximity",
	"minWordSizefor1Typo",
	"minWordSizefor2Typos",
	"minimumAroundRadius",
	"mode",
	"naturalLanguages",
	"numericFilters",
	"offset",
	"optionalFilters",
	"optionalWords",
	"page",
	"percentileComputation",
	"personalizationImpact",
	"query",
	"queryLanguages",
	"queryType",
	"ranking",
	"reRankingApplyFilter",
	"relevancyStrictness",
	"removeStopWords",
	"removeWordsIfNoResults",
	"renderingContent",
	"replaceSynonymsInHighlight",
	"responseFields",
	"restrictHighlightAndSnippetArrays",
	"restrictSearchableAttributes",
	"ruleContexts",
	"semanticSearch",
	"similarQuery",
	"snippetEllipsisText",
	"sortFacetValuesBy",
	"sumOrFiltersScores",
	"synonyms",
	"tagFilters",
	"typoTolerance",
	"userToken",
}

func AddBrowseParamsObjectFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Whether to support phrase matching and excluding words from search queries.

Use the `+"`"+`advancedSyntaxFeatures`+"`"+` parameter to control which feature is supported.
`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Advanced search syntax features you want to support.

<dl>
<dt><code>exactPhrase</code></dt>
<dd>

Phrases in quotes must match exactly.
For example, `+"`"+`sparkly blue "iPhone case"`+"`"+` only returns records with the exact string "iPhone case".

</dd>
<dt><code>excludeWords</code></dt>
<dd>

Query words prefixed with a `+"`"+`-`+"`"+` must not occur in a record.
For example, `+"`"+`search -engine`+"`"+` matches records that contain "search" but not "engine".

</dd>
</dl>

This setting only has an effect if `+"`"+`advancedSyntax`+"`"+` is true.
`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers in the search query.

Turn off this setting to reduce the number of irrelevant matches
when searching in large sets of similar numbers.
`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives of query words that should be considered as exact matches by the Exact ranking criterion.

<dl>
<dt><code>ignorePlurals</code></dt>
<dd>

Plurals and similar declensions added by the `+"`"+`ignorePlurals`+"`"+` setting are considered exact matches.

</dd>
<dt><code>singleWordSynonym</code></dt>
<dd>
Single-word synonyms, such as "NY/NYC" are considered exact matches.
</dd>
<dt><code>multiWordsSynonym</code></dt>
<dd>
Multi-word synonyms, such as "NY/New York" are considered exact matches.
</dd>
</dl>.
`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`Whether this search will be included in Analytics.`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`Tags to apply to the query for [segmenting analytics data](https://www.algolia.com/doc/guides/search-analytics/guides/segments/).`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Coordinates for the center of a circle, expressed as a comma-separated string of latitude and longitude.

Only records included within circle around this central location are included in the results.
The radius of the circle is determined by the `+"`"+`aroundRadius`+"`"+` and `+"`"+`minimumAroundRadius`+"`"+` settings.
This parameter is ignored if you also specify `+"`"+`insidePolygon`+"`"+` or `+"`"+`insideBoundingBox`+"`"+`.
`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`Whether to obtain the coordinates from the request's IP address.`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	aroundPrecision := NewJSONVar([]string{"integer", "array"}...)
	cmd.Flags().Var(aroundPrecision, "aroundPrecision", heredoc.Doc(`Precision of a coordinate-based search in meters to group results with similar distances.

The Geo ranking criterion considers all matches within the same range of distances to be equal.
`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Maximum radius for a search around a central location.

This parameter works in combination with the `+"`"+`aroundLatLng`+"`"+` and `+"`"+`aroundLatLngViaIP`+"`"+` parameters.
By default, the search radius is determined automatically from the density of hits around the central location.
The search radius is small if there are many hits close to the central coordinates.
`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`Whether the best matching attribute should be determined by minimum proximity.

This setting only affects ranking if the Attribute ranking criterion comes before Proximity in the `+"`"+`ranking`+"`"+` setting.
If true, the best matching attribute is selected based on the minimum proximity of multiple matches.
Otherwise, the best matching attribute is determined by the order in the `+"`"+`searchableAttributes`+"`"+` setting.
`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight.

By default, all searchable attributes are highlighted.
Use `+"`"+`*`+"`"+` to highlight all attributes or use an empty array `+"`"+`[]`+"`"+` to turn off highlighting.

With highlighting, strings that match the search query are surrounded by HTML tags defined by `+"`"+`highlightPreTag`+"`"+` and `+"`"+`highlightPostTag`+"`"+`.
You can use this to visually highlight matching parts of a search query in your UI.

For more information, see [Highlighting and snippeting](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/highlighting-snippeting/js/).
`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response.

To reduce the size of your response, you can retrieve only some of the attributes.

- `+"`"+`*`+"`"+` retrieves all attributes, except attributes included in the `+"`"+`customRanking`+"`"+` and `+"`"+`unretrievableAttributes`+"`"+` settings.
- To retrieve all attributes except a specific one, prefix the attribute with a dash and combine it with the `+"`"+`*`+"`"+`: `+"`"+`["*", "-ATTRIBUTE"]`+"`"+`.
- The `+"`"+`objectID`+"`"+` attribute is always included.
`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes for which to enable snippets.

Snippets provide additional context to matched words.
If you enable snippets, they include 10 words, including the matched word.
The matched word will also be wrapped by HTML tags for highlighting.
You can adjust the number of words with the following notation: `+"`"+`ATTRIBUTE:NUMBER`+"`"+`,
where `+"`"+`NUMBER`+"`"+` is the number of words to be extracted.
`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`Whether to include a `+"`"+`queryID`+"`"+` attribute in the response.

The query ID is a unique identifier for a search query and is required for tracking [click and conversion events](https://www.algolia.com/guides/sending-events/getting-started/).
`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().String("cursor", "", heredoc.Doc(`Cursor to get the next page of the response.

The parameter must match the value returned in the response of a previous request.
The last page of the response does not return a `+"`"+`cursor`+"`"+` attribute.
`))
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Attributes to use as [custom ranking](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/).

The custom ranking attributes decide which items are shown first if the other ranking criteria are equal.

Records with missing values for your selected custom ranking attributes are always sorted last.
Boolean attributes are sorted based on their alphabetical order.

**Modifiers**

<dl>
<dt><code>asc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in ascending order.</dd>
<dt><code>desc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in descending order.</dd>
</dl>

If you use two or more custom ranking attributes, [reduce the precision](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/how-to/controlling-custom-ranking-metrics-precision/) of your first attributes,
or the other attributes will never be applied.
`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Whether to split compound words into their building blocks.

For more information, see [Word segmentation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/#splitting-compound-words).
Word segmentation is supported for these languages: German, Dutch, Finnish, Swedish, and Norwegian.
`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Searchable attributes for which you want to [turn off the Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).

This can be useful for attributes with long values, where the likelyhood of an exact match is high,
such as product descriptions.
Turning off the Exact ranking criterion for these attributes favors exact matching on other attributes.
This reduces the impact of individual attributes with a lot of content on ranking.
`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/).

Returning only exact matches can help when:

- [Searching in hyphenated attributes](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/how-to/how-to-search-in-hyphenated-attributes/).
- Reducing the number of matches when you have too many.
  This can happen with attributes that are long blocks of text, such as product descriptions.

Consider alternatives such as `+"`"+`disableTypoToleranceOnWords`+"`"+` or adding synonyms if your attributes have intentional unusual spellings that might look like typos.
`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Determines how many records of a group are included in the search results.

Records with the same value for the `+"`"+`attributeForDistinct`+"`"+` attribute are considered a group.
The `+"`"+`distinct`+"`"+` setting controls how many members of the group are returned.
This is useful for [deduplication and grouping](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/#introducing-algolias-distinct-feature).

The `+"`"+`distinct`+"`"+` setting is ignored if `+"`"+`attributeForDistinct`+"`"+` is not set.
`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`Whether to enable A/B testing for this search.`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Whether to enable Personalization.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Whether this search will use [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/).

This setting only has an effect if you activated Dynamic Re-Ranking for this index in the Algolia dashboard.
`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether to enable rules.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the [Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes) is computed when the search query has only one word.

<dl>
<dt><code>attribute</code></dt>
<dd>
The Exact ranking criterion is 1 if the query word and attribute value are the same.
For example, a search for "road" will match the value "road", but not "road trip".
</dd>
<dt><code>none</code></dt>
<dd>
The Exact ranking criterion is ignored on single-word searches.
</dd>
<dt><code>word</code></dt>
<dd>
The Exact ranking criterion is 1 if the query word is found in the attribute value.
The query word must have at least 3 characters and must not be a stop word.
</dd>
</dl>

If `+"`"+`exactOnSingleWordQuery`+"`"+` is `+"`"+`word`+"`"+`, only exact matches will be highlighted, partial and prefix matches won't.
 One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter the search by facet values, so that only records with the same facet values are retrieved.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

- `+"`"+`[filter1, filter2]`+"`"+` is interpreted as `+"`"+`filter1 AND filter2`+"`"+`.
- `+"`"+`[[filter1, filter2], filter3]`+"`"+` is interpreted as `+"`"+`filter1 OR filter2 AND filter3`+"`"+`.
- `+"`"+`facet:-value`+"`"+` is interpreted as `+"`"+`NOT facet:value`+"`"+`.

While it's best to avoid attributes that start with a `+"`"+`-`+"`"+`, you can still filter them by escaping with a backslash:
`+"`"+`facet:\-value`+"`"+`.
`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`Whether faceting should be applied after deduplication with `+"`"+`distinct`+"`"+`.

This leads to accurate facet counts when using faceting in combination with `+"`"+`distinct`+"`"+`.
It's usually better to use `+"`"+`afterDistinct`+"`"+` modifiers in the `+"`"+`attributesForFaceting`+"`"+` setting,
as `+"`"+`facetingAfterDistinct`+"`"+` only computes correct facet counts if all records have the same facet values for the `+"`"+`attributeForDistinct`+"`"+`.
`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`Facets for which to retrieve facet values that match the search criteria and the number of matching facet values.

To retrieve all facets, use the wildcard character `+"`"+`*`+"`"+`.
For more information, see [facets](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#contextual-facet-values-and-counts).
`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter the search so that only records with matching values are included in the results.

These filters are supported:

- **Numeric filters.** `+"`"+`<facet> <op> <number>`+"`"+`, where `+"`"+`<op>`+"`"+` is one of `+"`"+`<`+"`"+`, `+"`"+`<=`+"`"+`, `+"`"+`=`+"`"+`, `+"`"+`!=`+"`"+`, `+"`"+`>`+"`"+`, `+"`"+`>=`+"`"+`.
- **Ranges.** `+"`"+`<facet>:<lower> TO <upper>`+"`"+` where `+"`"+`<lower>`+"`"+` and `+"`"+`<upper>`+"`"+` are the lower and upper limits of the range (inclusive).
- **Facet filters.** `+"`"+`<facet>:<value>`+"`"+` where `+"`"+`<facet>`+"`"+` is a facet attribute (case-sensitive) and `+"`"+`<value>`+"`"+` a facet value.
- **Tag filters.** `+"`"+`_tags:<value>`+"`"+` or just `+"`"+`<value>`+"`"+` (case-sensitive).
- **Boolean filters.** `+"`"+`<facet>: true | false`+"`"+`.

You can combine filters with `+"`"+`AND`+"`"+`, `+"`"+`OR`+"`"+`, and `+"`"+`NOT`+"`"+` operators with the following restrictions:

- You can only combine filters of the same type with `+"`"+`OR`+"`"+`.
  **Not supported:** `+"`"+`facet:value OR num > 3`+"`"+`.
- You can't use `+"`"+`NOT`+"`"+` with combinations of filters.
  **Not supported:** `+"`"+`NOT(facet:value OR facet:value)`+"`"+`
- You can't combine conjunctions (`+"`"+`AND`+"`"+`) with `+"`"+`OR`+"`"+`.
  **Not supported:** `+"`"+`facet:value OR (facet:value AND facet:value)`+"`"+`

Use quotes around your filters, if the facet attribute name or facet value has spaces, keywords (`+"`"+`OR`+"`"+`, `+"`"+`AND`+"`"+`, `+"`"+`NOT`+"`"+`), or quotes.
If a facet attribute is an array, the filter matches if it matches at least one element of the array.

For more information, see [Filters](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/).
`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("getRankingInfo", false, heredoc.Doc(`Whether the search response should include detailed ranking information.`))
	cmd.Flags().SetAnnotation("getRankingInfo", "Categories", []string{"Advanced"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`HTML tag to insert after the highlighted parts in all highlighted results and snippets.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`HTML tag to insert before the highlighted parts in all highlighted results and snippets.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treat singular, plurals, and other forms of declensions as equivalent.
You should only use this feature for the languages used in your index.
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters for which diacritics should be preserved.

By default, Algolia removes diacritics from letters.
For example, `+"`"+`é`+"`"+` becomes `+"`"+`e`+"`"+`. If this causes issues in your search,
you can specify characters that should keep their diacritics.
`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`Number of hits to retrieve (used in combination with `+"`"+`offset`+"`"+`).`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet values to return when [searching for facet values](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#search-for-facet-values).`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Minimum proximity score for two matching words.

This adjusts the [Proximity ranking criterion](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#proximity)
by equally scoring matches that are farther apart.

For example, if `+"`"+`minProximity`+"`"+` is 2, neighboring matches and matches with one word between them would have the same score.
`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with [one typo](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with [two typos](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`Minimum radius (in meters) for a search around a location when `+"`"+`aroundRadius`+"`"+` isn't set.`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results.

This setting only applies to indices, for which Algolia enabled NeuralSearch for you.
 One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`ISO language codes that adjust settings that are useful for processing natural language queries (as opposed to keyword searches):

- Sets `+"`"+`removeStopWords`+"`"+` and `+"`"+`ignorePlurals`+"`"+` to the list of provided languages.
- Sets `+"`"+`removeWordsIfNoResults`+"`"+` to `+"`"+`allOptional`+"`"+`.
- Adds a `+"`"+`natural_language`+"`"+` attribute to `+"`"+`ruleContexts`+"`"+` and `+"`"+`analyticsTags`+"`"+`.
`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter by numeric facets.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

You can use numeric comparison operators: `+"`"+`<`+"`"+`, `+"`"+`<=`+"`"+`, `+"`"+`=`+"`"+`, `+"`"+`!=`+"`"+`, `+"`"+`>`+"`"+`, `+"`"+`>=`+"`"+`. Comparsions are precise up to 3 decimals.
You can also provide ranges: `+"`"+`facet:<lower> TO <upper>`+"`"+`. The range includes the lower and upper boundaries.
The same combination rules apply as for `+"`"+`facetFilters`+"`"+`.
`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`Position of the first hit to retrieve.`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`Filters to promote or demote records in the search results.

Optional filters work like facet filters, but they don't exclude records from the search results.
Records that match the optional filter rank before records that don't match.
If you're using a negative filter `+"`"+`facet:-value`+"`"+`, matching records rank after records that don't match.

- Optional filters don't work on virtual replicas.
- Optional filters are applied _after_ sort-by attributes.
- Optional filters don't work with numeric attributes.
`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words that should be considered optional when found in the query.

By default, records must match all words in the search query to be included in the search results.
Adding optional words can help to increase the number of search results by running an additional search query that doesn't include the optional words.
For example, if the search query is "action video" and "video" is an optional word,
the search engine runs two queries. One for "action video" and one for "action".
Records that match all words are ranked higher.

For a search query with 4 or more words **and** all its words are optional,
the number of matched words required for a record to be included in the search results increases for every 1,000 records:

- If `+"`"+`optionalWords`+"`"+` has less than 10 words, the required number of matched words increases by 1:
  results 1 to 1,000 require 1 matched word, results 1,001 to 2000 need 2 matched words.
- If `+"`"+`optionalWords`+"`"+` has 10 or more words, the number of required matched words increases by the number of optional words dividied by 5 (rounded down).
  For example, with 18 optional words: results 1 to 1,000 require 1 matched word, results 1,001 to 2000 need 4 matched words.

For more information, see [Optional words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/#creating-a-list-of-optional-words).
`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`Page of search results to retrieve.`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`Whether to include this search when calculating processing-time percentiles.`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`Impact that Personalization should have on this search.

The higher this value is, the more Personalization determines the ranking compared to other factors.
For more information, see [Understanding Personalization impact](https://www.algolia.com/doc/guides/personalization/personalizing-results/in-depth/configuring-personalization/#understanding-personalization-impact).
`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`Search query.`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Languages for language-specific query processing steps such as plurals, stop-word removal, and word-detection dictionaries.

This setting sets a default list of languages used by the `+"`"+`removeStopWords`+"`"+` and `+"`"+`ignorePlurals`+"`"+` settings.
This setting also sets a dictionary for word detection in the logogram-based [CJK](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/#normalization-for-logogram-based-languages-cjk) languages.
To support this, you must place the CJK language **first**.

**You should always specify a query language.**
If you don't specify an indexing language, the search engine uses all [supported languages](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/),
or the languages you specified with the `+"`"+`ignorePlurals`+"`"+` or `+"`"+`removeStopWords`+"`"+` parameters.
This can lead to unexpected search results.
For more information, see [Language-specific configuration](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/).
`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines if and how query words are interpreted as prefixes.

By default, only the last query word is treated as prefix (`+"`"+`prefixLast`+"`"+`).
To turn off prefix search, use `+"`"+`prefixNone`+"`"+`.
Avoid `+"`"+`prefixAll`+"`"+`, which treats all query words as prefixes.
This might lead to counterintuitive results and makes your search slower.

For more information, see [Prefix searching](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/prefix-searching/).
 One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia returns your results.

By default, each entry corresponds to a [ranking criteria](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/).
The tie-breaking algorithm sequentially applies each criterion in the order they're specified.
If you configure a replica index for [sorting by an attribute](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/how-to/sort-by-attribute/),
you put the sorting attribute at the top of the list.

**Modifiers**

<dl>
<dt><code>asc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in ascending order.</dd>
<dt><code>desc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in descending order.</dd>
</dl>

Before you modify the default setting,
you should test your changes in the dashboard,
and by [A/B testing](https://www.algolia.com/doc/guides/ab-testing/what-is-ab-testing/).
`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string", "null"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`Restrict [Dynamic Re-ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/) to records that match these filters.
`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.

You can only set `+"`"+`relevancyStrictness`+"`"+` on [virtual replica indices](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/in-depth/replicas/#what-are-virtual-replicas).
Use this setting to strike a balance between the relevance and number of returned results.
`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop words from the search query.

Stop words are common words like articles, conjunctions, prepositions, or pronouns that have little or no meaning on their own.
In English, "the", "a", or "and" are stop words.

You should only use this feature for the languages used in your index.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy for removing words from the query when it doesn't return any results.
This helps to avoid returning empty search results.

<dl>
<dt><code>none</code></dt>
<dd>No words are removed when a query doesn't return results.</dd>
<dt><code>lastWords</code></dt>
<dd>Treat the last (then second to last, then third to last) word as optional, until there are results or at most 5 words have been removed.</dd>
<dt><code>firstWords</code></dt>
<dd>Treat the first (then second, then third) word as optional, until there are results or at most 5 words have been removed.</dd>
<dt><code>allOptional</code></dt>
<dd>Treat all words as optional.</dd>
</dl>

For more information, see [Remove words to improve results](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/in-depth/why-use-remove-words-if-no-results/).
 One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra data that can be used in the search UI.

You can use this to control aspects of your search UI, such as, the order of facet names and values
without changing your frontend code.
`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to replace a highlighted word with the matched synonym.

By default, the original words are highlighted even if a synonym matches.
For example, with `+"`"+`home`+"`"+` as a synonym for `+"`"+`house`+"`"+` and a search for `+"`"+`home`+"`"+`,
records matching either "home" or "house" are included in the search results,
and either "home" or "house" are highlighted.

With `+"`"+`replaceSynonymsInHighlight`+"`"+` set to `+"`"+`true`+"`"+`, a search for `+"`"+`home`+"`"+` still matches the same records,
but all occurences of "house" are replaced by "home" in the highlighted response.
`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{"*"}, heredoc.Doc(`Properties to include in the API response of `+"`"+`search`+"`"+` and `+"`"+`browse`+"`"+` requests.

By default, all response properties are included.
To reduce the response size, you can select, which attributes should be included.

You can't exclude these properties:
`+"`"+`message`+"`"+`, `+"`"+`warning`+"`"+`, `+"`"+`cursor`+"`"+`, `+"`"+`serverUsed`+"`"+`, `+"`"+`indexUsed`+"`"+`,
`+"`"+`abTestVariantID`+"`"+`, `+"`"+`parsedQuery`+"`"+`, or any property triggered by the `+"`"+`getRankingInfo`+"`"+` parameter.

Don't exclude properties that you might need in your search UI.
`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Whether to restrict highlighting and snippeting to items that at least partially matched the search query.
By default, all items are highlighted and snippeted.
`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a search to a subset of your searchable attributes.`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`Assigns a rule context to the search query.

[Rule contexts](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/how-to/customize-search-results-by-platform/#whats-a-context) are strings that you can use to trigger matching rules.
`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch.
Only used when `+"`"+`mode`+"`"+` is `+"`"+`neuralSearch`+"`"+`.
`))
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`Keywords to be used instead of the search query to conduct a more broader search.

Using the `+"`"+`similarQuery`+"`"+` parameter changes other settings:

- `+"`"+`queryType`+"`"+` is set to `+"`"+`prefixNone`+"`"+`.
- `+"`"+`removeStopWords`+"`"+` is set to true.
- `+"`"+`words`+"`"+` is set as the first ranking criterion.
- All remaining words are treated as `+"`"+`optionalWords`+"`"+`.

Since the `+"`"+`similarQuery`+"`"+` is supposed to do a broad search, they usually return many results.
Combine it with `+"`"+`filters`+"`"+` to narrow down the list of results.
`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Order in which to retrieve facet values.

<dl>
<dt><code>count</code></dt>
<dd>
Facet values are retrieved by decreasing count.
The count is the number of matching records containing this facet value.
</dd>
<dt><code>alpha</code></dt>
<dd>Retrieve facet values alphabetically.</dd>
</dl>

This setting doesn't influence how facet values are displayed in your UI (see `+"`"+`renderingContent`+"`"+`).
For more information, see [facet value display](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/facet-display/js/).
`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`Whether to sum all filter scores.

If true, all filter scores are summed.
Otherwise, the maximum filter score is kept.
For more information, see [filter scores](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/in-depth/filter-scoring/#accumulating-scores-with-sumorfiltersscores).
`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for this search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter the search by values of the special `+"`"+`_tags`+"`"+` attribute.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

Different from regular facets, `+"`"+`_tags`+"`"+` can only be used for filtering (including or excluding records).
You won't get a facet count.
The same combination and escaping rules apply as for `+"`"+`facetFilters`+"`"+`.
`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Whether [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/) is enabled and how it is applied.

If typo tolerance is true, `+"`"+`min`+"`"+`, or `+"`"+`strict`+"`"+`, [word splitting and concetenation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/splitting-and-concatenation/) is also active.
`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`Unique pseudonymous or anonymous user identifier.

This helps with analytics and click and conversion events.
For more information, see [user token](https://www.algolia.com/doc/guides/sending-events/concepts/usertoken/).
`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

func AddDeleteByParamsFlags(cmd *cobra.Command) {
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Coordinates for the center of a circle, expressed as a comma-separated string of latitude and longitude.

Only records included within circle around this central location are included in the results.
The radius of the circle is determined by the `+"`"+`aroundRadius`+"`"+` and `+"`"+`minimumAroundRadius`+"`"+` settings.
This parameter is ignored if you also specify `+"`"+`insidePolygon`+"`"+` or `+"`"+`insideBoundingBox`+"`"+`.
`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Maximum radius for a search around a central location.

This parameter works in combination with the `+"`"+`aroundLatLng`+"`"+` and `+"`"+`aroundLatLngViaIP`+"`"+` parameters.
By default, the search radius is determined automatically from the density of hits around the central location.
The search radius is small if there are many hits close to the central coordinates.
`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter the search by facet values, so that only records with the same facet values are retrieved.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

- `+"`"+`[filter1, filter2]`+"`"+` is interpreted as `+"`"+`filter1 AND filter2`+"`"+`.
- `+"`"+`[[filter1, filter2], filter3]`+"`"+` is interpreted as `+"`"+`filter1 OR filter2 AND filter3`+"`"+`.
- `+"`"+`facet:-value`+"`"+` is interpreted as `+"`"+`NOT facet:value`+"`"+`.

While it's best to avoid attributes that start with a `+"`"+`-`+"`"+`, you can still filter them by escaping with a backslash:
`+"`"+`facet:\-value`+"`"+`.
`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter the search so that only records with matching values are included in the results.

These filters are supported:

- **Numeric filters.** `+"`"+`<facet> <op> <number>`+"`"+`, where `+"`"+`<op>`+"`"+` is one of `+"`"+`<`+"`"+`, `+"`"+`<=`+"`"+`, `+"`"+`=`+"`"+`, `+"`"+`!=`+"`"+`, `+"`"+`>`+"`"+`, `+"`"+`>=`+"`"+`.
- **Ranges.** `+"`"+`<facet>:<lower> TO <upper>`+"`"+` where `+"`"+`<lower>`+"`"+` and `+"`"+`<upper>`+"`"+` are the lower and upper limits of the range (inclusive).
- **Facet filters.** `+"`"+`<facet>:<value>`+"`"+` where `+"`"+`<facet>`+"`"+` is a facet attribute (case-sensitive) and `+"`"+`<value>`+"`"+` a facet value.
- **Tag filters.** `+"`"+`_tags:<value>`+"`"+` or just `+"`"+`<value>`+"`"+` (case-sensitive).
- **Boolean filters.** `+"`"+`<facet>: true | false`+"`"+`.

You can combine filters with `+"`"+`AND`+"`"+`, `+"`"+`OR`+"`"+`, and `+"`"+`NOT`+"`"+` operators with the following restrictions:

- You can only combine filters of the same type with `+"`"+`OR`+"`"+`.
  **Not supported:** `+"`"+`facet:value OR num > 3`+"`"+`.
- You can't use `+"`"+`NOT`+"`"+` with combinations of filters.
  **Not supported:** `+"`"+`NOT(facet:value OR facet:value)`+"`"+`
- You can't combine conjunctions (`+"`"+`AND`+"`"+`) with `+"`"+`OR`+"`"+`.
  **Not supported:** `+"`"+`facet:value OR (facet:value AND facet:value)`+"`"+`

Use quotes around your filters, if the facet attribute name or facet value has spaces, keywords (`+"`"+`OR`+"`"+`, `+"`"+`AND`+"`"+`, `+"`"+`NOT`+"`"+`), or quotes.
If a facet attribute is an array, the filter matches if it matches at least one element of the array.

For more information, see [Filters](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/).
`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter by numeric facets.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

You can use numeric comparison operators: `+"`"+`<`+"`"+`, `+"`"+`<=`+"`"+`, `+"`"+`=`+"`"+`, `+"`"+`!=`+"`"+`, `+"`"+`>`+"`"+`, `+"`"+`>=`+"`"+`. Comparsions are precise up to 3 decimals.
You can also provide ranges: `+"`"+`facet:<lower> TO <upper>`+"`"+`. The range includes the lower and upper boundaries.
The same combination rules apply as for `+"`"+`facetFilters`+"`"+`.
`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter the search by values of the special `+"`"+`_tags`+"`"+` attribute.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

Different from regular facets, `+"`"+`_tags`+"`"+` can only be used for filtering (including or excluding records).
You won't get a facet count.
The same combination and escaping rules apply as for `+"`"+`facetFilters`+"`"+`.
`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
}

func AddIndexSettingsFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Whether to support phrase matching and excluding words from search queries.

Use the `+"`"+`advancedSyntaxFeatures`+"`"+` parameter to control which feature is supported.
`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Advanced search syntax features you want to support.

<dl>
<dt><code>exactPhrase</code></dt>
<dd>

Phrases in quotes must match exactly.
For example, `+"`"+`sparkly blue "iPhone case"`+"`"+` only returns records with the exact string "iPhone case".

</dd>
<dt><code>excludeWords</code></dt>
<dd>

Query words prefixed with a `+"`"+`-`+"`"+` must not occur in a record.
For example, `+"`"+`search -engine`+"`"+` matches records that contain "search" but not "engine".

</dd>
</dl>

This setting only has an effect if `+"`"+`advancedSyntax`+"`"+` is true.
`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowCompressionOfIntegerArray", false, heredoc.Doc(`Whether arrays with exclusively non-negative integers should be compressed for better performance.
If true, the compressed arrays may be reordered.
`))
	cmd.Flags().SetAnnotation("allowCompressionOfIntegerArray", "Categories", []string{"Performance"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers in the search query.

Turn off this setting to reduce the number of irrelevant matches
when searching in large sets of similar numbers.
`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives of query words that should be considered as exact matches by the Exact ranking criterion.

<dl>
<dt><code>ignorePlurals</code></dt>
<dd>

Plurals and similar declensions added by the `+"`"+`ignorePlurals`+"`"+` setting are considered exact matches.

</dd>
<dt><code>singleWordSynonym</code></dt>
<dd>
Single-word synonyms, such as "NY/NYC" are considered exact matches.
</dd>
<dt><code>multiWordsSynonym</code></dt>
<dd>
Multi-word synonyms, such as "NY/New York" are considered exact matches.
</dd>
</dl>.
`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`Whether the best matching attribute should be determined by minimum proximity.

This setting only affects ranking if the Attribute ranking criterion comes before Proximity in the `+"`"+`ranking`+"`"+` setting.
If true, the best matching attribute is selected based on the minimum proximity of multiple matches.
Otherwise, the best matching attribute is determined by the order in the `+"`"+`searchableAttributes`+"`"+` setting.
`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().String("attributeForDistinct", "", heredoc.Doc(`Attribute that should be used to establish groups of results.

All records with the same value for this attribute are considered a group.
You can combine `+"`"+`attributeForDistinct`+"`"+` with the `+"`"+`distinct`+"`"+` search parameter to control
how many items per group are included in the search results.

If you want to use the same attribute also for faceting, use the `+"`"+`afterDistinct`+"`"+` modifier of the `+"`"+`attributesForFaceting`+"`"+` setting.
This applies faceting _after_ deduplication, which will result in accurate facet counts.
`))
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`Attributes used for [faceting](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/).

Facets are ways to categorize search results based on attributes.
Facets can be used to let user filter search results.
By default, no attribute is used for faceting.

**Modifiers**

<dl>
<dt><code>filterOnly("ATTRIBUTE")</code></dt>
<dd>Allows using this attribute as a filter, but doesn't evalue the facet values.</dd>
<dt><code>searchable("ATTRIBUTE")</code></dt>
<dd>Allows searching for facet values.</dd>
<dt><code>afterDistinct("ATTRIBUTE")</code></dt>
<dd>

Evaluates the facet count _after_ deduplication with `+"`"+`distinct`+"`"+`.
This ensures accurate facet counts.
You can apply this modifier to searchable facets: `+"`"+`afterDistinct(searchable(ATTRIBUTE))`+"`"+`.

</dd>
</dl>

Without modifiers, the attribute is used as a regular facet.
`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight.

By default, all searchable attributes are highlighted.
Use `+"`"+`*`+"`"+` to highlight all attributes or use an empty array `+"`"+`[]`+"`"+` to turn off highlighting.

With highlighting, strings that match the search query are surrounded by HTML tags defined by `+"`"+`highlightPreTag`+"`"+` and `+"`"+`highlightPostTag`+"`"+`.
You can use this to visually highlight matching parts of a search query in your UI.

For more information, see [Highlighting and snippeting](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/highlighting-snippeting/js/).
`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response.

To reduce the size of your response, you can retrieve only some of the attributes.

- `+"`"+`*`+"`"+` retrieves all attributes, except attributes included in the `+"`"+`customRanking`+"`"+` and `+"`"+`unretrievableAttributes`+"`"+` settings.
- To retrieve all attributes except a specific one, prefix the attribute with a dash and combine it with the `+"`"+`*`+"`"+`: `+"`"+`["*", "-ATTRIBUTE"]`+"`"+`.
- The `+"`"+`objectID`+"`"+` attribute is always included.
`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes for which to enable snippets.

Snippets provide additional context to matched words.
If you enable snippets, they include 10 words, including the matched word.
The matched word will also be wrapped by HTML tags for highlighting.
You can adjust the number of words with the following notation: `+"`"+`ATTRIBUTE:NUMBER`+"`"+`,
where `+"`"+`NUMBER`+"`"+` is the number of words to be extracted.
`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToTransliterate", []string{}, heredoc.Doc(`Attributes, for which you want to support [Japanese transliteration](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/#japanese-transliteration-and-type-ahead).

Transliteration supports searching in any of the Japanese writing systems.
To support transliteration, you must set the indexing language to Japanese.
`))
	cmd.Flags().SetAnnotation("attributesToTransliterate", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("camelCaseAttributes", []string{}, heredoc.Doc(`Attributes for which to split [camel case](https://wikipedia.org/wiki/Camel_case) words.`))
	cmd.Flags().SetAnnotation("camelCaseAttributes", "Categories", []string{"Languages"})
	customNormalization := NewJSONVar([]string{}...)
	cmd.Flags().Var(customNormalization, "customNormalization", heredoc.Doc(`Characters and their normalized replacements.
This overrides Algolia's default [normalization](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/).
`))
	cmd.Flags().SetAnnotation("customNormalization", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Attributes to use as [custom ranking](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/).

The custom ranking attributes decide which items are shown first if the other ranking criteria are equal.

Records with missing values for your selected custom ranking attributes are always sorted last.
Boolean attributes are sorted based on their alphabetical order.

**Modifiers**

<dl>
<dt><code>asc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in ascending order.</dd>
<dt><code>desc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in descending order.</dd>
</dl>

If you use two or more custom ranking attributes, [reduce the precision](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/how-to/controlling-custom-ranking-metrics-precision/) of your first attributes,
or the other attributes will never be applied.
`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Whether to split compound words into their building blocks.

For more information, see [Word segmentation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/#splitting-compound-words).
Word segmentation is supported for these languages: German, Dutch, Finnish, Swedish, and Norwegian.
`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	decompoundedAttributes := NewJSONVar([]string{}...)
	cmd.Flags().Var(decompoundedAttributes, "decompoundedAttributes", heredoc.Doc(`Searchable attributes to which Algolia should apply [word segmentation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/how-to/customize-segmentation/) (decompounding).

Compound words are formed by combining two or more individual words,
and are particularly prevalent in Germanic languages—for example, "firefighter".
With decompounding, the individual components are indexed separately.

You can specify different lists for different languages.
Decompounding is supported for these languages:
Dutch (`+"`"+`nl`+"`"+`), German (`+"`"+`de`+"`"+`), Finnish (`+"`"+`fi`+"`"+`), Danish (`+"`"+`da`+"`"+`), Swedish (`+"`"+`sv`+"`"+`), and Norwegian (`+"`"+`no`+"`"+`).
`))
	cmd.Flags().SetAnnotation("decompoundedAttributes", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Searchable attributes for which you want to [turn off the Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).

This can be useful for attributes with long values, where the likelyhood of an exact match is high,
such as product descriptions.
Turning off the Exact ranking criterion for these attributes favors exact matching on other attributes.
This reduces the impact of individual attributes with a lot of content on ranking.
`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disablePrefixOnAttributes", []string{}, heredoc.Doc(`Searchable attributes for which you want to turn off [prefix matching](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/#adjusting-prefix-search).`))
	cmd.Flags().SetAnnotation("disablePrefixOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/).

Returning only exact matches can help when:

- [Searching in hyphenated attributes](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/how-to/how-to-search-in-hyphenated-attributes/).
- Reducing the number of matches when you have too many.
  This can happen with attributes that are long blocks of text, such as product descriptions.

Consider alternatives such as `+"`"+`disableTypoToleranceOnWords`+"`"+` or adding synonyms if your attributes have intentional unusual spellings that might look like typos.
`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("disableTypoToleranceOnWords", []string{}, heredoc.Doc(`Words for which you want to turn off [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/).
This also turns off [word splitting and concatenation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/splitting-and-concatenation/) for the specified words.
`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnWords", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Determines how many records of a group are included in the search results.

Records with the same value for the `+"`"+`attributeForDistinct`+"`"+` attribute are considered a group.
The `+"`"+`distinct`+"`"+` setting controls how many members of the group are returned.
This is useful for [deduplication and grouping](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/#introducing-algolias-distinct-feature).

The `+"`"+`distinct`+"`"+` setting is ignored if `+"`"+`attributeForDistinct`+"`"+` is not set.
`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Whether to enable Personalization.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Whether this search will use [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/).

This setting only has an effect if you activated Dynamic Re-Ranking for this index in the Algolia dashboard.
`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether to enable rules.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the [Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes) is computed when the search query has only one word.

<dl>
<dt><code>attribute</code></dt>
<dd>
The Exact ranking criterion is 1 if the query word and attribute value are the same.
For example, a search for "road" will match the value "road", but not "road trip".
</dd>
<dt><code>none</code></dt>
<dd>
The Exact ranking criterion is ignored on single-word searches.
</dd>
<dt><code>word</code></dt>
<dd>
The Exact ranking criterion is 1 if the query word is found in the attribute value.
The query word must have at least 3 characters and must not be a stop word.
</dd>
</dl>

If `+"`"+`exactOnSingleWordQuery`+"`"+` is `+"`"+`word`+"`"+`, only exact matches will be highlighted, partial and prefix matches won't.
 One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`HTML tag to insert after the highlighted parts in all highlighted results and snippets.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`HTML tag to insert before the highlighted parts in all highlighted results and snippets.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treat singular, plurals, and other forms of declensions as equivalent.
You should only use this feature for the languages used in your index.
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("indexLanguages", []string{}, heredoc.Doc(`Languages for language-specific processing steps, such as word detection and dictionary settings.

**You should always specify an indexing language.**
If you don't specify an indexing language, the search engine uses all [supported languages](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/),
or the languages you specified with the `+"`"+`ignorePlurals`+"`"+` or `+"`"+`removeStopWords`+"`"+` parameters.
This can lead to unexpected search results.
For more information, see [Language-specific configuration](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/).
`))
	cmd.Flags().SetAnnotation("indexLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters for which diacritics should be preserved.

By default, Algolia removes diacritics from letters.
For example, `+"`"+`é`+"`"+` becomes `+"`"+`e`+"`"+`. If this causes issues in your search,
you can specify characters that should keep their diacritics.
`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet values to return when [searching for facet values](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#search-for-facet-values).`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Minimum proximity score for two matching words.

This adjusts the [Proximity ranking criterion](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#proximity)
by equally scoring matches that are farther apart.

For example, if `+"`"+`minProximity`+"`"+` is 2, neighboring matches and matches with one word between them would have the same score.
`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with [one typo](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with [two typos](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results.

This setting only applies to indices, for which Algolia enabled NeuralSearch for you.
 One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("numericAttributesForFiltering", []string{}, heredoc.Doc(`Numeric attributes that can be used as [numerical filters](https://www.algolia.com/doc/guides/managing-results/rules/detecting-intent/how-to/applying-a-custom-filter-for-a-specific-query/#numerical-filters).

By default, all numeric attributes are available as numerical filters.
For faster indexing, reduce the number of numeric attributes.

If you want to turn off filtering for all numeric attributes, specifiy an attribute that doesn't exist in your index, such as `+"`"+`NO_NUMERIC_FILTERING`+"`"+`.

**Modifier**

<dl>
<dt><code>equalOnly("ATTRIBUTE")</code></dt>
<dd>

Support only filtering based on equality comparisons `+"`"+`=`+"`"+` and `+"`"+`!=`+"`"+`.

</dd>
</dl>

Without modifier, all numeric comparisons are supported.
`))
	cmd.Flags().SetAnnotation("numericAttributesForFiltering", "Categories", []string{"Performance"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words that should be considered optional when found in the query.

By default, records must match all words in the search query to be included in the search results.
Adding optional words can help to increase the number of search results by running an additional search query that doesn't include the optional words.
For example, if the search query is "action video" and "video" is an optional word,
the search engine runs two queries. One for "action video" and one for "action".
Records that match all words are ranked higher.

For a search query with 4 or more words **and** all its words are optional,
the number of matched words required for a record to be included in the search results increases for every 1,000 records:

- If `+"`"+`optionalWords`+"`"+` has less than 10 words, the required number of matched words increases by 1:
  results 1 to 1,000 require 1 matched word, results 1,001 to 2000 need 2 matched words.
- If `+"`"+`optionalWords`+"`"+` has 10 or more words, the number of required matched words increases by the number of optional words dividied by 5 (rounded down).
  For example, with 18 optional words: results 1 to 1,000 require 1 matched word, results 1,001 to 2000 need 4 matched words.

For more information, see [Optional words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/#creating-a-list-of-optional-words).
`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("paginationLimitedTo", 1000, heredoc.Doc(`Maximum number of search results that can be obtained through pagination.

Higher pagination limits might slow down your search.
For pagination limits above 1,000, the sorting of results beyond the 1,000th hit can't be guaranteed.
`))
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Languages for language-specific query processing steps such as plurals, stop-word removal, and word-detection dictionaries.

This setting sets a default list of languages used by the `+"`"+`removeStopWords`+"`"+` and `+"`"+`ignorePlurals`+"`"+` settings.
This setting also sets a dictionary for word detection in the logogram-based [CJK](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/#normalization-for-logogram-based-languages-cjk) languages.
To support this, you must place the CJK language **first**.

**You should always specify a query language.**
If you don't specify an indexing language, the search engine uses all [supported languages](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/),
or the languages you specified with the `+"`"+`ignorePlurals`+"`"+` or `+"`"+`removeStopWords`+"`"+` parameters.
This can lead to unexpected search results.
For more information, see [Language-specific configuration](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/).
`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines if and how query words are interpreted as prefixes.

By default, only the last query word is treated as prefix (`+"`"+`prefixLast`+"`"+`).
To turn off prefix search, use `+"`"+`prefixNone`+"`"+`.
Avoid `+"`"+`prefixAll`+"`"+`, which treats all query words as prefixes.
This might lead to counterintuitive results and makes your search slower.

For more information, see [Prefix searching](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/prefix-searching/).
 One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia returns your results.

By default, each entry corresponds to a [ranking criteria](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/).
The tie-breaking algorithm sequentially applies each criterion in the order they're specified.
If you configure a replica index for [sorting by an attribute](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/how-to/sort-by-attribute/),
you put the sorting attribute at the top of the list.

**Modifiers**

<dl>
<dt><code>asc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in ascending order.</dd>
<dt><code>desc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in descending order.</dd>
</dl>

Before you modify the default setting,
you should test your changes in the dashboard,
and by [A/B testing](https://www.algolia.com/doc/guides/ab-testing/what-is-ab-testing/).
`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string", "null"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`Restrict [Dynamic Re-ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/) to records that match these filters.
`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.

You can only set `+"`"+`relevancyStrictness`+"`"+` on [virtual replica indices](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/in-depth/replicas/#what-are-virtual-replicas).
Use this setting to strike a balance between the relevance and number of returned results.
`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop words from the search query.

Stop words are common words like articles, conjunctions, prepositions, or pronouns that have little or no meaning on their own.
In English, "the", "a", or "and" are stop words.

You should only use this feature for the languages used in your index.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy for removing words from the query when it doesn't return any results.
This helps to avoid returning empty search results.

<dl>
<dt><code>none</code></dt>
<dd>No words are removed when a query doesn't return results.</dd>
<dt><code>lastWords</code></dt>
<dd>Treat the last (then second to last, then third to last) word as optional, until there are results or at most 5 words have been removed.</dd>
<dt><code>firstWords</code></dt>
<dd>Treat the first (then second, then third) word as optional, until there are results or at most 5 words have been removed.</dd>
<dt><code>allOptional</code></dt>
<dd>Treat all words as optional.</dd>
</dl>

For more information, see [Remove words to improve results](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/in-depth/why-use-remove-words-if-no-results/).
 One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra data that can be used in the search UI.

You can use this to control aspects of your search UI, such as, the order of facet names and values
without changing your frontend code.
`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to replace a highlighted word with the matched synonym.

By default, the original words are highlighted even if a synonym matches.
For example, with `+"`"+`home`+"`"+` as a synonym for `+"`"+`house`+"`"+` and a search for `+"`"+`home`+"`"+`,
records matching either "home" or "house" are included in the search results,
and either "home" or "house" are highlighted.

With `+"`"+`replaceSynonymsInHighlight`+"`"+` set to `+"`"+`true`+"`"+`, a search for `+"`"+`home`+"`"+` still matches the same records,
but all occurences of "house" are replaced by "home" in the highlighted response.
`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("replicas", []string{}, heredoc.Doc(`Creates [replica indices](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/in-depth/replicas/).

Replicas are copies of a primary index with the same records but different settings, synonyms, or rules.
If you want to offer a different ranking or sorting of your search results, you'll use replica indices.
All index operations on a primary index are automatically forwarded to its replicas.
To add a replica index, you must provide the complete set of replicas to this parameter.
If you omit a replica from this list, the replica turns into a regular, standalone index that will no longer by synced with the primary index.

**Modifier**

<dl>
<dt><code>virtual("REPLICA")</code></dt>
<dd>

Create a virtual replica,
Virtual replicas don't increase the number of records and are optimized for [Relevant sorting](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/in-depth/relevant-sort/).

</dd>
</dl>

Without modifier, a standard replica is created, which duplicates your record count and is used for strict, or [exhaustive sorting](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/in-depth/exhaustive-sort/).
`))
	cmd.Flags().SetAnnotation("replicas", "Categories", []string{"Ranking"})
	cmd.Flags().StringSlice("responseFields", []string{"*"}, heredoc.Doc(`Properties to include in the API response of `+"`"+`search`+"`"+` and `+"`"+`browse`+"`"+` requests.

By default, all response properties are included.
To reduce the response size, you can select, which attributes should be included.

You can't exclude these properties:
`+"`"+`message`+"`"+`, `+"`"+`warning`+"`"+`, `+"`"+`cursor`+"`"+`, `+"`"+`serverUsed`+"`"+`, `+"`"+`indexUsed`+"`"+`,
`+"`"+`abTestVariantID`+"`"+`, `+"`"+`parsedQuery`+"`"+`, or any property triggered by the `+"`"+`getRankingInfo`+"`"+` parameter.

Don't exclude properties that you might need in your search UI.
`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Whether to restrict highlighting and snippeting to items that at least partially matched the search query.
By default, all items are highlighted and snippeted.
`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("searchableAttributes", []string{}, heredoc.Doc(`Attributes used for searching.

By default, all attributes are searchable and the [Attribute](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#attribute) ranking criterion is turned off.
With a non-empty list, Algolia only returns results with matches in the selected attributes.
In addition, the Attribute ranking criterion is turned on: matches in attributes that are higher in the list of `+"`"+`searchableAttributes`+"`"+` rank first.
To make matches in two attributes rank equally, include them in a comma-separated string, such as `+"`"+`"title,alternate_title"`+"`"+`.
Attributes with the same priority are always unordered.

For more information, see [Searchable attributes](https://www.algolia.com/doc/guides/sending-and-managing-data/prepare-your-data/how-to/setting-searchable-attributes/).

**Modifier**

<dl>
<dt><code>unordered("ATTRIBUTE")</code></dt>
<dd>
Ignore the position of a match within the attribute.
</dd>
</dl>

Without modifier, matches at the beginning of an attribute rank higer than matches at the end.
`))
	cmd.Flags().SetAnnotation("searchableAttributes", "Categories", []string{"Attributes"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch.
Only used when `+"`"+`mode`+"`"+` is `+"`"+`neuralSearch`+"`"+`.
`))
	cmd.Flags().String("separatorsToIndex", "", heredoc.Doc(`Controls which separators are indexed.

Separators are all non-letter characters except spaces and currency characters, such as $€£¥.
By default, separator characters aren't indexed.
With `+"`"+`separatorsToIndex`+"`"+`, Algolia treats separator characters as separate words.
For example, a search for `+"`"+`C#`+"`"+` would report two matches.
`))
	cmd.Flags().SetAnnotation("separatorsToIndex", "Categories", []string{"Typos"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Order in which to retrieve facet values.

<dl>
<dt><code>count</code></dt>
<dd>
Facet values are retrieved by decreasing count.
The count is the number of matching records containing this facet value.
</dd>
<dt><code>alpha</code></dt>
<dd>Retrieve facet values alphabetically.</dd>
</dl>

This setting doesn't influence how facet values are displayed in your UI (see `+"`"+`renderingContent`+"`"+`).
For more information, see [facet value display](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/facet-display/js/).
`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Whether [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/) is enabled and how it is applied.

If typo tolerance is true, `+"`"+`min`+"`"+`, or `+"`"+`strict`+"`"+`, [word splitting and concetenation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/splitting-and-concatenation/) is also active.
`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("unretrievableAttributes", []string{}, heredoc.Doc(`Attributes that can't be retrieved at query time.

This can be useful if you want to use an attribute for ranking or to [restrict access](https://www.algolia.com/doc/guides/security/api-keys/how-to/user-restricted-access-to-data/),
but don't want to include it in the search results.
`))
	cmd.Flags().SetAnnotation("unretrievableAttributes", "Categories", []string{"Attributes"})
	userData := NewJSONVar([]string{}...)
	cmd.Flags().Var(userData, "userData", heredoc.Doc(`An object with custom data.

You can store up to 32&nbsp;kB as custom data.
`))
	cmd.Flags().SetAnnotation("userData", "Categories", []string{"Advanced"})
}

func AddSearchParamsObjectFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Whether to support phrase matching and excluding words from search queries.

Use the `+"`"+`advancedSyntaxFeatures`+"`"+` parameter to control which feature is supported.
`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Advanced search syntax features you want to support.

<dl>
<dt><code>exactPhrase</code></dt>
<dd>

Phrases in quotes must match exactly.
For example, `+"`"+`sparkly blue "iPhone case"`+"`"+` only returns records with the exact string "iPhone case".

</dd>
<dt><code>excludeWords</code></dt>
<dd>

Query words prefixed with a `+"`"+`-`+"`"+` must not occur in a record.
For example, `+"`"+`search -engine`+"`"+` matches records that contain "search" but not "engine".

</dd>
</dl>

This setting only has an effect if `+"`"+`advancedSyntax`+"`"+` is true.
`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers in the search query.

Turn off this setting to reduce the number of irrelevant matches
when searching in large sets of similar numbers.
`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives of query words that should be considered as exact matches by the Exact ranking criterion.

<dl>
<dt><code>ignorePlurals</code></dt>
<dd>

Plurals and similar declensions added by the `+"`"+`ignorePlurals`+"`"+` setting are considered exact matches.

</dd>
<dt><code>singleWordSynonym</code></dt>
<dd>
Single-word synonyms, such as "NY/NYC" are considered exact matches.
</dd>
<dt><code>multiWordsSynonym</code></dt>
<dd>
Multi-word synonyms, such as "NY/New York" are considered exact matches.
</dd>
</dl>.
`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`Whether this search will be included in Analytics.`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`Tags to apply to the query for [segmenting analytics data](https://www.algolia.com/doc/guides/search-analytics/guides/segments/).`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Coordinates for the center of a circle, expressed as a comma-separated string of latitude and longitude.

Only records included within circle around this central location are included in the results.
The radius of the circle is determined by the `+"`"+`aroundRadius`+"`"+` and `+"`"+`minimumAroundRadius`+"`"+` settings.
This parameter is ignored if you also specify `+"`"+`insidePolygon`+"`"+` or `+"`"+`insideBoundingBox`+"`"+`.
`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`Whether to obtain the coordinates from the request's IP address.`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	aroundPrecision := NewJSONVar([]string{"integer", "array"}...)
	cmd.Flags().Var(aroundPrecision, "aroundPrecision", heredoc.Doc(`Precision of a coordinate-based search in meters to group results with similar distances.

The Geo ranking criterion considers all matches within the same range of distances to be equal.
`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Maximum radius for a search around a central location.

This parameter works in combination with the `+"`"+`aroundLatLng`+"`"+` and `+"`"+`aroundLatLngViaIP`+"`"+` parameters.
By default, the search radius is determined automatically from the density of hits around the central location.
The search radius is small if there are many hits close to the central coordinates.
`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`Whether the best matching attribute should be determined by minimum proximity.

This setting only affects ranking if the Attribute ranking criterion comes before Proximity in the `+"`"+`ranking`+"`"+` setting.
If true, the best matching attribute is selected based on the minimum proximity of multiple matches.
Otherwise, the best matching attribute is determined by the order in the `+"`"+`searchableAttributes`+"`"+` setting.
`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight.

By default, all searchable attributes are highlighted.
Use `+"`"+`*`+"`"+` to highlight all attributes or use an empty array `+"`"+`[]`+"`"+` to turn off highlighting.

With highlighting, strings that match the search query are surrounded by HTML tags defined by `+"`"+`highlightPreTag`+"`"+` and `+"`"+`highlightPostTag`+"`"+`.
You can use this to visually highlight matching parts of a search query in your UI.

For more information, see [Highlighting and snippeting](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/highlighting-snippeting/js/).
`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response.

To reduce the size of your response, you can retrieve only some of the attributes.

- `+"`"+`*`+"`"+` retrieves all attributes, except attributes included in the `+"`"+`customRanking`+"`"+` and `+"`"+`unretrievableAttributes`+"`"+` settings.
- To retrieve all attributes except a specific one, prefix the attribute with a dash and combine it with the `+"`"+`*`+"`"+`: `+"`"+`["*", "-ATTRIBUTE"]`+"`"+`.
- The `+"`"+`objectID`+"`"+` attribute is always included.
`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes for which to enable snippets.

Snippets provide additional context to matched words.
If you enable snippets, they include 10 words, including the matched word.
The matched word will also be wrapped by HTML tags for highlighting.
You can adjust the number of words with the following notation: `+"`"+`ATTRIBUTE:NUMBER`+"`"+`,
where `+"`"+`NUMBER`+"`"+` is the number of words to be extracted.
`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`Whether to include a `+"`"+`queryID`+"`"+` attribute in the response.

The query ID is a unique identifier for a search query and is required for tracking [click and conversion events](https://www.algolia.com/guides/sending-events/getting-started/).
`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Attributes to use as [custom ranking](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/).

The custom ranking attributes decide which items are shown first if the other ranking criteria are equal.

Records with missing values for your selected custom ranking attributes are always sorted last.
Boolean attributes are sorted based on their alphabetical order.

**Modifiers**

<dl>
<dt><code>asc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in ascending order.</dd>
<dt><code>desc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in descending order.</dd>
</dl>

If you use two or more custom ranking attributes, [reduce the precision](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/how-to/controlling-custom-ranking-metrics-precision/) of your first attributes,
or the other attributes will never be applied.
`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Whether to split compound words into their building blocks.

For more information, see [Word segmentation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/#splitting-compound-words).
Word segmentation is supported for these languages: German, Dutch, Finnish, Swedish, and Norwegian.
`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Searchable attributes for which you want to [turn off the Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).

This can be useful for attributes with long values, where the likelyhood of an exact match is high,
such as product descriptions.
Turning off the Exact ranking criterion for these attributes favors exact matching on other attributes.
This reduces the impact of individual attributes with a lot of content on ranking.
`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/).

Returning only exact matches can help when:

- [Searching in hyphenated attributes](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/how-to/how-to-search-in-hyphenated-attributes/).
- Reducing the number of matches when you have too many.
  This can happen with attributes that are long blocks of text, such as product descriptions.

Consider alternatives such as `+"`"+`disableTypoToleranceOnWords`+"`"+` or adding synonyms if your attributes have intentional unusual spellings that might look like typos.
`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Determines how many records of a group are included in the search results.

Records with the same value for the `+"`"+`attributeForDistinct`+"`"+` attribute are considered a group.
The `+"`"+`distinct`+"`"+` setting controls how many members of the group are returned.
This is useful for [deduplication and grouping](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/#introducing-algolias-distinct-feature).

The `+"`"+`distinct`+"`"+` setting is ignored if `+"`"+`attributeForDistinct`+"`"+` is not set.
`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`Whether to enable A/B testing for this search.`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Whether to enable Personalization.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Whether this search will use [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/).

This setting only has an effect if you activated Dynamic Re-Ranking for this index in the Algolia dashboard.
`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether to enable rules.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the [Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes) is computed when the search query has only one word.

<dl>
<dt><code>attribute</code></dt>
<dd>
The Exact ranking criterion is 1 if the query word and attribute value are the same.
For example, a search for "road" will match the value "road", but not "road trip".
</dd>
<dt><code>none</code></dt>
<dd>
The Exact ranking criterion is ignored on single-word searches.
</dd>
<dt><code>word</code></dt>
<dd>
The Exact ranking criterion is 1 if the query word is found in the attribute value.
The query word must have at least 3 characters and must not be a stop word.
</dd>
</dl>

If `+"`"+`exactOnSingleWordQuery`+"`"+` is `+"`"+`word`+"`"+`, only exact matches will be highlighted, partial and prefix matches won't.
 One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter the search by facet values, so that only records with the same facet values are retrieved.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

- `+"`"+`[filter1, filter2]`+"`"+` is interpreted as `+"`"+`filter1 AND filter2`+"`"+`.
- `+"`"+`[[filter1, filter2], filter3]`+"`"+` is interpreted as `+"`"+`filter1 OR filter2 AND filter3`+"`"+`.
- `+"`"+`facet:-value`+"`"+` is interpreted as `+"`"+`NOT facet:value`+"`"+`.

While it's best to avoid attributes that start with a `+"`"+`-`+"`"+`, you can still filter them by escaping with a backslash:
`+"`"+`facet:\-value`+"`"+`.
`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`Whether faceting should be applied after deduplication with `+"`"+`distinct`+"`"+`.

This leads to accurate facet counts when using faceting in combination with `+"`"+`distinct`+"`"+`.
It's usually better to use `+"`"+`afterDistinct`+"`"+` modifiers in the `+"`"+`attributesForFaceting`+"`"+` setting,
as `+"`"+`facetingAfterDistinct`+"`"+` only computes correct facet counts if all records have the same facet values for the `+"`"+`attributeForDistinct`+"`"+`.
`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`Facets for which to retrieve facet values that match the search criteria and the number of matching facet values.

To retrieve all facets, use the wildcard character `+"`"+`*`+"`"+`.
For more information, see [facets](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#contextual-facet-values-and-counts).
`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter the search so that only records with matching values are included in the results.

These filters are supported:

- **Numeric filters.** `+"`"+`<facet> <op> <number>`+"`"+`, where `+"`"+`<op>`+"`"+` is one of `+"`"+`<`+"`"+`, `+"`"+`<=`+"`"+`, `+"`"+`=`+"`"+`, `+"`"+`!=`+"`"+`, `+"`"+`>`+"`"+`, `+"`"+`>=`+"`"+`.
- **Ranges.** `+"`"+`<facet>:<lower> TO <upper>`+"`"+` where `+"`"+`<lower>`+"`"+` and `+"`"+`<upper>`+"`"+` are the lower and upper limits of the range (inclusive).
- **Facet filters.** `+"`"+`<facet>:<value>`+"`"+` where `+"`"+`<facet>`+"`"+` is a facet attribute (case-sensitive) and `+"`"+`<value>`+"`"+` a facet value.
- **Tag filters.** `+"`"+`_tags:<value>`+"`"+` or just `+"`"+`<value>`+"`"+` (case-sensitive).
- **Boolean filters.** `+"`"+`<facet>: true | false`+"`"+`.

You can combine filters with `+"`"+`AND`+"`"+`, `+"`"+`OR`+"`"+`, and `+"`"+`NOT`+"`"+` operators with the following restrictions:

- You can only combine filters of the same type with `+"`"+`OR`+"`"+`.
  **Not supported:** `+"`"+`facet:value OR num > 3`+"`"+`.
- You can't use `+"`"+`NOT`+"`"+` with combinations of filters.
  **Not supported:** `+"`"+`NOT(facet:value OR facet:value)`+"`"+`
- You can't combine conjunctions (`+"`"+`AND`+"`"+`) with `+"`"+`OR`+"`"+`.
  **Not supported:** `+"`"+`facet:value OR (facet:value AND facet:value)`+"`"+`

Use quotes around your filters, if the facet attribute name or facet value has spaces, keywords (`+"`"+`OR`+"`"+`, `+"`"+`AND`+"`"+`, `+"`"+`NOT`+"`"+`), or quotes.
If a facet attribute is an array, the filter matches if it matches at least one element of the array.

For more information, see [Filters](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/).
`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("getRankingInfo", false, heredoc.Doc(`Whether the search response should include detailed ranking information.`))
	cmd.Flags().SetAnnotation("getRankingInfo", "Categories", []string{"Advanced"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`HTML tag to insert after the highlighted parts in all highlighted results and snippets.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`HTML tag to insert before the highlighted parts in all highlighted results and snippets.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treat singular, plurals, and other forms of declensions as equivalent.
You should only use this feature for the languages used in your index.
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters for which diacritics should be preserved.

By default, Algolia removes diacritics from letters.
For example, `+"`"+`é`+"`"+` becomes `+"`"+`e`+"`"+`. If this causes issues in your search,
you can specify characters that should keep their diacritics.
`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`Number of hits to retrieve (used in combination with `+"`"+`offset`+"`"+`).`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet values to return when [searching for facet values](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#search-for-facet-values).`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Minimum proximity score for two matching words.

This adjusts the [Proximity ranking criterion](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#proximity)
by equally scoring matches that are farther apart.

For example, if `+"`"+`minProximity`+"`"+` is 2, neighboring matches and matches with one word between them would have the same score.
`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with [one typo](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with [two typos](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`Minimum radius (in meters) for a search around a location when `+"`"+`aroundRadius`+"`"+` isn't set.`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results.

This setting only applies to indices, for which Algolia enabled NeuralSearch for you.
 One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`ISO language codes that adjust settings that are useful for processing natural language queries (as opposed to keyword searches):

- Sets `+"`"+`removeStopWords`+"`"+` and `+"`"+`ignorePlurals`+"`"+` to the list of provided languages.
- Sets `+"`"+`removeWordsIfNoResults`+"`"+` to `+"`"+`allOptional`+"`"+`.
- Adds a `+"`"+`natural_language`+"`"+` attribute to `+"`"+`ruleContexts`+"`"+` and `+"`"+`analyticsTags`+"`"+`.
`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter by numeric facets.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

You can use numeric comparison operators: `+"`"+`<`+"`"+`, `+"`"+`<=`+"`"+`, `+"`"+`=`+"`"+`, `+"`"+`!=`+"`"+`, `+"`"+`>`+"`"+`, `+"`"+`>=`+"`"+`. Comparsions are precise up to 3 decimals.
You can also provide ranges: `+"`"+`facet:<lower> TO <upper>`+"`"+`. The range includes the lower and upper boundaries.
The same combination rules apply as for `+"`"+`facetFilters`+"`"+`.
`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`Position of the first hit to retrieve.`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`Filters to promote or demote records in the search results.

Optional filters work like facet filters, but they don't exclude records from the search results.
Records that match the optional filter rank before records that don't match.
If you're using a negative filter `+"`"+`facet:-value`+"`"+`, matching records rank after records that don't match.

- Optional filters don't work on virtual replicas.
- Optional filters are applied _after_ sort-by attributes.
- Optional filters don't work with numeric attributes.
`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words that should be considered optional when found in the query.

By default, records must match all words in the search query to be included in the search results.
Adding optional words can help to increase the number of search results by running an additional search query that doesn't include the optional words.
For example, if the search query is "action video" and "video" is an optional word,
the search engine runs two queries. One for "action video" and one for "action".
Records that match all words are ranked higher.

For a search query with 4 or more words **and** all its words are optional,
the number of matched words required for a record to be included in the search results increases for every 1,000 records:

- If `+"`"+`optionalWords`+"`"+` has less than 10 words, the required number of matched words increases by 1:
  results 1 to 1,000 require 1 matched word, results 1,001 to 2000 need 2 matched words.
- If `+"`"+`optionalWords`+"`"+` has 10 or more words, the number of required matched words increases by the number of optional words dividied by 5 (rounded down).
  For example, with 18 optional words: results 1 to 1,000 require 1 matched word, results 1,001 to 2000 need 4 matched words.

For more information, see [Optional words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/#creating-a-list-of-optional-words).
`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`Page of search results to retrieve.`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`Whether to include this search when calculating processing-time percentiles.`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`Impact that Personalization should have on this search.

The higher this value is, the more Personalization determines the ranking compared to other factors.
For more information, see [Understanding Personalization impact](https://www.algolia.com/doc/guides/personalization/personalizing-results/in-depth/configuring-personalization/#understanding-personalization-impact).
`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`Search query.`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Languages for language-specific query processing steps such as plurals, stop-word removal, and word-detection dictionaries.

This setting sets a default list of languages used by the `+"`"+`removeStopWords`+"`"+` and `+"`"+`ignorePlurals`+"`"+` settings.
This setting also sets a dictionary for word detection in the logogram-based [CJK](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/#normalization-for-logogram-based-languages-cjk) languages.
To support this, you must place the CJK language **first**.

**You should always specify a query language.**
If you don't specify an indexing language, the search engine uses all [supported languages](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/),
or the languages you specified with the `+"`"+`ignorePlurals`+"`"+` or `+"`"+`removeStopWords`+"`"+` parameters.
This can lead to unexpected search results.
For more information, see [Language-specific configuration](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/).
`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines if and how query words are interpreted as prefixes.

By default, only the last query word is treated as prefix (`+"`"+`prefixLast`+"`"+`).
To turn off prefix search, use `+"`"+`prefixNone`+"`"+`.
Avoid `+"`"+`prefixAll`+"`"+`, which treats all query words as prefixes.
This might lead to counterintuitive results and makes your search slower.

For more information, see [Prefix searching](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/prefix-searching/).
 One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia returns your results.

By default, each entry corresponds to a [ranking criteria](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/).
The tie-breaking algorithm sequentially applies each criterion in the order they're specified.
If you configure a replica index for [sorting by an attribute](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/how-to/sort-by-attribute/),
you put the sorting attribute at the top of the list.

**Modifiers**

<dl>
<dt><code>asc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in ascending order.</dd>
<dt><code>desc("ATTRIBUTE")</code></dt>
<dd>Sort the index by the values of an attribute, in descending order.</dd>
</dl>

Before you modify the default setting,
you should test your changes in the dashboard,
and by [A/B testing](https://www.algolia.com/doc/guides/ab-testing/what-is-ab-testing/).
`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string", "null"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`Restrict [Dynamic Re-ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/) to records that match these filters.
`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.

You can only set `+"`"+`relevancyStrictness`+"`"+` on [virtual replica indices](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/in-depth/replicas/#what-are-virtual-replicas).
Use this setting to strike a balance between the relevance and number of returned results.
`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop words from the search query.

Stop words are common words like articles, conjunctions, prepositions, or pronouns that have little or no meaning on their own.
In English, "the", "a", or "and" are stop words.

You should only use this feature for the languages used in your index.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy for removing words from the query when it doesn't return any results.
This helps to avoid returning empty search results.

<dl>
<dt><code>none</code></dt>
<dd>No words are removed when a query doesn't return results.</dd>
<dt><code>lastWords</code></dt>
<dd>Treat the last (then second to last, then third to last) word as optional, until there are results or at most 5 words have been removed.</dd>
<dt><code>firstWords</code></dt>
<dd>Treat the first (then second, then third) word as optional, until there are results or at most 5 words have been removed.</dd>
<dt><code>allOptional</code></dt>
<dd>Treat all words as optional.</dd>
</dl>

For more information, see [Remove words to improve results](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/in-depth/why-use-remove-words-if-no-results/).
 One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra data that can be used in the search UI.

You can use this to control aspects of your search UI, such as, the order of facet names and values
without changing your frontend code.
`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to replace a highlighted word with the matched synonym.

By default, the original words are highlighted even if a synonym matches.
For example, with `+"`"+`home`+"`"+` as a synonym for `+"`"+`house`+"`"+` and a search for `+"`"+`home`+"`"+`,
records matching either "home" or "house" are included in the search results,
and either "home" or "house" are highlighted.

With `+"`"+`replaceSynonymsInHighlight`+"`"+` set to `+"`"+`true`+"`"+`, a search for `+"`"+`home`+"`"+` still matches the same records,
but all occurences of "house" are replaced by "home" in the highlighted response.
`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{"*"}, heredoc.Doc(`Properties to include in the API response of `+"`"+`search`+"`"+` and `+"`"+`browse`+"`"+` requests.

By default, all response properties are included.
To reduce the response size, you can select, which attributes should be included.

You can't exclude these properties:
`+"`"+`message`+"`"+`, `+"`"+`warning`+"`"+`, `+"`"+`cursor`+"`"+`, `+"`"+`serverUsed`+"`"+`, `+"`"+`indexUsed`+"`"+`,
`+"`"+`abTestVariantID`+"`"+`, `+"`"+`parsedQuery`+"`"+`, or any property triggered by the `+"`"+`getRankingInfo`+"`"+` parameter.

Don't exclude properties that you might need in your search UI.
`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Whether to restrict highlighting and snippeting to items that at least partially matched the search query.
By default, all items are highlighted and snippeted.
`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a search to a subset of your searchable attributes.`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`Assigns a rule context to the search query.

[Rule contexts](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/how-to/customize-search-results-by-platform/#whats-a-context) are strings that you can use to trigger matching rules.
`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch.
Only used when `+"`"+`mode`+"`"+` is `+"`"+`neuralSearch`+"`"+`.
`))
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`Keywords to be used instead of the search query to conduct a more broader search.

Using the `+"`"+`similarQuery`+"`"+` parameter changes other settings:

- `+"`"+`queryType`+"`"+` is set to `+"`"+`prefixNone`+"`"+`.
- `+"`"+`removeStopWords`+"`"+` is set to true.
- `+"`"+`words`+"`"+` is set as the first ranking criterion.
- All remaining words are treated as `+"`"+`optionalWords`+"`"+`.

Since the `+"`"+`similarQuery`+"`"+` is supposed to do a broad search, they usually return many results.
Combine it with `+"`"+`filters`+"`"+` to narrow down the list of results.
`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Order in which to retrieve facet values.

<dl>
<dt><code>count</code></dt>
<dd>
Facet values are retrieved by decreasing count.
The count is the number of matching records containing this facet value.
</dd>
<dt><code>alpha</code></dt>
<dd>Retrieve facet values alphabetically.</dd>
</dl>

This setting doesn't influence how facet values are displayed in your UI (see `+"`"+`renderingContent`+"`"+`).
For more information, see [facet value display](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/facet-display/js/).
`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`Whether to sum all filter scores.

If true, all filter scores are summed.
Otherwise, the maximum filter score is kept.
For more information, see [filter scores](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/in-depth/filter-scoring/#accumulating-scores-with-sumorfiltersscores).
`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for this search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter the search by values of the special `+"`"+`_tags`+"`"+` attribute.

**Prefer using the `+"`"+`filters`+"`"+` parameter, which supports all filter types and combinations with boolean operators.**

Different from regular facets, `+"`"+`_tags`+"`"+` can only be used for filtering (including or excluding records).
You won't get a facet count.
The same combination and escaping rules apply as for `+"`"+`facetFilters`+"`"+`.
`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Whether [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/) is enabled and how it is applied.

If typo tolerance is true, `+"`"+`min`+"`"+`, or `+"`"+`strict`+"`"+`, [word splitting and concetenation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/splitting-and-concatenation/) is also active.
`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`Unique pseudonymous or anonymous user identifier.

This helps with analytics and click and conversion events.
For more information, see [user token](https://www.algolia.com/doc/guides/sending-events/concepts/usertoken/).
`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

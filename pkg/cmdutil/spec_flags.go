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
	"attributesForFaceting",
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
	"explain",
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
	"attributesForFaceting",
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
	"explain",
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
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Enables the [advanced query syntax](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/#advanced-syntax).`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Allows you to specify which advanced syntax features are active when `+"`"+`advancedSyntax`+"`"+` is enabled.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers ("numeric tokens") in the query string.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives that should be considered an exact match by [the exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`Indicates whether this query will be included in [analytics](https://www.algolia.com/doc/guides/search-analytics/guides/exclude-queries/).`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`Tags to apply to the query for [segmenting analytics data](https://www.algolia.com/doc/guides/search-analytics/guides/segments/).`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Search for entries [around a central location](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filter-around-a-central-point), enabling a geographical search within a circular area.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`Search for entries around a location. The location is automatically computed from the requester's IP address.`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	aroundPrecision := NewJSONVar([]string{"integer", "array"}...)
	cmd.Flags().Var(aroundPrecision, "aroundPrecision", heredoc.Doc(`Precision of a geographical search (in meters), to [group results that are more or less the same distance from a central point](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/in-depth/geo-ranking-precision/).`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`[Maximum radius](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#increase-the-search-radius) for a geographical search (in meters).
`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`When the [Attribute criterion is ranked above Proximity](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#attribute-and-proximity-combinations) in your ranking formula, Proximity is used to select which searchable attribute is matched in the Attribute ranking stage.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`Attributes used for [faceting](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/) and the [modifiers](https://www.algolia.com/doc/api-reference/api-parameters/attributesForFaceting/#modifiers) that can be applied: `+"`"+`filterOnly`+"`"+`, `+"`"+`searchable`+"`"+`, and `+"`"+`afterDistinct`+"`"+`.
`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight. Strings that match the search query in the attributes are highlighted by surrounding them with HTML tags (`+"`"+`highlightPreTag`+"`"+` and `+"`"+`highlightPostTag`+"`"+`).`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response. To reduce the size of your response, you can retrieve only some of the attributes. By default, the response includes all attributes.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes to _snippet_. 'Snippeting' is shortening the attribute to a certain number of words. If not specified, the attribute is shortened to the 10 words around the matching string but you can specify the number. For example: `+"`"+`body:20`+"`"+`.
`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`Indicates whether a query ID parameter is included in the search response. This is required for [tracking click and conversion events](https://www.algolia.com/doc/guides/sending-events/concepts/event-types/#events-related-to-algolia-requests).`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().String("cursor", "", heredoc.Doc(`Cursor indicating the location to resume browsing from. Must match the value returned by the previous call.
Pass this value to the subsequent browse call to get the next page of results.
When the end of the index has been reached, `+"`"+`cursor`+"`"+` is absent from the response.
`))
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Specifies the [Custom ranking criterion](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/). Use the `+"`"+`asc`+"`"+` and `+"`"+`desc`+"`"+` modifiers to specify the ranking order: ascending or descending.
`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`[Splits compound words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/#splitting-compound-words) into their component word parts in the query.
`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to [turn off the exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/).`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Enables [deduplication or grouping of results (Algolia's _distinct_ feature](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/#introducing-algolias-distinct-feature)).`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`Incidates whether this search will be considered in A/B testing.`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Incidates whether [Personalization](https://www.algolia.com/doc/guides/personalization/what-is-personalization/) is enabled.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Indicates whether this search will use [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/).`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Incidates whether [Rules](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/) are enabled.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the [Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes) is computed when the query contains only one word. One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("explain", []string{}, heredoc.Doc(`Enriches the API's response with information about how the query was processed.`))
	cmd.Flags().SetAnnotation("explain", "Categories", []string{"Advanced"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`[Filter hits by facet value](https://www.algolia.com/doc/api-reference/api-parameters/facetFilters/).
`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`Forces faceting to be applied after [de-duplication](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/) (with the distinct feature). Alternatively, the `+"`"+`afterDistinct`+"`"+` [modifier](https://www.algolia.com/doc/api-reference/api-parameters/attributesForFaceting/#modifiers) of `+"`"+`attributesForFaceting`+"`"+` allows for more granular control.
`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`Returns [facets](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#contextual-facet-values-and-counts), their facet values, and the number of matching facet values.`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`[Filter](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/) the query with numeric, facet, or tag filters.
`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("getRankingInfo", false, heredoc.Doc(`Incidates whether the search response includes [detailed ranking information](https://www.algolia.com/doc/guides/building-search-ui/going-further/backend-search/in-depth/understanding-the-api-response/#ranking-information).`))
	cmd.Flags().SetAnnotation("getRankingInfo", "Categories", []string{"Advanced"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`HTML string to insert after the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`HTML string to insert before the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treats singular, plurals, and other forms of declensions as matching terms.
`+"`"+`ignorePlurals`+"`"+` is used in conjunction with the `+"`"+`queryLanguages`+"`"+` setting.
_list_: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in `+"`"+`queryLanguages`+"`"+`. _true_: enables the ignore plurals feature, where singulars and plurals are considered equivalent ("foot" = "feet"). The languages supported here are either [every language](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/) (this is the default) or those set by `+"`"+`queryLanguages`+"`"+`. _false_: turns off the ignore plurals feature, so that singulars and plurals aren't considered to be the same ("foot" will not find "feet").
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().Float64Slice("insideBoundingBox", []float64{}, heredoc.Doc(`Search inside a [rectangular area](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filtering-inside-rectangular-or-polygonal-areas) (in geographical coordinates).`))
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().Float64Slice("insidePolygon", []float64{}, heredoc.Doc(`Search inside a [polygon](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filtering-inside-rectangular-or-polygonal-areas) (in geographical coordinates).`))
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters that the engine shouldn't automatically [normalize](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/).`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`Sets the number of hits to retrieve (for use with `+"`"+`offset`+"`"+`).
> **Note**: Using `+"`"+`page`+"`"+` and `+"`"+`hitsPerPage`+"`"+` is the recommended method for [paging results](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/pagination/js/). However, you can use `+"`"+`offset`+"`"+` and `+"`"+`length`+"`"+` to implement [an alternative approach to paging](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/pagination/js/#retrieving-a-subset-of-records-with-offset-and-length).
`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet hits to return when [searching for facet values](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#search-for-facet-values).`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Precision of the [proximity ranking criterion](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#proximity).`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with [one typo](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with [two typos](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`Minimum radius (in meters) used for a geographical search when `+"`"+`aroundRadius`+"`"+` isn't set.`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results. One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`Changes the default values of parameters that work best for a natural language query, such as `+"`"+`ignorePlurals`+"`"+`, `+"`"+`removeStopWords`+"`"+`, `+"`"+`removeWordsIfNoResults`+"`"+`, `+"`"+`analyticsTags`+"`"+`, and `+"`"+`ruleContexts`+"`"+`. These parameters work well together when the query consists of fuller natural language strings instead of keywords, for example when processing voice search queries.`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`[Filter on numeric attributes](https://www.algolia.com/doc/api-reference/api-parameters/numericFilters/).
`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`Specifies the offset of the first hit to return.
> **Note**: Using `+"`"+`page`+"`"+` and `+"`"+`hitsPerPage`+"`"+` is the recommended method for [paging results](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/pagination/js/). However, you can use `+"`"+`offset`+"`"+` and `+"`"+`length`+"`"+` to implement [an alternative approach to paging](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/pagination/js/#retrieving-a-subset-of-records-with-offset-and-length).
`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`Create filters to boost or demote records. 

Records that match the filter are ranked higher for positive and lower for negative optional filters. In contrast to regular filters, records that don't match the optional filter are still included in the results, only their ranking is affected.
`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words which should be considered [optional](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/#creating-a-list-of-optional-words) when found in a query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`Page to retrieve (the first page is `+"`"+`0`+"`"+`, not `+"`"+`1`+"`"+`).`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`Whether to include or exclude a query from the processing-time percentile computation.`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`Defines how much [Personalization affects results](https://www.algolia.com/doc/guides/personalization/personalizing-results/in-depth/configuring-personalization/#understanding-personalization-impact).`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`Text to search for in an index.`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Sets your user's search language. This adjusts language-specific settings and features such as `+"`"+`ignorePlurals`+"`"+`, `+"`"+`removeStopWords`+"`"+`, and [CJK](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/#normalization-for-logogram-based-languages-cjk) word detection.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines how query words are [interpreted as prefixes](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/prefix-searching/). One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia [returns your results](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/).`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`When [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/) is enabled, only records that match these filters will be affected by Dynamic Re-Ranking.`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop (common) words from the query before executing it.
`+"`"+`removeStopWords`+"`"+` is used in conjunction with the `+"`"+`queryLanguages`+"`"+` setting.
_list_: language ISO codes for which stop words should be enabled. This list will override any values that you may have set in `+"`"+`queryLanguages`+"`"+`. _true_: enables the stop words feature, ensuring that stop words are removed from consideration in a search. The languages supported here are either [every language](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/) (this is the default) or those set by `+"`"+`queryLanguages`+"`"+`. _false_: turns off the stop words feature, allowing stop words to be taken into account in a search.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy to [remove words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/in-depth/why-use-remove-words-if-no-results/) from the query when it doesn't match any hits. One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra content for the search UI, for example, to control the [ordering and display of facets](https://www.algolia.com/doc/guides/managing-results/rules/merchandising-and-promoting/how-to/merchandising-facets/#merchandise-facets-and-their-values-in-the-manual-editor). You can set a default value and dynamically override it with [Rules](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/).`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to highlight and snippet the original word that matches the synonym or the synonym itself.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`Attributes to include in the API response for search and browse queries.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Restrict highlighting and snippeting to items that matched the query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a query to only look at a subset of your [searchable attributes](https://www.algolia.com/doc/guides/managing-results/must-do/searchable-attributes/).`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`Assigns [rule contexts](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/how-to/customize-search-results-by-platform/#whats-a-context) to search queries.`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch. Only used when `+"`"+`mode`+"`"+` is _neuralSearch_.
`))
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`Overrides the query parameter and performs a more generic search.`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Controls how facet values are fetched.`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`Determines how to calculate [filter scores](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/in-depth/filter-scoring/#accumulating-scores-with-sumorfiltersscores).
If `+"`"+`false`+"`"+`, maximum score is kept.
If `+"`"+`true`+"`"+`, score is summed.
`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for a particular search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`[Filter hits by tags](https://www.algolia.com/doc/api-reference/api-parameters/tagFilters/).
`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Controls whether [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/) is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`Associates a [user token](https://www.algolia.com/doc/guides/sending-events/concepts/usertoken/) with the current search.`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

func AddDeleteByParamsFlags(cmd *cobra.Command) {
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Search for entries [around a central location](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filter-around-a-central-point), enabling a geographical search within a circular area.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`[Maximum radius](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#increase-the-search-radius) for a geographical search (in meters).
`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`[Filter hits by facet value](https://www.algolia.com/doc/api-reference/api-parameters/facetFilters/).
`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().String("filters", "", heredoc.Doc(`[Filter](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/) the query with numeric, facet, or tag filters.
`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Float64Slice("insideBoundingBox", []float64{}, heredoc.Doc(`Search inside a [rectangular area](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filtering-inside-rectangular-or-polygonal-areas) (in geographical coordinates).`))
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().Float64Slice("insidePolygon", []float64{}, heredoc.Doc(`Search inside a [polygon](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filtering-inside-rectangular-or-polygonal-areas) (in geographical coordinates).`))
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`[Filter on numeric attributes](https://www.algolia.com/doc/api-reference/api-parameters/numericFilters/).
`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`[Filter hits by tags](https://www.algolia.com/doc/api-reference/api-parameters/tagFilters/).
`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
}

func AddIndexSettingsFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Enables the [advanced query syntax](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/#advanced-syntax).`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Allows you to specify which advanced syntax features are active when `+"`"+`advancedSyntax`+"`"+` is enabled.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowCompressionOfIntegerArray", false, heredoc.Doc(`Incidates whether the engine compresses arrays with exclusively non-negative integers.
When enabled, the compressed arrays may be reordered.
`))
	cmd.Flags().SetAnnotation("allowCompressionOfIntegerArray", "Categories", []string{"Performance"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers ("numeric tokens") in the query string.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives that should be considered an exact match by [the exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`When the [Attribute criterion is ranked above Proximity](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#attribute-and-proximity-combinations) in your ranking formula, Proximity is used to select which searchable attribute is matched in the Attribute ranking stage.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().String("attributeForDistinct", "", heredoc.Doc(`Name of the deduplication attribute to be used with Algolia's [_distinct_ feature](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/#introducing-algolias-distinct-feature).`))
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`Attributes used for [faceting](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/) and the [modifiers](https://www.algolia.com/doc/api-reference/api-parameters/attributesForFaceting/#modifiers) that can be applied: `+"`"+`filterOnly`+"`"+`, `+"`"+`searchable`+"`"+`, and `+"`"+`afterDistinct`+"`"+`.
`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight. Strings that match the search query in the attributes are highlighted by surrounding them with HTML tags (`+"`"+`highlightPreTag`+"`"+` and `+"`"+`highlightPostTag`+"`"+`).`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response. To reduce the size of your response, you can retrieve only some of the attributes. By default, the response includes all attributes.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes to _snippet_. 'Snippeting' is shortening the attribute to a certain number of words. If not specified, the attribute is shortened to the 10 words around the matching string but you can specify the number. For example: `+"`"+`body:20`+"`"+`.
`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToTransliterate", []string{}, heredoc.Doc(`Attributes in your index to which [Japanese transliteration](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/#japanese-transliteration-and-type-ahead) applies. This will ensure that words indexed in Katakana or Kanji can also be searched in Hiragana.`))
	cmd.Flags().SetAnnotation("attributesToTransliterate", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("camelCaseAttributes", []string{}, heredoc.Doc(`Attributes on which to split [camel case](https://wikipedia.org/wiki/Camel_case) words.`))
	cmd.Flags().SetAnnotation("camelCaseAttributes", "Categories", []string{"Languages"})
	customNormalization := NewJSONVar([]string{}...)
	cmd.Flags().Var(customNormalization, "customNormalization", heredoc.Doc(`A list of characters and their normalized replacements to override Algolia's default [normalization](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/).`))
	cmd.Flags().SetAnnotation("customNormalization", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Specifies the [Custom ranking criterion](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/). Use the `+"`"+`asc`+"`"+` and `+"`"+`desc`+"`"+` modifiers to specify the ranking order: ascending or descending.
`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`[Splits compound words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/#splitting-compound-words) into their component word parts in the query.
`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	decompoundedAttributes := NewJSONVar([]string{}...)
	cmd.Flags().Var(decompoundedAttributes, "decompoundedAttributes", heredoc.Doc(`Attributes in your index to which [word segmentation](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/how-to/customize-segmentation/) (decompounding) applies.`))
	cmd.Flags().SetAnnotation("decompoundedAttributes", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to [turn off the exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disablePrefixOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off [prefix matching](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/#adjusting-prefix-search).`))
	cmd.Flags().SetAnnotation("disablePrefixOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/).`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("disableTypoToleranceOnWords", []string{}, heredoc.Doc(`Words for which you want to turn off [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/).`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnWords", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Enables [deduplication or grouping of results (Algolia's _distinct_ feature](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/#introducing-algolias-distinct-feature)).`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Incidates whether [Personalization](https://www.algolia.com/doc/guides/personalization/what-is-personalization/) is enabled.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Indicates whether this search will use [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/).`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Incidates whether [Rules](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/) are enabled.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the [Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes) is computed when the query contains only one word. One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`HTML string to insert after the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`HTML string to insert before the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treats singular, plurals, and other forms of declensions as matching terms.
`+"`"+`ignorePlurals`+"`"+` is used in conjunction with the `+"`"+`queryLanguages`+"`"+` setting.
_list_: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in `+"`"+`queryLanguages`+"`"+`. _true_: enables the ignore plurals feature, where singulars and plurals are considered equivalent ("foot" = "feet"). The languages supported here are either [every language](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/) (this is the default) or those set by `+"`"+`queryLanguages`+"`"+`. _false_: turns off the ignore plurals feature, so that singulars and plurals aren't considered to be the same ("foot" will not find "feet").
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("indexLanguages", []string{}, heredoc.Doc(`Set the languages of your index, for language-specific processing steps such as [tokenization](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/tokenization/) and [normalization](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/).`))
	cmd.Flags().SetAnnotation("indexLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters that the engine shouldn't automatically [normalize](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/).`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet hits to return when [searching for facet values](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#search-for-facet-values).`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Precision of the [proximity ranking criterion](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#proximity).`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with [one typo](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with [two typos](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results. One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("numericAttributesForFiltering", []string{}, heredoc.Doc(`Numeric attributes that can be used as [numerical filters](https://www.algolia.com/doc/guides/managing-results/rules/detecting-intent/how-to/applying-a-custom-filter-for-a-specific-query/#numerical-filters).`))
	cmd.Flags().SetAnnotation("numericAttributesForFiltering", "Categories", []string{"Performance"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words which should be considered [optional](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/#creating-a-list-of-optional-words) when found in a query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("paginationLimitedTo", 1000, heredoc.Doc(`Maximum number of hits accessible through pagination.`))
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Sets your user's search language. This adjusts language-specific settings and features such as `+"`"+`ignorePlurals`+"`"+`, `+"`"+`removeStopWords`+"`"+`, and [CJK](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/#normalization-for-logogram-based-languages-cjk) word detection.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines how query words are [interpreted as prefixes](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/prefix-searching/). One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia [returns your results](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/).`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`When [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/) is enabled, only records that match these filters will be affected by Dynamic Re-Ranking.`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop (common) words from the query before executing it.
`+"`"+`removeStopWords`+"`"+` is used in conjunction with the `+"`"+`queryLanguages`+"`"+` setting.
_list_: language ISO codes for which stop words should be enabled. This list will override any values that you may have set in `+"`"+`queryLanguages`+"`"+`. _true_: enables the stop words feature, ensuring that stop words are removed from consideration in a search. The languages supported here are either [every language](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/) (this is the default) or those set by `+"`"+`queryLanguages`+"`"+`. _false_: turns off the stop words feature, allowing stop words to be taken into account in a search.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy to [remove words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/in-depth/why-use-remove-words-if-no-results/) from the query when it doesn't match any hits. One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra content for the search UI, for example, to control the [ordering and display of facets](https://www.algolia.com/doc/guides/managing-results/rules/merchandising-and-promoting/how-to/merchandising-facets/#merchandise-facets-and-their-values-in-the-manual-editor). You can set a default value and dynamically override it with [Rules](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/).`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to highlight and snippet the original word that matches the synonym or the synonym itself.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("replicas", []string{}, heredoc.Doc(`Creates [replicas](https://www.algolia.com/doc/guides/managing-results/refine-results/sorting/in-depth/replicas/), which are copies of a primary index with the same records but different settings.`))
	cmd.Flags().SetAnnotation("replicas", "Categories", []string{"Ranking"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`Attributes to include in the API response for search and browse queries.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Restrict highlighting and snippeting to items that matched the query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("searchableAttributes", []string{}, heredoc.Doc(`[Attributes used for searching](https://www.algolia.com/doc/guides/managing-results/must-do/searchable-attributes/), including determining [if matches at the beginning of a word are important (ordered) or not (unordered)](https://www.algolia.com/doc/guides/managing-results/must-do/searchable-attributes/how-to/configuring-searchable-attributes-the-right-way/#understanding-word-position).
`))
	cmd.Flags().SetAnnotation("searchableAttributes", "Categories", []string{"Attributes"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch. Only used when `+"`"+`mode`+"`"+` is _neuralSearch_.
`))
	cmd.Flags().String("separatorsToIndex", "", heredoc.Doc(`Controls which separators are added to an Algolia index as part of [normalization](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/#what-does-normalization-mean). Separators are all non-letter characters except spaces and currency characters, such as $€£¥.`))
	cmd.Flags().SetAnnotation("separatorsToIndex", "Categories", []string{"Typos"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Controls how facet values are fetched.`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Controls whether [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/) is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("unretrievableAttributes", []string{}, heredoc.Doc(`Attributes that can't be retrieved at query time.`))
	cmd.Flags().SetAnnotation("unretrievableAttributes", "Categories", []string{"Attributes"})
	userData := NewJSONVar([]string{}...)
	cmd.Flags().Var(userData, "userData", heredoc.Doc(`Lets you store custom data in your indices.`))
	cmd.Flags().SetAnnotation("userData", "Categories", []string{"Advanced"})
}

func AddSearchParamsObjectFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Enables the [advanced query syntax](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/#advanced-syntax).`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Allows you to specify which advanced syntax features are active when `+"`"+`advancedSyntax`+"`"+` is enabled.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers ("numeric tokens") in the query string.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives that should be considered an exact match by [the exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`Indicates whether this query will be included in [analytics](https://www.algolia.com/doc/guides/search-analytics/guides/exclude-queries/).`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`Tags to apply to the query for [segmenting analytics data](https://www.algolia.com/doc/guides/search-analytics/guides/segments/).`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Search for entries [around a central location](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filter-around-a-central-point), enabling a geographical search within a circular area.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`Search for entries around a location. The location is automatically computed from the requester's IP address.`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	aroundPrecision := NewJSONVar([]string{"integer", "array"}...)
	cmd.Flags().Var(aroundPrecision, "aroundPrecision", heredoc.Doc(`Precision of a geographical search (in meters), to [group results that are more or less the same distance from a central point](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/in-depth/geo-ranking-precision/).`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`[Maximum radius](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#increase-the-search-radius) for a geographical search (in meters).
`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`When the [Attribute criterion is ranked above Proximity](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#attribute-and-proximity-combinations) in your ranking formula, Proximity is used to select which searchable attribute is matched in the Attribute ranking stage.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`Attributes used for [faceting](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/) and the [modifiers](https://www.algolia.com/doc/api-reference/api-parameters/attributesForFaceting/#modifiers) that can be applied: `+"`"+`filterOnly`+"`"+`, `+"`"+`searchable`+"`"+`, and `+"`"+`afterDistinct`+"`"+`.
`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight. Strings that match the search query in the attributes are highlighted by surrounding them with HTML tags (`+"`"+`highlightPreTag`+"`"+` and `+"`"+`highlightPostTag`+"`"+`).`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response. To reduce the size of your response, you can retrieve only some of the attributes. By default, the response includes all attributes.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes to _snippet_. 'Snippeting' is shortening the attribute to a certain number of words. If not specified, the attribute is shortened to the 10 words around the matching string but you can specify the number. For example: `+"`"+`body:20`+"`"+`.
`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`Indicates whether a query ID parameter is included in the search response. This is required for [tracking click and conversion events](https://www.algolia.com/doc/guides/sending-events/concepts/event-types/#events-related-to-algolia-requests).`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Specifies the [Custom ranking criterion](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/). Use the `+"`"+`asc`+"`"+` and `+"`"+`desc`+"`"+` modifiers to specify the ranking order: ascending or descending.
`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`[Splits compound words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/language-specific-configurations/#splitting-compound-words) into their component word parts in the query.
`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to [turn off the exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes).`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/).`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Enables [deduplication or grouping of results (Algolia's _distinct_ feature](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/#introducing-algolias-distinct-feature)).`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`Incidates whether this search will be considered in A/B testing.`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Incidates whether [Personalization](https://www.algolia.com/doc/guides/personalization/what-is-personalization/) is enabled.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Indicates whether this search will use [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/).`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Incidates whether [Rules](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/) are enabled.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the [Exact ranking criterion](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/adjust-exact-settings/#turn-off-exact-for-some-attributes) is computed when the query contains only one word. One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("explain", []string{}, heredoc.Doc(`Enriches the API's response with information about how the query was processed.`))
	cmd.Flags().SetAnnotation("explain", "Categories", []string{"Advanced"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`[Filter hits by facet value](https://www.algolia.com/doc/api-reference/api-parameters/facetFilters/).
`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`Forces faceting to be applied after [de-duplication](https://www.algolia.com/doc/guides/managing-results/refine-results/grouping/) (with the distinct feature). Alternatively, the `+"`"+`afterDistinct`+"`"+` [modifier](https://www.algolia.com/doc/api-reference/api-parameters/attributesForFaceting/#modifiers) of `+"`"+`attributesForFaceting`+"`"+` allows for more granular control.
`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`Returns [facets](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#contextual-facet-values-and-counts), their facet values, and the number of matching facet values.`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`[Filter](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/) the query with numeric, facet, or tag filters.
`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("getRankingInfo", false, heredoc.Doc(`Incidates whether the search response includes [detailed ranking information](https://www.algolia.com/doc/guides/building-search-ui/going-further/backend-search/in-depth/understanding-the-api-response/#ranking-information).`))
	cmd.Flags().SetAnnotation("getRankingInfo", "Categories", []string{"Advanced"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`HTML string to insert after the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`HTML string to insert before the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treats singular, plurals, and other forms of declensions as matching terms.
`+"`"+`ignorePlurals`+"`"+` is used in conjunction with the `+"`"+`queryLanguages`+"`"+` setting.
_list_: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in `+"`"+`queryLanguages`+"`"+`. _true_: enables the ignore plurals feature, where singulars and plurals are considered equivalent ("foot" = "feet"). The languages supported here are either [every language](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/) (this is the default) or those set by `+"`"+`queryLanguages`+"`"+`. _false_: turns off the ignore plurals feature, so that singulars and plurals aren't considered to be the same ("foot" will not find "feet").
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().Float64Slice("insideBoundingBox", []float64{}, heredoc.Doc(`Search inside a [rectangular area](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filtering-inside-rectangular-or-polygonal-areas) (in geographical coordinates).`))
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().Float64Slice("insidePolygon", []float64{}, heredoc.Doc(`Search inside a [polygon](https://www.algolia.com/doc/guides/managing-results/refine-results/geolocation/#filtering-inside-rectangular-or-polygonal-areas) (in geographical coordinates).`))
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters that the engine shouldn't automatically [normalize](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/).`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`Sets the number of hits to retrieve (for use with `+"`"+`offset`+"`"+`).
> **Note**: Using `+"`"+`page`+"`"+` and `+"`"+`hitsPerPage`+"`"+` is the recommended method for [paging results](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/pagination/js/). However, you can use `+"`"+`offset`+"`"+` and `+"`"+`length`+"`"+` to implement [an alternative approach to paging](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/pagination/js/#retrieving-a-subset-of-records-with-offset-and-length).
`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet hits to return when [searching for facet values](https://www.algolia.com/doc/guides/managing-results/refine-results/faceting/#search-for-facet-values).`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Precision of the [proximity ranking criterion](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/#proximity).`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with [one typo](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with [two typos](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/in-depth/configuring-typo-tolerance/#configuring-word-length-for-typos).`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`Minimum radius (in meters) used for a geographical search when `+"`"+`aroundRadius`+"`"+` isn't set.`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results. One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`Changes the default values of parameters that work best for a natural language query, such as `+"`"+`ignorePlurals`+"`"+`, `+"`"+`removeStopWords`+"`"+`, `+"`"+`removeWordsIfNoResults`+"`"+`, `+"`"+`analyticsTags`+"`"+`, and `+"`"+`ruleContexts`+"`"+`. These parameters work well together when the query consists of fuller natural language strings instead of keywords, for example when processing voice search queries.`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`[Filter on numeric attributes](https://www.algolia.com/doc/api-reference/api-parameters/numericFilters/).
`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`Specifies the offset of the first hit to return.
> **Note**: Using `+"`"+`page`+"`"+` and `+"`"+`hitsPerPage`+"`"+` is the recommended method for [paging results](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/pagination/js/). However, you can use `+"`"+`offset`+"`"+` and `+"`"+`length`+"`"+` to implement [an alternative approach to paging](https://www.algolia.com/doc/guides/building-search-ui/ui-and-ux-patterns/pagination/js/#retrieving-a-subset-of-records-with-offset-and-length).
`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`Create filters to boost or demote records. 

Records that match the filter are ranked higher for positive and lower for negative optional filters. In contrast to regular filters, records that don't match the optional filter are still included in the results, only their ranking is affected.
`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words which should be considered [optional](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/#creating-a-list-of-optional-words) when found in a query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`Page to retrieve (the first page is `+"`"+`0`+"`"+`, not `+"`"+`1`+"`"+`).`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`Whether to include or exclude a query from the processing-time percentile computation.`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`Defines how much [Personalization affects results](https://www.algolia.com/doc/guides/personalization/personalizing-results/in-depth/configuring-personalization/#understanding-personalization-impact).`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`Text to search for in an index.`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Sets your user's search language. This adjusts language-specific settings and features such as `+"`"+`ignorePlurals`+"`"+`, `+"`"+`removeStopWords`+"`"+`, and [CJK](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/normalization/#normalization-for-logogram-based-languages-cjk) word detection.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines how query words are [interpreted as prefixes](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/override-search-engine-defaults/in-depth/prefix-searching/). One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia [returns your results](https://www.algolia.com/doc/guides/managing-results/relevance-overview/in-depth/ranking-criteria/).`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`When [Dynamic Re-Ranking](https://www.algolia.com/doc/guides/algolia-ai/re-ranking/) is enabled, only records that match these filters will be affected by Dynamic Re-Ranking.`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop (common) words from the query before executing it.
`+"`"+`removeStopWords`+"`"+` is used in conjunction with the `+"`"+`queryLanguages`+"`"+` setting.
_list_: language ISO codes for which stop words should be enabled. This list will override any values that you may have set in `+"`"+`queryLanguages`+"`"+`. _true_: enables the stop words feature, ensuring that stop words are removed from consideration in a search. The languages supported here are either [every language](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/handling-natural-languages-nlp/in-depth/supported-languages/) (this is the default) or those set by `+"`"+`queryLanguages`+"`"+`. _false_: turns off the stop words feature, allowing stop words to be taken into account in a search.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy to [remove words](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/empty-or-insufficient-results/in-depth/why-use-remove-words-if-no-results/) from the query when it doesn't match any hits. One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra content for the search UI, for example, to control the [ordering and display of facets](https://www.algolia.com/doc/guides/managing-results/rules/merchandising-and-promoting/how-to/merchandising-facets/#merchandise-facets-and-their-values-in-the-manual-editor). You can set a default value and dynamically override it with [Rules](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/).`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to highlight and snippet the original word that matches the synonym or the synonym itself.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`Attributes to include in the API response for search and browse queries.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Restrict highlighting and snippeting to items that matched the query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a query to only look at a subset of your [searchable attributes](https://www.algolia.com/doc/guides/managing-results/must-do/searchable-attributes/).`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`Assigns [rule contexts](https://www.algolia.com/doc/guides/managing-results/rules/rules-overview/how-to/customize-search-results-by-platform/#whats-a-context) to search queries.`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch. Only used when `+"`"+`mode`+"`"+` is _neuralSearch_.
`))
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`Overrides the query parameter and performs a more generic search.`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Controls how facet values are fetched.`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`Determines how to calculate [filter scores](https://www.algolia.com/doc/guides/managing-results/refine-results/filtering/in-depth/filter-scoring/#accumulating-scores-with-sumorfiltersscores).
If `+"`"+`false`+"`"+`, maximum score is kept.
If `+"`"+`true`+"`"+`, score is summed.
`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for a particular search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`[Filter hits by tags](https://www.algolia.com/doc/api-reference/api-parameters/tagFilters/).
`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Controls whether [typo tolerance](https://www.algolia.com/doc/guides/managing-results/optimize-search-results/typo-tolerance/) is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`Associates a [user token](https://www.algolia.com/doc/guides/sending-events/concepts/usertoken/) with the current search.`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

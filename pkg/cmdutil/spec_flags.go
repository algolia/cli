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
	"attributeForDistinct",
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
	"enableRules",
	"exactOnSingleWordQuery",
	"highlightPostTag",
	"highlightPreTag",
	"hitsPerPage",
	"ignorePlurals",
	"indexLanguages",
	"keepDiacriticsOnCharacters",
	"maxFacetHits",
	"minProximity",
	"minWordSizefor1Typo",
	"minWordSizefor2Typos",
	"numericAttributesForFiltering",
	"optionalWords",
	"paginationLimitedTo",
	"queryLanguages",
	"queryType",
	"ranking",
	"relevancyStrictness",
	"removeStopWords",
	"removeWordsIfNoResults",
	"renderingContent",
	"replaceSynonymsInHighlight",
	"replicas",
	"responseFields",
	"restrictHighlightAndSnippetArrays",
	"restrictSearchableAttributes",
	"searchableAttributes",
	"separatorsToIndex",
	"snippetEllipsisText",
	"synonyms",
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
	"attributeForDistinct",
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
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Enables the advanced query syntax.`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Allows you to specify which advanced syntax features are active when ‘advancedSyntax' is enabled.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers ("numeric tokens") in the query string.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`List of alternatives that should be considered an exact match by the exact ranking criterion.`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`Whether the current query will be taken into account in the Analytics.`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`List of tags to apply to the query for analytics purposes.`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Search for entries around a central geolocation, enabling a geo search within a circular area.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`Search for entries around a given location automatically computed from the requester's IP address.`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	cmd.Flags().Int("aroundPrecision", 10, heredoc.Doc(`Precision of geo search (in meters), to add grouping by geo location to the ranking formula.`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Define the maximum radius for a geo search (in meters).`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`When attribute is ranked above proximity in your ranking formula, proximity is used to select which searchable attribute is matched in the attribute ranking stage.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().String("attributeForDistinct", "", heredoc.Doc(`Name of the de-duplication attribute to be used with the distinct feature.`))
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`The complete list of attributes that will be used for faceting.`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`List of attributes to highlight.`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`This parameter controls which attributes to retrieve and which not to retrieve.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`List of attributes to snippet, with an optional maximum number of words to snippet.`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`Enable the Click Analytics feature.`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().String("cursor", "", heredoc.Doc(`Cursor indicating the location to resume browsing from. Must match the value returned by the previous call.`))
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Specifies the custom ranking criterion.`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Splits compound words into their composing atoms in the query.`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`List of attributes on which you want to disable the exact ranking criterion.`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`List of attributes on which you want to disable typo tolerance.`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Enables de-duplication or grouping of results.`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`Whether this search should participate in running AB tests.`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Enable the Personalization feature.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Whether this search should use AI Re-Ranking.`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether Rules should be globally enabled.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Controls how the exact ranking criterion is computed when the query contains only one word. One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter hits by facet value.`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`Force faceting to be applied after de-duplication (via the Distinct setting).`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`Retrieve facets and their facet values.`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter the query with numeric, facet and/or tag filters.`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("getRankingInfo", false, heredoc.Doc(`Retrieve detailed ranking information.`))
	cmd.Flags().SetAnnotation("getRankingInfo", "Categories", []string{"Advanced"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`The HTML string to insert after the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`The HTML string to insert before the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Set the number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treats singular, plurals, and other forms of declensions as matching terms.
ignorePlurals is used in conjunction with the queryLanguages setting.
list: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in queryLanguages. true: enables the ignore plurals functionality, where singulars and plurals are considered equivalent (foot = feet). The languages supported here are either every language (this is the default, see list of languages below), or those set by queryLanguages. false: disables ignore plurals, where singulars and plurals are not considered the same for matching purposes (foot will not find feet).
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().Float64Slice("insideBoundingBox", []float64{}, heredoc.Doc(`Search inside a rectangular area (in geo coordinates).`))
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().Float64Slice("insidePolygon", []float64{}, heredoc.Doc(`Search inside a polygon (in geo coordinates).`))
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`List of characters that the engine shouldn't automatically normalize.`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`Set the number of hits to retrieve (used only with offset).`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet hits to return during a search for facet values. For performance reasons, the maximum allowed number of returned values is 100.`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet during a regular search.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Precision of the proximity ranking criterion.`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with 1 typo.`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with 2 typos.`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`Minimum radius (in meters) used for a geo search when aroundRadius is not set.`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`This parameter changes the default values of certain parameters and settings that work best for a natural language query, such as ignorePlurals, removeStopWords, removeWordsIfNoResults, analyticsTags and ruleContexts. These parameters and settings work well together when the query is formatted in natural language instead of keywords, for example when your user performs a voice search.`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter on numeric attributes.`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`Specify the offset of the first hit to return.`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`Create filters for ranking purposes, where records that match the filter are ranked higher, or lower in the case of a negative optional filter.`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`A list of words that should be considered as optional when found in the query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`Specify the page to retrieve.`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`Whether to include or exclude a query from the processing-time percentile computation.`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`Define the impact of the Personalization feature.`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`The text to search in the index.`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Sets the languages to be used by language-specific settings and functionalities such as ignorePlurals, removeStopWords, and CJK word-detection.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Controls if and how query words are interpreted as prefixes. One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Controls how Algolia should sort your results.`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`When Dynamic Re-Ranking is enabled, only records that match these filters will be impacted by Dynamic Re-Ranking.`))
	cmd.Flags().SetAnnotation("reRankingApplyFilter", "Categories", []string{"Filtering"})
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Controls the relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop (common) words from the query before executing it.
removeStopWords is used in conjunction with the queryLanguages setting.
list: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in queryLanguages. true: enables the stop word functionality, ensuring that stop words are removed from consideration in a search. The languages supported here are either every language, or those set by queryLanguages. false: disables stop word functionality, allowing stop words to be taken into account in a search.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Selects a strategy to remove words from the query when it doesn't match any hits. One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Content defining how the search interface should be rendered. Can be set via the settings for a default value and can be overridden via rules.`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to highlight and snippet the original word that matches the synonym or the synonym itself.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`Choose which fields to return in the API response. This parameters applies to search and browse queries.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Restrict highlighting and snippeting to items that matched the query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a given query to look in only a subset of your searchable attributes.`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`Enables contextual rules.`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`Overrides the query parameter and performs a more generic search that can be used to find "similar" results.`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Controls how facet values are fetched.`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`Determines how to calculate the total score for filtering.`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for a particular search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter hits by tags.`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Controls whether typo tolerance is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`Associates a certain user token with the current search.`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

func AddDeleteByParamsFlags(cmd *cobra.Command) {
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Search for entries around a central geolocation, enabling a geo search within a circular area.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Define the maximum radius for a geo search (in meters).`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter hits by facet value.`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter the query with numeric, facet and/or tag filters.`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Float64Slice("insideBoundingBox", []float64{}, heredoc.Doc(`Search inside a rectangular area (in geo coordinates).`))
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().Float64Slice("insidePolygon", []float64{}, heredoc.Doc(`Search inside a polygon (in geo coordinates).`))
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter on numeric attributes.`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter hits by tags.`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
}

func AddIndexSettingsFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Enables the advanced query syntax.`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Allows you to specify which advanced syntax features are active when ‘advancedSyntax' is enabled.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowCompressionOfIntegerArray", false, heredoc.Doc(`Enables compression of large integer arrays.`))
	cmd.Flags().SetAnnotation("allowCompressionOfIntegerArray", "Categories", []string{"Performance"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers ("numeric tokens") in the query string.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`List of alternatives that should be considered an exact match by the exact ranking criterion.`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`When attribute is ranked above proximity in your ranking formula, proximity is used to select which searchable attribute is matched in the attribute ranking stage.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().String("attributeForDistinct", "", heredoc.Doc(`Name of the de-duplication attribute to be used with the distinct feature.`))
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`The complete list of attributes that will be used for faceting.`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`List of attributes to highlight.`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`This parameter controls which attributes to retrieve and which not to retrieve.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`List of attributes to snippet, with an optional maximum number of words to snippet.`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToTransliterate", []string{}, heredoc.Doc(`Specify on which attributes in your index Algolia should apply Japanese transliteration to make words indexed in Katakana or Kanji searchable in Hiragana.`))
	cmd.Flags().SetAnnotation("attributesToTransliterate", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("camelCaseAttributes", []string{}, heredoc.Doc(`List of attributes on which to do a decomposition of camel case words.`))
	cmd.Flags().SetAnnotation("camelCaseAttributes", "Categories", []string{"Languages"})
	customNormalization := NewJSONVar([]string{}...)
	cmd.Flags().Var(customNormalization, "customNormalization", heredoc.Doc(`Overrides Algolia's default normalization.`))
	cmd.Flags().SetAnnotation("customNormalization", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Specifies the custom ranking criterion.`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Splits compound words into their composing atoms in the query.`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	decompoundedAttributes := NewJSONVar([]string{}...)
	cmd.Flags().Var(decompoundedAttributes, "decompoundedAttributes", heredoc.Doc(`Specify on which attributes in your index Algolia should apply word segmentation, also known as decompounding.`))
	cmd.Flags().SetAnnotation("decompoundedAttributes", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`List of attributes on which you want to disable the exact ranking criterion.`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disablePrefixOnAttributes", []string{}, heredoc.Doc(`List of attributes on which you want to disable prefix matching.`))
	cmd.Flags().SetAnnotation("disablePrefixOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`List of attributes on which you want to disable typo tolerance.`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("disableTypoToleranceOnWords", []string{}, heredoc.Doc(`A list of words for which you want to turn off typo tolerance.`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnWords", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Enables de-duplication or grouping of results.`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Enable the Personalization feature.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether Rules should be globally enabled.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Controls how the exact ranking criterion is computed when the query contains only one word. One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`The HTML string to insert after the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`The HTML string to insert before the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Set the number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treats singular, plurals, and other forms of declensions as matching terms.
ignorePlurals is used in conjunction with the queryLanguages setting.
list: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in queryLanguages. true: enables the ignore plurals functionality, where singulars and plurals are considered equivalent (foot = feet). The languages supported here are either every language (this is the default, see list of languages below), or those set by queryLanguages. false: disables ignore plurals, where singulars and plurals are not considered the same for matching purposes (foot will not find feet).
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("indexLanguages", []string{}, heredoc.Doc(`Sets the languages at the index level for language-specific processing such as tokenization and normalization.`))
	cmd.Flags().SetAnnotation("indexLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`List of characters that the engine shouldn't automatically normalize.`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet hits to return during a search for facet values. For performance reasons, the maximum allowed number of returned values is 100.`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Precision of the proximity ranking criterion.`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with 1 typo.`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with 2 typos.`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("numericAttributesForFiltering", []string{}, heredoc.Doc(`List of numeric attributes that can be used as numerical filters.`))
	cmd.Flags().SetAnnotation("numericAttributesForFiltering", "Categories", []string{"Performance"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`A list of words that should be considered as optional when found in the query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("paginationLimitedTo", 1000, heredoc.Doc(`Set the maximum number of hits accessible via pagination.`))
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Sets the languages to be used by language-specific settings and functionalities such as ignorePlurals, removeStopWords, and CJK word-detection.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Controls if and how query words are interpreted as prefixes. One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Controls how Algolia should sort your results.`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Controls the relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop (common) words from the query before executing it.
removeStopWords is used in conjunction with the queryLanguages setting.
list: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in queryLanguages. true: enables the stop word functionality, ensuring that stop words are removed from consideration in a search. The languages supported here are either every language, or those set by queryLanguages. false: disables stop word functionality, allowing stop words to be taken into account in a search.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Selects a strategy to remove words from the query when it doesn't match any hits. One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Content defining how the search interface should be rendered. Can be set via the settings for a default value and can be overridden via rules.`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to highlight and snippet the original word that matches the synonym or the synonym itself.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("replicas", []string{}, heredoc.Doc(`Creates replicas, exact copies of an index.`))
	cmd.Flags().SetAnnotation("replicas", "Categories", []string{"Ranking"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`Choose which fields to return in the API response. This parameters applies to search and browse queries.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Restrict highlighting and snippeting to items that matched the query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a given query to look in only a subset of your searchable attributes.`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("searchableAttributes", []string{}, heredoc.Doc(`The complete list of attributes used for searching.`))
	cmd.Flags().SetAnnotation("searchableAttributes", "Categories", []string{"Attributes"})
	cmd.Flags().String("separatorsToIndex", "", heredoc.Doc(`Control which separators are indexed.`))
	cmd.Flags().SetAnnotation("separatorsToIndex", "Categories", []string{"Typos"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for a particular search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Controls whether typo tolerance is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("unretrievableAttributes", []string{}, heredoc.Doc(`List of attributes that can't be retrieved at query time.`))
	cmd.Flags().SetAnnotation("unretrievableAttributes", "Categories", []string{"Attributes"})
	userData := NewJSONVar([]string{}...)
	cmd.Flags().Var(userData, "userData", heredoc.Doc(`Lets you store custom data in your indices.`))
	cmd.Flags().SetAnnotation("userData", "Categories", []string{"Advanced"})
}

func AddSearchParamsObjectFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Enables the advanced query syntax.`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Allows you to specify which advanced syntax features are active when ‘advancedSyntax' is enabled.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers ("numeric tokens") in the query string.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`List of alternatives that should be considered an exact match by the exact ranking criterion.`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`Whether the current query will be taken into account in the Analytics.`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`List of tags to apply to the query for analytics purposes.`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Search for entries around a central geolocation, enabling a geo search within a circular area.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`Search for entries around a given location automatically computed from the requester's IP address.`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	cmd.Flags().Int("aroundPrecision", 10, heredoc.Doc(`Precision of geo search (in meters), to add grouping by geo location to the ranking formula.`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Define the maximum radius for a geo search (in meters).`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`When attribute is ranked above proximity in your ranking formula, proximity is used to select which searchable attribute is matched in the attribute ranking stage.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().String("attributeForDistinct", "", heredoc.Doc(`Name of the de-duplication attribute to be used with the distinct feature.`))
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`The complete list of attributes that will be used for faceting.`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`List of attributes to highlight.`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`This parameter controls which attributes to retrieve and which not to retrieve.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`List of attributes to snippet, with an optional maximum number of words to snippet.`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`Enable the Click Analytics feature.`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Specifies the custom ranking criterion.`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Splits compound words into their composing atoms in the query.`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`List of attributes on which you want to disable the exact ranking criterion.`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`List of attributes on which you want to disable typo tolerance.`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Enables de-duplication or grouping of results.`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`Whether this search should participate in running AB tests.`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Enable the Personalization feature.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Whether this search should use AI Re-Ranking.`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether Rules should be globally enabled.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Controls how the exact ranking criterion is computed when the query contains only one word. One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter hits by facet value.`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`Force faceting to be applied after de-duplication (via the Distinct setting).`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`Retrieve facets and their facet values.`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter the query with numeric, facet and/or tag filters.`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("getRankingInfo", false, heredoc.Doc(`Retrieve detailed ranking information.`))
	cmd.Flags().SetAnnotation("getRankingInfo", "Categories", []string{"Advanced"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`The HTML string to insert after the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`The HTML string to insert before the highlighted parts in all highlight and snippet results.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Set the number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treats singular, plurals, and other forms of declensions as matching terms.
ignorePlurals is used in conjunction with the queryLanguages setting.
list: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in queryLanguages. true: enables the ignore plurals functionality, where singulars and plurals are considered equivalent (foot = feet). The languages supported here are either every language (this is the default, see list of languages below), or those set by queryLanguages. false: disables ignore plurals, where singulars and plurals are not considered the same for matching purposes (foot will not find feet).
`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().Float64Slice("insideBoundingBox", []float64{}, heredoc.Doc(`Search inside a rectangular area (in geo coordinates).`))
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().Float64Slice("insidePolygon", []float64{}, heredoc.Doc(`Search inside a polygon (in geo coordinates).`))
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`List of characters that the engine shouldn't automatically normalize.`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`Set the number of hits to retrieve (used only with offset).`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet hits to return during a search for facet values. For performance reasons, the maximum allowed number of returned values is 100.`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet during a regular search.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Precision of the proximity ranking criterion.`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with 1 typo.`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the query string must contain to accept matches with 2 typos.`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`Minimum radius (in meters) used for a geo search when aroundRadius is not set.`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`This parameter changes the default values of certain parameters and settings that work best for a natural language query, such as ignorePlurals, removeStopWords, removeWordsIfNoResults, analyticsTags and ruleContexts. These parameters and settings work well together when the query is formatted in natural language instead of keywords, for example when your user performs a voice search.`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter on numeric attributes.`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`Specify the offset of the first hit to return.`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`Create filters for ranking purposes, where records that match the filter are ranked higher, or lower in the case of a negative optional filter.`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`A list of words that should be considered as optional when found in the query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`Specify the page to retrieve.`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`Whether to include or exclude a query from the processing-time percentile computation.`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`Define the impact of the Personalization feature.`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`The text to search in the index.`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Sets the languages to be used by language-specific settings and functionalities such as ignorePlurals, removeStopWords, and CJK word-detection.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Controls if and how query words are interpreted as prefixes. One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Controls how Algolia should sort your results.`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`When Dynamic Re-Ranking is enabled, only records that match these filters will be impacted by Dynamic Re-Ranking.`))
	cmd.Flags().SetAnnotation("reRankingApplyFilter", "Categories", []string{"Filtering"})
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Controls the relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop (common) words from the query before executing it.
removeStopWords is used in conjunction with the queryLanguages setting.
list: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in queryLanguages. true: enables the stop word functionality, ensuring that stop words are removed from consideration in a search. The languages supported here are either every language, or those set by queryLanguages. false: disables stop word functionality, allowing stop words to be taken into account in a search.
`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Selects a strategy to remove words from the query when it doesn't match any hits. One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Content defining how the search interface should be rendered. Can be set via the settings for a default value and can be overridden via rules.`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to highlight and snippet the original word that matches the synonym or the synonym itself.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`Choose which fields to return in the API response. This parameters applies to search and browse queries.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Restrict highlighting and snippeting to items that matched the query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a given query to look in only a subset of your searchable attributes.`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`Enables contextual rules.`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`Overrides the query parameter and performs a more generic search that can be used to find "similar" results.`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Controls how facet values are fetched.`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`Determines how to calculate the total score for filtering.`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for a particular search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter hits by tags.`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Controls whether typo tolerance is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`Associates a certain user token with the current search.`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

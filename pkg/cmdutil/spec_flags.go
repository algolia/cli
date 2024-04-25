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
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Whether to support phrase matching and excluding words from search queries.`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Advanced search syntax features you want to support.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers in the search query.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives of query words that should be considered as exact matches by the Exact ranking criterion.`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`Whether this search will be included in Analytics.`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`Tags to apply to the query for segmenting analytics data.`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Coordinates for the center of a circle, expressed as a comma-separated string of latitude and longitude.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`Whether to obtain the coordinates from the request's IP address.`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	aroundPrecision := NewJSONVar([]string{"integer", "array"}...)
	cmd.Flags().Var(aroundPrecision, "aroundPrecision", heredoc.Doc(`Precision of a coordinate-based search in meters to group results with similar distances.`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Maximum radius for a search around a central location.`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`Whether the best matching attribute should be determined by minimum proximity.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight.`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes for which to enable snippets.`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`Whether to include a queryID attribute in the response.`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().String("cursor", "", heredoc.Doc(`Cursor to get the next page of the response.`))
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Attributes to use as custom ranking.`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Whether to split compound words into their building blocks.`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Searchable attributes for which you want to turn off the Exact ranking criterion.`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off typo tolerance.`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Determines how many records of a group are included in the search results.`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`Whether to enable A/B testing for this search.`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Whether to enable Personalization.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Whether this search will use Dynamic Re-Ranking.`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether to enable rules.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the Exact ranking criterion is computed when the search query has only one word. One of: attribute, none, word.`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter the search by facet values, so that only records with the same facet values are retrieved.`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`Whether faceting should be applied after deduplication with distinct.`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`Facets for which to retrieve facet values that match the search criteria and the number of matching facet values.`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter expression to only include items that match the filter criteria in the response.`))
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
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treat singular, plurals, and other forms of declensions as equivalent.`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters for which diacritics should be preserved.`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`Number of hits to retrieve (used in combination with offset).`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet values to return when searching for facet values.`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Minimum proximity score for two matching words.`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with one typo.`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with two typos.`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`Minimum radius (in meters) for a search around a location when aroundRadius isn't set.`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results. One of: neuralSearch, keywordSearch.`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`ISO language codes that adjust settings that are useful for processing natural language queries (as opposed to keyword searches).`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter by numeric facets.`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`Position of the first hit to retrieve.`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`Filters to promote or demote records in the search results.`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words that should be considered optional when found in the query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`Page of search results to retrieve.`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`Whether to include this search when calculating processing-time percentiles.`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`Impact that Personalization should have on this search.`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`Search query.`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Languages for language-specific query processing steps such as plurals, stop-word removal, and word-detection dictionaries.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines if and how query words are interpreted as prefixes. One of: prefixLast, prefixAll, prefixNone.`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia returns your results.`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string", "null"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`Restrict Dynamic Re-Ranking to records that match these filters.`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop words from the search query.`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy for removing words from the query when it doesn't return any results. One of: none, lastWords, firstWords, allOptional.`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra data that can be used in the search UI.`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to replace a highlighted word with the matched synonym.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{"*"}, heredoc.Doc(`Properties to include in the API response of search and browse requests.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Whether to restrict highlighting and snippeting to items that at least partially matched the search query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a search to a subset of your searchable attributes.`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`Assigns a rule context to the search query.`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch.`))
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`Keywords to be used instead of the search query to conduct a more broader search.`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Order in which to retrieve facet values.`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`Whether to sum all filter scores.`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for this search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter the search by values of the special _tags attribute.`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Whether typo tolerance is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`Unique pseudonymous or anonymous user identifier.`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

func AddDeleteByParamsFlags(cmd *cobra.Command) {
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Coordinates for the center of a circle, expressed as a comma-separated string of latitude and longitude.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Maximum radius for a search around a central location.`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter the search by facet values, so that only records with the same facet values are retrieved.`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter expression to only include items that match the filter criteria in the response.`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter by numeric facets.`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter the search by values of the special _tags attribute.`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
}

func AddIndexSettingsFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Whether to support phrase matching and excluding words from search queries.`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Advanced search syntax features you want to support.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowCompressionOfIntegerArray", false, heredoc.Doc(`Whether arrays with exclusively non-negative integers should be compressed for better performance.`))
	cmd.Flags().SetAnnotation("allowCompressionOfIntegerArray", "Categories", []string{"Performance"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers in the search query.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives of query words that should be considered as exact matches by the Exact ranking criterion.`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`Whether the best matching attribute should be determined by minimum proximity.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().String("attributeForDistinct", "", heredoc.Doc(`Attribute that should be used to establish groups of results.`))
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`Attributes used for faceting.`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight.`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes for which to enable snippets.`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToTransliterate", []string{}, heredoc.Doc(`Attributes, for which you want to support Japanese transliteration.`))
	cmd.Flags().SetAnnotation("attributesToTransliterate", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("camelCaseAttributes", []string{}, heredoc.Doc(`Attributes for which to split camel case words.`))
	cmd.Flags().SetAnnotation("camelCaseAttributes", "Categories", []string{"Languages"})
	customNormalization := NewJSONVar([]string{}...)
	cmd.Flags().Var(customNormalization, "customNormalization", heredoc.Doc(`Characters and their normalized replacements.`))
	cmd.Flags().SetAnnotation("customNormalization", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Attributes to use as custom ranking.`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Whether to split compound words into their building blocks.`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	decompoundedAttributes := NewJSONVar([]string{}...)
	cmd.Flags().Var(decompoundedAttributes, "decompoundedAttributes", heredoc.Doc(`Searchable attributes to which Algolia should apply word segmentation (decompounding).`))
	cmd.Flags().SetAnnotation("decompoundedAttributes", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Searchable attributes for which you want to turn off the Exact ranking criterion.`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disablePrefixOnAttributes", []string{}, heredoc.Doc(`Searchable attributes for which you want to turn off prefix matching.`))
	cmd.Flags().SetAnnotation("disablePrefixOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off typo tolerance.`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("disableTypoToleranceOnWords", []string{}, heredoc.Doc(`Words for which you want to turn off typo tolerance.`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnWords", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Determines how many records of a group are included in the search results.`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Whether to enable Personalization.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Whether this search will use Dynamic Re-Ranking.`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether to enable rules.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the Exact ranking criterion is computed when the search query has only one word. One of: attribute, none, word.`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`HTML tag to insert after the highlighted parts in all highlighted results and snippets.`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`HTML tag to insert before the highlighted parts in all highlighted results and snippets.`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`Number of hits per page.`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treat singular, plurals, and other forms of declensions as equivalent.`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("indexLanguages", []string{}, heredoc.Doc(`Languages for language-specific processing steps, such as word detection and dictionary settings.`))
	cmd.Flags().SetAnnotation("indexLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters for which diacritics should be preserved.`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet values to return when searching for facet values.`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Minimum proximity score for two matching words.`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with one typo.`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with two typos.`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results. One of: neuralSearch, keywordSearch.`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("numericAttributesForFiltering", []string{}, heredoc.Doc(`Numeric attributes that can be used as numerical filters.`))
	cmd.Flags().SetAnnotation("numericAttributesForFiltering", "Categories", []string{"Performance"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words that should be considered optional when found in the query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("paginationLimitedTo", 1000, heredoc.Doc(`Maximum number of search results that can be obtained through pagination.`))
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Languages for language-specific query processing steps such as plurals, stop-word removal, and word-detection dictionaries.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines if and how query words are interpreted as prefixes. One of: prefixLast, prefixAll, prefixNone.`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia returns your results.`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string", "null"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`Restrict Dynamic Re-Ranking to records that match these filters.`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop words from the search query.`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy for removing words from the query when it doesn't return any results. One of: none, lastWords, firstWords, allOptional.`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra data that can be used in the search UI.`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to replace a highlighted word with the matched synonym.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("replicas", []string{}, heredoc.Doc(`Creates replica indices.`))
	cmd.Flags().SetAnnotation("replicas", "Categories", []string{"Ranking"})
	cmd.Flags().StringSlice("responseFields", []string{"*"}, heredoc.Doc(`Properties to include in the API response of search and browse requests.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Whether to restrict highlighting and snippeting to items that at least partially matched the search query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("searchableAttributes", []string{}, heredoc.Doc(`Attributes used for searching. Attribute names are case-sensitive.`))
	cmd.Flags().SetAnnotation("searchableAttributes", "Categories", []string{"Attributes"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch.`))
	cmd.Flags().String("separatorsToIndex", "", heredoc.Doc(`Controls which separators are indexed.`))
	cmd.Flags().SetAnnotation("separatorsToIndex", "Categories", []string{"Typos"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Order in which to retrieve facet values.`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Whether typo tolerance is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("unretrievableAttributes", []string{}, heredoc.Doc(`Attributes that can't be retrieved at query time.`))
	cmd.Flags().SetAnnotation("unretrievableAttributes", "Categories", []string{"Attributes"})
	userData := NewJSONVar([]string{}...)
	cmd.Flags().Var(userData, "userData", heredoc.Doc(`An object with custom data.`))
	cmd.Flags().SetAnnotation("userData", "Categories", []string{"Advanced"})
}

func AddSearchParamsObjectFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`Whether to support phrase matching and excluding words from search queries.`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`Advanced search syntax features you want to support.`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`Whether to allow typos on numbers in the search query.`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`Alternatives of query words that should be considered as exact matches by the Exact ranking criterion.`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`Whether this search will be included in Analytics.`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`Tags to apply to the query for segmenting analytics data.`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`Coordinates for the center of a circle, expressed as a comma-separated string of latitude and longitude.`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`Whether to obtain the coordinates from the request's IP address.`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	aroundPrecision := NewJSONVar([]string{"integer", "array"}...)
	cmd.Flags().Var(aroundPrecision, "aroundPrecision", heredoc.Doc(`Precision of a coordinate-based search in meters to group results with similar distances.`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`Maximum radius for a search around a central location.`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`Whether the best matching attribute should be determined by minimum proximity.`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`Attributes to highlight.`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`Attributes to include in the API response.`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`Attributes for which to enable snippets.`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`Whether to include a queryID attribute in the response.`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`Attributes to use as custom ranking.`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`Whether to split compound words into their building blocks.`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`Searchable attributes for which you want to turn off the Exact ranking criterion.`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`Attributes for which you want to turn off typo tolerance.`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`Determines how many records of a group are included in the search results.`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`Whether to enable A/B testing for this search.`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`Whether to enable Personalization.`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`Whether this search will use Dynamic Re-Ranking.`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`Whether to enable rules.`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`Determines how the Exact ranking criterion is computed when the search query has only one word. One of: attribute, none, word.`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`Filter the search by facet values, so that only records with the same facet values are retrieved.`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`Whether faceting should be applied after deduplication with distinct.`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`Facets for which to retrieve facet values that match the search criteria and the number of matching facet values.`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`Filter expression to only include items that match the filter criteria in the response.`))
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
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`Treat singular, plurals, and other forms of declensions as equivalent.`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`Characters for which diacritics should be preserved.`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`Number of hits to retrieve (used in combination with offset).`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`Maximum number of facet values to return when searching for facet values.`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`Maximum number of facet values to return for each facet.`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`Minimum proximity score for two matching words.`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with one typo.`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`Minimum number of characters a word in the search query must contain to accept matches with two typos.`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`Minimum radius (in meters) for a search around a location when aroundRadius isn't set.`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`Search mode the index will use to query for results. One of: neuralSearch, keywordSearch.`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`ISO language codes that adjust settings that are useful for processing natural language queries (as opposed to keyword searches).`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`Filter by numeric facets.`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`Position of the first hit to retrieve.`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`Filters to promote or demote records in the search results.`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`Words that should be considered optional when found in the query.`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`Page of search results to retrieve.`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`Whether to include this search when calculating processing-time percentiles.`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`Impact that Personalization should have on this search.`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`Search query.`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`Languages for language-specific query processing steps such as plurals, stop-word removal, and word-detection dictionaries.`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`Determines if and how query words are interpreted as prefixes. One of: prefixLast, prefixAll, prefixNone.`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`Determines the order in which Algolia returns your results.`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string", "null"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`Restrict Dynamic Re-Ranking to records that match these filters.`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`Relevancy threshold below which less relevant results aren't included in the results.`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`Removes stop words from the search query.`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`Strategy for removing words from the query when it doesn't return any results. One of: none, lastWords, firstWords, allOptional.`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`Extra data that can be used in the search UI.`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`Whether to replace a highlighted word with the matched synonym.`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{"*"}, heredoc.Doc(`Properties to include in the API response of search and browse requests.`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`Whether to restrict highlighting and snippeting to items that at least partially matched the search query.`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`Restricts a search to a subset of your searchable attributes.`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`Assigns a rule context to the search query.`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`Settings for the semantic search part of NeuralSearch.`))
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`Keywords to be used instead of the search query to conduct a more broader search.`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`String used as an ellipsis indicator when a snippet is truncated.`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`Order in which to retrieve facet values.`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`Whether to sum all filter scores.`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`Whether to take into account an index's synonyms for this search.`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`Filter the search by values of the special _tags attribute.`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`Whether typo tolerance is enabled and how it is applied.`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`Unique pseudonymous or anonymous user identifier.`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

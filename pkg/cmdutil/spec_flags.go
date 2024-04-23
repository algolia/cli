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
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/advancedSyntax/`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/advancedSyntaxFeatures/`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/allowTyposOnNumericTokens/`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/alternativesAsExact/`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/analytics/`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/analyticsTags/`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundLatLng/`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundLatLngViaIP/`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	aroundPrecision := NewJSONVar([]string{"integer", "array"}...)
	cmd.Flags().Var(aroundPrecision, "aroundPrecision", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundPrecision/`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundRadius/`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributeCriteriaComputedByMinProximity/`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesForFaceting/`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToHighlight/`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToRetrieve/`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToSnippet/`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/clickAnalytics/`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().String("cursor", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/cursor/`))
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/customRanking/`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/decompoundQuery/`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/disableExactOnAttributes/`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/disableTypoToleranceOnAttributes/`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/distinct/`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enableABTest/`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enablePersonalization/`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enableReRanking/`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enableRules/`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/exactOnSingleWordQuery/ One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("explain", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/explain/`))
	cmd.Flags().SetAnnotation("explain", "Categories", []string{"Advanced"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/facetFilters/`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/facetingAfterDistinct/`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/facets/`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/filters/`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("getRankingInfo", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/getRankingInfo/`))
	cmd.Flags().SetAnnotation("getRankingInfo", "Categories", []string{"Advanced"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/highlightPostTag/`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/highlightPreTag/`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/hitsPerPage/`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/ignorePlurals/`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/keepDiacriticsOnCharacters/`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/length/`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/maxFacetHits/`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/maxValuesPerFacet/`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minProximity/`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minWordSizefor1Typo/`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minWordSizefor2Typos/`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minimumAroundRadius/`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/mode/ One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/naturalLanguages/`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/numericFilters/`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/offset/`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/optionalFilters/`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/optionalWords/`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/page/`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/percentileComputation/`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/personalizationImpact/`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/query/`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/queryLanguages/`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/queryType/ One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/ranking/`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/reRankingApplyFilter/`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/relevancyStrictness/`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/removeStopWords/`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/removeWordsIfNoResults/ One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/renderingContent/`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/replaceSynonymsInHighlight/`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/responseFields/`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/restrictHighlightAndSnippetArrays/`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/restrictSearchableAttributes/`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/ruleContexts/`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/semanticSearch/`))
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/similarQuery/`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/snippetEllipsisText/`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/sortFacetValuesBy/`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/sumOrFiltersScores/`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/synonyms/`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/tagFilters/`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/typoTolerance/`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/userToken/`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

func AddDeleteByParamsFlags(cmd *cobra.Command) {
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundLatLng/`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundRadius/`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/facetFilters/`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().String("filters", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/filters/`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/numericFilters/`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/tagFilters/`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
}

func AddIndexSettingsFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/advancedSyntax/`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/advancedSyntaxFeatures/`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowCompressionOfIntegerArray", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/allowCompressionOfIntegerArray/`))
	cmd.Flags().SetAnnotation("allowCompressionOfIntegerArray", "Categories", []string{"Performance"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/allowTyposOnNumericTokens/`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/alternativesAsExact/`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributeCriteriaComputedByMinProximity/`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().String("attributeForDistinct", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributeForDistinct/`))
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesForFaceting/`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToHighlight/`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToRetrieve/`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToSnippet/`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToTransliterate", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToTransliterate/`))
	cmd.Flags().SetAnnotation("attributesToTransliterate", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("camelCaseAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/camelCaseAttributes/`))
	cmd.Flags().SetAnnotation("camelCaseAttributes", "Categories", []string{"Languages"})
	customNormalization := NewJSONVar([]string{}...)
	cmd.Flags().Var(customNormalization, "customNormalization", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/customNormalization/`))
	cmd.Flags().SetAnnotation("customNormalization", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/customRanking/`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/decompoundQuery/`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	decompoundedAttributes := NewJSONVar([]string{}...)
	cmd.Flags().Var(decompoundedAttributes, "decompoundedAttributes", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/decompoundedAttributes/`))
	cmd.Flags().SetAnnotation("decompoundedAttributes", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/disableExactOnAttributes/`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disablePrefixOnAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/disablePrefixOnAttributes/`))
	cmd.Flags().SetAnnotation("disablePrefixOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/disableTypoToleranceOnAttributes/`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("disableTypoToleranceOnWords", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/disableTypoToleranceOnWords/`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnWords", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/distinct/`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enablePersonalization/`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enableReRanking/`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enableRules/`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/exactOnSingleWordQuery/ One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/highlightPostTag/`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/highlightPreTag/`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/hitsPerPage/`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/ignorePlurals/`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("indexLanguages", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/indexLanguages/`))
	cmd.Flags().SetAnnotation("indexLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/keepDiacriticsOnCharacters/`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/maxFacetHits/`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/maxValuesPerFacet/`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minProximity/`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minWordSizefor1Typo/`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minWordSizefor2Typos/`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/mode/ One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("numericAttributesForFiltering", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/numericAttributesForFiltering/`))
	cmd.Flags().SetAnnotation("numericAttributesForFiltering", "Categories", []string{"Performance"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/optionalWords/`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("paginationLimitedTo", 1000, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/paginationLimitedTo/`))
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/queryLanguages/`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/queryType/ One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/ranking/`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/reRankingApplyFilter/`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/relevancyStrictness/`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/removeStopWords/`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/removeWordsIfNoResults/ One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/renderingContent/`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/replaceSynonymsInHighlight/`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("replicas", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/replicas/`))
	cmd.Flags().SetAnnotation("replicas", "Categories", []string{"Ranking"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/responseFields/`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/restrictHighlightAndSnippetArrays/`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("searchableAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/searchableAttributes/`))
	cmd.Flags().SetAnnotation("searchableAttributes", "Categories", []string{"Attributes"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/semanticSearch/`))
	cmd.Flags().String("separatorsToIndex", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/separatorsToIndex/`))
	cmd.Flags().SetAnnotation("separatorsToIndex", "Categories", []string{"Typos"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/snippetEllipsisText/`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/sortFacetValuesBy/`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/typoTolerance/`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("unretrievableAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/unretrievableAttributes/`))
	cmd.Flags().SetAnnotation("unretrievableAttributes", "Categories", []string{"Attributes"})
	userData := NewJSONVar([]string{}...)
	cmd.Flags().Var(userData, "userData", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/userData/`))
	cmd.Flags().SetAnnotation("userData", "Categories", []string{"Advanced"})
}

func AddSearchParamsObjectFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("advancedSyntax", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/advancedSyntax/`))
	cmd.Flags().SetAnnotation("advancedSyntax", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("advancedSyntaxFeatures", []string{"exactPhrase", "excludeWords"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/advancedSyntaxFeatures/`))
	cmd.Flags().SetAnnotation("advancedSyntaxFeatures", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("allowTyposOnNumericTokens", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/allowTyposOnNumericTokens/`))
	cmd.Flags().SetAnnotation("allowTyposOnNumericTokens", "Categories", []string{"Typos"})
	cmd.Flags().StringSlice("alternativesAsExact", []string{"ignorePlurals", "singleWordSynonym"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/alternativesAsExact/`))
	cmd.Flags().SetAnnotation("alternativesAsExact", "Categories", []string{"Query strategy"})
	cmd.Flags().Bool("analytics", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/analytics/`))
	cmd.Flags().SetAnnotation("analytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("analyticsTags", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/analyticsTags/`))
	cmd.Flags().SetAnnotation("analyticsTags", "Categories", []string{"Analytics"})
	cmd.Flags().String("aroundLatLng", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundLatLng/`))
	cmd.Flags().SetAnnotation("aroundLatLng", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("aroundLatLngViaIP", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundLatLngViaIP/`))
	cmd.Flags().SetAnnotation("aroundLatLngViaIP", "Categories", []string{"Geo-Search"})
	aroundPrecision := NewJSONVar([]string{"integer", "array"}...)
	cmd.Flags().Var(aroundPrecision, "aroundPrecision", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundPrecision/`))
	cmd.Flags().SetAnnotation("aroundPrecision", "Categories", []string{"Geo-Search"})
	aroundRadius := NewJSONVar([]string{"integer", "string"}...)
	cmd.Flags().Var(aroundRadius, "aroundRadius", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/aroundRadius/`))
	cmd.Flags().SetAnnotation("aroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().Bool("attributeCriteriaComputedByMinProximity", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributeCriteriaComputedByMinProximity/`))
	cmd.Flags().SetAnnotation("attributeCriteriaComputedByMinProximity", "Categories", []string{"Advanced"})
	cmd.Flags().StringSlice("attributesForFaceting", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesForFaceting/`))
	cmd.Flags().SetAnnotation("attributesForFaceting", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("attributesToHighlight", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToHighlight/`))
	cmd.Flags().SetAnnotation("attributesToHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("attributesToRetrieve", []string{"*"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToRetrieve/`))
	cmd.Flags().SetAnnotation("attributesToRetrieve", "Categories", []string{"Attributes"})
	cmd.Flags().StringSlice("attributesToSnippet", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/attributesToSnippet/`))
	cmd.Flags().SetAnnotation("attributesToSnippet", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Bool("clickAnalytics", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/clickAnalytics/`))
	cmd.Flags().SetAnnotation("clickAnalytics", "Categories", []string{"Analytics"})
	cmd.Flags().StringSlice("customRanking", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/customRanking/`))
	cmd.Flags().SetAnnotation("customRanking", "Categories", []string{"Ranking"})
	cmd.Flags().Bool("decompoundQuery", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/decompoundQuery/`))
	cmd.Flags().SetAnnotation("decompoundQuery", "Categories", []string{"Languages"})
	cmd.Flags().StringSlice("disableExactOnAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/disableExactOnAttributes/`))
	cmd.Flags().SetAnnotation("disableExactOnAttributes", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("disableTypoToleranceOnAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/disableTypoToleranceOnAttributes/`))
	cmd.Flags().SetAnnotation("disableTypoToleranceOnAttributes", "Categories", []string{"Typos"})
	distinct := NewJSONVar([]string{"boolean", "integer"}...)
	cmd.Flags().Var(distinct, "distinct", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/distinct/`))
	cmd.Flags().SetAnnotation("distinct", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enableABTest", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enableABTest/`))
	cmd.Flags().SetAnnotation("enableABTest", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("enablePersonalization", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enablePersonalization/`))
	cmd.Flags().SetAnnotation("enablePersonalization", "Categories", []string{"Personalization"})
	cmd.Flags().Bool("enableReRanking", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enableReRanking/`))
	cmd.Flags().SetAnnotation("enableReRanking", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("enableRules", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/enableRules/`))
	cmd.Flags().SetAnnotation("enableRules", "Categories", []string{"Rules"})
	cmd.Flags().String("exactOnSingleWordQuery", "attribute", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/exactOnSingleWordQuery/ One of: (attribute, none, word).`))
	cmd.Flags().SetAnnotation("exactOnSingleWordQuery", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("explain", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/explain/`))
	cmd.Flags().SetAnnotation("explain", "Categories", []string{"Advanced"})
	facetFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(facetFilters, "facetFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/facetFilters/`))
	cmd.Flags().SetAnnotation("facetFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("facetingAfterDistinct", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/facetingAfterDistinct/`))
	cmd.Flags().SetAnnotation("facetingAfterDistinct", "Categories", []string{"Faceting"})
	cmd.Flags().StringSlice("facets", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/facets/`))
	cmd.Flags().SetAnnotation("facets", "Categories", []string{"Faceting"})
	cmd.Flags().String("filters", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/filters/`))
	cmd.Flags().SetAnnotation("filters", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("getRankingInfo", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/getRankingInfo/`))
	cmd.Flags().SetAnnotation("getRankingInfo", "Categories", []string{"Advanced"})
	cmd.Flags().String("highlightPostTag", "</em>", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/highlightPostTag/`))
	cmd.Flags().SetAnnotation("highlightPostTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("highlightPreTag", "<em>", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/highlightPreTag/`))
	cmd.Flags().SetAnnotation("highlightPreTag", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().Int("hitsPerPage", 20, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/hitsPerPage/`))
	cmd.Flags().SetAnnotation("hitsPerPage", "Categories", []string{"Pagination"})
	ignorePlurals := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(ignorePlurals, "ignorePlurals", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/ignorePlurals/`))
	cmd.Flags().SetAnnotation("ignorePlurals", "Categories", []string{"Languages"})
	cmd.Flags().SetAnnotation("insideBoundingBox", "Categories", []string{"Geo-Search"})
	cmd.Flags().SetAnnotation("insidePolygon", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("keepDiacriticsOnCharacters", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/keepDiacriticsOnCharacters/`))
	cmd.Flags().SetAnnotation("keepDiacriticsOnCharacters", "Categories", []string{"Languages"})
	cmd.Flags().Int("length", 0, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/length/`))
	cmd.Flags().SetAnnotation("length", "Categories", []string{"Pagination"})
	cmd.Flags().Int("maxFacetHits", 10, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/maxFacetHits/`))
	cmd.Flags().SetAnnotation("maxFacetHits", "Categories", []string{"Advanced"})
	cmd.Flags().Int("maxValuesPerFacet", 100, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/maxValuesPerFacet/`))
	cmd.Flags().SetAnnotation("maxValuesPerFacet", "Categories", []string{"Faceting"})
	cmd.Flags().Int("minProximity", 1, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minProximity/`))
	cmd.Flags().SetAnnotation("minProximity", "Categories", []string{"Advanced"})
	cmd.Flags().Int("minWordSizefor1Typo", 4, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minWordSizefor1Typo/`))
	cmd.Flags().SetAnnotation("minWordSizefor1Typo", "Categories", []string{"Typos"})
	cmd.Flags().Int("minWordSizefor2Typos", 8, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minWordSizefor2Typos/`))
	cmd.Flags().SetAnnotation("minWordSizefor2Typos", "Categories", []string{"Typos"})
	cmd.Flags().Int("minimumAroundRadius", 0, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/minimumAroundRadius/`))
	cmd.Flags().SetAnnotation("minimumAroundRadius", "Categories", []string{"Geo-Search"})
	cmd.Flags().String("mode", "keywordSearch", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/mode/ One of: (neuralSearch, keywordSearch).`))
	cmd.Flags().SetAnnotation("mode", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("naturalLanguages", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/naturalLanguages/`))
	cmd.Flags().SetAnnotation("naturalLanguages", "Categories", []string{"Languages"})
	numericFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(numericFilters, "numericFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/numericFilters/`))
	cmd.Flags().SetAnnotation("numericFilters", "Categories", []string{"Filtering"})
	cmd.Flags().Int("offset", 0, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/offset/`))
	cmd.Flags().SetAnnotation("offset", "Categories", []string{"Pagination"})
	optionalFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(optionalFilters, "optionalFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/optionalFilters/`))
	cmd.Flags().SetAnnotation("optionalFilters", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("optionalWords", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/optionalWords/`))
	cmd.Flags().SetAnnotation("optionalWords", "Categories", []string{"Query strategy"})
	cmd.Flags().Int("page", 0, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/page/`))
	cmd.Flags().SetAnnotation("page", "Categories", []string{"Pagination"})
	cmd.Flags().Bool("percentileComputation", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/percentileComputation/`))
	cmd.Flags().SetAnnotation("percentileComputation", "Categories", []string{"Advanced"})
	cmd.Flags().Int("personalizationImpact", 100, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/personalizationImpact/`))
	cmd.Flags().SetAnnotation("personalizationImpact", "Categories", []string{"Personalization"})
	cmd.Flags().String("query", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/query/`))
	cmd.Flags().SetAnnotation("query", "Categories", []string{"Search"})
	cmd.Flags().StringSlice("queryLanguages", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/queryLanguages/`))
	cmd.Flags().SetAnnotation("queryLanguages", "Categories", []string{"Languages"})
	cmd.Flags().String("queryType", "prefixLast", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/queryType/ One of: (prefixLast, prefixAll, prefixNone).`))
	cmd.Flags().SetAnnotation("queryType", "Categories", []string{"Query strategy"})
	cmd.Flags().StringSlice("ranking", []string{"typo", "geo", "words", "filters", "proximity", "attribute", "exact", "custom"}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/ranking/`))
	cmd.Flags().SetAnnotation("ranking", "Categories", []string{"Ranking"})
	reRankingApplyFilter := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(reRankingApplyFilter, "reRankingApplyFilter", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/reRankingApplyFilter/`))
	cmd.Flags().Int("relevancyStrictness", 100, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/relevancyStrictness/`))
	cmd.Flags().SetAnnotation("relevancyStrictness", "Categories", []string{"Ranking"})
	removeStopWords := NewJSONVar([]string{"array", "boolean"}...)
	cmd.Flags().Var(removeStopWords, "removeStopWords", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/removeStopWords/`))
	cmd.Flags().SetAnnotation("removeStopWords", "Categories", []string{"Languages"})
	cmd.Flags().String("removeWordsIfNoResults", "none", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/removeWordsIfNoResults/ One of: (none, lastWords, firstWords, allOptional).`))
	cmd.Flags().SetAnnotation("removeWordsIfNoResults", "Categories", []string{"Query strategy"})
	renderingContent := NewJSONVar([]string{}...)
	cmd.Flags().Var(renderingContent, "renderingContent", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/renderingContent/`))
	cmd.Flags().SetAnnotation("renderingContent", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("replaceSynonymsInHighlight", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/replaceSynonymsInHighlight/`))
	cmd.Flags().SetAnnotation("replaceSynonymsInHighlight", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("responseFields", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/responseFields/`))
	cmd.Flags().SetAnnotation("responseFields", "Categories", []string{"Advanced"})
	cmd.Flags().Bool("restrictHighlightAndSnippetArrays", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/restrictHighlightAndSnippetArrays/`))
	cmd.Flags().SetAnnotation("restrictHighlightAndSnippetArrays", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().StringSlice("restrictSearchableAttributes", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/restrictSearchableAttributes/`))
	cmd.Flags().SetAnnotation("restrictSearchableAttributes", "Categories", []string{"Filtering"})
	cmd.Flags().StringSlice("ruleContexts", []string{}, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/ruleContexts/`))
	cmd.Flags().SetAnnotation("ruleContexts", "Categories", []string{"Rules"})
	semanticSearch := NewJSONVar([]string{}...)
	cmd.Flags().Var(semanticSearch, "semanticSearch", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/semanticSearch/`))
	cmd.Flags().String("similarQuery", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/similarQuery/`))
	cmd.Flags().SetAnnotation("similarQuery", "Categories", []string{"Search"})
	cmd.Flags().String("snippetEllipsisText", "…", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/snippetEllipsisText/`))
	cmd.Flags().SetAnnotation("snippetEllipsisText", "Categories", []string{"Highlighting and Snippeting"})
	cmd.Flags().String("sortFacetValuesBy", "count", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/sortFacetValuesBy/`))
	cmd.Flags().SetAnnotation("sortFacetValuesBy", "Categories", []string{"Faceting"})
	cmd.Flags().Bool("sumOrFiltersScores", false, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/sumOrFiltersScores/`))
	cmd.Flags().SetAnnotation("sumOrFiltersScores", "Categories", []string{"Filtering"})
	cmd.Flags().Bool("synonyms", true, heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/synonyms/`))
	cmd.Flags().SetAnnotation("synonyms", "Categories", []string{"Advanced"})
	tagFilters := NewJSONVar([]string{"array", "string"}...)
	cmd.Flags().Var(tagFilters, "tagFilters", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/tagFilters/`))
	cmd.Flags().SetAnnotation("tagFilters", "Categories", []string{"Filtering"})
	typoTolerance := NewJSONVar([]string{"boolean", "string"}...)
	cmd.Flags().Var(typoTolerance, "typoTolerance", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/typoTolerance/`))
	cmd.Flags().SetAnnotation("typoTolerance", "Categories", []string{"Typos"})
	cmd.Flags().String("userToken", "", heredoc.Doc(`https://www.algolia.com/doc/api-reference/api-parameters/userToken/`))
	cmd.Flags().SetAnnotation("userToken", "Categories", []string{"Personalization"})
}

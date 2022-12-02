// This file is generated; DO NOT EDIT.

import { Flags } from "../../../Terminal/types";

const flags: Flags = {
  advancedSyntax: {
    type: "boolean",
    shortDesc: `Enables the advanced query syntax.`,
    multiple: false,
  },
  advancedSyntaxFeatures: {
    type: "string",
    shortDesc: `Allows you to specify which advanced syntax features are active when ‘advancedSyntax' is enabled.`,
    multiple: true,
  },
  allowTyposOnNumericTokens: {
    type: "boolean",
    shortDesc: `Whether to allow typos on numbers ("numeric tokens") in the query string.`,
    multiple: false,
  },
  alternativesAsExact: {
    type: "string",
    shortDesc: `List of alternatives that should be considered an exact match by the exact ranking criterion.`,
    multiple: true,
  },
  analytics: {
    type: "boolean",
    shortDesc: `Whether the current query will be taken into account in the Analytics.`,
    multiple: false,
  },
  analyticsTags: {
    type: "string",
    shortDesc: `List of tags to apply to the query for analytics purposes.`,
    multiple: true,
  },
  aroundLatLng: {
    type: "string",
    shortDesc: `Search for entries around a central geolocation, enabling a geo search within a circular area.`,
    multiple: false,
  },
  aroundLatLngViaIP: {
    type: "boolean",
    shortDesc: `Search for entries around a given location automatically computed from the requester's IP address.`,
    multiple: false,
  },
  aroundPrecision: {
    type: "number",
    shortDesc: `Precision of geo search (in meters), to add grouping by geo location to the ranking formula.`,
    multiple: false,
  },
  aroundRadius: {
    type: "string",
    shortDesc: `Define the maximum radius for a geo search (in meters).`,
    multiple: false,
  },
  attributeCriteriaComputedByMinProximity: {
    type: "boolean",
    shortDesc: `When attribute is ranked above proximity in your ranking formula, proximity is used to select which searchable attribute is matched in the attribute ranking stage.`,
    multiple: false,
  },
  attributeForDistinct: {
    type: "string",
    shortDesc: `Name of the de-duplication attribute to be used with the distinct feature.`,
    multiple: false,
  },
  attributesForFaceting: {
    type: "string",
    shortDesc: `The complete list of attributes that will be used for faceting.`,
    multiple: true,
  },
  attributesToHighlight: {
    type: "string",
    shortDesc: `List of attributes to highlight.`,
    multiple: true,
  },
  attributesToRetrieve: {
    type: "string",
    shortDesc: `This parameter controls which attributes to retrieve and which not to retrieve.`,
    multiple: true,
  },
  attributesToSnippet: {
    type: "string",
    shortDesc: `List of attributes to snippet, with an optional maximum number of words to snippet.`,
    multiple: true,
  },
  clickAnalytics: {
    type: "boolean",
    shortDesc: `Enable the Click Analytics feature.`,
    multiple: false,
  },
  customRanking: {
    type: "string",
    shortDesc: `Specifies the custom ranking criterion.`,
    multiple: true,
  },
  decompoundQuery: {
    type: "boolean",
    shortDesc: `Splits compound words into their composing atoms in the query.`,
    multiple: false,
  },
  disableExactOnAttributes: {
    type: "string",
    shortDesc: `List of attributes on which you want to disable the exact ranking criterion.`,
    multiple: true,
  },
  disableTypoToleranceOnAttributes: {
    type: "string",
    shortDesc: `List of attributes on which you want to disable typo tolerance.`,
    multiple: true,
  },
  distinct: {
    type: "string",
    shortDesc: `Enables de-duplication or grouping of results.`,
    multiple: false,
  },
  enableABTest: {
    type: "boolean",
    shortDesc: `Whether this search should participate in running AB tests.`,
    multiple: false,
  },
  enablePersonalization: {
    type: "boolean",
    shortDesc: `Enable the Personalization feature.`,
    multiple: false,
  },
  enableReRanking: {
    type: "boolean",
    shortDesc: `Whether this search should use AI Re-Ranking.`,
    multiple: false,
  },
  enableRules: {
    type: "boolean",
    shortDesc: `Whether Rules should be globally enabled.`,
    multiple: false,
  },
  exactOnSingleWordQuery: {
    type: "string",
    shortDesc: `Controls how the exact ranking criterion is computed when the query contains only one word. One of: (attribute, none, word).`,
    multiple: false,
  },
  facetFilters: {
    type: "string",
    shortDesc: `Filter hits by facet value.`,
    multiple: false,
  },
  facetingAfterDistinct: {
    type: "boolean",
    shortDesc: `Force faceting to be applied after de-duplication (via the Distinct setting).`,
    multiple: false,
  },
  facets: {
    type: "string",
    shortDesc: `Retrieve facets and their facet values.`,
    multiple: true,
  },
  filters: {
    type: "string",
    shortDesc: `Filter the query with numeric, facet and/or tag filters.`,
    multiple: false,
  },
  getRankingInfo: {
    type: "boolean",
    shortDesc: `Retrieve detailed ranking information.`,
    multiple: false,
  },
  highlightPostTag: {
    type: "string",
    shortDesc: `The HTML string to insert after the highlighted parts in all highlight and snippet results.`,
    multiple: false,
  },
  highlightPreTag: {
    type: "string",
    shortDesc: `The HTML string to insert before the highlighted parts in all highlight and snippet results.`,
    multiple: false,
  },
  hitsPerPage: {
    type: "number",
    shortDesc: `Set the number of hits per page.`,
    multiple: false,
  },
  ignorePlurals: {
    type: "string",
    shortDesc: `Treats singular, plurals, and other forms of declensions as matching terms.
ignorePlurals is used in conjunction with the queryLanguages setting.
list: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in queryLanguages. true: enables the ignore plurals functionality, where singulars and plurals are considered equivalent (foot = feet). The languages supported here are either every language (this is the default, see list of languages below), or those set by queryLanguages. false: disables ignore plurals, where singulars and plurals are not considered the same for matching purposes (foot will not find feet).
`,
    multiple: false,
  },
  insideBoundingBox: {
    type: "string",
    shortDesc: `Search inside a rectangular area (in geo coordinates).`,
    multiple: true,
  },
  insidePolygon: {
    type: "string",
    shortDesc: `Search inside a polygon (in geo coordinates).`,
    multiple: true,
  },
  keepDiacriticsOnCharacters: {
    type: "string",
    shortDesc: `List of characters that the engine shouldn't automatically normalize.`,
    multiple: false,
  },
  length: {
    type: "number",
    shortDesc: `Set the number of hits to retrieve (used only with offset).`,
    multiple: false,
  },
  maxFacetHits: {
    type: "number",
    shortDesc: `Maximum number of facet hits to return during a search for facet values. For performance reasons, the maximum allowed number of returned values is 100.`,
    multiple: false,
  },
  maxValuesPerFacet: {
    type: "number",
    shortDesc: `Maximum number of facet values to return for each facet during a regular search.`,
    multiple: false,
  },
  minProximity: {
    type: "number",
    shortDesc: `Precision of the proximity ranking criterion.`,
    multiple: false,
  },
  minWordSizefor1Typo: {
    type: "number",
    shortDesc: `Minimum number of characters a word in the query string must contain to accept matches with 1 typo.`,
    multiple: false,
  },
  minWordSizefor2Typos: {
    type: "number",
    shortDesc: `Minimum number of characters a word in the query string must contain to accept matches with 2 typos.`,
    multiple: false,
  },
  minimumAroundRadius: {
    type: "number",
    shortDesc: `Minimum radius (in meters) used for a geo search when aroundRadius is not set.`,
    multiple: false,
  },
  naturalLanguages: {
    type: "string",
    shortDesc: `This parameter changes the default values of certain parameters and settings that work best for a natural language query, such as ignorePlurals, removeStopWords, removeWordsIfNoResults, analyticsTags and ruleContexts. These parameters and settings work well together when the query is formatted in natural language instead of keywords, for example when your user performs a voice search.`,
    multiple: true,
  },
  numericFilters: {
    type: "string",
    shortDesc: `Filter on numeric attributes.`,
    multiple: false,
  },
  offset: {
    type: "number",
    shortDesc: `Specify the offset of the first hit to return.`,
    multiple: false,
  },
  optionalFilters: {
    type: "string",
    shortDesc: `Create filters for ranking purposes, where records that match the filter are ranked higher, or lower in the case of a negative optional filter.`,
    multiple: false,
  },
  optionalWords: {
    type: "string",
    shortDesc: `A list of words that should be considered as optional when found in the query.`,
    multiple: true,
  },
  page: {
    type: "number",
    shortDesc: `Specify the page to retrieve.`,
    multiple: false,
  },
  percentileComputation: {
    type: "boolean",
    shortDesc: `Whether to include or exclude a query from the processing-time percentile computation.`,
    multiple: false,
  },
  personalizationImpact: {
    type: "number",
    shortDesc: `Define the impact of the Personalization feature.`,
    multiple: false,
  },
  query: {
    type: "string",
    shortDesc: `The text to search in the index.`,
    multiple: false,
  },
  queryLanguages: {
    type: "string",
    shortDesc: `Sets the languages to be used by language-specific settings and functionalities such as ignorePlurals, removeStopWords, and CJK word-detection.`,
    multiple: true,
  },
  queryType: {
    type: "string",
    shortDesc: `Controls if and how query words are interpreted as prefixes. One of: (prefixLast, prefixAll, prefixNone).`,
    multiple: false,
  },
  ranking: {
    type: "string",
    shortDesc: `Controls how Algolia should sort your results.`,
    multiple: true,
  },
  reRankingApplyFilter: {
    type: "string",
    shortDesc: `When Dynamic Re-Ranking is enabled, only records that match these filters will be impacted by Dynamic Re-Ranking.`,
    multiple: false,
  },
  relevancyStrictness: {
    type: "number",
    shortDesc: `Controls the relevancy threshold below which less relevant results aren't included in the results.`,
    multiple: false,
  },
  removeStopWords: {
    type: "string",
    shortDesc: `Removes stop (common) words from the query before executing it.
removeStopWords is used in conjunction with the queryLanguages setting.
list: language ISO codes for which ignoring plurals should be enabled. This list will override any values that you may have set in queryLanguages. true: enables the stop word functionality, ensuring that stop words are removed from consideration in a search. The languages supported here are either every language, or those set by queryLanguages. false: disables stop word functionality, allowing stop words to be taken into account in a search.
`,
    multiple: false,
  },
  removeWordsIfNoResults: {
    type: "string",
    shortDesc: `Selects a strategy to remove words from the query when it doesn't match any hits. One of: (none, lastWords, firstWords, allOptional).`,
    multiple: false,
  },
  renderingContent: {
    type: "string",
    shortDesc: `Content defining how the search interface should be rendered. Can be set via the settings for a default value and can be overridden via rules.`,
    multiple: false,
  },
  replaceSynonymsInHighlight: {
    type: "boolean",
    shortDesc: `Whether to highlight and snippet the original word that matches the synonym or the synonym itself.`,
    multiple: false,
  },
  responseFields: {
    type: "string",
    shortDesc: `Choose which fields to return in the API response. This parameters applies to search and browse queries.`,
    multiple: true,
  },
  restrictHighlightAndSnippetArrays: {
    type: "boolean",
    shortDesc: `Restrict highlighting and snippeting to items that matched the query.`,
    multiple: false,
  },
  restrictSearchableAttributes: {
    type: "string",
    shortDesc: `Restricts a given query to look in only a subset of your searchable attributes.`,
    multiple: true,
  },
  ruleContexts: {
    type: "string",
    shortDesc: `Enables contextual rules.`,
    multiple: true,
  },
  similarQuery: {
    type: "string",
    shortDesc: `Overrides the query parameter and performs a more generic search that can be used to find "similar" results.`,
    multiple: false,
  },
  snippetEllipsisText: {
    type: "string",
    shortDesc: `String used as an ellipsis indicator when a snippet is truncated.`,
    multiple: false,
  },
  sortFacetValuesBy: {
    type: "string",
    shortDesc: `Controls how facet values are fetched.`,
    multiple: false,
  },
  sumOrFiltersScores: {
    type: "boolean",
    shortDesc: `Determines how to calculate the total score for filtering.`,
    multiple: false,
  },
  synonyms: {
    type: "boolean",
    shortDesc: `Whether to take into account an index's synonyms for a particular search.`,
    multiple: false,
  },
  tagFilters: {
    type: "string",
    shortDesc: `Filter hits by tags.`,
    multiple: false,
  },
  typoTolerance: {
    type: "string",
    shortDesc: `Controls whether typo tolerance is enabled and how it is applied.`,
    multiple: false,
  },
  userToken: {
    type: "string",
    shortDesc: `Associates a certain user token with the current search.`,
    multiple: false,
  },
};

export default flags;

package shared

import "github.com/algolia/algoliasearch-client-go/v4/algolia/search"

// DictionaryTypes returns the allowed dictionary types as strings
func DictionaryTypes() []string {
	var types []string
	for _, d := range search.AllowedDictionaryTypeEnumValues {
		types = append(types, string(d))
	}
	return types
}

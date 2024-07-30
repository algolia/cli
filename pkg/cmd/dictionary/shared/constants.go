package shared

import (
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
)

// EntryType represents the type of an entry in a dictionary.
// It can be either a custom entry or a standard entry.
type (
	EntryType      string
	DictionaryType int
)

// DictionaryNames returns the list of available dictionaries.
var DictionaryNames = func() []string {
	return []string{
		string(search.DICTIONARY_TYPE_STOPWORDS),
		string(search.DICTIONARY_TYPE_COMPOUNDS),
		string(search.DICTIONARY_TYPE_PLURALS),
	}
}

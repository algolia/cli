package shared

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

// EntryType represents the type of an entry in a dictionnary.
// It can be either a custom entry or a standard entry.
type EntryType string
type DictionaryType int

// DictionaryEntry can be plural, compound or stopword entry.
type DictionaryEntry struct {
	Type          EntryType
	Word          string   `json:"word,omitempty"`
	Words         []string `json:"words,omitempty"`
	Decomposition []string `json:"decomposition,omitempty"`
	ObjectID      string
	Language      string
	State         string
}

const (
	// CustomEntryType is the type of a custom entry in a dictionnary (i.e. added by the user).
	CustomEntryType EntryType = "custom"
	// StandardEntryType is the type of a standard entry in a dictionnary (i.e. added by Algolia).
	StandardEntryType EntryType = "standard"
)

var (
	// DictionaryNames returns the list of available dictionnaries.
	DictionaryNames = func() []string {
		return []string{
			string(search.Stopwords),
			string(search.Compounds),
			string(search.Plurals),
		}
	}
)

package config

import (
	"fmt"
	"io"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func GetSynonyms(srcIndex *search.Index) ([]search.Synonym, error) {
	it, err := srcIndex.BrowseSynonyms()
	if err != nil {
		return nil, fmt.Errorf("cannot browse source index synonyms: %v", err)
	}

	var synonyms []search.Synonym

	for {
		synonym, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, fmt.Errorf("error while iterating source index synonyms: %v", err)
			}
		}
		synonyms = append(synonyms, synonym)
	}

	return synonyms, nil
}

func GetRules(srcIndex *search.Index) ([]search.Rule, error) {
	it, err := srcIndex.BrowseRules()
	if err != nil {
		return nil, fmt.Errorf("cannot browse source index rules: %v", err)
	}

	var rules []search.Rule

	for {
		rule, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, fmt.Errorf("error while iterating source index rules: %v", err)
			}
		}
		rules = append(rules, *rule)
	}

	return rules, nil
}

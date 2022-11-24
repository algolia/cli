package indeximport

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	config "github.com/algolia/cli/pkg/cmd/shared/handler/indices"
	"github.com/algolia/cli/pkg/cmd/synonyms/shared"
)

func SynonymsToSearchSynonyms(synonyms []config.Synonym) ([]search.Synonym, error) {
	var searchSynonyms []search.Synonym
	for _, synonym := range synonyms {
		searchSynonym, err := synonymToSearchSynonm(synonym)
		if err != nil {
			return nil, err
		}
		searchSynonyms = append(searchSynonyms, searchSynonym)
	}

	return searchSynonyms, nil
}

func synonymToSearchSynonm(synonym config.Synonym) (search.Synonym, error) {
	switch synonym.Type {
	case shared.OneWay:
		return search.NewOneWaySynonym(
			synonym.ObjectID,
			synonym.Input,
			synonym.Synonyms...,
		), nil
	case shared.AltCorrection1:
		return search.NewAltCorrection1(
			synonym.ObjectID,
			synonym.Word,
			synonym.Corrections...,
		), nil
	case shared.AltCorrection2:
		return search.NewAltCorrection2(
			synonym.ObjectID,
			synonym.Word,
			synonym.Corrections...,
		), nil
	case shared.Placeholder:
		return search.NewPlaceholder(
			synonym.ObjectID,
			synonym.Placeholder,
			synonym.Replacements...,
		), nil
	case "", shared.Regular:
		return search.NewRegularSynonym(
			synonym.ObjectID,
			synonym.Synonyms...,
		), nil
	}

	return nil, fmt.Errorf("invalid synonym type for object id %s", synonym.ObjectID)
}

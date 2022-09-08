package shared

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type SynonymFlags struct {
	SynonymID           string
	SynonymInput        string
	SynonymWord         string
	SynonymPlaceholder  string
	SynonymType         string
	Synonyms            []string
	SynonymCorrections  []string
	SynonymReplacements []string
}

// Defining new type that implements pflag.Value interface with String, Set and Type
// https://stackoverflow.com/questions/50824554/permitted-flag-values-for-cobra
type SynonymType string

const (
	// Matching API https://www.algolia.com/doc/api-reference/api-methods/save-synonym/#method-param-type
	Regular        string = "synonym"
	OneWay         string = "oneWaySynonym"
	AltCorrection1 string = "altCorrection1"
	AltCorrection2 string = "altCorrection2"
	Placeholder    string = "placeholder"
)

func (e *SynonymType) String() string {
	return string(*e)
}

func (e *SynonymType) Set(v string) error {
	if v == "" {
		*e = SynonymType(v)
		return nil
	}

	switch v {
	case Regular, OneWay, AltCorrection1, AltCorrection2, Placeholder:
		*e = SynonymType(v)
		return nil
	default:
		return fmt.Errorf(`must be one of "regular", "one-way", "alt-correction1", "alt-correction2" or "placeholder"`)
	}
}

func (e *SynonymType) Type() string {
	return "SynonymType"
}

func FlagsToSynonym(flags SynonymFlags) (search.Synonym, error) {
	if flags.SynonymID == "" {
		return nil, fmt.Errorf("a unique synonym id is required")
	}

	// Default case
	if flags.SynonymType == "" || flags.SynonymType == Regular {
		if len(flags.Synonyms) < 1 {
			return nil, fmt.Errorf("at least 1 synonym is required")
		}
		return search.NewRegularSynonym(
			flags.SynonymID,
			flags.Synonyms...,
		), nil
	}

	switch flags.SynonymType {
	case OneWay:
		if len(flags.Synonyms) < 1 {
			return nil, fmt.Errorf("at least 1 synonym is required")
		}
		if flags.SynonymInput == "" {
			return nil, fmt.Errorf("a synonym input is required for one way synonyms")
		}
		return search.NewOneWaySynonym(
			flags.SynonymID,
			flags.SynonymInput,
			flags.Synonyms...,
		), nil
	case AltCorrection1:
		if flags.SynonymWord == "" {
			return nil, fmt.Errorf("synonym word is required for alt correction 1 synonyms")
		}
		if len(flags.SynonymCorrections) < 1 {
			return nil, fmt.Errorf("synonym corrections are required for alt correction 1 synonyms")
		}
		return search.NewAltCorrection1(
			flags.SynonymID,
			flags.SynonymWord,
			flags.SynonymCorrections...,
		), nil
	case AltCorrection2:
		if flags.SynonymWord == "" {
			return nil, fmt.Errorf("synonym word is required for alt correction 2 synonyms")
		}
		if len(flags.SynonymCorrections) < 1 {
			return nil, fmt.Errorf("synonym corrections are required for alt correction 2 synonyms")
		}
		return search.NewAltCorrection2(
			flags.SynonymID,
			flags.SynonymWord,
			flags.SynonymCorrections...,
		), nil
	case Placeholder:
		if flags.SynonymPlaceholder == "" {
			return nil, fmt.Errorf("a synonym placeholder is required for placeholder synonyms")
		}
		if len(flags.SynonymReplacements) < 1 {
			return nil, fmt.Errorf("synonym replacements are required for placeholder synonyms")
		}
		return search.NewPlaceholder(
			flags.SynonymID,
			flags.SynonymPlaceholder,
			flags.SynonymReplacements...,
		), nil
	case Regular:
		return search.NewRegularSynonym(
			flags.SynonymID,
			flags.Synonyms...,
		), nil
	}

	return nil, fmt.Errorf("invalid synonym type")
}

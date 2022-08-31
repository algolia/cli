package validator

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type SaveOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice              string
	SynonymID           string
	ForwardToReplicas   bool
	SynonymType         SynonymType
	SynonymInput        string
	SynonymWord         string
	SynonymPlaceholder  string
	Synonyms            []string
	SynonymCorrections  []string
	SynonymReplacements []string
	Synonym             search.Synonym
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

func ValidateFlags(options SaveOptions) (search.Synonym, error) {
	if options.SynonymID == "" {
		return nil, fmt.Errorf("a unique synonym id is required")
	}

	// Default case
	if options.SynonymType == "" || options.SynonymType == SynonymType(Regular) {
		return search.NewRegularSynonym(
			options.SynonymID,
			options.Synonyms...,
		), nil
	}

	switch options.SynonymType {
	case SynonymType(OneWay):
		if options.SynonymInput == "" {
			return nil, fmt.Errorf("a synonym input is required for one way synonyms")
		}
		return search.NewOneWaySynonym(
			options.SynonymID,
			options.SynonymInput,
			options.Synonyms...,
		), nil
	case SynonymType(AltCorrection1):
		if options.SynonymWord == "" {
			return nil, fmt.Errorf("synonym word is required for alt correction 1 synonyms")
		}
		if len(options.SynonymCorrections) < 1 {
			return nil, fmt.Errorf("synonym corrections are required for alt correction 1 synonyms")
		}
		return search.NewAltCorrection1(
			options.SynonymID,
			options.SynonymWord,
			options.SynonymCorrections...,
		), nil
	case SynonymType(AltCorrection2):
		if options.SynonymWord == "" {
			return nil, fmt.Errorf("synonym word is required for alt correction 2 synonyms")
		}
		if len(options.SynonymCorrections) < 1 {
			return nil, fmt.Errorf("synonym corrections are required for alt correction 2 synonyms")
		}
		return search.NewAltCorrection2(
			options.SynonymID,
			options.SynonymWord,
			options.SynonymCorrections...,
		), nil
	case SynonymType(Placeholder):
		if options.SynonymPlaceholder == "" {
			return nil, fmt.Errorf("a synonym placeholder is required for placeholder synonyms")
		}
		if len(options.SynonymReplacements) < 1 {
			return nil, fmt.Errorf("synonym replacements are required for placeholder synonyms")
		}
		return search.NewPlaceholder(
			options.SynonymID,
			options.SynonymPlaceholder,
			options.SynonymReplacements...,
		), nil
	case SynonymType(Regular):
		return search.NewRegularSynonym(
			options.SynonymID,
			options.Synonyms...,
		), nil
	}

	return nil, fmt.Errorf("invalid synonym type")
}

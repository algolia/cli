package shared

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
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
	Regular string = string(search.SYNONYM_TYPE_SYNONYM)
	// oneWaySynonym
	OneWay string = string(search.SYNONYM_TYPE_ONE_WAY_SYNONYM)
	// onewaysynonym
	AltOneWay string = string(search.SYNONYM_TYPE_ONEWAYSYNONYM)
	// altCorrection1
	AltCorrection1 string = string(search.SYNONYM_TYPE_ALT_CORRECTION1)
	// altcorrection1
	AltAltCorrection1 string = string(search.SYNONYM_TYPE_ALTCORRECTION1)
	// altCorrection2
	AltCorrection2 string = string(search.SYNONYM_TYPE_ALT_CORRECTION2)
	// altcorrection2
	AltAltCorrection2 string = string(search.SYNONYM_TYPE_ALTCORRECTION2)
	Placeholder       string = string(search.SYNONYM_TYPE_PLACEHOLDER)
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
	case Regular,
		OneWay,
		AltCorrection1,
		AltCorrection2,
		Placeholder,
		AltOneWay,
		AltAltCorrection1,
		AltAltCorrection2:
		*e = SynonymType(v)
		return nil
	default:
		return fmt.Errorf(
			`must be one of "regular", "one-way", "alt-correction1", "alt-correction2" or "placeholder"`,
		)
	}
}

func (e *SynonymType) Type() string {
	return "SynonymType"
}

func FlagsToSynonym(flags SynonymFlags) (*search.SynonymHit, error) {
	switch flags.SynonymType {
	case OneWay, AltOneWay:
		return search.NewEmptySynonymHit().
				SetType(search.SYNONYM_TYPE_ONE_WAY_SYNONYM).
				SetObjectID(flags.SynonymID).
				SetInput(flags.SynonymInput).
				SetSynonyms(flags.Synonyms),
			nil
	case AltCorrection1, AltAltCorrection1:
		return search.NewEmptySynonymHit().
				SetType(search.SYNONYM_TYPE_ALT_CORRECTION1).
				SetObjectID(flags.SynonymID).
				SetWord(flags.SynonymWord).
				SetCorrections(flags.SynonymCorrections),
			nil
	case AltCorrection2, AltAltCorrection2:
		return search.NewEmptySynonymHit().
				SetType(search.SYNONYM_TYPE_ALT_CORRECTION2).
				SetObjectID(flags.SynonymID).
				SetWord(flags.SynonymWord).
				SetCorrections(flags.SynonymCorrections),
			nil
	case Placeholder:
		return search.NewEmptySynonymHit().
				SetType(search.SYNONYM_TYPE_PLACEHOLDER).
				SetObjectID(flags.SynonymID).
				SetPlaceholder(flags.SynonymPlaceholder).
				SetReplacements(flags.SynonymReplacements),
			nil
	case "", Regular:
		return search.NewEmptySynonymHit().
				SetType(search.SYNONYM_TYPE_SYNONYM).
				SetObjectID(flags.SynonymID).
				SetSynonyms(flags.Synonyms),
			nil
	}

	return nil, fmt.Errorf("invalid synonym type")
}

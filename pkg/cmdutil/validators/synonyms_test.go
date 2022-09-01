package validator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name        string
		saveOptions SaveOptions
		synonymType string
		wantsErr    bool
	}{
		{
			name:     "Regular synonym",
			wantsErr: false,
			saveOptions: SaveOptions{
				SynonymID: "23",
				Synonyms:  []string{"mj", "goat"}},
			synonymType: "search.RegularSynonym",
		},
		{
			name:     "Regular synonym explicit type",
			wantsErr: false,
			saveOptions: SaveOptions{
				SynonymType: SynonymType(Regular),
				SynonymID:   "23",
				Synonyms:    []string{"mj", "goat"}},
			synonymType: "search.RegularSynonym",
		},
		{
			name:     "One way synonym",
			wantsErr: false,
			saveOptions: SaveOptions{
				SynonymType:  SynonymType(OneWay),
				SynonymID:    "23",
				Synonyms:     []string{"mj", "goat"},
				SynonymInput: "michael",
			},
			synonymType: "search.OneWaySynonym",
		},
		{
			name:     "One way synonym without input",
			wantsErr: true,
			saveOptions: SaveOptions{
				SynonymType: SynonymType(OneWay),
				SynonymID:   "23",
				Synonyms:    []string{"mj", "goat"},
			},
		},
		{
			name:     "AltCorrection1 synonym",
			wantsErr: false,
			saveOptions: SaveOptions{
				SynonymType:        SynonymType(AltCorrection1),
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
				SynonymWord:        "michael",
			},
			synonymType: "search.AltCorrection1",
		},
		{
			name:     "AltCorrection1 synonym without word",
			wantsErr: true,
			saveOptions: SaveOptions{
				SynonymType:        SynonymType(AltCorrection1),
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
			},
		},
		{
			name:     "AltCorrection2 synonym",
			wantsErr: false,
			saveOptions: SaveOptions{
				SynonymType:        SynonymType(AltCorrection2),
				SynonymID:          "24",
				SynonymCorrections: []string{"bryant", "mamba"},
				SynonymWord:        "kobe",
			},
			synonymType: "search.AltCorrection2",
		},
		{
			name:     "AltCorrection2 synonym without correction",
			wantsErr: true,
			saveOptions: SaveOptions{
				SynonymType: SynonymType(AltCorrection2),
				SynonymID:   "24",
				SynonymWord: "kobe",
			},
		},
		{
			name:     "Placeholder synonym",
			wantsErr: false,
			saveOptions: SaveOptions{
				SynonymType:         SynonymType(Placeholder),
				SynonymID:           "23",
				SynonymReplacements: []string{"james", "lebron"},
				SynonymPlaceholder:  "king",
			},
			synonymType: "search.Placeholder",
		},
		{
			name:     "Placeholder synonym without placeholder",
			wantsErr: true,
			saveOptions: SaveOptions{
				SynonymType:         SynonymType(Placeholder),
				SynonymID:           "23",
				SynonymReplacements: []string{"james", "lebron"},
			},
		},
		{
			name:     "Placeholder synonym without replacements",
			wantsErr: true,
			saveOptions: SaveOptions{
				SynonymType:        SynonymType(Placeholder),
				SynonymID:          "23",
				SynonymPlaceholder: "king",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synonym, err := ValidateFlags(tt.saveOptions)

			if tt.wantsErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, reflect.TypeOf(synonym).String(), tt.synonymType)
		})
	}
}

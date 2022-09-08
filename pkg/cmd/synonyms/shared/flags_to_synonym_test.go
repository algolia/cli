package shared

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FlagsToSynonym(t *testing.T) {
	tests := []struct {
		name         string
		synonymFlags SynonymFlags
		synonymType  string
		wantsErr     bool
		wantsErrMsg  string
	}{
		// Regular type
		{
			name:     "Regular synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymID: "23",
				Synonyms:  []string{"mj", "goat"}},
			synonymType: "search.RegularSynonym",
		},
		{
			name:     "Regular synonym explicit type",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType: Regular,
				SynonymID:   "23",
				Synonyms:    []string{"mj", "goat"}},
			synonymType: "search.RegularSynonym",
		},
		{
			name:        "Regular synonym without id",
			wantsErr:    true,
			wantsErrMsg: "a unique synonym id is required",
			synonymFlags: SynonymFlags{
				Synonyms: []string{"mj", "goat"}},
		},
		// One way type
		{
			name:     "One way synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType:  OneWay,
				SynonymID:    "23",
				Synonyms:     []string{"mj", "goat"},
				SynonymInput: "michael",
			},
			synonymType: "search.OneWaySynonym",
		},
		{
			name:        "One way synonym without input",
			wantsErr:    true,
			wantsErrMsg: "a synonym input is required for one way synonyms",
			synonymFlags: SynonymFlags{
				SynonymType: OneWay,
				SynonymID:   "23",
				Synonyms:    []string{"mj", "goat"},
			},
		},
		// Alt correction type
		{
			name:     "AltCorrection1 synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType:        AltCorrection1,
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
				SynonymWord:        "michael",
			},
			synonymType: "search.AltCorrection1",
		},
		{
			name:        "AltCorrection1 synonym without word",
			wantsErr:    true,
			wantsErrMsg: "synonym word is required for alt correction 1 synonyms",
			synonymFlags: SynonymFlags{
				SynonymType:        AltCorrection1,
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
			},
		},
		{
			name:     "AltCorrection2 synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType:        AltCorrection2,
				SynonymID:          "24",
				SynonymCorrections: []string{"bryant", "mamba"},
				SynonymWord:        "kobe",
			},
			synonymType: "search.AltCorrection2",
		},
		{
			name:        "AltCorrection2 synonym without correction",
			wantsErr:    true,
			wantsErrMsg: "synonym corrections are required for alt correction 2 synonyms",
			synonymFlags: SynonymFlags{
				SynonymType: AltCorrection2,
				SynonymID:   "24",
				SynonymWord: "kobe",
			},
		},
		// Placeholder type
		{
			name:     "Placeholder synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType:         Placeholder,
				SynonymID:           "23",
				SynonymReplacements: []string{"james", "lebron"},
				SynonymPlaceholder:  "king",
			},
			synonymType: "search.Placeholder",
		},
		{
			name:        "Placeholder synonym without placeholder",
			wantsErr:    true,
			wantsErrMsg: "a synonym placeholder is required for placeholder synonyms",
			synonymFlags: SynonymFlags{
				SynonymType:         Placeholder,
				SynonymID:           "23",
				SynonymReplacements: []string{"james", "lebron"},
			},
		},
		{
			name:        "Placeholder synonym without replacements",
			wantsErr:    true,
			wantsErrMsg: "synonym replacements are required for placeholder synonyms",
			synonymFlags: SynonymFlags{
				SynonymType:        Placeholder,
				SynonymID:          "23",
				SynonymPlaceholder: "king",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synonym, err := FlagsToSynonym(tt.synonymFlags)

			if tt.wantsErr {
				assert.EqualError(t, err, tt.wantsErrMsg)
				return
			}

			assert.Equal(t, reflect.TypeOf(synonym).String(), tt.synonymType)
		})
	}
}

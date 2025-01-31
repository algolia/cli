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
				Synonyms:  []string{"mj", "goat"},
			},
			synonymType: "search.RegularSynonym",
		},
		{
			name:     "Regular synonym explicit type",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType: Regular,
				SynonymID:   "23",
				Synonyms:    []string{"mj", "goat"},
			},
			synonymType: "search.RegularSynonym",
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
		// Wrong type
		{
			name:        "Wrong synonym type",
			wantsErr:    true,
			wantsErrMsg: "invalid synonym type",
			synonymFlags: SynonymFlags{
				SynonymType:         "wrongType",
				SynonymID:           "23",
				SynonymReplacements: []string{"james", "lebron"},
				SynonymPlaceholder:  "king",
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

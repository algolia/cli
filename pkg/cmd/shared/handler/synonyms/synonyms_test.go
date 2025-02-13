package synonms

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/cmd/synonyms/shared"
)

func Test_ValidateSynonymFlags(t *testing.T) {
	tests := []struct {
		name         string
		synonymFlags shared.SynonymFlags
		synonymType  string
		wantsErr     bool
		wantsErrMsg  string
	}{
		// Regular type
		{
			name:     "Regular synonym",
			wantsErr: false,
			synonymFlags: shared.SynonymFlags{
				SynonymID: "23",
				Synonyms:  []string{"mj", "goat"},
			},
		},
		{
			name:     "Regular synonym explicit type",
			wantsErr: false,
			synonymFlags: shared.SynonymFlags{
				SynonymType: shared.Regular,
				SynonymID:   "23",
				Synonyms:    []string{"mj", "goat"},
			},
		},
		{
			name:        "Regular synonym without id",
			wantsErr:    true,
			wantsErrMsg: "a unique synonym id is required",
			synonymFlags: shared.SynonymFlags{
				Synonyms: []string{"mj", "goat"},
			},
		},
		// One way type
		{
			name:     "One way synonym",
			wantsErr: false,
			synonymFlags: shared.SynonymFlags{
				SynonymType:  shared.OneWay,
				SynonymID:    "23",
				Synonyms:     []string{"mj", "goat"},
				SynonymInput: "michael",
			},
		},
		{
			name:        "One way synonym without input",
			wantsErr:    true,
			wantsErrMsg: "a synonym input is required for one way synonyms",
			synonymFlags: shared.SynonymFlags{
				SynonymType: shared.OneWay,
				SynonymID:   "23",
				Synonyms:    []string{"mj", "goat"},
			},
		},
		// Alt correction type
		{
			name:     "AltCorrection1 synonym",
			wantsErr: false,
			synonymFlags: shared.SynonymFlags{
				SynonymType:        shared.AltCorrection1,
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
				SynonymWord:        "michael",
			},
		},
		{
			name:        "AltCorrection1 synonym without word",
			wantsErr:    true,
			wantsErrMsg: "synonym word is required for alt correction 1 synonyms",
			synonymFlags: shared.SynonymFlags{
				SynonymType:        shared.AltCorrection1,
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
			},
		},
		{
			name:     "AltCorrection2 synonym",
			wantsErr: false,
			synonymFlags: shared.SynonymFlags{
				SynonymType:        shared.AltCorrection2,
				SynonymID:          "24",
				SynonymCorrections: []string{"bryant", "mamba"},
				SynonymWord:        "kobe",
			},
		},
		{
			name:        "AltCorrection2 synonym without correction",
			wantsErr:    true,
			wantsErrMsg: "synonym corrections are required for alt correction 2 synonyms",
			synonymFlags: shared.SynonymFlags{
				SynonymType: shared.AltCorrection2,
				SynonymID:   "24",
				SynonymWord: "kobe",
			},
		},
		// Placeholder type
		{
			name:     "Placeholder synonym",
			wantsErr: false,
			synonymFlags: shared.SynonymFlags{
				SynonymType:         shared.Placeholder,
				SynonymID:           "23",
				SynonymReplacements: []string{"james", "lebron"},
				SynonymPlaceholder:  "king",
			},
		},
		{
			name:        "Placeholder synonym without placeholder",
			wantsErr:    true,
			wantsErrMsg: "a synonym placeholder is required for placeholder synonyms",
			synonymFlags: shared.SynonymFlags{
				SynonymType:         shared.Placeholder,
				SynonymID:           "23",
				SynonymReplacements: []string{"james", "lebron"},
			},
		},
		{
			name:        "Placeholder synonym without replacements",
			wantsErr:    true,
			wantsErrMsg: "synonym replacements are required for placeholder synonyms",
			synonymFlags: shared.SynonymFlags{
				SynonymType:        shared.Placeholder,
				SynonymID:          "23",
				SynonymPlaceholder: "king",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSynonymFlags(tt.synonymFlags)

			if tt.wantsErr {
				assert.EqualError(t, err, tt.wantsErrMsg)
				return
			}
		})
	}
}

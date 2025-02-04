package shared

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"
)

func Test_FlagsToSynonym(t *testing.T) {
	tests := []struct {
		name         string
		synonymFlags SynonymFlags
		synonymType  search.SynonymType
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
			synonymType: search.SYNONYM_TYPE_SYNONYM,
		},
		{
			name:     "Regular synonym explicit type",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType: "synonym",
				SynonymID:   "23",
				Synonyms:    []string{"mj", "goat"},
			},
			synonymType: search.SYNONYM_TYPE_SYNONYM,
		},
		// One way type
		{
			name:     "One way synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType:  "oneWaySynonym",
				SynonymID:    "23",
				Synonyms:     []string{"mj", "goat"},
				SynonymInput: "michael",
			},
			synonymType: search.SYNONYM_TYPE_ONE_WAY_SYNONYM,
		},
		// Alt correction type
		{
			name:     "AltCorrection1 synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType:        "altCorrection1",
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
				SynonymWord:        "michael",
			},
			synonymType: search.SYNONYM_TYPE_ALT_CORRECTION1,
		},
		{
			name:     "AltCorrection2 synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType:        "altCorrection2",
				SynonymID:          "24",
				SynonymCorrections: []string{"bryant", "mamba"},
				SynonymWord:        "kobe",
			},
			synonymType: search.SYNONYM_TYPE_ALT_CORRECTION2,
		},
		// Placeholder type
		{
			name:     "Placeholder synonym",
			wantsErr: false,
			synonymFlags: SynonymFlags{
				SynonymType:         string(search.SYNONYM_TYPE_PLACEHOLDER),
				SynonymID:           "23",
				SynonymReplacements: []string{"james", "lebron"},
				SynonymPlaceholder:  "king",
			},
			synonymType: search.SYNONYM_TYPE_PLACEHOLDER,
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

			assert.Equal(t, synonym.Type, tt.synonymType)
		})
	}
}

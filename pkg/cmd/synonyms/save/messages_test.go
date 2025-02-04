package save

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	shared "github.com/algolia/cli/pkg/cmd/synonyms/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

func Test_GetSynonymSuccessMessage(t *testing.T) {
	tests := []struct {
		name         string
		synonymFlags shared.SynonymFlags
		saveOptions  SaveOptions
		wantsOutput  string
		saveWording  string
	}{
		{
			name: "Save regular synonym",
			synonymFlags: shared.SynonymFlags{
				SynonymID: "23",
				Synonyms:  []string{"mj", "goat"},
			},
			saveOptions: SaveOptions{
				Index: "legends",
			},
			wantsOutput: "✓ Synonym '23' successfully saved with 2 synonyms (mj, goat) to legends\n",
		},
		{
			name: "Save one way synonym",
			synonymFlags: shared.SynonymFlags{
				SynonymType:  shared.OneWay,
				SynonymID:    "23",
				Synonyms:     []string{"mj", "goat"},
				SynonymInput: "michael",
			},
			saveOptions: SaveOptions{
				Index: "legends",
			},
			wantsOutput: "✓ One way synonym '23' successfully saved with input 'michael' and 2 synonyms (mj, goat) to legends\n",
		},
		{
			name: "Save placeholder synonym",
			synonymFlags: shared.SynonymFlags{
				SynonymType:         shared.Placeholder,
				SynonymID:           "23",
				SynonymReplacements: []string{"mj", "goat"},
				SynonymPlaceholder:  "michael",
			},
			saveOptions: SaveOptions{
				Index: "legends",
			},
			wantsOutput: "✓ Placeholder synonym '23' successfully saved with placeholder 'michael' and 2 replacements (mj, goat) to legends\n",
		},
		{
			name: "Save alt correction 1 synonym",
			synonymFlags: shared.SynonymFlags{
				SynonymType:        shared.AltCorrection1,
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
				SynonymWord:        "michael",
			},
			saveOptions: SaveOptions{
				Index: "legends",
			},
			wantsOutput: "✓ Alt correction 1 synonym '23' successfully saved with word 'michael' and 2 corrections (mj, goat) to legends\n",
		},
		{
			name: "Save alt correction 2 synonym",
			synonymFlags: shared.SynonymFlags{
				SynonymType:        shared.AltCorrection2,
				SynonymID:          "23",
				SynonymCorrections: []string{"mj", "goat"},
				SynonymWord:        "michael",
			},
			saveOptions: SaveOptions{
				Index: "legends",
			},
			wantsOutput: "✓ Alt correction 2 synonym '23' successfully saved with word 'michael' and 2 corrections (mj, goat) to legends\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			f := &cmdutil.Factory{
				IOStreams: io,
			}

			err, message := GetSuccessMessage(tt.synonymFlags, tt.saveOptions.Index)

			assert.Equal(t, err, nil)
			assert.Equal(
				t,
				tt.wantsOutput,
				fmt.Sprintf("%s %s", f.IOStreams.ColorScheme().SuccessIcon(), message),
			)
		})
	}
}

package indiceimport

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/iostreams"
)

func TestGetConfirmMessage(t *testing.T) {
	tests := []struct {
		name string

		cs                    *iostreams.ColorScheme
		scope                 []string
		clearExistingRules    bool
		clearExistingSynonyms bool

		wantsOutput string
	}{
		{
			name:                  "full scope",
			scope:                 []string{"settings", "rules", "synonyms"},
			clearExistingRules:    false,
			clearExistingSynonyms: false,
			wantsOutput:           "! Your settings will be CLEARED and REPLACED.\n! Your rules and synonyms will be UPDATED\n",
		},
		{
			name:                  "full scope, --clearExistingSynonyms",
			scope:                 []string{"settings", "rules", "synonyms"},
			clearExistingRules:    false,
			clearExistingSynonyms: true,
			wantsOutput:           "! Your settings and synonyms will be CLEARED and REPLACED.\n! Your rules will be UPDATED\n",
		},
		{
			name:                  "full scope, --clearExistingRules --clearExistingSynonyms",
			scope:                 []string{"settings", "rules", "synonyms"},
			clearExistingRules:    true,
			clearExistingSynonyms: true,
			wantsOutput:           "! Your settings, rules and synonyms will be CLEARED and REPLACED.\n",
		},
		{
			name:                  "rules and synonyms scope",
			scope:                 []string{"rules", "synonyms"},
			clearExistingRules:    false,
			clearExistingSynonyms: false,
			wantsOutput:           "! Your rules and synonyms will be UPDATED\n",
		},
		{
			name:                  "rules and synonyms scope --clearExistingSynonyms",
			scope:                 []string{"rules", "synonyms"},
			clearExistingRules:    false,
			clearExistingSynonyms: true,
			wantsOutput:           "! Your synonyms will be CLEARED and REPLACED.\n! Your rules will be UPDATED\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			tt.cs = io.ColorScheme()

			assert.Equal(t, tt.wantsOutput, GetConfirmMessage(tt.cs, tt.scope, tt.clearExistingRules, tt.clearExistingSynonyms))
		})
	}
}

package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_runStringSliceCompletion(t *testing.T) {
	allowedMap := map[string]string{
		"settings": "settings",
		"synonyms": "synonyms",
		"rules":    "rules",
	}
	prefixDescription := "copy only"

	tests := []struct {
		name       string
		toComplete string
		results    []string
	}{
		{
			name:       "first input, no letter",
			toComplete: "",
			results:    []string{"rules\tcopy only rules", "settings\tcopy only settings", "synonyms\tcopy only synonyms"},
		},
		{
			name:       "second input (settings already passed), no letter",
			toComplete: "settings,s",
			results:    []string{"settings,synonyms\tcopy only settings and synonyms"},
		},
		{
			name:       "first input, first letter",
			toComplete: "s",
			results:    []string{"settings\tcopy only settings", "synonyms\tcopy only synonyms"},
		},
		{
			name:       "second input (settings already passed), first letter",
			toComplete: "settings,s",
			results:    []string{"settings,synonyms\tcopy only settings and synonyms"},
		},
		{
			name:       "third input (settings and synonyms already passed), no letter",
			toComplete: "settings,synonyms,",
			results:    []string{"settings,synonyms,rules\tcopy only settings, synonyms and rules"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, rule := runStringSliceCompletion(
				allowedMap,
				tt.toComplete,
				prefixDescription,
			)

			assert.Equal(t, tt.results, results)
			assert.Equal(t, cobra.ShellCompDirectiveNoSpace, rule)
		})

	}
}

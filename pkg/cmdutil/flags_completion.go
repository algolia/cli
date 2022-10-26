package cmdutil

import (
	"fmt"
	"strings"

	"github.com/algolia/cli/pkg/utils"
	"github.com/spf13/cobra"
)

func ConfiguredProfilesCompletionFunc(f *Factory) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		profiles := f.Config.ConfiguredProfiles()
		completions := make([]string, 0, len(profiles))

		// We want to show the profile name and the Application ID as the description.
		// https://github.com/spf13/cobra/blob/master/shell_completions.md#descriptions-for-completions
		for _, profile := range profiles {
			completions = append(completions, fmt.Sprintf("%s\t%s", profile.Name, profile.ApplicationID))
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}
}

func StringCompletionFunc(allowedMap map[string]string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		allowedValues := make([]string, 0, len(allowedMap))
		for name, description := range allowedMap {
			allowedValues = append(allowedValues, fmt.Sprintf("%s\t%s", name, description))
		}
		return allowedValues, cobra.ShellCompDirectiveNoSpace
	}
}

// Inspired from https://github.com/cli/cli/blob/trunk/pkg/cmdutil/json_flags.go#L26
func StringSliceCompletionFunc(allowedMap map[string]string, prefixAllDescription string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var results []string
		var prefix string

		if idx := strings.LastIndexByte(toComplete, ','); idx >= 0 {
			prefix = toComplete[:idx+1]
			toComplete = toComplete[idx+1:]
		}
		toComplete = strings.ToLower(toComplete)

		for name, description := range allowedMap {
			prefixSlice := utils.StringToSlice(prefix)

			// Build dynamic description with previous selected values
			dynamicSliceDescriptions := []string{}
			for _, prefixName := range prefixSlice {
				prefixDescription := allowedMap[prefixName]
				if prefixDescription != "" {
					dynamicSliceDescriptions = append(dynamicSliceDescriptions, prefixDescription)
				}
			}
			// If current value isn't already selected and if prefix matches
			if !utils.Contains(prefixSlice, name) && strings.HasPrefix(strings.ToLower(name), toComplete) {
				// Add description of current value
				dynamicSliceDescriptions = append(dynamicSliceDescriptions, description)
				results = append(results, fmt.Sprintf("%s%s\t%s",
					prefix,
					name,
					fmt.Sprintf("%s %s", prefixAllDescription, utils.SliceToReadableString(dynamicSliceDescriptions))))
			}
		}

		return results, cobra.ShellCompDirectiveNoSpace
	}
}

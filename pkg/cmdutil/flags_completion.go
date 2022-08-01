package cmdutil

import (
	"fmt"

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

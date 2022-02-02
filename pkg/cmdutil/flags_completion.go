package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/config"
)

func ConfiguredApplicationsCompletionFunc(cfg *config.Config) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		applications := cfg.Applications()
		completions := make([]string, 0, len(applications))

		// We want to show the application name and the ID as the description.
		// https://github.com/spf13/cobra/blob/master/shell_completions.md#descriptions-for-completions
		for profileName, AppID := range applications {
			completions = append(completions, fmt.Sprintf("%s\t%s", profileName, AppID))
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}
}

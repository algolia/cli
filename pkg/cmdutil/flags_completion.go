package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ConfiguredApplicationsCompletionFunc(f *Factory) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		applications := f.Config.ConfiguredApplications()
		completions := make([]string, 0, len(applications))

		// We want to show the application name and the ID as the description.
		// https://github.com/spf13/cobra/blob/master/shell_completions.md#descriptions-for-completions
		for _, application := range applications {
			completions = append(completions, fmt.Sprintf("%s\t%s", application.Name, application.ID))
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}
}

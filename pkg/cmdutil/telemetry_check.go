package cmdutil

import "github.com/spf13/cobra"

func ShouldTrackUsage(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "help", "powershell", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
		return false
	}

	if cmd.Parent() != nil && cmd.Parent().Name() == "completion" {
		return false
	}

	return true
}

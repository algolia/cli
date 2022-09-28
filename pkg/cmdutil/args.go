package cmdutil

import (
	"github.com/spf13/cobra"
)

func ExactArgs(n int, msg string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > n {
			return FlagErrorf("too many arguments")
		}

		if len(args) < n {
			return FlagErrorf("%s", msg)
		}

		return nil
	}
}

package validators

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/algolia/cli/pkg/cmdutil"
)

// ExactArgs is a validator for commands to print an error with a custom message
// followed by usage, flags and available commands when too few/much arguments are provided
func ExactArgs(n int, msg string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return cmdutil.FlagErrorf(msg)
		}

		return nil
	}
}

// NoArgs is a validator for commands to print an error when an argument is provided
func NoArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		extractArgs := ExactArgs(0, fmt.Sprintf(
			"`%s` does not take any positional arguments.",
			cmd.CommandPath(),
		))

		return extractArgs(cmd, args)
	}
}

// ExactArgsWithDefaultRequiredMsg is the same as ExactArgs but displays
// a default error message
func ExactArgsWithDefaultRequiredMsg(n int) cobra.PositionalArgs {
	argument := "argument"
	if n > 1 {
		argument = argument + "s"
	}

	return func(cmd *cobra.Command, args []string) error {
		extractArgs := ExactArgs(n, fmt.Sprintf("`%s` requires exactly %d %s.",
			cmd.CommandPath(),
			n,
			argument,
		))

		return extractArgs(cmd, args)
	}

}

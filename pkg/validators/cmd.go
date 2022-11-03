package validators

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/algolia/cli/pkg/cmdutil"
)

// ExactArgs is a validator for commands to print an error with a custom message
// followed by usage, flags and available commands when too few/much arguments are provided
func ExactArgsWithMsg(n int, msg string) cobra.PositionalArgs {
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
		extractArgs := ExactArgsWithMsg(0, fmt.Sprintf(
			"`%s` does not take any positional arguments.",
			cmd.CommandPath(),
		))

		return extractArgs(cmd, args)
	}
}

// ExactArgs is the same as ExactArgsWithMsg but displays
// a default error message
func ExactArgs(n int) cobra.PositionalArgs {
	argument := "argument"
	if n > 1 {
		argument = argument + "s"
	}

	return func(cmd *cobra.Command, args []string) error {
		extractArgs := ExactArgsWithMsg(n, fmt.Sprintf("`%s` requires exactly %d %s.",
			cmd.CommandPath(),
			n,
			argument,
		))

		return extractArgs(cmd, args)
	}

}

// AtLeastArgs is a validator for commands to print an error with a custom message
// followed by usage, flags and available commands when too few argument(s) are provided
func AtLeastArgs(n int) cobra.PositionalArgs {
	argument := "argument"
	if n > 1 {
		argument = argument + "s"
	}

	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return cmdutil.FlagErrorf(
				fmt.Sprintf("`%s` requires at least %d %s.", cmd.CommandPath(), n, argument))
		}

		return nil
	}

}

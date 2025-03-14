package root

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

func rootUsageFunc(
	IOStreams *iostreams.IOStreams,
	command *cobra.Command,
) func(cmd *cobra.Command) error {
	return cmdutil.UsageFuncDefault(IOStreams, command)
}

func rootFlagErrorFunc(cmd *cobra.Command, err error) error {
	if err == pflag.ErrHelp {
		return err
	}
	return cmdutil.FlagErrorWrap(err)
}

var hasFailed bool

// HasFailed signals that the main process should exit with non-zero status
func HasFailed() bool {
	return hasFailed
}

// Display helpful error message in case subcommand name was mistyped.
// This matches Cobra's behavior for root command, which Cobra
// confusingly doesn't apply to nested commands.
func nestedSuggestFunc(IOStreams *iostreams.IOStreams, command *cobra.Command, arg string) {
	w := IOStreams.ErrOut
	fmt.Fprintf(w, "unknown command %q for %q\n", arg, command.CommandPath())

	var candidates []string
	if arg == "help" {
		candidates = []string{"--help"}
	} else {
		if command.SuggestionsMinimumDistance <= 0 {
			command.SuggestionsMinimumDistance = 2
		}
		candidates = command.SuggestionsFor(arg)
	}

	if len(candidates) > 0 {
		fmt.Fprint(w, "\nDid you mean this?\n")
		for _, c := range candidates {
			fmt.Fprintf(w, "\t%s\n", c)
		}
	}

	fmt.Fprint(w, "\n")
	_ = rootUsageFunc(IOStreams, command)(command)
}

func isRootCmd(command *cobra.Command) bool {
	return command != nil && !command.HasParent()
}

func rootHelpFunc(f *cmdutil.Factory, command *cobra.Command, args []string) {
	cs := f.IOStreams.ColorScheme()

	if isRootCmd(command.Parent()) && len(args) >= 2 && args[1] != "--help" && args[1] != "-h" {
		nestedSuggestFunc(f.IOStreams, command, args[1])
		hasFailed = true
		return
	}

	longText := command.Long
	if longText == "" {
		longText = command.Short
	}

	helpEntries := cmdutil.UsageEntries{}
	if longText != "" {
		helpEntries.AddEntry(cmdutil.UsageEntry{Title: "", Body: longText})
	}

	helpEntries.AddBasicUsage(f.IOStreams, command)

	categoryFlagSet := cmdutil.NewCategoryFlagSet(command.LocalFlags())
	if len(categoryFlagSet.Categories) > 0 {
		for _, categoryName := range categoryFlagSet.SortedCategoryNames() {
			groupName := fmt.Sprintf("%s Flags", categoryName)
			helpEntries.AddEntry(
				cmdutil.UsageEntry{
					Title: groupName,
					Body:  cmdutil.Dedent(categoryFlagSet.Categories[categoryName].FlagUsages()),
				},
			)
		}
		if categoryFlagSet.Others.FlagUsages() != "" {
			helpEntries.AddEntry(
				cmdutil.UsageEntry{
					Title: cs.Bold("Other Flags"),
					Body:  cmdutil.Dedent(categoryFlagSet.Others.FlagUsages()),
				},
			)
		}
	} else {
		helpEntries.AddFlags(f.IOStreams, command, cmdutil.Dedent(categoryFlagSet.Others.FlagUsages()))
	}

	printFlagUsages := categoryFlagSet.Print.FlagUsages()
	if printFlagUsages != "" {
		helpEntries.AddEntry(
			cmdutil.UsageEntry{
				Title: cs.Bold("Output Formatting Flags"),
				Body:  cmdutil.Dedent(printFlagUsages),
			},
		)
	}
	inheritedFlagUsages := command.InheritedFlags().FlagUsages()
	if inheritedFlagUsages != "" {
		helpEntries.AddEntry(
			cmdutil.UsageEntry{
				Title: cs.Bold("Inherited Flags"),
				Body:  cmdutil.Dedent(inheritedFlagUsages),
			},
		)
	}
	if command.Example != "" {
		helpEntries.AddEntry(cmdutil.UsageEntry{Title: cs.Bold("Examples"), Body: command.Example})
	}
	if _, ok := command.Annotations["help:see-also"]; ok {
		helpEntries.AddEntry(
			cmdutil.UsageEntry{
				Title: cs.Bold("See also"),
				Body:  command.Annotations["help:see-also"],
			},
		)
	}
	helpEntries.AddEntry(cmdutil.UsageEntry{Title: cs.Bold("Learn More"), Body: `
Use 'algolia <command> <subcommand> --help' for more information about a command.
Read the documentation at https://algolia.com/doc/tools/cli/`})

	helpEntries.DisplayEntries(f.IOStreams.Out)
}

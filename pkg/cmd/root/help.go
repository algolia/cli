package root

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/text"
)

func rootUsageFunc(w io.Writer, command *cobra.Command) error {
	fmt.Fprintf(w, "Usage:  %s", command.UseLine())

	subcommands := command.Commands()
	if len(subcommands) > 0 {
		fmt.Fprint(w, "\n\nAvailable commands:\n")
		for _, c := range subcommands {
			if c.Hidden {
				continue
			}
			fmt.Fprintf(w, "  %s\n", c.Name())
		}
		return nil
	}

	flagUsages := command.LocalFlags().FlagUsages()
	if flagUsages != "" {
		fmt.Fprintln(w, "\n\nFlags:")
		fmt.Fprint(w, text.Indent(dedent(flagUsages), "  "))
	}
	return nil
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
func nestedSuggestFunc(w io.Writer, command *cobra.Command, arg string) {
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
	_ = rootUsageFunc(w, command)
}

func isRootCmd(command *cobra.Command) bool {
	return command != nil && !command.HasParent()
}

func rootHelpFunc(f *cmdutil.Factory, command *cobra.Command, args []string) {
	cs := f.IOStreams.ColorScheme()

	if isRootCmd(command.Parent()) && len(args) >= 2 && args[1] != "--help" && args[1] != "-h" {
		nestedSuggestFunc(f.IOStreams.ErrOut, command, args[1])
		hasFailed = true
		return
	}

	namePadding := 12
	commands := []string{}
	for _, c := range command.Commands() {
		if c.Short == "" {
			continue
		}
		if c.Hidden {
			continue
		}

		s := rpad(c.Name()+":", namePadding) + c.Short
		commands = append(commands, s)
	}

	type helpEntry struct {
		Title string
		Body  string
	}

	longText := command.Long
	if longText == "" {
		longText = command.Short
	}

	helpEntries := []helpEntry{}
	if longText != "" {
		helpEntries = append(helpEntries, helpEntry{"", longText})
	}
	helpEntries = append(helpEntries, helpEntry{cs.Bold("Usage"), command.UseLine()})
	if len(commands) > 0 {
		helpEntries = append(helpEntries, helpEntry{cs.Bold("Commands"), strings.Join(commands, "\n")})
	}

	categoryFlagSet := cmdutil.NewCategoryFlagSet(command.LocalFlags())
	if len(categoryFlagSet.Categories) > 0 {
		for _, categoryName := range categoryFlagSet.SortedCategoryNames() {
			groupName := fmt.Sprintf("%s Flags", categoryName)
			helpEntries = append(helpEntries, helpEntry{groupName, dedent(categoryFlagSet.Categories[categoryName].FlagUsages())})
		}
		if categoryFlagSet.Others.FlagUsages() != "" {
			helpEntries = append(helpEntries, helpEntry{cs.Bold("Other Flags"), dedent(categoryFlagSet.Others.FlagUsages())})
		}
	} else {
		if categoryFlagSet.Others.FlagUsages() != "" {
			helpEntries = append(helpEntries, helpEntry{cs.Bold("Flags"), dedent(categoryFlagSet.Others.FlagUsages())})
		}
	}

	printFlagUsages := categoryFlagSet.Print.FlagUsages()
	if printFlagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{cs.Bold("Output Formatting Flags"), dedent(printFlagUsages)})
	}
	inheritedFlagUsages := command.InheritedFlags().FlagUsages()
	if inheritedFlagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{cs.Bold("Inherited Flags"), dedent(inheritedFlagUsages)})
	}
	if command.Example != "" {
		helpEntries = append(helpEntries, helpEntry{cs.Bold("Examples"), command.Example})
	}
	helpEntries = append(helpEntries, helpEntry{cs.Bold("Learn More"), `
Use 'algolia <command> <subcommand> --help' for more information about a command.`})

	out := f.IOStreams.Out
	for _, e := range helpEntries {
		if e.Title != "" {
			// If there is a title, add indentation to each line in the body
			fmt.Fprintln(out, cs.Bold(e.Title))
			fmt.Fprintln(out, text.Indent(strings.Trim(e.Body, "\r\n"), "  "))
		} else {
			// If there is no title print the body as is
			fmt.Fprintln(out, e.Body)
		}
		fmt.Fprintln(out)
	}
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

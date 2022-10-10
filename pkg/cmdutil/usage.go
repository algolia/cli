package cmdutil

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/text"
	"github.com/algolia/cli/pkg/utils"
)

type UsageEntry struct {
	Title string
	Body  string
}

type UsageEntries struct {
	entries []UsageEntry
}

func (u *UsageEntries) AddEntry(entry UsageEntry) {
	u.entries = append(u.entries, entry)
}

func (u UsageEntries) AddBasicUsage(IOStreams *iostreams.IOStreams, command *cobra.Command) {
	cs := IOStreams.ColorScheme()
	u.AddEntry(UsageEntry{cs.Bold("Usage:"), command.UseLine()})
	subcommands := command.Commands()

	if len(subcommands) > 0 {
		namePadding := 12
		commands := []string{}
		for _, c := range subcommands {
			if c.Short == "" {
				continue
			}
			if c.Hidden {
				continue
			}
			commands = append(commands, rpad(c.Name()+":", namePadding)+c.Short)
		}

		u.AddEntry(UsageEntry{
			cs.Bold("Available commands:"),
			strings.Join(commands, "\n")},
		)
	}
}

func (u UsageEntries) AddFlags(IOStreams *iostreams.IOStreams, command *cobra.Command, flagUsages string) {
	cs := IOStreams.ColorScheme()

	if flagUsages != "" {
		u.AddEntry(UsageEntry{cs.Bold("Flags:"), flagUsages})
	}
}

func (u UsageEntries) AddAllFlags(IOStreams *iostreams.IOStreams, command *cobra.Command, flagUsages string) {
	u.AddFlags(IOStreams, command, command.LocalFlags().FlagUsages())
}

func (u UsageEntries) AddFilteredFlags(IOStreams *iostreams.IOStreams, command *cobra.Command, flagsToDisplay []string) {
	filteredFlags := filterFlagSet(*command.LocalFlags(), flagsToDisplay)

	u.AddFlags(IOStreams, command, filteredFlags.FlagUsages())
}

func (u UsageEntries) AddInheritedFlags(IOStreams *iostreams.IOStreams, command *cobra.Command) {
	cs := IOStreams.ColorScheme()
	inheritedFlagUsages := command.InheritedFlags().FlagUsages()

	if inheritedFlagUsages != "" {
		dedentedInheritedFlagUsages := Dedent(inheritedFlagUsages)
		u.AddEntry(UsageEntry{
			cs.Bold("Inherited Flags:"),
			fmt.Sprintln(text.Indent(strings.Trim(dedentedInheritedFlagUsages, "\r\n"), "  ")),
		})
	}
}

func (u UsageEntries) DisplayEntries(out io.Writer) {
	for _, e := range u.entries {
		if e.Title != "" {
			// If there is a title, add indentation to each line in the body
			fmt.Fprintln(out, e.Title)
			fmt.Fprintln(out, text.Indent(strings.Trim(e.Body, "\r\n"), "  "))
		} else {
			// If there is no title print the body as is
			fmt.Fprintln(out, e.Body)
		}
		fmt.Fprintln(out)
	}
}

func UsageFunc(IOStreams *iostreams.IOStreams, command *cobra.Command, flagUsages string) func(cmd *cobra.Command) error {
	return func(cmd *cobra.Command) error {
		entries := UsageEntries{}

		entries.AddBasicUsage(IOStreams, command)
		entries.AddFlags(IOStreams, command, flagUsages)

		entries.DisplayEntries(IOStreams.Out)
		return nil
	}
}

func UsageFuncDefault(IOStreams *iostreams.IOStreams, command *cobra.Command) func(cmd *cobra.Command) error {
	return UsageFunc(IOStreams, command, command.LocalFlags().FlagUsages())
}

func UsageFuncWithFilteredFlags(IOStreams *iostreams.IOStreams, command *cobra.Command, flagsToDisplay []string) func(cmd *cobra.Command) error {
	filteredFlags := filterFlagSet(*command.LocalFlags(), flagsToDisplay)

	return UsageFunc(IOStreams, command, filteredFlags.FlagUsages())

}

func UsageFuncWithFilteredAndInheritedFlags(IOStreams *iostreams.IOStreams, command *cobra.Command, flagsToDisplay []string) func(cmd *cobra.Command) error {
	return func(cmd *cobra.Command) error {

		entries := UsageEntries{}
		entries.AddBasicUsage(IOStreams, command)
		entries.AddFilteredFlags(IOStreams, command, flagsToDisplay)
		entries.AddInheritedFlags(IOStreams, command)

		entries.DisplayEntries(IOStreams.Out)
		return nil
	}
}

func UsageFuncWithInheritedFlagsOnly(IOStreams *iostreams.IOStreams, command *cobra.Command) func(cmd *cobra.Command) error {
	return func(cmd *cobra.Command) error {
		entries := UsageEntries{}
		entries.AddBasicUsage(IOStreams, command)
		entries.AddInheritedFlags(IOStreams, command)

		entries.DisplayEntries(IOStreams.Out)
		return nil
	}

}

func Dedent(s string) string {
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

func filterFlagSet(f pflag.FlagSet, flagsToDisplay []string) pflag.FlagSet {
	filteredFlags := pflag.NewFlagSet("flags", pflag.ContinueOnError)

	f.VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden && utils.Contains(flagsToDisplay, flag.Name) {
			filteredFlags.AddFlag(flag)
		}
	})

	return *filteredFlags
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

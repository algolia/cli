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

func UsageFunc(IOStreams *iostreams.IOStreams, command *cobra.Command, flagUsages string) error {
	entries := []UsageEntry{}
	err := AddBasicUsage(&entries, IOStreams, command)
	if err != nil {
		return err
	}

	if flagUsages != "" {
		AddFlags(&entries, IOStreams, command, flagUsages)
	}
	DisplayUsageEntry(IOStreams.Out, entries)

	return nil
}

func UsageFuncDefault(IOStreams *iostreams.IOStreams, command *cobra.Command) error {
	return UsageFunc(IOStreams, command, command.LocalFlags().FlagUsages())
}

func UsageFuncWithFilteredFlags(IOStreams *iostreams.IOStreams, command *cobra.Command, flagsToDisplay []string) error {
	filteredFlags := filterFlagSet(*command.LocalFlags(), flagsToDisplay)

	return UsageFunc(IOStreams, command, filteredFlags.FlagUsages())
}

func UsageFuncWithFilteredAndInheritedFlags(IOStreams *iostreams.IOStreams, command *cobra.Command, flagsToDisplay []string) error {
	entries := []UsageEntry{}
	err := AddBasicUsage(&entries, IOStreams, command)
	if err != nil {
		return err
	}

	err = AddFilteredFlags(&entries, IOStreams, command, flagsToDisplay)
	if err != nil {
		return err
	}

	err = AddInheritedFlags(&entries, IOStreams, command)
	if err != nil {
		return err
	}

	DisplayUsageEntry(IOStreams.Out, entries)
	return nil
}

func UsageFuncWithInheritedFlagsOnly(IOStreams *iostreams.IOStreams, command *cobra.Command) error {
	entries := []UsageEntry{}
	err := AddBasicUsage(&entries, IOStreams, command)
	if err != nil {
		return err
	}

	err = AddInheritedFlags(&entries, IOStreams, command)
	if err != nil {
		return err
	}

	DisplayUsageEntry(IOStreams.Out, entries)
	return nil
}

func DisplayUsageEntry(out io.Writer, entries []UsageEntry) {
	for _, e := range entries {
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

func AddBasicUsage(entries *[]UsageEntry, IOStreams *iostreams.IOStreams, command *cobra.Command) error {
	cs := IOStreams.ColorScheme()
	*entries = append(*entries, UsageEntry{cs.Bold("Usage:"), command.UseLine()})

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

		*entries = append(*entries, UsageEntry{
			cs.Bold("Available commands:"),
			strings.Join(commands, "\n")},
		)
		return nil
	}

	return nil
}

func AddFlags(entries *[]UsageEntry, IOStreams *iostreams.IOStreams, command *cobra.Command, flagUsages string) error {
	cs := IOStreams.ColorScheme()

	if flagUsages != "" {
		*entries = append(*entries, UsageEntry{cs.Bold("Flags:"), flagUsages})
	}
	return nil
}

func AddAllFlags(entries *[]UsageEntry, IOStreams *iostreams.IOStreams, command *cobra.Command) error {
	return AddFlags(entries, IOStreams, command, command.LocalFlags().FlagUsages())
}

func AddFilteredFlags(entries *[]UsageEntry, IOStreams *iostreams.IOStreams, command *cobra.Command, flagsToDisplay []string) error {
	filteredFlags := filterFlagSet(*command.LocalFlags(), flagsToDisplay)

	return AddFlags(entries, IOStreams, command, filteredFlags.FlagUsages())
}

func AddInheritedFlags(entries *[]UsageEntry, IOStreams *iostreams.IOStreams, command *cobra.Command) error {
	cs := IOStreams.ColorScheme()
	inheritedFlagUsages := command.InheritedFlags().FlagUsages()

	if inheritedFlagUsages != "" {
		dedentedInheritedFlagUsages := Dedent(inheritedFlagUsages)
		*entries = append(*entries, UsageEntry{
			cs.Bold("Inherited Flags:"),
			fmt.Sprintln(text.Indent(strings.Trim(dedentedInheritedFlagUsages, "\r\n"), "  ")),
		})
	}
	return nil
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

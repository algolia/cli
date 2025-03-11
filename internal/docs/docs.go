package docs

import (
	"strings"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Aliases     []string
	Examples    string
	Slug        string
	RunInWebCLI bool

	Flags map[string][]Flag

	SubCommands []Command
}

type Flag struct {
	Name        string
	Shorthand   string
	Description string
	Default     string
}

func newFlag(flag *pflag.Flag) Flag {
	return Flag{
		Name:      flag.Name,
		Shorthand: flag.Shorthand,
		// Add two spaces before each newline to force a newline in markdown
		Description: strings.ReplaceAll(flag.Usage, "\n", "  \n"),
		Default:     flag.DefValue,
	}
}

func newFlags(flagSet *pflag.FlagSet) []Flag {
	flags := make([]Flag, 0)
	flagSet.VisitAll(func(flag *pflag.Flag) {
		flags = append(flags, newFlag(flag))
	})
	return flags
}

func newCommand(cmd *cobra.Command) Command {
	categoryFlagSet := cmdutil.NewCategoryFlagSet(cmd.NonInheritedFlags())
	// Make sure the command description ends with a period.
	if !strings.HasSuffix(cmd.Short, ".") {
		cmd.Short += "."
	}

	command := Command{
		Name:        cmd.CommandPath(),
		Description: cmd.Short,
		Usage:       cmd.UseLine(),
		Aliases:     cmd.Aliases,
		Examples:    cmd.Example,
		RunInWebCLI: false,
	}
	if value, ok := cmd.Annotations["runInWebCLI"]; ok && value != "" {
		command.RunInWebCLI = true
	}

	flags := make(map[string][]Flag)

	if len(categoryFlagSet.Categories) > 0 {
		for _, categoryName := range categoryFlagSet.SortedCategoryNames() {
			flags[categoryName] = newFlags(categoryFlagSet.Categories[categoryName])
		}
		if categoryFlagSet.Others.HasAvailableFlags() {
			flags["Other flags"] = newFlags(categoryFlagSet.Others)
		}
	} else {
		if categoryFlagSet.Others.HasAvailableFlags() {
			flags["Flags"] = newFlags(categoryFlagSet.Others)
		}
	}

	if categoryFlagSet.Print.HasAvailableFlags() {
		flags["Output formatting flags"] = newFlags(categoryFlagSet.Print)

	}
	command.Flags = flags

	return command
}

func getCommands(cmd *cobra.Command) []Command {
	var commands []Command
	for _, c := range cmd.Commands() {
		if c.Hidden || c.Name() == "help" {
			continue
		}
		command := newCommand(c)
		if c.HasAvailableSubCommands() {
			for _, s := range c.Commands() {
				sub := newCommand(s)
				if s.HasAvailableSubCommands() {
					for _, sus := range s.Commands() {
						sub.SubCommands = append(sub.SubCommands, newCommand(sus))
					}
				}
				command.SubCommands = append(command.SubCommands, sub)
			}
		}
		commands = append(commands, command)
	}

	return commands
}

type Example struct {
	Desc          string
	Code          string
	WebCLICommand string
}

func (cmd Command) ExamplesList() []Example {
	var examples []Example
	examplesRaw := strings.Split(cmd.Examples, "#")
	for _, example := range examplesRaw {
		example = strings.TrimSpace(example)
		if len(example) == 0 {
			continue
		}
		exampleLines := strings.Split(example, "\n")

		code := strings.ReplaceAll(exampleLines[1], "$", "")
		code = strings.TrimSpace(code)

		formattedExample := Example{
			Desc: strings.Replace(exampleLines[0], "#", "", 1),
			Code: code,
		}
		if cmd.RunInWebCLI &&
			!strings.Contains(code, ">") &&
			!strings.Contains(code, "<") &&
			!strings.Contains(code, "output") &&
			!strings.Contains(code, "-F") &&
			!strings.Contains(code, "|") {
			formattedExample.WebCLICommand = strings.ReplaceAll(code, `"`, `\"`)
		}

		examples = append(examples, formattedExample)
	}
	return examples
}

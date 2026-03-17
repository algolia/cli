package docs

import (
	"strings"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Command struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Usage       string   `json:"usage"`
	Aliases     []string `json:"aliases,omitempty"`
	Examples    string   `json:"examples,omitempty"`
	Slug        string   `json:"slug,omitempty"`
	RunInWebCLI bool     `json:"runInWebCLI,omitempty"`
	CommandType string   `json:"commandType"`

	Annotations map[string]string `json:"annotations,omitempty"`

	Flags map[string][]Flag `json:"flags,omitempty"`

	SubCommands []Command `json:"subCommands,omitempty"`
}

type Flag struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Description string `json:"description"`
	Default     string `json:"default"`
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
		CommandType: commandType(cmd),
		Annotations: cmd.Annotations,
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

func commandType(cmd *cobra.Command) string {
	switch cmd.Name() {
	case "tail":
		return "stream"
	case "search", "browse", "get", "list", "stats", "analyze", "describe", "schema", "open":
		return "read"
	case "create", "delete", "clear", "import", "update", "set", "save", "move", "copy", "crawl", "run", "pause", "reindex", "unblock", "remove", "add", "setdefault":
		return "write"
	default:
		if cmd.HasAvailableSubCommands() {
			return "namespace"
		}
		return "other"
	}
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

func describeCommand(cmd *cobra.Command) Command {
	command := newCommand(cmd)
	if cmd.HasAvailableSubCommands() {
		command.SubCommands = getCommands(cmd)
	}
	return command
}

// DescribeCommand returns a machine-readable description of a command.
func DescribeCommand(cmd *cobra.Command) Command {
	return describeCommand(cmd)
}

// FindCommand resolves a command path against the provided root command.
func FindCommand(root *cobra.Command, args []string) (*cobra.Command, error) {
	current := root
	for _, arg := range args {
		next := findChildCommand(current, arg)
		if next == nil {
			return nil, cmdutil.FlagErrorf("unknown command %q for %q", arg, current.CommandPath())
		}
		current = next
	}
	return current, nil
}

func findChildCommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" {
			continue
		}
		if child.Name() == name {
			return child
		}
		for _, alias := range child.Aliases {
			if alias == name {
				return child
			}
		}
	}
	return nil
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

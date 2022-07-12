package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Aliases     []string
	Examples    string

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
	command := Command{
		Name:        cmd.CommandPath(),
		Description: cmd.Short,
		Usage:       cmd.UseLine(),
		Aliases:     cmd.Aliases,
		Examples:    cmd.Example,
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

func GenYamlTree(cmd *cobra.Command, dir string) error {
	var commands []Command
	for _, c := range cmd.Commands() {
		if c.Hidden || c.Name() == "help" {
			continue
		}
		command := newCommand(c)
		if c.HasAvailableSubCommands() {
			for _, sub := range c.Commands() {
				command.SubCommands = append(command.SubCommands, newCommand(sub))
			}
		}
		commands = append(commands, command)
	}

	for i, c := range commands {
		command_path := strings.ReplaceAll(c.Name, " ", "_")
		filename := filepath.Join(dir, fmt.Sprintf("0%d-%s.yml", i+1, command_path))
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		// Write a comment to the file, ignoring vale errors
		_, err = f.WriteString("# <!-- vale off -->\n")
		if err != nil {
			return err
		}

		encoder := yaml.NewEncoder(f)
		defer encoder.Close()

		err = encoder.Encode(c)
		if err != nil {
			return err
		}
	}

	return nil
}

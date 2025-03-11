package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/algolia/cli/internal/docs"
	"github.com/algolia/cli/pkg/cmd/root"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	dir := flags.StringP("app_data-path", "", "", "Path directory where you want generate documentation data files")
	help := flags.BoolP("help", "h", false, "Help about any command")
	target := flags.StringP("target", "T", "old", "target old or new documentation website")

	if *target != "old" && *target != "new" {
		return fmt.Errorf("error: --destination can only be 'old' or 'new' ('old' by default)")
	}

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n%s", filepath.Base(args[0]), flags.FlagUsages())
		return nil
	}

	if *dir == "" {
		return fmt.Errorf("error: --app_data-path not set")
	}

	ios, _, _, _ := iostreams.Test()
	rootCmd := root.NewRootCmd(&cmdutil.Factory{
		IOStreams: ios,
		Config:    &config.Config{},
	})
	rootCmd.InitDefaultHelpCmd()

	if err := os.MkdirAll(*dir, 0755); err != nil {
		return err
	}

	if *target == "old" {
		if err := docs.GenYamlTree(rootCmd, *dir); err != nil {
			return err
		}
	} else {
		if err := docs.GenMdxTree(rootCmd, *dir); err != nil {
			return err
		}
	}

	return nil

}

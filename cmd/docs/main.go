package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	dir := flags.StringP("doc-path", "", "", "Path directory where you want generate doc files")
	help := flags.BoolP("help", "h", false, "Help about any command")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n%s", filepath.Base(args[0]), flags.FlagUsages())
		return nil
	}

	if *dir == "" {
		return fmt.Errorf("error: --doc-path not set")
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

	if err := docs.GenMarkdownTree(rootCmd, *dir, linkHandler); err != nil {
		return err
	}

	return nil

}

func linkHandler(name string) string {
	return fmt.Sprintf("./%s", strings.TrimSuffix(name, ".md"))
}

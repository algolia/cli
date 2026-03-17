package describe

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/internal/docs"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

const schemaVersion = "v1"

type DescribeOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CommandPath []string
	PrintFlags  *cmdutil.PrintFlags
}

type DescribeResponse struct {
	SchemaVersion string       `json:"schemaVersion"`
	Command       docs.Command `json:"command"`
}

// NewDescribeCmd creates and returns a describe command.
func NewDescribeCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DescribeOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:     "describe [command] [subcommand]...",
		Aliases: []string{"schema"},
		Args:    cobra.ArbitraryArgs,
		Short:   "Describe commands and flags as JSON.",
		Long: heredoc.Doc(`
			Describe the CLI's command tree in a machine-readable format.
			With no arguments, this command describes the root command.
		`),
		Example: heredoc.Doc(`
			# Describe the root command
			$ algolia describe

			# Describe the search command
			$ algolia describe search

			# Describe the objects browse command
			$ algolia describe objects browse
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CommandPath = args
			return runDescribeCmd(cmd, opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runDescribeCmd(cmd *cobra.Command, opts *DescribeOptions) error {
	target, err := docs.FindCommand(cmd.Root(), opts.CommandPath)
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	return p.Print(opts.IO, DescribeResponse{
		SchemaVersion: schemaVersion,
		Command:       docs.DescribeCommand(target),
	})
}

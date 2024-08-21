package list

import (
	"github.com/algolia/cli/api/provisionning"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	ProvisionningClient func() (*provisionning.Client, error)

	PrintFlags *cmdutil.PrintFlags
}

// NewListCmd creates and returns a list command for Applications.
func NewListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                  f.IOStreams,
		Config:              f.Config,
		ProvisionningClient: f.ProvisionningClient,
		PrintFlags:          cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all applications",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// runListCmd executes the list command
func runListCmd(opts *ListOptions) error {
	client, err := opts.ProvisionningClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching applications")
	res, err := client.ListApplications()
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	return p.Print(opts.IO, res)
}

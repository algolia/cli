package domains

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type CreateOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
	AgentID, Domain   string
	DryRun            bool
	OutputChanged     bool
}

func newCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "create <agent-id> --domain <pattern>",
		Short: "Add a single allowed domain to an agent",
		Example: heredoc.Doc(`
			$ algolia agents domains create <agent-id> --domain https://app.example.com
			$ algolia agents domains create <agent-id> --domain "*.example.com"
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.Domain == "" {
				return cmdutil.FlagErrorf("--domain is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}
	cmd.Flags().StringVar(&opts.Domain, "domain", "", "Domain or pattern (required)")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Print the would-be request without sending")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out,
			"Dry run: would POST /1/agents/%s/allowed-domains\n  body: {\"domain\":%q}\n",
			opts.AgentID, opts.Domain)
		return nil
	}
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Adding allowed domain")
	d, err := client.CreateAllowedDomain(ctxOrBackground(opts.Ctx), opts.AgentID, opts.Domain)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	return opts.PrintFlags.Print(opts.IO, d)
}

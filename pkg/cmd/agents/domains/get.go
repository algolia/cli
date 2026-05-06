package domains

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// nowFn is overridable for deterministic time-based output in tests.
var nowFn = time.Now

func nowFnOrTime() time.Time { return nowFn() }

type GetOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
	AgentID, DomainID string
}

func newGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "get <agent-id> <domain-id>",
		Short: "Get a single allowed domain",
		Args:  validators.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID, opts.DomainID = args[0], args[1]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" || opts.DomainID == "" {
				return cmdutil.FlagErrorf("agent-id and domain-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runGetCmd(opts)
		},
	}
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching allowed domain")
	d, err := client.GetAllowedDomain(ctxOrBackground(opts.Ctx), opts.AgentID, opts.DomainID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	return opts.PrintFlags.Print(opts.IO, d)
}

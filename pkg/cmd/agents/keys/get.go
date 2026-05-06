package keys

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

var nowFn = time.Now

func nowFnOrTime() time.Time { return nowFn() }

type GetOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
	ID                string
	ShowSecret        bool
}

func newGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a single secret key",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			opts.Ctx = cmd.Context()
			if opts.ID == "" {
				return cmdutil.FlagErrorf("id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runGetCmd(opts)
		},
	}
	cmd.Flags().BoolVar(&opts.ShowSecret, "show-secret", false, "Reveal raw key value (default redacted as ***)")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching secret key")
	k, err := client.GetSecretKey(ctxOrBackground(opts.Ctx), opts.ID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	masked := maskKey(*k, opts.ShowSecret)
	return opts.PrintFlags.Print(opts.IO, &masked)
}

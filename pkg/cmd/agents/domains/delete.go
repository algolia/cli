package domains

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	AgentID, DomainID string
	DryRun            bool
	DoConfirm         bool
}

func newDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	var confirm bool
	cmd := &cobra.Command{
		Use:   "delete <agent-id> <domain-id> [--confirm]",
		Short: "Remove an allowed domain by ID",
		Args:  validators.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID, opts.DomainID = args[0], args[1]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" || opts.DomainID == "" {
				return cmdutil.FlagErrorf("agent-id and domain-id must not be empty")
			}
			doConfirm, err := shared.ResolveConfirm(opts.IO, confirm, opts.DryRun)
			if err != nil {
				return err
			}
			opts.DoConfirm = doConfirm
			if runF != nil {
				return runF(opts)
			}
			return runDeleteCmd(opts)
		},
	}
	shared.AddConfirmFlag(cmd, &confirm)
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Print what would be deleted without calling the API")
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out,
			"Dry run: would DELETE /1/agents/%s/allowed-domains/%s\n",
			opts.AgentID, opts.DomainID)
		return nil
	}
	if opts.DoConfirm {
		ok, err := shared.Confirm(fmt.Sprintf("Delete allowed domain %s?", opts.DomainID))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Deleting allowed domain")
	err = client.DeleteAllowedDomain(shared.OrBackground(opts.Ctx), opts.AgentID, opts.DomainID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted allowed domain %s\n", cs.SuccessIcon(), opts.DomainID)
	}
	return nil
}

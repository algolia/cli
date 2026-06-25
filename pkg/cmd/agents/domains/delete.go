package domains

import (
	"context"
	"fmt"

	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	IO                   *iostreams.IOStreams
	Ctx                  context.Context
	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	AgentID, DomainID    string
	DoConfirm            bool
}

func newDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
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
			doConfirm, err := shared.ResolveConfirm(opts.IO, confirm)
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
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	if opts.DoConfirm {
		ok, err := shared.Confirm(fmt.Sprintf("Delete allowed domain %s?", opts.DomainID))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Deleting allowed domain")
	err = client.DeleteAllowedDomain(
		client.NewApiDeleteAllowedDomainRequest(opts.DomainID, opts.AgentID),
		agentStudio.WithContext(shared.OrBackground(opts.Ctx)),
	)
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

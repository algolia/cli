package domains

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
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
			if !confirm && !opts.DryRun {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}
			if runF != nil {
				return runF(opts)
			}
			return runDeleteCmd(opts)
		},
	}
	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
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
		var ok bool
		if err := prompt.Confirm(fmt.Sprintf("Delete allowed domain %s?", opts.DomainID), &ok); err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
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
	err = client.DeleteAllowedDomain(ctxOrBackground(opts.Ctx), opts.AgentID, opts.DomainID)
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

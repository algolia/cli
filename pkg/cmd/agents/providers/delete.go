package providers

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// DeleteOptions collects inputs for `agents providers delete`.
type DeleteOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	ProviderID string
	DoConfirm  bool
}

func newDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <provider-id> [--confirm]",
		Short: "Delete an LLM provider authentication",
		Long: heredoc.Doc(`
			Delete a provider authentication. The backend may 409 if any
			agent still references this provider; the structured detail
			surfaces verbatim. Detach affected agents first (or update
			them to point at a different provider) before retrying.

			Like "agents delete", interactive use prompts to confirm and
			non-interactive use requires --confirm.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers delete <id>           # interactive
			$ algolia agents providers delete <id> -y        # CI
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]
			opts.Ctx = cmd.Context()
			if opts.ProviderID == "" {
				return cmdutil.FlagErrorf("provider-id must not be empty")
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
	cs := opts.IO.ColorScheme()
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := shared.OrBackground(opts.Ctx)

	// Pre-fetch so the prompt can show name+providerName,
	// matching `agents delete`'s contract.
	p, err := client.GetProvider(ctx, opts.ProviderID)
	if err != nil {
		return err
	}

	if opts.DoConfirm {
		ok, err := shared.Confirm(
			fmt.Sprintf("Delete provider %q (%s, %s)?", p.Name, p.ProviderName, opts.ProviderID),
		)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	opts.IO.StartProgressIndicatorWithLabel("Deleting provider")
	err = client.DeleteProvider(ctx, opts.ProviderID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted provider %s\n", cs.SuccessIcon(), opts.ProviderID)
	}
	return nil
}

package keys

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
	ID                string
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
		Use:   "delete <id> [--confirm]",
		Short: "Delete a secret key (admin key required)",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			opts.Ctx = cmd.Context()
			if opts.ID == "" {
				return cmdutil.FlagErrorf("id must not be empty")
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
		fmt.Fprintf(opts.IO.Out, "Dry run: would DELETE /1/secret-keys/%s\n", opts.ID)
		return nil
	}
	if opts.DoConfirm {
		ok, err := shared.Confirm(fmt.Sprintf("Delete secret key %s?", opts.ID))
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
	opts.IO.StartProgressIndicatorWithLabel("Deleting secret key")
	err = client.DeleteSecretKey(shared.OrBackground(opts.Ctx), opts.ID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted secret key %s\n", cs.SuccessIcon(), opts.ID)
	}
	return nil
}

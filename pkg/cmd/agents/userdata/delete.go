package userdata

import (
	"context"
	"fmt"
	"strings"

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
	UserToken         string
	DoConfirm         bool
}

func newDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	var confirm bool
	cmd := &cobra.Command{
		Use:   "delete <user-token> [--confirm]",
		Short: "Erase ALL conversations + memories tied to a user token (irreversible)",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.UserToken = args[0]
			opts.Ctx = cmd.Context()
			if opts.UserToken == "" {
				return cmdutil.FlagErrorf("user-token must not be empty")
			}
			if strings.Contains(opts.UserToken, "/") {
				return cmdutil.FlagErrorf("%s", rejectSlashMsg)
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
		msg := fmt.Sprintf(
			"Erase ALL conversations and memories for user token %q across every agent in this app? This cannot be undone.",
			opts.UserToken,
		)
		ok, err := shared.Confirm(msg)
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
	opts.IO.StartProgressIndicatorWithLabel("Erasing user data")
	err = client.DeleteUserData(shared.OrBackground(opts.Ctx), opts.UserToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Erased user data for %s\n", cs.SuccessIcon(), opts.UserToken)
	}
	return nil
}

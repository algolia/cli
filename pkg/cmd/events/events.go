package events

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/events/tail"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewEventsCmd returns a new command for events.
func NewEventsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Manage your Algolia events",
	}

	cmd.AddCommand(tail.NewTailCmd(f, nil))

	return cmd
}

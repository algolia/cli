package app

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/app/copy"
	"github.com/algolia/cli/pkg/cmd/app/move"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type AppOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	PrintFlags *cmdutil.PrintFlags
}

// NewAppCmd creates and returns an app command
func NewAppCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		Aliases: []string{"apps"},
		Short:   "Manage your Algolia apps",
	}

	cmd.AddCommand(copy.NewCopyCmd(f))
	cmd.AddCommand(move.NewMoveCmd(f))

	return cmd
}

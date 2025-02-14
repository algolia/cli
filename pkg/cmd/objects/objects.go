package objects

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/objects/browse"
	"github.com/algolia/cli/pkg/cmd/objects/delete"
	importObjects "github.com/algolia/cli/pkg/cmd/objects/import"
	"github.com/algolia/cli/pkg/cmd/objects/operations"
	updateObjects "github.com/algolia/cli/pkg/cmd/objects/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewObjectsCmd returns a new command for indices objects.
func NewObjectsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "objects",
		Short:   "Manage your indices' objects",
		Aliases: []string{"records"},
	}

	cmd.AddCommand(browse.NewBrowseCmd(f))
	cmd.AddCommand(importObjects.NewImportCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))
	cmd.AddCommand(updateObjects.NewUpdateCmd(f, nil))
	cmd.AddCommand(operations.NewOperationsCmd(f, nil))

	return cmd
}

package feedback

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
)

// NewFeedbackCmd is the parent for `algolia agents feedback <verb>`.
func NewFeedbackCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feedback",
		Short: "Submit user feedback on agent messages",
	}
	cmd.AddCommand(newCreateCmd(f, nil))
	return cmd
}

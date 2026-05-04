package export

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type ExportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID string
	Output  string
}

func NewExportCmd(f *cmdutil.Factory, runF func(*ExportOptions) error) *cobra.Command {
	opts := &ExportOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}
	cmd := &cobra.Command{
		Use:   "export <agent_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Export an agent's conversations as raw bytes",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# Print to stdout
			$ algolia agentstudio conversations export a1b2

			# Save to a file
			$ algolia agentstudio conversations export a1b2 -o conversations.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			if runF != nil {
				return runF(opts)
			}
			return runExportCmd(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.Output, "output-file", "o", "", "Write the export body to the named file instead of stdout")
	return cmd
}

func runExportCmd(opts *ExportOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Exporting conversations")
	body, err := client.ExportConversations(opts.AgentID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.Output != "" {
		if err := os.WriteFile(opts.Output, body, 0o644); err != nil {
			return err
		}
		if opts.IO.IsStdoutTTY() {
			cs := opts.IO.ColorScheme()
			fmt.Fprintf(opts.IO.Out, "%s Wrote %d bytes to %s\n", cs.SuccessIcon(), len(body), opts.Output)
		}
		return nil
	}

	_, err = opts.IO.Out.Write(body)
	return err
}

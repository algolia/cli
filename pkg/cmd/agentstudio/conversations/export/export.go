package export

import (
	"fmt"
	"io"
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
	defer body.Close()

	if opts.Output != "" {
		// 0600: conversation exports may contain user PII / prompts.
		f, err := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			return err
		}
		n, copyErr := io.Copy(f, body)
		closeErr := f.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		if opts.IO.IsStdoutTTY() {
			cs := opts.IO.ColorScheme()
			fmt.Fprintf(opts.IO.Out, "%s Wrote %d bytes to %s\n", cs.SuccessIcon(), n, opts.Output)
		}
		return nil
	}

	_, err = io.Copy(opts.IO.Out, body)
	return err
}

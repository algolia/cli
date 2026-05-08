package create

import (
	"context"
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type CreateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	File          string
	OutputChanged bool
}

func NewCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "create -F <file>",
		Short: "Create an Agent Studio agent from a JSON file",
		Long: heredoc.Doc(`
			Create a new Agent Studio agent from a JSON file describing the
			AgentConfigCreate body (name, instructions, model, providerId,
			tools, config, …). The file is sent verbatim to the backend; the
			CLI only validates that it's well-formed JSON. Field-level
			validation is the backend's job and surfaces as a 422 error.
		`),
		Example: heredoc.Doc(`
			# Create from a file
			$ algolia agents create -F spec.json

			# Create from stdin
			$ cat spec.json | algolia agents create -F -
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")
			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "JSON file with the agent body (use \"-\" for stdin)")
	_ = cmd.MarkFlagRequired("file")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	body, err := shared.ReadJSONFile(opts.IO.In, opts.File)
	if err != nil {
		return err
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Creating agent")
	agent, err := client.CreateAgent(ctx, json.RawMessage(body))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return opts.PrintFlags.Print(opts.IO, agent)
}

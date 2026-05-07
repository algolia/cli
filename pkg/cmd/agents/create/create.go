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
	DryRun        bool
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

			Use --dry-run to print the resolved request body without sending
			it — useful in CI to lint generated agent specs.
		`),
		Example: heredoc.Doc(`
			# Create from a file
			$ algolia agents create -F spec.json

			# Create from stdin
			$ cat spec.json | algolia agents create -F -

			# Preview the request without sending it
			$ algolia agents create -F spec.json --dry-run
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
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Validate and print the resolved request body without calling the API")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	body, err := shared.ReadJSONFile(opts.IO.In, opts.File)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, opts.PrintFlags, opts.OutputChanged,
			"create_agent", "POST /1/agents", opts.File, body, nil)
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

package get

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type GetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	ObjectID string
}

// NewGetCmd creates and returns a get command for GenAI responses.
func NewGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
	}

	cmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get a GenAI response",
		Long: heredoc.Doc(`
			Get a GenAI response by ID.
		`),
		Example: heredoc.Doc(`
			# Get a response by ID
			$ algolia genai response get b4e52d1a-2509-49ea-ba36-f6f5c3a83ba9
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ObjectID = args[0]

			if runF != nil {
				return runF(opts)
			}

			return runGetCmd(opts)
		},
	}

	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Getting response")

	response, err := client.GetResponse(opts.ObjectID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("ID:"), cs.Bold(response.ObjectID))
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Query:"), response.Query)
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Data Source:"), response.DataSourceID)
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Prompt:"), response.PromptID)

		if response.AdditionalFilters != "" {
			fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Additional Filters:"), response.AdditionalFilters)
		}

		fmt.Fprintf(opts.IO.Out, "%s %t\n", cs.Bold("Save:"), response.Save)
		fmt.Fprintf(opts.IO.Out, "%s %t\n", cs.Bold("Use Cache:"), response.UseCache)
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Origin:"), response.Origin)
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Created:"), formatTime(response.CreatedAt))
		fmt.Fprintf(opts.IO.Out, "%s %s\n\n", cs.Bold("Updated:"), formatTime(response.UpdatedAt))

		if response.Response != "" {
			fmt.Fprintf(opts.IO.Out, "%s\n", response.Response)
		}
	} else {
		fmt.Fprintf(opts.IO.Out, "%s", response.Response)
	}

	return nil
}

// formatTime formats time to a readable format
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

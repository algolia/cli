package get

import (
	"bytes"
	"encoding/json"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/jsoncolor"
	"github.com/algolia/cli/pkg/validators"
)

type GetOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice string
}

// NewGetCmd creates and returns a get command for settings
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:               "get <indice>",
		Args:              validators.ExactArgs(1),
		Short:             "Get settings",
		Long:              `Settings for the specified index.`,
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			return runListCmd(opts)
		},
	}

	return cmd
}

func runListCmd(opts *GetOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching settings")
	res, err := client.InitIndex(opts.Indice).GetSettings()
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(res)

	if opts.IO.ColorEnabled() {
		jsoncolor.Write(opts.IO.Out, &buf, "  ")
	} else {
		opts.IO.Out.Write(buf.Bytes())
	}
	return nil
}

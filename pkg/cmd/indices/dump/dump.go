package dump

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algolia-cli/pkg/cmdutil"
	"github.com/algolia/algolia-cli/pkg/config"
	"github.com/algolia/algolia-cli/pkg/iostreams"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"
)

type DumpOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice string
}

// NewDumpCmd creates and returns a delete command for indices
func NewDumpCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DumpOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:   "dump <index_1>...",
		Args:  cobra.ExactArgs(1),
		Short: "Dump the indice",
		Long: heredoc.Doc(`
			Dump the given indices.
			This command dump the records of the specified indice.
		`),
		Example: heredoc.Doc(`
			$ algolia indices dump TEST_PRODUCTS_1
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			client, err := opts.SearchClient()
			if err != nil {
				return err
			}

			// Test that all the provided indices exist
			exists, err := client.InitIndex(opts.Indice).Exists()
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("the index %s does not exist", opts.Indice)
			}

			return runDumpCmd(opts)
		},
	}

	return cmd
}

func runDumpCmd(opts *DumpOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	res, err := indice.BrowseObjects()
	if err != nil {
		return err
	}
	for {
		iObject, err := res.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		obj, err := json.Marshal(iObject)
		if err != nil {
			return err
		}
		opts.IO.Out.Write([]byte(fmt.Sprintf("%s\n", obj)))
	}
}

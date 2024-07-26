package browse

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type BrowseOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index        string
	BrowseParams map[string]interface{}
	Retries      int

	PrintFlags *cmdutil.PrintFlags
}

// NewBrowseCmd creates and returns a browse command for index objects
func NewBrowseCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &BrowseOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4_SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "browse <index>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.V4_IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "browse",
		},
		Short: "Browse the index objects",
		Long: heredoc.Doc(`
			This command browse the objects of the specified index.
		`),
		Example: heredoc.Doc(`
			# Browse the objects from the "MOVIES" index
			$ algolia objects browse MOVIES

			# Browse the objects from the "MOVIES" index and select which attributes to retrieve
			$ algolia objects browse MOVIES --attributesToRetrieve title,overview

			# Browse the objects from the "MOVIES" index with filters
			$ algolia objects browse MOVIES --filters "genres:Drama"

			# Browse the objects from the "MOVIES" and export the results to a new line delimited JSON (ndjson) file
			$ algolia objects browse MOVIES > movies.ndjson
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			browseParams, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.BrowseParamsObject...)
			if err != nil {
				return err
			}
			opts.BrowseParams = browseParams

			return runBrowseCmd(opts)
		},
	}

	cmd.Flags().
		IntVarP(&opts.Retries, "retries", "r", 1000, "Max. number of browse requests. Each request retrieves up to 1,000 records.")

	cmd.SetUsageFunc(cmdutil.UsageFuncWithInheritedFlagsOnly(f.IOStreams, cmd))

	cmdutil.AddSearchParamsObjectFlags(cmd)
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runBrowseCmd(opts *BrowseOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	browseParams := search.NewEmptyBrowseParamsObject()
	// Convert `v3` options to `v4`
	MapToStruct(opts.BrowseParams, browseParams)

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	_, err = client.BrowseObjects(
		opts.Index,
		*browseParams,
		search.WithAggregator(func(res any, _ error) {
			for _, hit := range res.(*search.BrowseResponse).Hits {
				p.Print(opts.IO, hit)
			}
		}),
		search.WithMaxRetries(opts.Retries),
	)
	if err != nil {
		return err
	}
	return nil
}

// Capitalize makes the first letter of a word uppercase
func Capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	firstRune := []rune(word)[0]
	rest := []rune(word)[1:]
	return string(unicode.ToUpper(firstRune)) + string(rest)
}

// MapToStruct converts a map into a struct
func MapToStruct(m map[string]any, s interface{}) error {
	val := reflect.ValueOf(s).Elem()

	for k, v := range m {
		// cmdline options are lowercase (`--query`),
		// but struct fields are capital (`Query`)
		field := val.FieldByName(Capitalize(k))
		if !field.IsValid() {
			return errors.New(fmt.Sprintf("No such parameter: %s for browse\n.", k))
		}

		if !field.CanSet() {
			return errors.New(fmt.Sprintf("Can't set field: %s\n", field))
		}

		fieldValue := reflect.ValueOf(v)

		if field.Type().Kind() == reflect.Ptr &&
			fieldValue.Type().ConvertibleTo(field.Type().Elem()) {
			newValue := reflect.New(fieldValue.Type()).Elem()
			newValue.Set(fieldValue)
			field.Set(newValue.Addr())
		} else if fieldValue.Type().ConvertibleTo(field.Type()) {
			field.Set(fieldValue.Convert(field.Type()))
		} else {
			return errors.New(fmt.Sprintf("Can't convert type of %s to %s\n", fieldValue.Type(), field.Type()))
		}
	}
	return nil
}

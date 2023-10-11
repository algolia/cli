package delete

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice       string
	ObjectIDs    []string
	DeleteParams map[string]interface{}

	DoConfirm bool
	Wait      bool
}

// NewDeleteCmd creates and returns a delete command for index objects
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	var confirm bool

	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "delete <index> [--object-ids <object-ids> | --filters  <filters>...] [--confirm] [--wait]",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Delete objects from an index",
		Long: heredoc.Doc(`
			This command deletes the objects from the specified index.

			You can either directly specify the objects to delete by theirs IDs and/or use the filters related flags to delete the matching objects.
		`),
		Example: heredoc.Doc(`
			# Delete one single object with the ID "1" from the "MOVIES" index
			$ algolia objects delete MOVIES --object-ids 1

			# Delete multiple objects with the IDs "1" and "2" from the "MOVIES" index
			$ algolia objects delete MOVIES --object-ids 1,2

			# Delete all objects matching the filters "type:Scripted" from the "MOVIES" index
			$ algolia objects delete MOVIES --filters "type:Scripted" --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]
			deleteParams, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.DeleteByParams...)
			if err != nil {
				return err
			}
			opts.DeleteParams = deleteParams

			if len(opts.ObjectIDs) == 0 && len(opts.DeleteParams) == 0 {
				return cmdutil.FlagErrorf("you must specify either --object-ids or a filter")
			}

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.ObjectIDs, "object-ids", "", nil, "Object IDs to delete")
	cmdutil.AddDeleteByParamsFlags(cmd)

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().BoolVar(&opts.Wait, "wait", false, "wait for all the operations to complete before returning")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	nbObjectsToDelete := len(opts.ObjectIDs)
	extra := "Operation aborted, no deletion action taken"

	// Tests if the provided object IDs exists.
	for _, objectID := range opts.ObjectIDs {
		var obj interface{}
		if err := indice.GetObject(objectID, &obj); err != nil {
			// The original error is not helpful, so we print a more helpful message
			if strings.Contains(err.Error(), "ObjectID does not exist") {
				return fmt.Errorf("object with ID '%s' does not exist. %s", objectID, extra)
			}
			return fmt.Errorf("%s. %s", err, extra)
		}
	}

	// We count the number of objects matching the filters if they are provided.
	// The count is used to display the confirmation message, but it is sometimes approximate.
	exactOrApproximate := "exactly"

	// If the user provided filters, we need to count the number of objects matching the filters
	if len(opts.DeleteParams) > 0 {
		res, err := indice.Search("", opt.ExtraOptions(opts.DeleteParams))
		if err != nil {
			return err
		}
		nbObjectsToDelete = nbObjectsToDelete + res.NbHits
		if !res.ExhaustiveNbHits {
			exactOrApproximate = "approximately"
		}
	}

	if nbObjectsToDelete == 0 {
		if _, err = fmt.Fprintf(opts.IO.Out, "%s No objects to delete. %s\n", cs.WarningIcon(), extra); err != nil {
			return err
		}
		return nil
	}

	objectNbMessage := fmt.Sprintf("%s %s from %s", exactOrApproximate, utils.Pluralize(nbObjectsToDelete, "object"), opts.Indice)

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(fmt.Sprintf("Delete %s?", objectNbMessage), &confirmed)
		if err != nil {
			return fmt.Errorf("%s Failed to prompt: %w", cs.FailureIcon(), err)
		}
		if !confirmed {
			return nil
		}
	}

	var taskIDs []int64

	// Delete the objects by their IDs
	if len(opts.ObjectIDs) > 0 {
		deleteByIDRes, err := indice.DeleteObjects(opts.ObjectIDs)
		if err != nil {
			return err
		}

		taskIDs = append(taskIDs, deleteByIDRes.TaskID)
	}

	// Delete the objects matching the filters
	if len(opts.DeleteParams) > 0 {
		deleteByOpts, err := deleteParamsToDeleteByOpts(opts.DeleteParams)
		if err != nil {
			return err
		}

		deleteByRes, err := indice.DeleteBy(deleteByOpts...)
		if err != nil {
			return err
		}

		taskIDs = append(taskIDs, deleteByRes.TaskID)
	}

	// Wait for the tasks to complete
	if opts.Wait {
		opts.IO.StartProgressIndicatorWithLabel("Waiting for all of the deletion tasks to complete")
		for _, taskID := range taskIDs {
			if err := indice.WaitTask(taskID); err != nil {
				return err
			}
		}
		opts.IO.StopProgressIndicator()
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully deleted %s\n", cs.SuccessIcon(), objectNbMessage)
	}

	return nil
}

// flagValueToOpts returns a given option from the provided flag.
// It is used to convert the flag value to the correct type expected by the `DeleteBy` method.
func flagValueToOpts(value interface{}, opt interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, opt); err != nil {
		return err
	}

	return nil
}

// deleteParamsToDeleteByOpts returns an array of deleteByOptions from the provided delete parameters.
func deleteParamsToDeleteByOpts(params map[string]interface{}) ([]interface{}, error) {
	var opts []interface{}

	for key, value := range params {
		switch key {
		case "filters":
			var filtersOpt opt.FiltersOption
			if err := flagValueToOpts(value, &filtersOpt); err != nil {
				return nil, err
			}

			opts = append(opts, &filtersOpt)

		case "facetFilters":
			var facetFiltersOpt opt.FacetFiltersOption
			if err := flagValueToOpts(value, &facetFiltersOpt); err != nil {
				return nil, err
			}

			opts = append(opts, &facetFiltersOpt)

		case "numericFilters":
			var numericFiltersOpt opt.NumericFiltersOption
			if err := flagValueToOpts(value, &numericFiltersOpt); err != nil {
				return nil, err
			}

			opts = append(opts, &numericFiltersOpt)

		case "tagFilters":
			var tagFiltersOpt opt.TagFiltersOption
			if err := flagValueToOpts(value, &tagFiltersOpt); err != nil {
				return nil, err
			}

			opts = append(opts, &tagFiltersOpt)

		case "aroundLatLng":
			var aroundLatLngOpt opt.AroundLatLngOption
			if err := flagValueToOpts(value, &aroundLatLngOpt); err != nil {
				return nil, err
			}

			opts = append(opts, &aroundLatLngOpt)
		}
	}

	return opts, nil
}

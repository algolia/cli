package delete

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
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

	SearchClient func() (*search.APIClient, error)

	Index         string
	ObjectIDs     []string
	DeleteParams  search.DeleteByParams
	NdeleteParams int

	DoConfirm bool
	Wait      bool
}

// NewDeleteCmd creates and returns a delete command for index objects
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	var confirm bool

	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "delete <index> [--object-ids <object-ids> | --filters  <filters>...] [--confirm] [--wait]",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.V4IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "deleteObject",
		},
		Short: "Delete objects from an index",
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
			opts.Index = args[0]
			deleteParams, err := cmdutil.FlagValuesMap(cmd.Flags(), cmdutil.DeleteByParams...)
			if err != nil {
				return err
			}
			opts.NdeleteParams = len(deleteParams)

			// Convert map into object
			tmp, err := json.Marshal(deleteParams)
			if err != nil {
				return err
			}
			err = json.Unmarshal(tmp, &opts.DeleteParams)
			if err != nil {
				return err
			}

			if len(opts.ObjectIDs) == 0 && opts.NdeleteParams == 0 {
				return cmdutil.FlagErrorf("you must specify either --object-ids or a filter")
			}

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf(
						"--confirm required when non-interactive shell is detected",
					)
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
	cmd.Flags().
		BoolVar(&opts.Wait, "wait", false, "wait for all the operations to complete before returning")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}
	nbObjectsToDelete := len(opts.ObjectIDs)
	extra := "Operation aborted, no deletion action taken"

	// Tests if the provided object IDs exists.
	for _, objectID := range opts.ObjectIDs {
		_, err := client.GetObject(client.NewApiGetObjectRequest(opts.Index, objectID))
		if err != nil {
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
	if opts.NdeleteParams > 0 {
		res, err := client.SearchSingleIndex(
			client.NewApiSearchSingleIndexRequest(opts.Index).
				WithSearchParams(search.SearchParamsObjectAsSearchParams(deleteByToSearchParams(&opts.DeleteParams))),
		)
		if err != nil {
			return err
		}
		nbObjectsToDelete = nbObjectsToDelete + int(*res.NbHits)
		if res.ExhaustiveNbHits != nil && !*res.ExhaustiveNbHits {
			exactOrApproximate = "approximately"
		}
	}

	if nbObjectsToDelete == 0 {
		if _, err = fmt.Fprintf(opts.IO.Out, "%s No objects to delete. %s\n", cs.WarningIcon(), extra); err != nil {
			return err
		}
		return nil
	}

	objectNbMessage := fmt.Sprintf(
		"%s %s from %s",
		exactOrApproximate,
		utils.Pluralize(nbObjectsToDelete, "object"),
		opts.Index,
	)

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
		batchRes, err := client.DeleteObjects(opts.Index, opts.ObjectIDs)
		if err != nil {
			return err
		}
		for _, res := range batchRes {
			taskIDs = append(taskIDs, res.TaskID)
		}
	}

	// Delete the objects matching the filters
	if opts.NdeleteParams > 0 {
		res, err := client.DeleteBy(client.NewApiDeleteByRequest(opts.Index, &opts.DeleteParams))
		if err != nil {
			return err
		}

		taskIDs = append(taskIDs, res.TaskID)
	}

	// Wait for the tasks to complete
	if opts.Wait {
		opts.IO.StartProgressIndicatorWithLabel("Waiting for all of the deletion tasks to complete")
		for _, taskID := range taskIDs {
			_, err := client.WaitForTask(opts.Index, taskID)
			if err != nil {
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

func deleteByToSearchParams(input *search.DeleteByParams) *search.SearchParamsObject {
	return &search.SearchParamsObject{
		Filters:           input.Filters,
		FacetFilters:      input.FacetFilters,
		NumericFilters:    input.NumericFilters,
		TagFilters:        input.TagFilters,
		AroundLatLng:      input.AroundLatLng,
		AroundRadius:      input.AroundRadius,
		InsideBoundingBox: input.InsideBoundingBox,
		InsidePolygon:     input.InsidePolygon,
	}
}

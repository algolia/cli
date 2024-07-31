package delete

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Indices         []string
	DoConfirm       bool
	IncludeReplicas bool
}

// NewDeleteCmd creates and returns a delete command for indices
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4_SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "delete <index>",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cmdutil.V4_IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "deleteIndex",
		},
		Short: "Delete one or multiple indices",
		Long: heredoc.Doc(`
			Delete one or multiples indices.
			This command permanently removes one or multiple indices from your application, and removes their metadata and configured settings.
		`),
		Example: heredoc.Doc(`
			# Delete the index named "MOVIES"
			$ algolia indices delete MOVIES

      # Delete the index named "MOVIES" and its replicas
      $ algolia indices delete MOVIES --includeReplicas

			# Delete the index named "MOVIES", skipping the confirmation prompt
			$ algolia indices delete MOVIES -y

			# Delete multiple indices
			$ algolia indices delete MOVIES SERIES ANIMES
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indices = args

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

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().
		BoolVarP(&opts.IncludeReplicas, "includeReplicas", "r", false, "delete replica indices too")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	if opts.DoConfirm {
		var confirmed bool
		msg := "Are you sure you want to delete"
		if len(opts.Indices) == 1 {
			msg += " the index %q"
		} else {
			msg += " the indices %q"
		}
		if opts.IncludeReplicas {
			msg += " and their replicas"
		}
		msg += "?"
		err := prompt.Confirm(fmt.Sprintf(msg, strings.Join(opts.Indices, ", ")), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	var deletedIndices []string
	for _, index := range opts.Indices {
		settings, err := client.GetSettings(client.NewApiGetSettingsRequest(index))
		if err != nil {
			er, ok := err.(*search.APIError)
			if ok && er.Status == 404 {
				return fmt.Errorf("Index '%s' does not exist\n", index)
			}
			return err
		}
		// Is it a replica index?
		if primary := settings.GetPrimary(); len(primary) > 0 {
			primarySettings, err := client.GetSettings(client.NewApiGetSettingsRequest(primary))
			oneRemoved := removeElement(primarySettings.GetReplicas(), index)
			// Detach replica index from primary index (keeping other replicas)
			res, err := client.SetSettings(
				client.NewApiSetSettingsRequest(
					primary,
					&search.IndexSettings{
						Replicas: oneRemoved,
					},
				),
			)
			if err != nil {
				return err
			}
			// Wait until the settings change has been made
			_, err = client.WaitForTask(primary, res.TaskID)
			if err != nil {
				return err
			}
		}

		deletedRes, err := client.DeleteIndex(client.NewApiDeleteIndexRequest(index))
		if err != nil {
			return fmt.Errorf("failed to delete index %s: %w", index, err)
		}

		deletedIndices = append(deletedIndices, index)
		if opts.IncludeReplicas && settings.HasReplicas() {
			client.WaitForTask(index, deletedRes.TaskID)
			// Construct batch request for deleting replicas of this index
			var requests []search.MultipleBatchRequest
			replicas := settings.GetReplicas()
			for _, index := range replicas {
				requests = append(
					requests,
					*search.NewMultipleBatchRequest(search.ACTION_DELETE, map[string]any{"indexName": index}, index),
				)
			}
			_, err := client.MultipleBatch(
				client.NewApiMultipleBatchRequest(search.NewBatchParams(requests)),
			)
			if err != nil {
				return err
			}
			deletedIndices = append(deletedIndices, replicas...)
		}
	}
	whatWasDeleted := "index"
	if len(deletedIndices) > 1 {
		whatWasDeleted = "indices"
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Deleted %s %s\n",
			cs.SuccessIcon(),
			whatWasDeleted,
			strings.Join(deletedIndices, ", "),
		)
	}
	return nil
}

// removeElement removes one element from a slice
func removeElement(slice []string, element string) []string {
	index := -1
	for i, v := range slice {
		if v == element || v == virtual(element) {
			index = i
			break
		}
	}

	if index == -1 {
		// Element not found, return the original slice
		return slice
	}

	return append(slice[:index], slice[index+1:]...)
}

// virtual wraps a string in the `virtual` modifier
func virtual(s string) string {
	return "virtual(" + s + ")"
}

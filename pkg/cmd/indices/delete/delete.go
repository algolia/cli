package delete

import (
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
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indices         []string
	DoConfirm       bool
	IncludeReplicas bool
}

// NewDeleteCmd creates and returns a delete command for indices
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "delete <index>",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Delete one or multiple indices",
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

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().BoolVarP(&opts.IncludeReplicas, "includeReplicas", "r", false, "delete replica indices too")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	if opts.DoConfirm {
		var confirmed bool
		msg := "Are you sure you want to delete the indices %q?"
		if opts.IncludeReplicas {
			msg = "Are you sure you want to delete the indices %q including their replicas?"
		}
		err := prompt.Confirm(fmt.Sprintf(msg, strings.Join(opts.Indices, ", ")), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	indices := make([]*search.Index, 0, len(opts.Indices))
	for _, indexName := range opts.Indices {
		index := client.InitIndex(indexName)
		exists, err := index.Exists()
		if err != nil || !exists {
			return fmt.Errorf("index %q does not exist", indexName)
		}
		indices = append(indices, index)

		if opts.IncludeReplicas {
			settings, err := index.GetSettings()

			if err != nil {
				return fmt.Errorf("can't get settings of index %q: %w", indexName, err)
			}

			replicas := settings.Replicas
			for _, replicaName := range replicas.Get() {
				replica := client.InitIndex(replicaName)
				indices = append(indices, replica)
			}
		}
	}

	for _, index := range indices {
		var mustWait bool

		if opts.IncludeReplicas {
			settings, err := index.GetSettings()
			if err != nil {
				return fmt.Errorf("failed to get settings of index %q: %w", index.GetName(), err)
			}
			if len(settings.Replicas.Get()) > 0 {
				mustWait = true
			}
		}

		res, err := index.Delete()

		// Otherwise, the replica indices might not be 'fully detached' yet.
		if mustWait {
			_ = res.Wait()
		}

		if err != nil {
			opts.IO.StartProgressIndicatorWithLabel(fmt.Sprint("Deleting replica index ", index.GetName()))
			err := deleteReplicaIndex(client, index)
			opts.IO.StopProgressIndicator()
			if err != nil {
				return fmt.Errorf("failed to delete index %q: %w", index.GetName(), err)
			}
		}
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted indices %s\n", cs.SuccessIcon(), strings.Join(opts.Indices, ", "))
	}

	return nil
}

// Delete a replica index.
func deleteReplicaIndex(client *search.Client, replicaIndex *search.Index) error {
	replicaName := replicaIndex.GetName()
	primaryName, err := findPrimaryIndex(replicaIndex)
	if err != nil {
		return fmt.Errorf("can't find primary index for %q: %w", replicaName, err)
	}

	err = detachReplicaIndex(replicaName, primaryName, client)
	if err != nil {
		return fmt.Errorf("can't unlink replica index %s from primary index %s: %w", replicaName, primaryName, err)
	}

	_, err = replicaIndex.Delete()
	if err != nil {
		return fmt.Errorf("can't delete replica index %q: %w", replicaName, err)
	}

	return nil
}

// Find the primary index of a replica index
func findPrimaryIndex(replicaIndex *search.Index) (string, error) {
	replicaName := replicaIndex.GetName()
	settings, err := replicaIndex.GetSettings()

	if err != nil {
		return "", fmt.Errorf("can't get settings of replica index %q: %w", replicaName, err)
	}

	primary := settings.Primary
	if primary == nil {
		return "", fmt.Errorf("index %s doesn't have a primary", replicaName)
	}

	return primary.Get(), nil
}

// Remove replica from `replicas` settings of the primary index
func detachReplicaIndex(replicaName string, primaryName string, client *search.Client) error {
	primaryIndex := client.InitIndex(primaryName)
	settings, err := primaryIndex.GetSettings()

	if err != nil {
		return fmt.Errorf("can't get settings of primary index %q: %w", primaryName, err)
	}

	replicas := settings.Replicas.Get()
	indexOfReplica := findIndex(replicas, replicaName)

	// Delete the replica at position `indexOfReplica` from the array
	replicas = append(replicas[:indexOfReplica], replicas[indexOfReplica+1:]...)

	res, err := primaryIndex.SetSettings(
		search.Settings{
			Replicas: opt.Replicas(replicas...),
		},
	)

	if err != nil {
		return fmt.Errorf("can't update settings of index %q: %w", primaryName, err)
	}

	// Wait until the settings are updated, else a subsequent `delete` will fail.
	_ = res.Wait()
	return nil
}

// Find the index of the string `target` in the array `arr`
func findIndex(arr []string, target string) int {
	for i, v := range arr {
		if v == target {
			return i
		}
	}
	return -1
}

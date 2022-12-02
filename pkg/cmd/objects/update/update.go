package updateObjects

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type UpdateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Index                          string
	FilePath                       string
	CreateIfNotExists              bool
	AutoGenerateObjectIDIfNotExist bool

	DoConfirm bool
}

// NewUpdateCmd creates and returns an update command for index objects
func NewUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	var confirm bool

	opts := &UpdateOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "update <index> -F <file>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Update objects from file to the specified index",
		Long: heredoc.Doc(`
			Update objects from JSON file to the specified index.
			The JSON file must contains an array of valid objects
		`),
		Example: heredoc.Doc(`
			# Update objects from the "objects.json" file to the "TEST_PRODUCTS" index
			$ algolia objects update TEST_PRODUCTS -F objects.json

			# Update objects (create if not exists) from the "objects.json" file to the "TEST_PRODUCTS" index
			$ algolia objects update TEST_PRODUCTS -F objects.json --create-if-not-exists

			# Update objects (create and auto generate objectID if not exists) from the "objects.json" file to the "TEST_PRODUCTS" index
			$ algolia objects update TEST_PRODUCTS -F objects.json --create-if-not-exists --auto-generate-object-id-if-not-exist
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cs := opts.IO.ColorScheme()
			opts.Index = args[0]

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if opts.AutoGenerateObjectIDIfNotExist && !opts.CreateIfNotExists {
				return fmt.Errorf("%s --auto-generate-object-id-if-not-exist flag can only be set if --create-if-not-exists flag is set", cs.FailureIcon())
			}

			if runF != nil {
				return runF(opts)
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")
	cmd.Flags().StringVarP(&opts.FilePath, "file", "F", "", "Directory path of the JSON that contains updated objects")
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().BoolVarP(&opts.CreateIfNotExists, "create-if-not-exists", "c", false, "Updating a nonexistent object will create a new object")
	cmd.Flags().BoolVarP(&opts.AutoGenerateObjectIDIfNotExist, "auto-generate-object-id-if-not-exist", "a", false, "The engine assigns an objectID to any object without objectID")

	return cmd
}

func runDeleteCmd(opts *UpdateOptions) error {
	cs := opts.IO.ColorScheme()

	jsonFile, err := os.Open(opts.FilePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	var objectsToUpdate ObjectsToUpdate
	err = json.Unmarshal(byteValue, &objectsToUpdate)
	if err != nil {
		return err
	}

	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	index := client.InitIndex(opts.Index)

	_, err = index.PartialUpdateObjects(objectsToUpdate,
		opt.CreateIfNotExists(opts.CreateIfNotExists), opt.AutoGenerateObjectIDIfNotExist(opts.AutoGenerateObjectIDIfNotExist))
	if err != nil {
		return err
	}

	fmt.Printf("%s %d objects successfully updated", cs.SuccessIcon(), len(objectsToUpdate))
	return nil
}

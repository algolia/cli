package add

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/provisionning"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

type AddOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	ProvisionningClient func() (*provisionning.Client, error)

	AppID  string
	Emails []string

	PrintFlags *cmdutil.PrintFlags
}

// NewAddCmd creates and returns a add command for Applications Members.
func NewAddCmd(f *cmdutil.Factory, runF func(*AddOptions) error) *cobra.Command {
	opts := &AddOptions{
		IO:                  f.IOStreams,
		Config:              f.Config,
		ProvisionningClient: f.ProvisionningClient,
	}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new member to an application",
		Long: heredoc.Doc(`
			Add a new member to an application.
		`),
		Example: heredoc.Doc(`
			# Add one new member to the application ABCD1234
			$ algolia applications members add --app-id ABCD1234 --email john.doe@algolia.com

			# Add multiple new members to the application ABCD1234
			$ algolia applications members add --app-id ABCD1234 --email john.doe@algolia.com,jane.doe@algolia.com
		`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runAddCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.AppID, "app-id", "a", "", "Application ID")
	_ = cmd.MarkFlagRequired("app-id")
	cmd.Flags().StringSliceVarP(&opts.Emails, "email", "e", nil, "Email address of the member to add")
	_ = cmd.MarkFlagRequired("email")

	return cmd
}

// runAddCmd executes the create command
func runAddCmd(opts *AddOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.ProvisionningClient()
	if err != nil {
		return err
	}

	emailsNb := len(opts.Emails)
	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Adding %d %s to application %s", emailsNb, utils.Pluralize(emailsNb, "member"), cs.Bold(opts.AppID)))
	errors := make([]error, 0)
	for _, email := range opts.Emails {
		err = client.AddApplicationMember(opts.AppID, email)
		if err != nil {
			errors = append(errors, err)
		}
	}
	opts.IO.StopProgressIndicator()

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintf(opts.IO.ErrOut, "%s\n", cs.Red(err.Error()))
		}
		return fmt.Errorf("failed to add %d member(s) to application %s", len(errors), opts.AppID)
	}

	fmt.Fprintf(opts.IO.Out, "%s %s added to application %s\n", cs.SuccessIcon(), utils.Pluralize(emailsNb, "member"), cs.Bold(opts.AppID))

	return nil
}

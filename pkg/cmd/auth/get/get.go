package get

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type GetOptions struct {
	IO                  *iostreams.IOStreams
	LoadToken           func() *auth.StoredToken
	PrintFlags          *cmdutil.PrintFlags
	NewDashboardClient  func(clientID string) *dashboard.Client
	EnsureAuthenticated func(io *iostreams.IOStreams, client *dashboard.Client) (string, error)

	WithAccessToken bool
}

type Identity struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	Name   string `json:"name,omitempty"`
	Token  string `json:"token,omitempty"`
}

func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:         f.IOStreams,
		LoadToken:  auth.LoadToken,
		PrintFlags: cmdutil.NewPrintFlags().WithDefaultOutput("json"),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
		EnsureAuthenticated: auth.EnsureAuthenticated,
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the authenticated user",
		Long: heredoc.Doc(`
			Get the identity of the authenticated user from the OAuth
			credentials stored in the local keychain.
		`),
		Example: heredoc.Doc(`
			# Get the authenticated user
			$ algolia auth get

			# Include the access token in the output
			$ algolia auth get --with-access-token
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCmd(opts)
		},
	}

	cmd.Flags().
		BoolVar(&opts.WithAccessToken, "with-access-token", false, "Include the OAuth access token in the output")
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client := opts.NewDashboardClient(auth.OAuthClientID())
	if _, err := opts.EnsureAuthenticated(opts.IO, client); err != nil {
		return err
	}

	stored := opts.LoadToken()

	identity := Identity{
		UserID: stored.UserID,
		Email:  stored.Email,
		Name:   stored.Name,
	}
	if opts.WithAccessToken {
		identity.Token = stored.AccessToken
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	return p.Print(opts.IO, identity)
}

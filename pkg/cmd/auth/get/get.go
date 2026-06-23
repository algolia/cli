package get

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// GetOptions represents the options for the get command.
type GetOptions struct {
	IO *iostreams.IOStreams

	LoadToken func() *auth.StoredToken

	PrintFlags *cmdutil.PrintFlags
}

// Identity is the authenticated user, without any token information.
type Identity struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	Name   string `json:"name,omitempty"`
}

// NewGetCmd returns a new instance of the get command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:         f.IOStreams,
		LoadToken:  auth.LoadToken,
		PrintFlags: cmdutil.NewPrintFlags().WithDefaultOutput("json"),
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
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// runGetCmd runs the get command.
func runGetCmd(opts *GetOptions) error {
	stored := opts.LoadToken()
	if stored == nil {
		return fmt.Errorf("you are not logged in — run `algolia auth login` first")
	}

	if stored.IsExpired() {
		return fmt.Errorf("your session has expired — run `algolia auth login` again")
	}

	identity := Identity{
		UserID: stored.UserID,
		Email:  stored.Email,
		Name:   stored.Name,
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	return p.Print(opts.IO, identity)
}

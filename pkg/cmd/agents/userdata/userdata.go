package userdata

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
)

// NewUserDataCmd is the parent for `algolia agents user-data <verb>`.
// Backs the GDPR right-to-access (get) and right-to-be-forgotten (delete).
func NewUserDataCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user-data",
		Aliases: []string{"userdata"},
		Short:   "Read/erase per-user-token conversation+memory data (GDPR)",
		Long: heredoc.Doc(`
			Inspect or erase all conversations and memories tied to a
			user token (X-Algolia-Secure-User-Token). Intended for
			GDPR data-subject requests; use with care — delete is
			irreversible and erases across every agent in the app.
		`),
	}
	cmd.AddCommand(newGetCmd(f, nil))
	cmd.AddCommand(newDeleteCmd(f, nil))
	return cmd
}

func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// rejectSlashMsg is the user-facing reason for rejecting tokens that
// contain "/". The CLI escapes the path segment correctly, but the
// staging gateway decodes %2F before routing — verified live, see
// Anya QA report F-1. Blocking with a clear message beats a misleading
// 404 from the backend.
const rejectSlashMsg = `user-token contains "/", which the Agent Studio gateway misroutes (decodes "%2F" before path matching, yielding a misleading 404). Use a token without "/" or contact support if a slash-bearing token is required.`

package providers

import (
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
)

var nowFn = time.Now

// NewProvidersCmd is the parent for `algolia agents providers <verb>`.
func NewProvidersCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "Manage Agent Studio LLM providers",
		Long: heredoc.Doc(`
			Manage LLM provider authentications (one per OpenAI /
			Anthropic / Azure / etc. account) used by Agent Studio agents.

			Agents reference a provider by ID via their "providerId"
			field; an agent without a working provider 4xxs at completion
			time. Use this group to bootstrap or audit your providers
			without leaving the terminal.
		`),
	}

	cmd.AddCommand(newListCmd(f, nil))
	cmd.AddCommand(newGetCmd(f, nil))
	cmd.AddCommand(newCreateCmd(f, nil))
	cmd.AddCommand(newUpdateCmd(f, nil))
	cmd.AddCommand(newDeleteCmd(f, nil))
	cmd.AddCommand(newModelsCmd(f, nil))
	cmd.AddCommand(newDefaultsCmd(f, nil))
	return cmd
}

func relTimeOrDash(t *time.Time, now time.Time) string {
	if t == nil || t.IsZero() {
		return "-"
	}
	return humanize.RelTime(now, *t, "from now", "ago")
}

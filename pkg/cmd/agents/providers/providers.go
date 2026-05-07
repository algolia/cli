package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
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

// readBody is the file-or-stdin → validated JSON pipeline used by create/update.
func readBody(file string, ios *iostreams.IOStreams) ([]byte, error) {
	body, err := cmdutil.ReadFile(file, ios.In)
	if err != nil {
		return nil, fmt.Errorf("failed to read provider body from %s: %w", shared.SourceLabel(file), err)
	}
	body = shared.TrimUTF8BOM(body)
	if !json.Valid(body) {
		return nil, cmdutil.FlagErrorf("provider body in %s is not valid JSON", shared.SourceLabel(file))
	}
	return body, nil
}

func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func relTimeOrDash(t *time.Time, now time.Time) string {
	if t == nil || t.IsZero() {
		return "-"
	}
	return humanize.RelTime(now, *t, "from now", "ago")
}

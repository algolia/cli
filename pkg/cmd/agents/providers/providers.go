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

// nowFn is overridable for deterministic time-based output in tests.
var nowFn = time.Now

// NewProvidersCmd is the parent for `algolia agents providers <verb>`.
//
// Provider records are LLM-credential bindings (one per
// OpenAI/Anthropic/Azure/etc. account the app talks to). Agents
// reference a provider by ID via their `providerId` field.
//
// Naming: `providers` not `provider` to match every other listable
// resource in the CLI tree (`apikeys`, `objects`, `rules`, ...).
//
// Verbs are split per file (list.go, get.go, create.go, update.go,
// delete.go, models.go) within this package — same package keeps
// the masking helper and the small formatters below accessible
// without exporting them. Cross-cutting helpers live here.
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
	return cmd
}

// readBody is the shared file-or-stdin → validated JSON pipeline used
// by create/update. Centralised so any future verb that takes a JSON
// body (e.g., a hypothetical bulk-create) gets identical UX:
// "<source> is not valid JSON" with the same source-label conventions.
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

// ctxOrBackground promotes a possibly-nil command Context to
// context.Background. Cobra always supplies one in production, but
// table-test invocations of run* helpers occasionally don't, and
// returning context.Background here keeps test setup boilerplate down.
func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// relTimeOrDash is the standard "relative time, or dash for unset"
// formatter used by the table renderer in list.go (and reused by any
// future verb whose output table includes timestamps).
func relTimeOrDash(t *time.Time, now time.Time) string {
	if t == nil || t.IsZero() {
		return "-"
	}
	return humanize.RelTime(now, *t, "from now", "ago")
}

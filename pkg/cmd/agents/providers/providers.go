package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
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

// ---------------------------------------------------------------------
// list
// ---------------------------------------------------------------------

type ListOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	Page    int
	PerPage int
	Show    bool
}

func newListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List configured LLM providers",
		Long: heredoc.Doc(`
			List provider authentications on Agent Studio for the active
			application.

			By default, structured output (--output json) masks the
			"apiKey" field (and similar secrets) with a "***" prefix.
			Pass --show-secret to render values verbatim — useful for
			scripted exports, dangerous in shared logs.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers list
			$ algolia agents providers list --output json --show-secret
			$ algolia agents providers list --page 2 --per-page 25
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			if runF != nil {
				return runF(opts)
			}
			return runListCmd(opts)
		},
	}

	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (1-indexed; 0 = backend default)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Items per page (0 = backend default, currently 10)")
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields (apiKey, ...) verbatim instead of masking")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Fetching providers")
	res, err := client.ListProviders(ctx, agentstudio.ListProvidersParams{
		Page:  opts.Page,
		Limit: opts.PerPage,
	})
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if !opts.Show {
		// Masking happens in-place on the slice we'll print. Doesn't
		// touch the cached *Provider on the backend; only what we hand
		// to the printer.
		for i := range res.Data {
			res.Data[i].Input = MaskInput(res.Data[i].Input)
		}
	}

	if opts.PrintFlags.HasStructuredOutput() {
		return opts.PrintFlags.Print(opts.IO, res)
	}

	now := nowFn()
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("NAME", nil, nil)
		table.AddField("PROVIDER", nil, nil)
		table.AddField("LAST USED", nil, nil)
		table.AddField("UPDATED", nil, nil)
		table.EndRow()
	}
	for _, p := range res.Data {
		table.AddField(p.ID, nil, nil)
		table.AddField(p.Name, nil, nil)
		table.AddField(p.ProviderName, nil, nil)
		table.AddField(relTimeOrDash(p.LastUsedAt, now), nil, nil)
		table.AddField(relTimeOrDash(&p.UpdatedAt, now), nil, nil)
		table.EndRow()
	}
	if err := table.Render(); err != nil {
		return err
	}
	if table.IsTTY() {
		fmt.Fprintf(opts.IO.Out,
			"\n%d provider(s) — page %d of %d (total %d).\n",
			len(res.Data),
			res.Pagination.Page, res.Pagination.TotalPages, res.Pagination.TotalCount)
	}
	return nil
}

// ---------------------------------------------------------------------
// get
// ---------------------------------------------------------------------

type GetOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	ProviderID string
	Show       bool
}

func newGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <provider-id>",
		Short: "Get an LLM provider authentication by ID",
		Long: heredoc.Doc(`
			Fetch a provider authentication by ID. By default, secret
			fields ("apiKey") are masked. Pass --show-secret to reveal.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers get 11111111-1111-1111-1111-111111111111
			$ algolia agents providers get <id> --show-secret --output json
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]
			opts.Ctx = cmd.Context()
			if opts.ProviderID == "" {
				return cmdutil.FlagErrorf("provider-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runGetCmd(opts)
		},
	}

	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields (apiKey, ...) verbatim instead of masking")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Fetching provider")
	p, err := client.GetProvider(ctx, opts.ProviderID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if !opts.Show {
		p.Input = MaskInput(p.Input)
	}
	return opts.PrintFlags.Print(opts.IO, p)
}

// ---------------------------------------------------------------------
// create
// ---------------------------------------------------------------------

type CreateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	File          string
	DryRun        bool
	Show          bool
	OutputChanged bool
}

func newCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "create -F <file>",
		Short: "Create an LLM provider authentication from a JSON file",
		Long: heredoc.Doc(`
			Create a provider authentication from a JSON file describing
			the ProviderAuthenticationCreate body (name, providerName,
			input). The "input" subobject's shape varies per providerName:

			  - openai / anthropic: {apiKey, baseUrl?}
			  - azure_openai:       {apiKey, azureEndpoint, azureDeployment, apiVersion?}
			  - openai_compatible:  {apiKey, baseUrl, defaultModel}
			  - google_genai / deepseek: {apiKey}

			The file is sent verbatim; field-level validation is the
			backend's job (a 4xx surfaces with the structured detail).

			Use --dry-run to preview the request without sending.

			By default the created provider in the success response is
			masked. Pass --show-secret to render the apiKey verbatim.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers create -F openai-prod.json
			$ cat spec.json | algolia agents providers create -F -
			$ algolia agents providers create -F spec.json --dry-run
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")
			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "JSON file with the provider body (use \"-\" for stdin)")
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Validate and print the resolved request body without calling the API")
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields verbatim in the success response")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	body, err := readBody(opts.File, opts.IO)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, opts.PrintFlags, opts.OutputChanged,
			"create_provider", "POST /1/providers", opts.File, body, nil)
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Creating provider")
	p, err := client.CreateProvider(ctx, json.RawMessage(body))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if !opts.Show {
		p.Input = MaskInput(p.Input)
	}
	return opts.PrintFlags.Print(opts.IO, p)
}

// ---------------------------------------------------------------------
// update
// ---------------------------------------------------------------------

type UpdateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	ProviderID    string
	File          string
	DryRun        bool
	Show          bool
	OutputChanged bool
}

func newUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "update <provider-id> -F <file>",
		Short: "Patch an LLM provider authentication from a JSON file",
		Long: heredoc.Doc(`
			Patch a provider authentication. PATCH semantics: only the
			fields in the file are updated. Pass {"name":"new-name"} to
			rename, or {"input":{"apiKey":"sk-NEW"}} to rotate the key.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers update <id> -F rename.json
			$ algolia agents providers update <id> -F rotate.json --dry-run
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")
			if opts.ProviderID == "" {
				return cmdutil.FlagErrorf("provider-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "JSON file with the provider patch body (use \"-\" for stdin)")
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Validate and print the resolved request body without calling the API")
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields verbatim in the success response")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	body, err := readBody(opts.File, opts.IO)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, opts.PrintFlags, opts.OutputChanged,
			"update_provider",
			fmt.Sprintf("PATCH /1/providers/%s", opts.ProviderID),
			opts.File, body, map[string]any{"providerId": opts.ProviderID})
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Updating provider")
	p, err := client.UpdateProvider(ctx, opts.ProviderID, json.RawMessage(body))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if !opts.Show {
		p.Input = MaskInput(p.Input)
	}
	return opts.PrintFlags.Print(opts.IO, p)
}

// ---------------------------------------------------------------------
// delete
// ---------------------------------------------------------------------

type DeleteOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	ProviderID string
	DryRun     bool
	DoConfirm  bool
}

func newDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <provider-id> [--confirm]",
		Short: "Delete an LLM provider authentication",
		Long: heredoc.Doc(`
			Delete a provider authentication. The backend may 409 if any
			agent still references this provider; the structured detail
			surfaces verbatim. Detach affected agents first (or update
			them to point at a different provider) before retrying.

			Like "agents delete", interactive use prompts to confirm and
			non-interactive use requires --confirm. Use --dry-run to
			preview without deleting.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers delete <id>           # interactive
			$ algolia agents providers delete <id> -y        # CI
			$ algolia agents providers delete <id> --dry-run # preview
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]
			opts.Ctx = cmd.Context()
			if opts.ProviderID == "" {
				return cmdutil.FlagErrorf("provider-id must not be empty")
			}
			if !confirm && !opts.DryRun {
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

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Fetch and preview the provider without deleting it")
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	// Pre-fetch so the prompt + dry-run output show name+providerName,
	// matching `agents delete`'s contract.
	p, err := client.GetProvider(ctx, opts.ProviderID)
	if err != nil {
		return err
	}

	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out, "Dry run: would DELETE /1/providers/%s\n", opts.ProviderID)
		fmt.Fprintf(opts.IO.Out, "  name:     %s\n", p.Name)
		fmt.Fprintf(opts.IO.Out, "  provider: %s\n", p.ProviderName)
		return nil
	}

	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(
			fmt.Sprintf("Delete provider %q (%s, %s)?", p.Name, p.ProviderName, opts.ProviderID),
			&confirmed,
		)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	opts.IO.StartProgressIndicatorWithLabel("Deleting provider")
	err = client.DeleteProvider(ctx, opts.ProviderID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted provider %s\n", cs.SuccessIcon(), opts.ProviderID)
	}
	return nil
}

// ---------------------------------------------------------------------
// models
// ---------------------------------------------------------------------

type ModelsOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	ProviderID string
}

func newModelsCmd(f *cmdutil.Factory, runF func(*ModelsOptions) error) *cobra.Command {
	opts := &ModelsOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
	}

	cmd := &cobra.Command{
		Use:   "models [--provider-id <id>]",
		Short: "List models available per provider type, or for a configured provider",
		Long: heredoc.Doc(`
			Two routes, one command:

			  - No --provider-id: GET /1/providers/models
			    Returns the static catalog of models per provider TYPE
			    (openai → [gpt-4o, gpt-4o-mini, ...], anthropic → ...).
			    Use this BEFORE creating a provider to pick a model.

			  - --provider-id <id>: GET /1/providers/<id>/models
			    Returns models exposed by a specific configured provider,
			    which can include account-specific ones (OpenAI fine-tunes,
			    Azure deployments, etc.).
		`),
		Example: heredoc.Doc(`
			# What can I configure?
			$ algolia agents providers models

			# What does my OpenAI provider actually expose (incl. fine-tunes)?
			$ algolia agents providers models --provider-id <id>
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			// An explicit empty --provider-id is almost certainly a
			// scripting bug ("$PROV_ID" before assignment); refuse it
			// instead of silently falling back to the catalog route.
			if cmd.Flags().Changed("provider-id") && opts.ProviderID == "" {
				return cmdutil.FlagErrorf("--provider-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runModelsCmd(opts)
		},
	}

	cmd.Flags().
		StringVar(&opts.ProviderID, "provider-id", "", "Use the configured-provider route instead of the static catalog")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runModelsCmd(opts *ModelsOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	if opts.ProviderID == "" {
		return runListCatalog(opts, client, ctx)
	}
	return runListForProvider(opts, client, ctx)
}

func runListCatalog(opts *ModelsOptions, client *agentstudio.Client, ctx context.Context) error {
	opts.IO.StartProgressIndicatorWithLabel("Fetching provider model catalog")
	models, err := client.ListProviderModels(ctx)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.HasStructuredOutput() {
		return opts.PrintFlags.Print(opts.IO, models)
	}

	// Sorted, table-friendly: provider name then model name.
	providerNames := make([]string, 0, len(models))
	for k := range models {
		providerNames = append(providerNames, k)
	}
	sort.Strings(providerNames)

	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("PROVIDER", nil, nil)
		table.AddField("MODEL", nil, nil)
		table.EndRow()
	}
	for _, prov := range providerNames {
		ms := models[prov]
		sort.Strings(ms)
		for _, m := range ms {
			table.AddField(prov, nil, nil)
			table.AddField(m, nil, nil)
			table.EndRow()
		}
	}
	return table.Render()
}

func runListForProvider(opts *ModelsOptions, client *agentstudio.Client, ctx context.Context) error {
	opts.IO.StartProgressIndicatorWithLabel("Fetching models for provider")
	raw, err := client.ListModelsForProvider(ctx, opts.ProviderID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	// Spec leaves the response shape unspecified. Empirically it's a
	// JSON array of strings, but we don't pin that — pretty-print
	// whatever arrives and let the user pipe to jq.
	if opts.PrintFlags.HasStructuredOutput() {
		// Round-trip through a generic decode so --output json
		// re-formats consistently rather than emitting the backend's
		// exact whitespace.
		var anyV any
		if err := json.Unmarshal(raw, &anyV); err != nil {
			return fmt.Errorf("decode models response: %w", err)
		}
		return opts.PrintFlags.Print(opts.IO, anyV)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, raw, "", "  "); err != nil {
		_, _ = opts.IO.Out.Write(raw)
		return nil
	}
	_, _ = opts.IO.Out.Write(pretty.Bytes())
	_, _ = opts.IO.Out.Write([]byte("\n"))
	return nil
}

// ---------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------

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

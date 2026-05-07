# AGENTS.md

Guidance for coding agents working in `github.com/algolia/cli`.

## Scope

- Applies to the whole repository.
- Prefer small, local changes over broad refactors.
- Follow existing Go + Cobra CLI patterns.

## Repository Facts

- Language: Go.
- Module: `github.com/algolia/cli`.
- Go version: `1.23.0`; toolchain: `go1.23.4`.
- Main binary entrypoint: `cmd/algolia/main.go`.
- Docs generator entrypoint: `cmd/docs/main.go`.
- Command tree: `pkg/cmd/...`.

## Cursor / Copilot Rules

- No `.cursor/rules/` directory was found.
- No `.cursorrules` file was found.
- No `.github/copilot-instructions.md` file was found.
- Use this file plus the existing codebase as the repo-specific source of truth.

## Tooling

- Preferred toolchain is listed in `devbox.json`.
- Common tools expected here: `go`, `task`, `golangci-lint`, `gofumpt`, `golines`, `gh`, `curl`.
- E2E tests require `ALGOLIA_APPLICATION_ID` and `ALGOLIA_API_KEY` in the environment or root `.env`.

## Build Commands

Preferred:

```sh
task build
```

```sh
go generate ./...
go build -ldflags "-s -w -X=github.com/algolia/cli/pkg/version.Version=main" -o algolia cmd/algolia/main.go
go build -v ./...
```

- `task build` runs generation first.
- CI also checks `go build -v ./...`.

## Test Commands

All unit tests:

```sh
task test
go test ./...
go test ./... -p 1
```

Single package / single test:

```sh
go test ./pkg/cmd/search
go test ./pkg/cmd/apikeys/list -run Test_runListCmd
go test ./pkg/cmd/apikeys/list -run 'Test_runListCmd/list_tty'
go test ./... -run Test_runListCmd
go test ./pkg/cmd/apikeys/list -run Test_runListCmd -v -count=1
```

E2E tests:

```sh
task e2e
go test ./e2e -tags=e2e
go test ./e2e -tags=e2e -run TestIndices
go test ./e2e -tags=e2e -run TestAgentReady -v
```

- E2E uses `github.com/cli/go-internal/testscript`.
- E2E makes real Algolia API requests.
- Keep E2E runs narrow when possible.

## Lint / Format Commands

```sh
task lint
golangci-lint run
task format
gofumpt -w pkg cmd test internal api e2e
golines -w pkg cmd test internal api e2e
```

- `gosec`
- `gofumpt`
- `stylecheck`

## Generation / Docs Commands

```sh
go generate ./...
go run ./cmd/docs --app_data-path tmp
go run ./cmd/docs --app_data-path tmp
```

Run generation when changing generated flags or API-spec-derived code.

## Fast Local Verification

For substantial changes, prefer this order:

```sh
task format
go test ./path/to/touched/package -run TestName
task lint
task build
```

Use narrower verification for small edits.

## Architecture Guidelines

- Add CLI commands under `pkg/cmd/<domain>`.
- Construct commands with `New...Cmd` functions.
- Keep option structs close to their commands.
- Inject dependencies via `*cmdutil.Factory`.
- Put shared command logic in focused helper packages, usually `pkg/cmdutil`.
- Keep docs-generation logic in `internal/docs` and `cmd/docs`.

## Agent Studio (`pkg/cmd/agents/...`, `api/agentstudio/`)

Top-level command group: `algolia agents`. Verbs: `list`, `get`, `create`, `update`, `delete`, `publish`, `unpublish`, `duplicate`, `try`, `run`. Sub-groups: `cache` (`invalidate`), `providers` (`list`/`get`/`create`/`update`/`delete`/`models`), `config` (`get`/`set`). Backend source of truth: `github.com/algolia/conversational-ai`.

**Naming note**: `try` (not `test`) — see "On `--dry-run`" below for why. All flat verbs are single-word lowercase to match the CLI-wide convention; no hyphenated subcommand names exist anywhere in the tree. Sub-groups (`cache`, `providers`, `config`, future `conversations` / `keys` / `domains`) read as noun-then-verb (`agents cache invalidate`, `agents providers create`) — also a single word per token. `providers` is plural to match every other listable resource group in the CLI tree (`apikeys`, `objects`, `rules`); `config` is singular because there's exactly one config record per app.

### File organisation

Both layers mirror the OpenAPI spec's tag boundaries — one source file (and one test file) per API tag. Adding a new resource is a single new file plus a single new test file:

```
api/agentstudio/
  client.go          ← Config, Client, NewClient, setHeaders, checkResponse, extractDetail, sentinelFor (infra only)
  agents.go          ← Agents tag (CRUD + lifecycle + InvalidateAgentCache)
  completions.go     ← Completions tag (Completions + boolToWire)
  providers.go       ← Providers tag (CRUD + 2 model-discovery routes)
  configuration.go   ← Configurations tag
  sse.go             ← cross-cutting helper (StreamEvent parser, both AI SDK v4 and v5 shapes)
  host.go            ← BaseURL resolution (profile → env → ldflag → cluster-proxy fallback)
  errors.go          ← APIError, sentinel errors
  types.go           ← shared response types (Agent, Provider, ApplicationConfig, PaginationMetadata, ProviderName constants)
```

Tests follow the same pairing — `agents_test.go`, `completions_test.go`, etc. `client_test.go` keeps the `newTestClient` harness plus infra-only tests (NewClient validation, `checkResponse` error mapping, ctx cancellation). The error-mapping suite uses `ListAgents` as a vehicle because it's the simplest GET; that's intentional and the file header calls it out.

For the cmd layer, top-level verbs each own a subpackage (`pkg/cmd/agents/{list,get,create,update,delete,publish,unpublish,duplicate,try,run}/`) — that pattern is inherited from the wider CLI codebase. Sub-groups (`cache/`, `providers/`, `config/`) keep all their verbs in **one package** but **one file per verb**:

```
pkg/cmd/agents/providers/
  providers.go   ← NewProvidersCmd parent + readBody/ctxOrBackground/relTimeOrDash helpers
  list.go        ← list verb
  get.go         ← get verb
  create.go      ← create verb
  update.go      ← update verb
  delete.go      ← delete verb
  models.go      ← models verb (handles both /1/providers/models and /1/providers/{id}/models)
  mask.go        ← MaskInput helper (shared by list/get/create/update)
  <verb>_test.go ← one test file per verb; helpers in providers_test.go
```

Same package keeps internal helpers (`MaskInput`, `readBody`, `ctxOrBackground`) accessible without exporting; per-file split keeps each verb under ~200 LOC. **Don't** promote sub-group verbs to per-verb subpackages — they're tightly coupled with shared internals; the directory churn buys nothing.

`conversations/` (5 verbs) follows the same pattern: `conversations.go` (parent + helpers), `list.go`, `get.go`, `delete.go`, `purge.go`, `export.go`, with paired `<verb>_test.go`. The agent ID is positional first arg on every verb (matches `agents publish/run/cache invalidate`).

For `cache/` (1 verb) and `config/` (2 verbs) the per-file split is unnecessary; they live in one file each.

### API client (`api/agentstudio/`)

- Auth: standard Algolia headers (`X-Algolia-Application-Id`, `X-Algolia-API-Key`). No bearer tokens. Comes from the active profile via `*cmdutil.Factory.AgentStudioClient`.
- Base URL resolution priority: per-profile `agent_studio_url` → env `ALGOLIA_AGENT_STUDIO_URL` → build-time `agentstudio.DefaultBaseURL` (set via `ldflags`, mirrors `dashboard.DefaultDashboardURL`) → cluster-proxy fallback `https://{appID}.algolia.net/agent-studio`. The cluster proxy already does region routing — don't add a `Region` field.
- Errors: `*APIError` with `StatusCode`, `Detail`, optional `Sentinel`. The detail extractor prefers structured FastAPI `detail[].msg` arrays over the generic `message` field — backends that return both pair them as `{"message":"Input is invalid, see detail/body:","detail":[{"msg":"..."}]}` and the structured form is the actionable one.
- `CreateAgent` / `UpdateAgent` accept `json.RawMessage` bodies on purpose. The backend's `AgentConfigCreate` schema is large, deeply validated, and evolves often. The CLI is a pass-through; the backend validates; our 422-detail surfacing makes errors actionable.
- `Completions(...)` returns the raw `*http.Response`. Caller checks `Content-Type` (`text/event-stream` → `ParseStream`; else copy verbatim). One method, two output shapes.
- `CompletionOptions.No*` fields (`NoCache`, `NoMemory`, `NoAnalytics`) are **inverted** from the backend's query polarity. Two reasons: the backend defaults all three to true (only the negative is interesting at the CLI), and `memory` in particular has an `anyOf [{const false}, {type null}]` schema — sending `memory=true` would 422. Therefore the wire form omits the param when the No* field is false, and sends `<param>=false` when true. Polarity is enforced end-to-end by `TestCompletions_QueryFlagsAndSecureUserToken` in `api/agentstudio/completions_test.go`.
- `CompletionOptions.SecureUserToken` populates the `X-Algolia-Secure-User-Token` header when non-empty. It carries a signed JWT scoping the conversation/memory/analytics partition to a specific end-user (see `rag/dependencies/secure_user_token.py` in the backend). Empty means no header — `X-Algolia-User-ID` fallback applies.
- `InvalidateAgentCache(id, before)` calls `DELETE /1/agents/{id}/cache?before=YYYY-MM-DD` (query omitted when `before` is empty). Date format validation is **deliberately not done client-side** — the backend's Pydantic parser is the source of truth, and our 422 surfacing turns malformed input into an actionable message verbatim. Mirroring the parser in Go would create silent skew.
- `ListProviders` / `GetProvider` / `CreateProvider` / `UpdateProvider` / `DeleteProvider` cover the `/1/providers` CRUD. Same `json.RawMessage` body convention as `CreateAgent` — the `input` subobject is a 6-way discriminated union (`openai` / `azure_openai` / `google_genai` / `deepseek` / `openai_compatible` / `anthropic`) with deeply-validated per-variant fields. Mirroring those structs in Go would lie about parity. The CLI passes through; the backend validates.
- `ListProviderModels()` returns `map[string][]string` — the static catalog of "what models can each provider type expose" (used as a discoverability primitive before creating a provider). `ListModelsForProvider(id)` returns `json.RawMessage` because the spec leaves the response shape unspecified — empirically `[]string` (incl. account-specific entries like OpenAI fine-tunes / Azure deployments) but we don't pin it.
- `GetConfiguration` / `UpdateConfiguration` cover `/1/configuration`. ACL is `logs`, **not** `settings` — the only field today (`maxRetentionDays`) governs log/conversation retention, hence the unusual ACL. Body shape kept as `json.RawMessage` for symmetry with the rest of the agents tree, even though the schema is a single int field — future fields will land here.
- `ListConversations` / `GetConversation` / `DeleteConversation` / `PurgeConversations` / `ExportConversations` cover `/1/agents/{id}/conversations*`. Note the per-agent scope — every endpoint takes `{agent_id}` in the path. `GetConversation` returns `json.RawMessage` because `ConversationFullResponse.messages` is a discriminated union over message roles (system/user/assistant/tool); `ListConversations` returns the typed `PaginatedConversationsResponse` because the lightweight base shape (no messages) is stable. `ListConversationsParams.FeedbackVote` is `*int` because nil = no filter while 0 (downvote) is a meaningful filter — pointer/nil distinction matters here. `ExportConversations` returns `json.RawMessage` because the spec leaves the export response shape unspecified.

### Conversations: `purge` vs `delete` (`agents conversations`)

Two distinct verbs share the underlying HTTP method (`DELETE`), but the blast radius is different by orders of magnitude:

- `delete <agent-id> <conv-id>` — surgical, one conversation. Mistype the conv ID and you nuke an unrelated conversation; same risk profile as `agents delete`.
- `purge <agent-id> --start-date | --end-date` — bulk. Wipes every conversation in the date range.

**Spec vs. reality on `purge`**: the OpenAPI spec marks both `startDate` and `endDate` as `required: false`, which reads as "dateless DELETE wipes everything." The live backend disagrees — it rejects dateless DELETE with `400 "At least one filter is required."` (caught during Phase 7 live vet against staging EU). The CLI mirrors backend reality, **not** the spec: at least one of `--start-date` / `--end-date` is required at the flag layer. If you genuinely want to wipe every conversation, pass an open-ended bound (`--start-date 1970-01-01` or `--end-date 9999-12-31`). The error message points users at the backend constraint so future spec drift is debuggable.

Both verbs flow through the same `--confirm` / non-TTY-refuses-without-it rule as `agents delete`. `--dry-run` previews the URL with its query string and labels the scope (`scope: between A and B`, `scope: from A onwards`, etc.).

### Streaming (`api/agentstudio/sse.go`)

The wire format is **not** standard SSE. Two protocols, both served as `text/event-stream`:

- **v5 (CLI default)**: standard SSE — `data: <json>\n\n`, `data: [DONE]` sentinel.
- **v4**: line-delimited bespoke — `<type-code>:<json>\n` per line, no terminator. Type codes: `0` = text, `9` = tool-call, `d` = finish-message, etc. (see `v4TypeNames` in `sse.go`).

`ParseStream` sniffs the line prefix and emits a normalized `StreamEvent{Type, Data, Raw}` for both. `compatibilityMode` is a **required** server-side query parameter — the CLI defaults to v5 and exposes `--compatibility v4|v5`.

Streaming output convention: TTY-attached stdout renders a flowing assistant transcript (text-deltas inline, tool calls/results as dim annotations, errors red). Non-TTY (piped, redirected) emits NDJSON, one `{"type":"...","data":{...}}` per line — stable contract for `jq -r 'select(.type=="text-delta") | .data.delta'` and similar pipelines. `--ndjson` forces NDJSON on a TTY for users who want machine output on screen. Branch lives in `shared.RenderCompletion`; if you add a new event type that should surface in the TTY render, extend `renderTTY`'s switch — leave NDJSON's verbatim pass-through alone (it's the wire-fidelity escape hatch).

### Completion runtime knobs (`agents try` / `agents run`)

Both commands expose the same set of completion-time flags, mapping directly to backend query params + headers:

| Flag | Wire | Default | Notes |
|---|---|---|---|
| `--no-stream` | `?stream=false` | stream | Buffered single-JSON response instead of SSE |
| `--compatibility v4\|v5` | `?compatibilityMode=ai-sdk-{4,5}` | v5 | **Required** server-side; CLI promotes empty → v5 |
| `--no-cache` | `?cache=false` | cache on | Bypasses backend completion cache for this call |
| `--no-memory` | `?memory=false` | memory on | Disables agent memory retrieval/write for this call |
| `--no-analytics` | `?analytics=false` | analytics on | Skips Agent Studio analytics for this call |
| `--secure-user-token <jwt>` | `X-Algolia-Secure-User-Token` header | (omitted) | Signed JWT, end-user scoping |

The flag set is intentionally duplicated across `try.go` and `run.go` rather than extracted into a `RegisterCompletionFlags` shared helper — there are exactly two consumers and the duplication is mechanical (8 lines per command). If a third consumer appears, extract following the "second use" rule (same as `PrintDryRun` / `NormalizeCompatibility`).

### Provider secret masking (`agents providers list/get/create/update`)

Provider responses include `apiKey` literally (per `OpenAIProviderInput-Output` et al. in the spec). Without masking, `--output json` writes raw API keys to stdout, which routinely lands in CI logs, terminal scrollback, and shared pastes. Convention:

- All four read/write commands mask `apiKey` to `"***"` by default in their success-path output.
- Pass `--show-secret` to render verbatim (scripted exports, debugging).
- Masking happens at the cmd layer (`pkg/cmd/agents/providers/mask.go:MaskInput`), not the client — the client returns the raw response so the cmd layer can opt out per-invocation.
- `--dry-run` does **not** mask: the user authored the file and is being shown what THEY are about to send. Hiding it would break the "what would be sent" contract.
- Three asterisks, no last-N preview. Goal is "impossible to copy by accident", not "allow last-4 lookup".
- `secretFieldNames` in `mask.go` is the closed set; today it's just `apiKey`. Extend alphabetically when new credential fields land. Same convention will land for Phase 8 (`agents keys`); when it does, lift `MaskInput` into a shared helper following the second-use rule.

### On `--dry-run`

Two distinct concepts share the name and they MUST NOT be conflated:

- **CLI-side `--dry-run` flag** on `agents create / update / delete / run`: standard CLI preview convention (matches `objects/update --dry-run`). Validates the request, prints what would be sent, makes no HTTP call. Two output modes:
  - **Human (default)**: `Dry run: would <METHOD> <PATH>` followed by the pretty-printed JSON body.
  - **Structured (only when `--output` is explicitly set)**: `{"action":"...","request":"...","source":"...","bytes":N,"body":<...>,"dryRun":true,...}`. Gate on `cmd.Flags().Changed("output")`, **not** on `PrintFlags.HasStructuredOutput()` — using the latter would let `WithDefaultOutput("json")` from the success path silently steal the human dry-run output.
- **Conversational-ai-side "dry-run" semantics** = run a real completion against a configuration that hasn't been persisted. The backend exposes this via the `agent_id="test"` route (`AgentTestConfiguration` in the body). The CLI surfaces it as the `agents try` command (NOT named `test` to avoid pytest/unit-test connotations and to read naturally as "experiment without commitment").

Why `agents try` has no `--dry-run` flag: the whole command IS the dry-run. Adding a CLI-level `--dry-run` on top would be "dry-run a dry-run." If you want to inspect the body the CLI would POST without calling the backend, build the JSON yourself — the wire shape is `{"messages":[...], "configuration":{...}}` and is documented in `pkg/cmd/agents/shared/completion.go:CompletionRequest`. The e2e in `e2e/testscripts/agents/dry-run.txtar` includes a regression assertion that `agents try --dry-run` is rejected as an unknown flag, so a future contributor doesn't add it back without seeing this rationale.

Shared helpers live in `pkg/cmd/agents/shared/` (`PrintDryRun`, `BuildMessages`, `ReadJSONFile`, `MarshalCompletionBody`, `RenderCompletion`, `NormalizeCompatibility`). Extract on second use, not pre-emptively.

### Telemetry

Existing `pkg/telemetry` model is **one event (`"Command Invoked"`) per invocation** from root, with `{command: cmd.CommandPath(), flags: [<changed flag names>]}`. That already attributes per-verb (`algolia agents create`) and surfaces `--dry-run` (it's in `flags`). Don't add bespoke per-verb telemetry events — it would diverge from convention for one feature only. Outcome (success/error) is a separate, all-commands refactor.

## Code Style

### Imports

- Use standard Go grouping: stdlib, third-party, local module.
- Let `gofumpt` handle ordering and spacing.
- Avoid aliases unless they prevent collisions or materially improve clarity.

### Formatting

- Run `gofumpt` on all modified Go files.
- Run `golines` if wrapping becomes awkward.
- Preserve existing multiline layout for structs, literals, and signatures.

### Types And Structs

- Prefer explicit structs for command options and helper state.
- Keep exported APIs minimal.
- Use `any` only where JSON-like dynamic values are genuinely needed.

### Naming

- Exported names: PascalCase.
- Unexported names: camelCase.
- Command constructors: `NewXCmd`.
- Command runners: `runXCmd`.
- Match surrounding test naming, commonly `Test_runXCmd`, `TestNewXCmd`, or `Test_Feature`.

### Cobra Conventions

- Use `RunE`, not `Run`, for command handlers.
- Validate args with `cobra.ExactArgs`, `cobra.MinimumNArgs`, or repo validators.
- Use `ValidArgsFunction` when completion helpers already exist.
- Reuse `cmdutil` helpers for usage text, print flags, JSON flags, and validations.
- Use heredocs for multiline examples and help text.

### Error Handling

- Return errors instead of exiting except in true entrypoints like `main()`.
- Wrap with `%w` when the original cause matters.
- Use plain `return err` when extra context adds no value.
- Prefer actionable CLI-facing error messages.
- Use `cmdutil.FlagErrorf` for invalid flag combinations and user input issues.
- Stop progress indicators on all error paths after starting them.

### I/O And UX

- Use factory-provided `IOStreams` for stdout, stderr, TTY checks, colors, and progress indicators.
- Keep non-TTY output deterministic and script-friendly.
- Use structured output helpers for commands that support `--output`.
- Preserve dry-run behavior: validate, summarize, and avoid side effects.

### Config And Clients

- Read config through `config.IConfig`.
- Acquire API clients from injected functions like `SearchClient` and `CrawlerClient`.
- Do not hardcode credentials, hosts, or profile logic.

### Testing Style

- Prefer table-driven tests for flags, output modes, and edge cases.
- Use `test.NewFactory(...)` and `test.Execute(...)` for command tests.
- Stub API calls with `pkg/httpmock`.
- Use `assert` / `require` from `testify` consistently with nearby tests.
- Use `t.Cleanup(...)` for restoring globals.
- For E2E, add new `txtar` cases under `e2e/testscripts/<area>` and register them in `e2e/e2e_test.go`.

## Change Guidance

- Check for an existing helper before adding a new utility.
- If flag surfaces or generated spec flags change, run `go generate ./...`.
- If command help or command trees change, consider whether docs generation should be rerun.
- Add or update tests when behavior changes.
- Avoid unrelated formatting churn.

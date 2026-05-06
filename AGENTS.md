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

Top-level command group: `algolia agents`. Verbs: `list`, `get`, `create`, `update`, `delete`, `publish`, `unpublish`, `duplicate`, `test`, `run`. Backend source of truth: `github.com/algolia/conversational-ai`.

### API client (`api/agentstudio/`)

- Auth: standard Algolia headers (`X-Algolia-Application-Id`, `X-Algolia-API-Key`). No bearer tokens. Comes from the active profile via `*cmdutil.Factory.AgentStudioClient`.
- Base URL resolution priority: per-profile `agent_studio_url` → env `ALGOLIA_AGENT_STUDIO_URL` → build-time `agentstudio.DefaultBaseURL` (set via `ldflags`, mirrors `dashboard.DefaultDashboardURL`) → cluster-proxy fallback `https://{appID}.algolia.net/agent-studio`. The cluster proxy already does region routing — don't add a `Region` field.
- Errors: `*APIError` with `StatusCode`, `Detail`, optional `Sentinel`. The detail extractor prefers structured FastAPI `detail[].msg` arrays over the generic `message` field — backends that return both pair them as `{"message":"Input is invalid, see detail/body:","detail":[{"msg":"..."}]}` and the structured form is the actionable one.
- `CreateAgent` / `UpdateAgent` accept `json.RawMessage` bodies on purpose. The backend's `AgentConfigCreate` schema is large, deeply validated, and evolves often. The CLI is a pass-through; the backend validates; our 422-detail surfacing makes errors actionable.
- `Completions(...)` returns the raw `*http.Response`. Caller checks `Content-Type` (`text/event-stream` → `ParseStream`; else copy verbatim). One method, two output shapes.

### Streaming (`api/agentstudio/sse.go`)

The wire format is **not** standard SSE. Two protocols, both served as `text/event-stream`:

- **v5 (CLI default)**: standard SSE — `data: <json>\n\n`, `data: [DONE]` sentinel.
- **v4**: line-delimited bespoke — `<type-code>:<json>\n` per line, no terminator. Type codes: `0` = text, `9` = tool-call, `d` = finish-message, etc. (see `v4TypeNames` in `sse.go`).

`ParseStream` sniffs the line prefix and emits a normalized `StreamEvent{Type, Data, Raw}` for both. `compatibilityMode` is a **required** server-side query parameter — the CLI defaults to v5 and exposes `--compatibility v4|v5`.

Streaming output convention: NDJSON to stdout regardless of TTY, one `{"type":"...","data":{...}}` per line. Plays well with `jq -r 'select(.type=="text-delta") | .data.delta'`. Don't fork rendering between TTY/non-TTY for streaming responses.

### Dry-run convention

`agents create / update / delete / test / run` all support `--dry-run`. Two output modes:

- **Human (default)**: print `Dry run: would <METHOD> <PATH>` followed by the resolved JSON body (pretty-printed for the body preview).
- **Structured (only when `--output` is explicitly set)**: emit `{"action":"...","request":"...","source":"...","bytes":N,"body":<...>,"dryRun":true,...}`. Gate this on `cmd.Flags().Changed("output")`, **not** on `PrintFlags.HasStructuredOutput()` — using the latter would let `WithDefaultOutput("json")` from the success path silently steal the human dry-run output.

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

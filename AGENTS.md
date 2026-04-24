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

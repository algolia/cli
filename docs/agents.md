# Agent Studio CLI

Reference for `algolia agents` and `api/agentstudio/`. The Go source files keep one-line godoc comments only; everything that needs more than that lives here.

## Command surface

```
algolia agents
  list | get | create | update | delete | publish | unpublish | duplicate
  try     # run a completion against an unsaved configuration
  run     # run a completion against a persisted agent
  cache invalidate
  config get | set
  providers     list | get | create | update | delete | models | defaults
  conversations list | get | delete | purge | export
  domains       list | get | create | delete | bulk-insert | bulk-delete
  keys          list | get | create | update | delete
  feedback      create
  user-data     get | delete
  internal      status | memorize | ponder | consolidate    (hidden)
```

Backend source of truth: `github.com/algolia/conversational-ai`.

## Layout

Both layers mirror the OpenAPI spec's tag boundaries — one source file (and one test file) per API tag.

```
api/agentstudio/
  client.go          infra: Config, NewClient, setHeaders, checkResponse, extractDetail, sentinelFor
  agents.go          Agents tag (CRUD + lifecycle + InvalidateAgentCache)
  completions.go     Completions tag
  providers.go       Providers tag (CRUD + model discovery)
  configuration.go   Configurations tag
  conversations.go   Conversations tag
  domains.go         Allowed-domains tag
  keys.go            Secret-keys tag
  feedback.go        Feedback tag
  userdata.go        User-data tag
  internal.go        hidden/internal endpoints
  sse.go             v4/v5 stream parser
  host.go            base URL resolution
  errors.go          APIError + sentinels
  types.go           shared response types
```

For the cmd layer, top-level verbs each own a subpackage. Sub-groups (`cache/`, `providers/`, `config/`, `conversations/`, `domains/`, `keys/`, `feedback/`, `userdata/`, `internal/`) keep all their verbs in one package, one file per verb. Don't promote sub-group verbs into per-verb subpackages — they share internal helpers.

## Auth + host resolution

- Standard Algolia headers: `X-Algolia-Application-Id`, `X-Algolia-API-Key`. No bearer tokens.
- `X-Algolia-User-ID` is **not** an authorization signal — it's a cleartext label for telemetry/rate-limiting. The signed equivalent is `X-Algolia-Secure-User-Token` (JWT), wired only into `/completions`.
- Base URL precedence: profile `agent_studio_url` → env `ALGOLIA_AGENT_STUDIO_URL` → ldflag `agentstudio.DefaultBaseURL` → cluster-proxy fallback `https://{appID}.algolia.net/agent-studio`.
- Admin API key is required for `keys create|update|delete`. Backend rejects with 403 `"Admin API key required."` otherwise.

## Pass-through bodies (`json.RawMessage`)

`CreateAgent`, `UpdateAgent`, `CreateProvider`, `UpdateProvider`, `Completions`, and the internal memory verbs accept `json.RawMessage`. The backend schemas are large discriminated unions that evolve frequently (provider `input` is a 6-way union, `messages` is a role union, agent `config`/`tools` are free-form). Mirroring them in Go would lie about parity and force a release on every backend bump. The CLI is a pass-through; backend validates; our 422 surfacing makes errors actionable.

## Streaming (`sse.go`)

Two protocols, both served as `text/event-stream`. `compatibilityMode` is a **required** server-side query param; CLI defaults to v5.

| Protocol | Frame | Terminator |
|----------|-------|------------|
| v5 (default) | `data: <json>\n\n` | `data: [DONE]` |
| v4           | `<type-code>:<json>\n` | body close |

`ParseStream` sniffs the prefix and emits a normalised `StreamEvent{Type, Data, Raw}`. v4 type codes (`0`=text, `9`=tool-call, `d`=finish-message, …) live in `v4TypeNames`.

Output rendering (`pkg/cmd/agents/shared/completion.go`):
- TTY default → human transcript (text-deltas inline, tool calls/results dim, errors red).
- Non-TTY → NDJSON `{"type":"...","data":{...}}` per line. Stable contract for `jq` pipelines.
- `--ndjson` forces NDJSON on a TTY.

## Completion runtime knobs

`agents try` and `agents run` share the same flags:

| Flag | Wire | Default | Notes |
|------|------|---------|-------|
| `--no-stream` | `?stream=false` | stream | buffered single JSON |
| `--compatibility v4\|v5` | `?compatibilityMode=ai-sdk-{4,5}` | v5 | required server-side |
| `--no-cache` | `?cache=false` | cache on | bypass completion cache |
| `--no-memory` | `?memory=false` | memory on | disable agent memory |
| `--no-analytics` | `?analytics=false` | analytics on | skip analytics writes |
| `--secure-user-token <jwt>` | `X-Algolia-Secure-User-Token` | omitted | per-end-user scoping |
| `--ndjson` | (output) | TTY: rich render | force NDJSON on TTY |

`No*` polarity is inverted from the wire because the backend defaults all three to true. `memory` in particular has an `anyOf [{const false}, {type null}]` schema — `memory=true` would 422. Wire form omits the param when `No*` is false, sends `<param>=false` when true.

## On `--dry-run`

Two distinct concepts share the name:

- **CLI `--dry-run`** on `agents create / update / delete / run`: validates the request, prints what would be sent, makes no HTTP call. Two output modes (human / structured) gated on whether `--output` was explicitly set — gating on `PrintFlags.HasStructuredOutput()` would let `WithDefaultOutput("json")` from the success path silently steal the human dry-run output.
- **Conversational-ai "dry-run"** = a real completion against a configuration that hasn't been persisted. Backend exposes this via the `agent_id="test"` route. CLI surfaces it as `agents try`.

`agents try` therefore has no `--dry-run` flag — the whole command IS the dry-run. To preview the wire body without calling the backend, marshal `{"messages":[...], "configuration":{...}}` yourself. The dry-run e2e regression-asserts that `agents try --dry-run` is rejected.

## Providers: `-F` vs flags

`agents providers create` accepts either:

- **`-F <file>`** — full `ProviderAuthenticationCreate` JSON (all `providerName` variants, including `azure_openai` and `openai_compatible`), or
- **Flags** — `--name`, `--provider` (`openai` \| `anthropic` \| `google_genai` \| `deepseek`), plus exactly one of `--api-key`, `--api-key-stdin`, or `--api-key-env <VAR>`. Optional `--base-url` only for `openai` / `anthropic`.

`-F` and the shortcut flags are **mutually exclusive**.

`agents providers update <id>` accepts **`-F`** (patch JSON) **or** shortcut flags: any non-empty combination of `--name`, `--api-key` / `--api-key-stdin` / `--api-key-env`, and `--base-url`, with the same exclusivity rule against `-F`.

Prefer **`--api-key-env`** or **`--api-key-stdin`** over **`--api-key`** (shell history). `--dry-run` still shows the resolved body unredacted so authors can verify what would be sent.

Team sign-off and **Anya** QA checklist: [`docs/qa/arg_friendly_providers_SIGNOFF.md`](qa/arg_friendly_providers_SIGNOFF.md).

## Secret masking

`apiKey` (provider input) and `value` (secret-keys) are masked to `"***"` by default. Pass `--show-secret` to render verbatim. Masking happens at the cmd layer (`pkg/cmd/agents/shared/mask.go`), not the client. `--dry-run` does **not** mask: the user authored the file and is being shown what THEY are about to send. Three asterisks, no last-N preview — goal is "impossible to copy by accident", not "allow last-4 lookup". `secretFieldNames` is the closed set; extend alphabetically when new credential fields land.

## Backend / spec gotchas

These are real and have all bitten us during live vetting against staging. Don't remove the CLI-side guardrails without re-checking on a deployed backend.

- **`conversations purge`** — OpenAPI marks `startDate`/`endDate` optional; live backend rejects dateless DELETE with `400 "At least one filter is required."`. CLI requires at least one of `--start-date` / `--end-date`. To wipe everything, pass an open-ended bound (`--start-date 1970-01-01`).
- **`user-data get|delete` with `/` in the token** — gateway decodes `%2F` before path matching, yielding a misleading 404. CLI rejects tokens containing `/` with an actionable error.
- **Memory ops path** — `memorize`/`ponder`/`consolidate` live under `/1/agents/agents/{id}/<verb>` (segment genuinely doubled in the live deployment). `/1/agents/{id}/memorize` 404s.
- **`InvalidateAgentCache`** — date format is **not** validated client-side. Backend's Pydantic parser is the source of truth; our 422 surfacing forwards the message verbatim. Mirroring in Go would create silent skew.
- **`bulk-delete` reporting** — backend returns 204 with no body, so it can't tell us which IDs actually existed. On TTY we pre-fetch the list to split requested IDs into removed-vs-already-absent. Non-TTY skips the GET (script perf).

## Conversations: `purge` vs `delete`

Same HTTP method, two orders of magnitude difference in blast radius:

- `delete <agent-id> <conv-id>` — surgical, one conversation.
- `purge <agent-id> --start-date|--end-date` — bulk, every conversation in the date range.

Both flow through the same `--confirm` / non-TTY-refuses-without-it rule as `agents delete`. `--dry-run` previews the URL and labels the scope (`scope: between A and B`, `scope: from A onwards`, …).

## Telemetry

Single `"Command Invoked"` event from root, with `{command: cmd.CommandPath(), flags: [<changed flag names>]}`. Already attributes per-verb and surfaces `--dry-run`. Don't add bespoke per-verb telemetry — diverging from convention for one feature would be a mistake. Outcome (success/error) is a separate, all-commands refactor.

## Shared helpers

`pkg/cmd/agents/shared/`:

- `BuildMessages`, `ReadJSONFile`, `MarshalCompletionBody` — completion body assembly.
- `RenderCompletion`, `renderTTY`, `renderNDJSON` — streaming output.
- `NormalizeCompatibility` — `v4`/`v5` aliases → wire form.
- `PrintDryRun` — shared `--dry-run` formatter.
- `MaskInput`, `MaskString` — secret redaction.
- `SourceLabel`, `TrimUTF8BOM`, `ReadJSONBody` — file/stdin plumbing.

Extract on second use, not pre-emptively.

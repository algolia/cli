# Agent Studio CLI

Reference for `algolia agents` and `api/agentstudio/`. The Go source files keep one-line godoc comments only; everything that needs more than that lives here.

## Command surface

```
algolia agents
  list | get | create | update | delete | publish | unpublish | duplicate
  try     # run a completion against an unsaved configuration
  run     # run a completion against a persisted agent
  tools   add-search-index   # merge algolia_search_index tool / index entry
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

## How to use

The CLI is a thin client over the Agent Studio HTTP API. **Exploration and help** start with `algolia agents --help` and `algolia agents <subcommand> --help` (flag names, examples, and defaults are maintained there).

### 1. Credentials and host

Use the **same application ID and API key** as for other Algolia APIs (CLI **profile** or **`ALGOLIA_APPLICATION_ID`** / **`ALGOLIA_API_KEY`**). The key must be allowed to call Agent Studio for your app (feature access is enforced server-side).

If requests hit the wrong cluster, set **`ALGOLIA_AGENT_STUDIO_URL`** or profile **`agent_studio_url`** (see **Auth + host resolution** below).

### 2. Read-only orientation

Inspect existing resources before you mutate anything:

```bash
algolia agents list
algolia agents get <agent-id> --output json
algolia agents providers list
```

**Agent IDs and provider IDs** are opaque strings (often UUID-shaped). **`publish`**, **`run`**, and most verbs that address a single agent expect **`agent-id`**, not the human-readable `name` field inside the agent body.

### 3. Providers (LLM backing)

Either pass **full JSON** with **`-F`** (all backend provider types — Azure, `openai_compatible`, …) or use **shortcut flags** for the common OpenAI-compatible single-key vendors — see **Providers: `-F` vs flags** below.

```bash
algolia agents providers create --name prod-openai --provider openai \
  --api-key-env OPENAI_API_KEY
```

Note the **`providerId`** returned; you embed it when authoring agent JSON (`providerId` in the create body).

### 4. Agents: JSON-first create/update

There is **no interactive wizard** or template subcommands — **create**/**update** consume **`-F`** files the backend validates:

```bash
algolia agents create -F agent.json
algolia agents update <agent-id> -F patch.json
```

Author JSON by aligning with **`agents get`** from an agent tuned in Dashboard, from your OpenAPI/SDK examples, or from internal integration docs — the schemas move with the backend.

### 5. Completions: try vs run

**`try`** sends an **unsaved** configuration file to **`/1/agents/test/completions`** (nothing persisted in the agents list):

```bash
algolia agents try -c draft.json --message "Smoke test message"
algolia agents try -c draft.json -m hi --no-stream
```

**`run`** targets a **persisted agent** **`agent-id`** and is the shape production usage follows:

```bash
algolia agents run <agent-id> -m "Hello"
```

Streaming, compatibility mode, cache/memory/analytics knobs, and secure-user token behave the same on both commands — see **Completion runtime knobs** below and **`--help`** on each command.

### 6. Lifecycle and cache

Publishing makes an agent reachable for **`run`**-style completions the way downstream apps expect:

```bash
algolia agents publish <agent-id>
algolia agents unpublish <agent-id>
```

Clear cached completions after config changes — non-interactive shells need **`-y` / `--confirm`** (same rule as other destructive `agents` verbs):

```bash
algolia agents cache invalidate <agent-id> -y
```

Destructive verbs (**`agents delete`**, **`conversations delete|purge`**, **`providers delete`**, …) require **`-y` / `--confirm`** in scripts; interactive sessions get a prompt unless you pass **`-y`**.

### 7. Output for scripts

Prefer **`--output json`** when you automate parsing (**`jq`**, CI). **`try`** / **`run`** on a non-TTY default to NDJSON streams; **`--ndjson`** forces stream JSON on an interactive terminal. Masked secrets (**`providers`**, **`keys`**) omit raw credentials unless **`--show-secret`** — see **Secret masking** below.

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
- Overrides (`ALGOLIA_AGENT_STUDIO_URL` / profile `agent_studio_url`) must be **HTTPS** URLs with a scheme and host. Plain **HTTP** is rejected unless **`ALGOLIA_AGENT_STUDIO_ALLOW_INSECURE_HTTP=1`** is set (local development only).
- Cluster-proxy fallback requires an **application id** that is **4–32 alphanumeric** characters (so it is safe to embed as a single DNS label). Invalid characters produce a CLI error before any request.
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

## Ephemeral completions (`agents try`)

Backend route **`/1/agents/test/completions`** runs a completion from an **unsaved** configuration (`agents try -c`). Nothing is written to the agents list—use it to iterate on prompts and tools before **`create` / `publish`**.

To inspect the JSON you would POST to **`run`**, assemble `{"messages":[...]}` locally (the persisted agent supplies configuration server-side for real agent ids).

## Providers: `-F` vs flags

`agents providers create` accepts either:

- **`-F <file>`** — full `ProviderAuthenticationCreate` JSON (all `providerName` variants, including `azure_openai` and `openai_compatible`), or
- **Flags** — `--name`, `--provider` (`openai` \| `anthropic` \| `google_genai` \| `deepseek`), plus exactly one of `--api-key`, `--api-key-stdin`, or `--api-key-env <VAR>`. Optional `--base-url` only for `openai` / `anthropic`.

`-F` and the shortcut flags are **mutually exclusive**.

`agents providers update <id>` accepts **`-F`** (patch JSON) **or** shortcut flags: any non-empty combination of `--name`, `--api-key` / `--api-key-stdin` / `--api-key-env`, and `--base-url`, with the same exclusivity rule against `-F`.

Prefer **`--api-key-env`** or **`--api-key-stdin`** over **`--api-key`** (shell history).

## Secret masking

`apiKey` (provider input) and `value` (secret-keys) are masked to `"***"` by default. Pass **`--show-secret`** to render verbatim. Masking happens at the cmd layer (`pkg/cmd/agents/shared/mask.go`), not the client. Three asterisks, no last-N preview — goal is "impossible to copy by accident", not "allow last-4 lookup". `secretFieldNames` is the closed set; extend alphabetically when new credential fields land.

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

Both flow through the same **`--confirm`** / non-TTY-refuses-without-it rule as `agents delete`.

## Telemetry

Single `"Command Invoked"` event from root, with `{command: cmd.CommandPath(), flags: [<changed flag names>]}`. Don't add bespoke per-verb telemetry — diverging from convention for one feature would be a mistake. Outcome (success/error) is a separate, all-commands refactor.

## Shared helpers

`pkg/cmd/agents/shared/`:

- `BuildMessages`, `ReadJSONFile`, `MarshalCompletionBody` — completion body assembly.
- `RenderCompletion`, `renderTTY`, `renderNDJSON` — streaming output.
- `NormalizeCompatibility` — `v4`/`v5` / `ai-sdk-4`/`ai-sdk-5` (case-insensitive) → wire form.
- `MaskInput`, `MaskString` — secret redaction.
- `SourceLabel`, `TrimUTF8BOM` — file/stdin plumbing.

Extract on second use, not pre-emptively.

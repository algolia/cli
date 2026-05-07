# Arg-friendly `agents providers` — team sign-off & QA sheet

**Feature branch:** `feat/agents-providers-arg-friendly` (stacked on `lab_week_3` / [PR #212](https://github.com/algolia/cli/pull/212))

**Scope:** Optional flags for `agents providers create|update` so the common “single API key” case does not require `-F`. Full JSON via `-F` unchanged; Azure / `openai_compatible` still require `-F`.

**Planning reference:** Lab Week 3 roles (“The team (5 people)”: Iris, Marek, Yuki, Diego, Priya) — often kept locally as `tmp/lab_week_3_plan.md` beside other lab notes.

> This file lives under `docs/qa/` so it ships with the repo. Some checkouts list `tmp/` in `.git/info/exclude`; use this path as the canonical copy for PRs.

---

## Engineering sign-off (planning team)

| Person | Role | Position on this change |
|--------|------|-------------------------|
| **Iris** | Auth & identity | **On board.** Out of scope for OAuth; provider keys stay user-supplied. Flags do not change header auth (`X-Algolia-Application-Id` / `X-Algolia-API-Key`). |
| **Marek** | CLI / DX | **On board.** Matches Cobra patterns: `-F` vs flags mutually exclusive; long-help documents shell-history risk for `--api-key` and recommends `--api-key-env` / `--api-key-stdin`. |
| **Yuki** | Agent Studio backend | **On board.** Wire body is identical to JSON create/patch; CLI only assembles JSON for `openai`, `anthropic`, `google_genai`, `deepseek`. |
| **Diego** | Security | **On board.** Same threat model as `-F`: secrets can still hit shell history if the user passes `--api-key`; env/stdin paths preferred. No new persistence. |
| **Priya** | Search UI (aspirational track) | **N/A / on board.** No UI work; no objection to CLI ergonomics that reduce temp files. |

---

## QA owner: **Anya**

**Exit criteria**

- [ ] `task build` (or `go build` per repo) produces a binary; all of the commands below use that binary as `./algolia`.
- [ ] Agents e2e suite runs with app credentials: `go test ./e2e -tags=e2e -run TestAgents -count=1` includes `testscripts/agents/providers.txtar` and passes (or document skip if creds unavailable).
- [ ] Unit tests: `go test ./pkg/cmd/agents/providers/...` green.
- [ ] Live smoke (staging / beta app): at least one **create** and one **update** via flags succeed; **dry-run** paths show expected JSON and do not call mutating APIs.

---

## Anya — testing command sheet

Run from the **repository root** after `task build`. Replace placeholders.

### 0. Binary and profile

```bash
cd /path/to/cli
task build
./algolia version
# Ensure profile or env has app + key (same as other `algolia agents` commands)
```

### 1. Unit tests (no network)

```bash
go test ./pkg/cmd/agents/providers/... -count=1
```

### 2. Contract tests (e2e harness; needs credentials)

```bash
export ALGOLIA_APPLICATION_ID="YOUR_APP_ID"
export ALGOLIA_API_KEY="YOUR_API_KEY"
go test ./e2e -tags=e2e -run TestAgents -count=1
```

### 3. Help / usage sanity

```bash
./algolia agents providers create --help
./algolia agents providers update --help
```

### 4. Create — **dry-run** (flags), no API call

```bash
./algolia agents providers create \
  --name "qa-cli-flags-smoke" \
  --provider anthropic \
  --api-key "sk-ant-REDACTED" \
  --dry-run
```

Expect: human or structured dry-run summary; body includes `providerName`, `name`, `input.apiKey` (unmasked in dry-run per masking rules in `docs/agents.md`).

### 5. Create — **dry-run** via env (no key on command line)

```bash
export QA_ANTHROPIC_KEY="sk-ant-your-real-test-key"
./algolia agents providers create \
  --name "qa-cli-flags-env" \
  --provider anthropic \
  --api-key-env QA_ANTHROPIC_KEY \
  --dry-run
unset QA_ANTHROPIC_KEY
```

### 6. Create — **live** (optional; consumes provider quota)

```bash
export QA_ANTHROPIC_KEY="sk-ant-..."
PID=$(./algolia agents providers create \
  --name "qa-live-$(date +%s)" \
  --provider anthropic \
  --api-key-env QA_ANTHROPIC_KEY \
  --output json | jq -r .id)
echo "PID=$PID"
```

### 7. Update — **dry-run** (rename only)

```bash
./algolia agents providers update "$PID" --name "qa-renamed" --dry-run
```

### 8. Update — **live** (rotate key or rename)

```bash
./algolia agents providers update "$PID" --name "qa-renamed-final"
# or rotate:
# ./algolia agents providers update "$PID" --api-key-env QA_ANTHROPIC_KEY
```

### 9. Negative cases

```bash
# Must error: cannot mix -F and flags (use a real JSON file path):
# echo '{"name":"x","providerName":"openai","input":{"apiKey":"sk"}}' > /tmp/prov.json
# ./algolia agents providers create -F /tmp/prov.json --name x --provider openai --api-key sk 2>&1 | head -5

# Must error: unsupported provider for flag path
./algolia agents providers create --name x --provider azure_openai --api-key sk --dry-run 2>&1 | head -5

# Must error: create without -F and without full flag set
./algolia agents providers create 2>&1 | head -5
```

### 10. Cleanup

```bash
./algolia agents providers delete "$PID" -y
```

---

## Report back

Add `docs/qa/anya_arg_friendly_providers_vet/REPORT.md` (date, commit SHA, pass/fail per section, backend 422 text if any) after live runs.

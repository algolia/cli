package providers

import (
	"encoding/json"
)

// MaskInput redacts secret fields in a provider's `input` payload
// before display.
//
// Why: the backend returns the apiKey field in responses (verified
// against the OpenAPI spec — the OpenAIProviderInput-Output et al.
// schemas all include apiKey). Without masking, `agents providers
// list/get/create/update --output json` writes raw API keys to stdout,
// which routinely lands in CI logs, terminal scrollback, and shared
// pastes.
//
// What: parses Input as a flat JSON object, replaces the value of any
// "apiKey" field with the literal string "***" (no prefix preview —
// the goal is to make the key impossible to copy by accident, not to
// allow last-4 lookup), and re-marshals. Non-object inputs and
// invalid JSON pass through unchanged so a future schema with a
// non-object `input` shape doesn't lose information silently.
//
// What NOT: this is best-effort defense in depth, not a substitute
// for a server-side redaction toggle. We don't try to walk arbitrarily
// nested structures because the current schemas are flat and adding
// recursive walking adds bug surface for no current value. If a
// future provider grows nested credential fields, expand here with
// tests that pin the new shape.
//
// The same convention will land in Phase 8 for `agents keys` (which
// vends per-agent secret keys with similar handling concerns); when
// it does, lift this into a shared helper.
func MaskInput(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return raw
	}
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		// Non-object input (string, array, null, malformed). Leave it
		// alone rather than throw away the data — see godoc.
		return raw
	}
	masked := false
	for _, k := range secretFieldNames {
		if _, ok := obj[k]; ok {
			obj[k] = maskSentinel
			masked = true
		}
	}
	if !masked {
		return raw
	}
	out, err := json.Marshal(obj)
	if err != nil {
		// Unreachable in practice (we just unmarshaled this map). If it
		// ever happens, return the original to fail closed — better to
		// preserve the user's data than emit garbage.
		return raw
	}
	return out
}

// secretFieldNames is the closed set of input-object keys we treat as
// secrets. Today it's just `apiKey` because that's the only field the
// backend's provider input schemas mark as a credential. Extend
// alphabetically if new credential fields land.
var secretFieldNames = []string{"apiKey"}

// maskSentinel is the literal we substitute for secret values. Three
// asterisks is short enough to read in tables and unambiguous as
// "redacted, not the real value." Avoid empty string (loses
// information that a key WAS present) and avoid prefix-of-original
// (encourages users to assume some characters leaked).
const maskSentinel = "***"

package shared

import "encoding/json"

// MaskInput redacts known secret fields ("apiKey") in a flat JSON
// object and returns the re-marshaled bytes. Non-object input,
// invalid JSON, and objects without a known secret field pass
// through unchanged.
func MaskInput(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return raw
	}
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return raw
	}
	masked := false
	for _, k := range secretFieldNames {
		if _, ok := obj[k]; ok {
			obj[k] = MaskSentinel
			masked = true
		}
	}
	if !masked {
		return raw
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return raw
	}
	return out
}

// MaskSecretsInValue marshals v to JSON, redacts known secret fields
// ("apiKey") anywhere in the resulting tree, and returns a generic value
// suitable for PrintFlags. Used for typed SDK responses (e.g. provider
// authentications) whose secret lives in a nested union the CLI can't mask
// in place. On any marshal/unmarshal error it returns v unchanged.
func MaskSecretsInValue(v any) any {
	raw, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return v
	}
	maskSecretsRecursive(decoded)
	return decoded
}

func maskSecretsRecursive(v any) {
	switch t := v.(type) {
	case map[string]any:
		for _, k := range secretFieldNames {
			if _, ok := t[k]; ok {
				t[k] = MaskSentinel
			}
		}
		for _, child := range t {
			maskSecretsRecursive(child)
		}
	case []any:
		for _, child := range t {
			maskSecretsRecursive(child)
		}
	}
}

// MaskString returns "***" if v is non-empty, otherwise "".
// Used for top-level scalar secrets like SecretKey.Value.
func MaskString(v string) string {
	if v == "" {
		return ""
	}
	return MaskSentinel
}

// secretFieldNames is the closed set of object keys treated as
// secrets. Extend alphabetically as new credential fields land.
var secretFieldNames = []string{"apiKey"}

// MaskSentinel is the literal substituted for secret values.
const MaskSentinel = "***"

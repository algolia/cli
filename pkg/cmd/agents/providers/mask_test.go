package providers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskInput(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty input passes through unchanged",
			in:   "",
			want: "",
		},
		{
			name: "openai shape masks apiKey, preserves baseUrl",
			in:   `{"apiKey":"sk-abc-123","baseUrl":"https://api.openai.com"}`,
			want: `{"apiKey":"***","baseUrl":"https://api.openai.com"}`,
		},
		{
			name: "anthropic shape masks apiKey",
			in:   `{"apiKey":"sk-ant-XYZ"}`,
			want: `{"apiKey":"***"}`,
		},
		{
			name: "azure shape masks apiKey, preserves endpoint+deployment+version (none of those are secrets)",
			in:   `{"apiKey":"k","azureEndpoint":"https://x.openai.azure.com","azureDeployment":"gpt-4","apiVersion":"2024-12-01-preview"}`,
			want: `{"apiKey":"***","azureDeployment":"gpt-4","apiVersion":"2024-12-01-preview","azureEndpoint":"https://x.openai.azure.com"}`,
		},
		{
			name: "no apiKey field: object passes through unchanged",
			in:   `{"baseUrl":"https://x"}`,
			want: `{"baseUrl":"https://x"}`,
		},
		{
			name: "non-object input passes through unchanged (string)",
			in:   `"hello"`,
			want: `"hello"`,
		},
		{
			name: "non-object input passes through unchanged (array)",
			in:   `[1,2,3]`,
			want: `[1,2,3]`,
		},
		{
			name: "malformed JSON passes through unchanged (don't lose data)",
			in:   `{not json`,
			want: `{not json`,
		},
		{
			name: "null apiKey value still gets masked (don't leak the absence-as-known shape)",
			in:   `{"apiKey":null}`,
			want: `{"apiKey":"***"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaskInput(json.RawMessage(tc.in))
			if tc.want == "" {
				assert.Empty(t, got)
				return
			}
			// JSON object key order is not preserved across re-marshal,
			// so compare as JSON for cases where both sides are valid
			// JSON objects. Fall back to exact string match for
			// pass-through cases (malformed JSON, primitives, arrays).
			if json.Valid([]byte(tc.in)) && json.Valid([]byte(tc.want)) &&
				len(tc.in) > 0 && tc.in[0] == '{' {
				assert.JSONEq(t, tc.want, string(got))
			} else {
				assert.Equal(t, tc.want, string(got))
			}
		})
	}
}

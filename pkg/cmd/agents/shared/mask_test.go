package shared

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskInput(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"empty", "", ""},
		{
			"openai",
			`{"apiKey":"sk-abc","baseUrl":"https://api.openai.com"}`,
			`{"apiKey":"***","baseUrl":"https://api.openai.com"}`,
		},
		{"anthropic", `{"apiKey":"sk-ant"}`, `{"apiKey":"***"}`},
		{
			"azure",
			`{"apiKey":"k","azureEndpoint":"https://x","azureDeployment":"gpt-4","apiVersion":"v1"}`,
			`{"apiKey":"***","azureEndpoint":"https://x","azureDeployment":"gpt-4","apiVersion":"v1"}`,
		},
		{"no-secret", `{"baseUrl":"https://x"}`, `{"baseUrl":"https://x"}`},
		{"primitive-string", `"hello"`, `"hello"`},
		{"primitive-array", `[1,2,3]`, `[1,2,3]`},
		{"malformed", `{not json`, `{not json`},
		{"null-apiKey-still-masked", `{"apiKey":null}`, `{"apiKey":"***"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaskInput(json.RawMessage(tc.in))
			if tc.want == "" {
				assert.Empty(t, got)
				return
			}
			if json.Valid([]byte(tc.in)) && json.Valid([]byte(tc.want)) && tc.in[0] == '{' {
				assert.JSONEq(t, tc.want, string(got))
			} else {
				assert.Equal(t, tc.want, string(got))
			}
		})
	}
}

func TestMaskString(t *testing.T) {
	assert.Equal(t, "", MaskString(""))
	assert.Equal(t, "***", MaskString("anything"))
}

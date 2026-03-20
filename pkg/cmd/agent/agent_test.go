package agent

import (
	"os"
	"strings"
	"testing"
)

func TestParseSSEStream(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantText  string
		wantMsgID string
	}{
		{
			name: "parses a complete stream with start and text-delta events",
			input: "data: {\"type\":\"start\",\"messageId\":\"msg_123\"}\n" +
				"data: {\"type\":\"text-delta\",\"delta\":\"Hello \"}\n" +
				"data: {\"type\":\"text-delta\",\"delta\":\"world\"}\n" +
				"data: [DONE]\n",
			wantText:  "Hello world",
			wantMsgID: "msg_123",
		},
		{
			name:      "returns empty on empty stream",
			input:     "",
			wantText:  "",
			wantMsgID: "",
		},
		{
			name: "ignores non-data lines",
			input: "event: message\n" +
				"id: 1\n" +
				"data: {\"type\":\"text-delta\",\"delta\":\"hi\"}\n" +
				"data: [DONE]\n",
			wantText:  "hi",
			wantMsgID: "",
		},
		{
			name: "skips malformed JSON",
			input: "data: not-json\n" +
				"data: {\"type\":\"text-delta\",\"delta\":\"ok\"}\n" +
				"data: [DONE]\n",
			wantText:  "ok",
			wantMsgID: "",
		},
		{
			name: "stops at DONE marker",
			input: "data: {\"type\":\"text-delta\",\"delta\":\"before\"}\n" +
				"data: [DONE]\n" +
				"data: {\"type\":\"text-delta\",\"delta\":\"after\"}\n",
			wantText:  "before",
			wantMsgID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			gotText, gotMsgID, err := parseSSEStream(r)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotText != tt.wantText {
				t.Errorf("text = %q, want %q", gotText, tt.wantText)
			}
			if gotMsgID != tt.wantMsgID {
				t.Errorf("messageID = %q, want %q", gotMsgID, tt.wantMsgID)
			}
		})
	}
}

func TestEnvOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		envValue   string
		defaultVal string
		want       string
	}{
		{
			name:       "returns default when env is not set",
			key:        "TEST_ENV_OR_DEFAULT_UNSET",
			defaultVal: "fallback",
			want:       "fallback",
		},
		{
			name:       "returns env value when set",
			key:        "TEST_ENV_OR_DEFAULT_SET",
			envValue:   "from-env",
			defaultVal: "fallback",
			want:       "from-env",
		},
		{
			name:       "returns default when env is empty string",
			key:        "TEST_ENV_OR_DEFAULT_EMPTY",
			envValue:   "",
			defaultVal: "fallback",
			want:       "fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := envOrDefault(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("envOrDefault(%q, %q) = %q, want %q", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

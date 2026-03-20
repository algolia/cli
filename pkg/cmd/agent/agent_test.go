package agent

import (
	"os"
	"strings"
	"testing"
)

func TestParseSSEStream(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantText    string
		wantMsgID   string
		wantCommand string
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
		{
			name: "parses suggestCommand tool call",
			input: "data: {\"type\":\"start\",\"messageId\":\"msg_456\"}\n" +
				"data: {\"type\":\"text-delta\",\"delta\":\"Try this:\"}\n" +
				"data: {\"type\":\"tool-input-available\",\"toolName\":\"suggestCommand\",\"input\":{\"command\":\"algolia search MOVIES\"}}\n" +
				"data: [DONE]\n",
			wantText:    "Try this:",
			wantMsgID:   "msg_456",
			wantCommand: "algolia search MOVIES",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			result, err := parseSSEStream(r)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Text != tt.wantText {
				t.Errorf("text = %q, want %q", result.Text, tt.wantText)
			}
			if result.MessageID != tt.wantMsgID {
				t.Errorf("messageID = %q, want %q", result.MessageID, tt.wantMsgID)
			}
			if result.Command != tt.wantCommand {
				t.Errorf("command = %q, want %q", result.Command, tt.wantCommand)
			}
		})
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{name: "simple command", command: "algolia indices list", wantErr: false},
		{name: "command with pipe", command: "algolia profile list | head -5", wantErr: false},
		{name: "command with quotes", command: `algolia profile list`, wantErr: false},
		{name: "blocks objects browse", command: "algolia objects browse MOVIES", wantErr: true},
		{name: "blocks search", command: "algolia search MOVIES --query test", wantErr: true},
		{name: "blocks rules browse", command: "algolia rules browse MOVIES", wantErr: true},
		{name: "blocks synonyms browse", command: "algolia synonyms browse MOVIES", wantErr: true},
		{name: "blocks dictionary entries browse", command: "algolia dictionary entries browse stopwords", wantErr: true},
		{name: "blocks semicolon", command: "algolia indices list; rm -rf /", wantErr: true},
		{name: "blocks double ampersand", command: "algolia indices list && echo pwned", wantErr: true},
		{name: "blocks double pipe", command: "algolia indices list || echo fallback", wantErr: true},
		{name: "blocks dollar paren", command: "algolia indices list $(whoami)", wantErr: true},
		{name: "blocks backtick", command: "algolia indices list `whoami`", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCommand(%q) error = %v, wantErr %v", tt.command, err, tt.wantErr)
			}
		})
	}
}

func TestIsSafeCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{name: "profile list is safe", command: "algolia profile list", want: true},
		{name: "search is not safe", command: "algolia search MOVIES --query test", want: false},
		{name: "objects browse is not safe", command: "algolia objects browse MOVIES", want: false},
		{name: "indices list is safe", command: "algolia indices list", want: true},
		{name: "describe is safe", command: "algolia describe search", want: true},
		{name: "delete is not safe", command: "algolia indices delete MOVIES -y", want: false},
		{name: "objects import is not safe", command: "algolia objects import MOVIES -F data.ndjson", want: false},
		{name: "non-algolia command is not safe", command: "rm -rf /", want: false},
		{name: "empty string is not safe", command: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSafeCommand(tt.command)
			if got != tt.want {
				t.Errorf("isSafeCommand(%q) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain text unchanged", input: "hello world", want: "hello world"},
		{name: "strips color codes", input: "\033[31mred\033[0m", want: "red"},
		{name: "strips bold", input: "\033[1mbold\033[0m", want: "bold"},
		{name: "handles spinner overwrite", input: "Loading ⣾\rLoading ⣽\rDone", want: "Done"},
		{name: "preserves newlines", input: "line1\nline2", want: "line1\nline2"},
		{name: "handles CRLF", input: "line1\r\nline2", want: "line1\nline2"},
		{name: "filters empty lines from spinner", input: "Fetching\rFetching\r\n\nresult", want: "Fetching\nresult"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripANSI(tt.input)
			if got != tt.want {
				t.Errorf("stripANSI(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestReplaceAlgoliaBinary(t *testing.T) {
	// Save and restore os.Args[0]
	origArg := os.Args[0]
	defer func() { os.Args[0] = origArg }()

	t.Run("no replacement when binary is algolia", func(t *testing.T) {
		os.Args[0] = "algolia"
		got := replaceAlgoliaBinary("algolia search MOVIES")
		if got != "algolia search MOVIES" {
			t.Errorf("got %q, want %q", got, "algolia search MOVIES")
		}
	})

	t.Run("replaces at start", func(t *testing.T) {
		os.Args[0] = "./algolia"
		got := replaceAlgoliaBinary("algolia search MOVIES")
		if got != "./algolia search MOVIES" {
			t.Errorf("got %q, want %q", got, "./algolia search MOVIES")
		}
	})

	t.Run("replaces after pipe", func(t *testing.T) {
		os.Args[0] = "./algolia"
		got := replaceAlgoliaBinary("algolia objects browse SRC | algolia objects import DST -F -")
		if got != "./algolia objects browse SRC | ./algolia objects import DST -F -" {
			t.Errorf("got %q, want %q", got, "./algolia objects browse SRC | ./algolia objects import DST -F -")
		}
	})

	t.Run("does not replace in arguments", func(t *testing.T) {
		os.Args[0] = "./algolia"
		got := replaceAlgoliaBinary("algolia search hello-algolia")
		if got != "./algolia search hello-algolia" {
			t.Errorf("got %q, want %q", got, "./algolia search hello-algolia")
		}
	})
}

func TestTruncateOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLines int
		want     string
	}{
		{name: "short output unchanged", input: "line1\nline2", maxLines: 10, want: "line1\nline2"},
		{name: "truncates at limit", input: "1\n2\n3\n4\n5", maxLines: 3, want: "1\n2\n3\n[... 2 more lines truncated]"},
		{name: "single line", input: "hello", maxLines: 10, want: "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateOutput(tt.input, tt.maxLines)
			if got != tt.want {
				t.Errorf("truncateOutput() = %q, want %q", got, tt.want)
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

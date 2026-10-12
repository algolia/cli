package shared

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// EnvSecureUserToken is read when neither --secure-user-token nor
// --secure-user-token-file is given, keeping the signed token out of argv.
const EnvSecureUserToken = "ALGOLIA_AGENT_STUDIO_SECURE_USER_TOKEN" //nolint:gosec // env var name, not a credential

// CompletionInputs is the shared flag surface for `agents try` /
// `agents run`. Embed in your command Options struct and call
// RegisterCompletionFlags to wire them onto the cobra command.
//
// Owners of the surrounding command still register their own
// `--config` (`try`) or positional agent-id (`run`) as usual.
type CompletionInputs struct {
	InputFile           string
	Message             string
	NoStream            bool
	Compatibility       string
	NoCache             bool
	NoMemory            bool
	NoAnalytics         bool
	SecureUserToken     string
	SecureUserTokenFile string
	NDJSON              bool
}

// RegisterCompletionFlags adds the standard completion-runtime flags
// (-i/--input, -m/--message, --no-stream, --compatibility,
// --no-cache, --no-memory, --no-analytics, --secure-user-token,
// --ndjson) and the input/message mutex.
func RegisterCompletionFlags(cmd *cobra.Command, in *CompletionInputs) {
	cmd.Flags().
		StringVarP(&in.InputFile, "input", "i", "", "JSON file with messages array (use \"-\" for stdin)")
	cmd.Flags().
		StringVarP(&in.Message, "message", "m", "", "Single user message (convenience for one-shot prompts)")
	cmd.Flags().BoolVar(&in.NoStream, "no-stream", false, "Request a buffered JSON response instead of SSE")
	cmd.Flags().
		StringVar(&in.Compatibility, "compatibility", "", "Streaming protocol: v4 / ai-sdk-4 or v5 / ai-sdk-5 (default); case-insensitive")
	cmd.Flags().BoolVar(&in.NoCache, "no-cache", false, "Bypass the backend completion cache (default: cache enabled)")
	cmd.Flags().
		BoolVar(&in.NoMemory, "no-memory", false, "Disable agent memory for this completion (default: memory enabled)")
	cmd.Flags().
		BoolVar(&in.NoAnalytics, "no-analytics", false, "Skip Agent Studio analytics for this completion (default: analytics enabled)")
	cmd.Flags().
		StringVar(&in.SecureUserToken, "secure-user-token", "", "Signed JWT scoping the conversation/memory/analytics partition to an end-user (X-Algolia-Secure-User-Token)")
	cmd.Flags().
		StringVar(&in.SecureUserTokenFile, "secure-user-token-file", "", "Read the secure user token from this file instead of passing it on the command line (also: "+EnvSecureUserToken+")")
	cmd.MarkFlagsMutuallyExclusive("secure-user-token", "secure-user-token-file")
	cmd.Flags().
		BoolVar(&in.NDJSON, "ndjson", false, "Force NDJSON output even on a TTY (default on non-TTY; use this when you want machine-parseable events but also want to see them on screen)")
	cmd.MarkFlagsMutuallyExclusive("input", "message")
}

// ResolveSecureUserToken returns the secure user token from, in order,
// --secure-user-token, --secure-user-token-file or the
// ALGOLIA_AGENT_STUDIO_SECURE_USER_TOKEN environment variable. Prefer the file
// or the environment variable: argv is visible in `ps` and shell history.
func (in *CompletionInputs) ResolveSecureUserToken() (string, error) {
	if in.SecureUserToken != "" {
		return in.SecureUserToken, nil
	}
	if in.SecureUserTokenFile != "" {
		b, err := os.ReadFile(in.SecureUserTokenFile)
		if err != nil {
			return "", fmt.Errorf("read --secure-user-token-file: %w", err)
		}
		tok := strings.TrimSpace(string(b))
		if tok == "" {
			return "", fmt.Errorf("--secure-user-token-file %s is empty", in.SecureUserTokenFile)
		}
		return tok, nil
	}
	return strings.TrimSpace(os.Getenv(EnvSecureUserToken)), nil
}

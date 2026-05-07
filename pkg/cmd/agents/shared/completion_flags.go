package shared

import "github.com/spf13/cobra"

// CompletionInputs is the shared flag surface for `agents try` /
// `agents run`. Embed in your command Options struct and call
// RegisterCompletionFlags to wire them onto the cobra command.
//
// Owners of the surrounding command still register their own
// `--config` / `--dry-run` / agent-id flags as usual.
type CompletionInputs struct {
	InputFile       string
	Message         string
	NoStream        bool
	Compatibility   string
	NoCache         bool
	NoMemory        bool
	NoAnalytics     bool
	SecureUserToken string
	NDJSON          bool
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
		StringVar(&in.Compatibility, "compatibility", "", "Streaming protocol: v4 (ai-sdk-4) or v5 (ai-sdk-5, default)")
	cmd.Flags().BoolVar(&in.NoCache, "no-cache", false, "Bypass the backend completion cache (default: cache enabled)")
	cmd.Flags().
		BoolVar(&in.NoMemory, "no-memory", false, "Disable agent memory for this completion (default: memory enabled)")
	cmd.Flags().
		BoolVar(&in.NoAnalytics, "no-analytics", false, "Skip Agent Studio analytics for this completion (default: analytics enabled)")
	cmd.Flags().
		StringVar(&in.SecureUserToken, "secure-user-token", "", "Signed JWT scoping the conversation/memory/analytics partition to an end-user (X-Algolia-Secure-User-Token)")
	cmd.Flags().
		BoolVar(&in.NDJSON, "ndjson", false, "Force NDJSON output even on a TTY (default on non-TTY; use this when you want machine-parseable events but also want to see them on screen)")
	cmd.MarkFlagsMutuallyExclusive("input", "message")
}

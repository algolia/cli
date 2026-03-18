package root

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/telemetry"
)

// runSummary is the JSON object written to stderr when DEBUG is set.
type runSummary struct {
	Event        string `json:"event"`
	InvocationID string `json:"invocation_id"`
	Command      string `json:"command"`
	Status       string `json:"status"`
	DurationMs   int64  `json:"duration_ms"`
	Error        string `json:"error,omitempty"`
}

func shouldEmitRunSummary() bool {
	return os.Getenv("DEBUG") != ""
}

// sensitiveFlagPattern matches flags that take a secret value; submatch 1 = flag name (e.g. --api-key= or -p ), submatch 2 = value.
var sensitiveFlagPattern = regexp.MustCompile(
	`(--(?:api-key|application-id|admin-api-key)=)(\S+)|(--(?:api-key|application-id|admin-api-key)\s+)(\S+)|(-p\s+)(\S+)`)

// sensitiveValuePattern matches key/value pairs in error messages; submatch 1 = label, submatch 2 = value.
var sensitiveValuePattern = regexp.MustCompile(
	`(?i)(api[_\s]?key|application[_\s]?id)\s*[:=]\s*([^\s]+)`)

// maskWithLast4 returns a masked string showing only the last 4 chars (e.g. ***c123). If len <= 4, returns ****.
func maskWithLast4(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return "***" + s[len(s)-4:]
}

// sanitizeRunSummaryCommand redacts sensitive flag values, keeping last 4 chars (e.g. --api-key=***c123).
func sanitizeRunSummaryCommand(cmd string) string {
	return sensitiveFlagPattern.ReplaceAllStringFunc(cmd, func(match string) string {
		subs := sensitiveFlagPattern.FindStringSubmatch(match)
		if len(subs) < 3 {
			return match
		}
		for i := 1; i < len(subs); i += 2 {
			if subs[i] != "" && i+1 < len(subs) && subs[i+1] != "" {
				return subs[i] + maskWithLast4(subs[i+1])
			}
		}
		return match
	})
}

// sanitizeRunSummaryError redacts sensitive values in error messages, keeping last 4 chars (e.g. api_key: ***c123).
func sanitizeRunSummaryError(errMsg string) string {
	s := sensitiveValuePattern.ReplaceAllStringFunc(errMsg, func(match string) string {
		subs := sensitiveValuePattern.FindStringSubmatch(match)
		if len(subs) >= 3 && subs[2] != "" {
			return subs[1] + ": " + maskWithLast4(subs[2])
		}
		return match
	})
	if len(s) > 500 {
		s = s[:497] + "..."
	}
	return s
}

func emitRunSummary(stderr io.Writer, ctx context.Context, cmd *cobra.Command, runErr error, duration time.Duration) {
	meta := telemetry.GetEventMetadata(ctx)
	invocationID := ""
	if meta != nil {
		invocationID = meta.InvocationID
	}
	commandPath := ""
	if meta != nil && meta.CommandPath != "" {
		commandPath = meta.CommandPath
	}
	if commandPath == "" && cmd != nil {
		commandPath = cmd.CommandPath()
	}
	commandPath = sanitizeRunSummaryCommand(commandPath)
	status := "ok"
	errMsg := ""
	if runErr != nil {
		status = "error"
		errMsg = sanitizeRunSummaryError(runErr.Error())
	}
	s := runSummary{
		Event:        "cli_run",
		InvocationID: invocationID,
		Command:      commandPath,
		Status:       status,
		DurationMs:   duration.Milliseconds(),
		Error:        errMsg,
	}
	enc := json.NewEncoder(stderr)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
}

package agent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
)

// blockedCommands lists command prefixes that cannot be executed from the agent.
// These commands will be shown as suggestions only.
var blockedCommands = []string{
	"objects browse",
	"search ",
	"rules browse",
	"synonyms browse",
	"dictionary entries browse",
}

// safeCommands lists read-only command prefixes that can be auto-run without confirmation.
var safeCommands = []string{
	"profile list",
	"application list",
	"indices list",
	"apikeys list",
	"settings get",
	"dictionary settings get",
	"describe",
	"open",
	"events tail",
	"crawler list",
	"crawler get",
	"crawler stats",
	"indices config export",
	"indices analyze",
}

// jsonOutputCommands lists command prefixes that support the -o json flag.
var jsonOutputCommands = []string{
	"profile list",
	"application list",
	"indices list",
	"apikeys list",
	"settings get",
	"crawler list",
	"crawler get",
	"crawler stats",
	"indices analyze",
}

// forceJSONOutput appends -o json to the command if it supports it and doesn't already have it.
func forceJSONOutput(command string) string {
	if strings.Contains(command, " -o ") || strings.Contains(command, " --output ") {
		return command
	}
	sub := strings.TrimPrefix(command, "algolia ")
	if sub == command {
		return command
	}
	for _, prefix := range jsonOutputCommands {
		if strings.HasPrefix(sub, prefix) {
			return command + " -o json"
		}
	}
	return command
}

// isSafeCommand checks if a command is read-only and can be auto-run.
func isSafeCommand(command string) bool {
	// Strip the leading "algolia " to get the subcommand.
	sub := strings.TrimPrefix(command, "algolia ")
	if sub == command {
		return false
	}
	for _, safe := range safeCommands {
		if strings.HasPrefix(sub, safe) {
			return true
		}
	}
	return false
}

// isBlockedCommand checks if a command is not allowed to run from the agent.
func isBlockedCommand(command string) bool {
	sub := strings.TrimPrefix(command, "algolia ")
	if sub == command {
		return false
	}
	for _, blocked := range blockedCommands {
		if strings.HasPrefix(sub, blocked) {
			return true
		}
	}
	return false
}

// requiredFlags maps command prefixes to flags that must be present when run from the agent.
var requiredFlags = map[string][]string{
	"auth login":         {"--no-browser", "--app-name"},
	"auth signup":        {"--no-browser", "--app-name"},
	"application create": {"--region"},
}

// validateRequiredFlags checks that commands include required flags when run from the agent.
func validateRequiredFlags(command string) error {
	sub := strings.TrimPrefix(command, "algolia ")
	if sub == command {
		return nil
	}
	for prefix, flags := range requiredFlags {
		if strings.HasPrefix(sub, prefix) {
			for _, flag := range flags {
				if !strings.Contains(command, flag) {
					return fmt.Errorf("%s requires %s when run from the agent", prefix, flag)
				}
			}
		}
	}
	return nil
}

// validateCommand checks that a command string does not contain dangerous shell metacharacters.
func validateCommand(command string) error {
	if isBlockedCommand(command) {
		return fmt.Errorf("command is not allowed from the agent")
	}
	if err := validateRequiredFlags(command); err != nil {
		return err
	}
	for _, pattern := range []string{"&&", "||", ";", "$(", "`"} {
		if strings.Contains(command, pattern) {
			return fmt.Errorf("command contains disallowed shell operator: %s", pattern)
		}
	}
	return nil
}

// executeCommand runs a command string inside a PTY so that the child process
// sees a real terminal (IsTerminal returns true). Output is tee'd to the user's
// terminal and captured for the agent context.
func executeCommand(command string) (string, error) {
	if command == "" {
		return "", fmt.Errorf("empty command")
	}
	if err := validateCommand(command); err != nil {
		return "", err
	}
	command = forceJSONOutput(command)
	command = replaceAlgoliaBinary(command)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Env = append(os.Environ(), "PAGER=cat", "ALGOLIA_NO_PROMPT=1")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}
	defer ptmx.Close()

	// Tee PTY output to both the real terminal and a buffer.
	var buf bytes.Buffer
	_, _ = io.Copy(io.MultiWriter(os.Stdout, &buf), ptmx)

	_ = cmd.Wait()
	return strings.TrimSpace(stripANSI(buf.String())), nil
}

// replaceAlgoliaBinary replaces "algolia" at command positions with the actual binary path
// (e.g. "./algolia" in dev). Only replaces at the start and after pipes.
func replaceAlgoliaBinary(command string) string {
	bin := os.Args[0]
	if bin == "algolia" {
		return command
	}
	if strings.HasPrefix(command, "algolia ") {
		command = bin + command[len("algolia"):]
	}
	command = strings.ReplaceAll(command, "| algolia ", "| "+bin+" ")
	return command
}

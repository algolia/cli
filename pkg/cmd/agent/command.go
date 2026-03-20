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

// validateCommand checks that a command string does not contain dangerous shell metacharacters.
func validateCommand(command string) error {
	if isBlockedCommand(command) {
		return fmt.Errorf("command is not allowed from the agent")
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
	command = replaceAlgoliaBinary(command)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Env = append(os.Environ(), "PAGER=cat")

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

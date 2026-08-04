package cmdutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/iostreams"
)

func Test_ApplyNonInteractive_DefaultsToJSON(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	printFlags := NewPrintFlags()

	ApplyNonInteractive(io, printFlags)

	assert.True(t, io.GetNeverPrompt())
	assert.False(t, io.CanPrompt())
	assert.Equal(t, "json", *printFlags.OutputFormat)
}

func Test_ApplyNonInteractive_KeepsExplicitOutput(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	printFlags := NewPrintFlags()
	*printFlags.OutputFormat = "jsonpath={.id}"

	ApplyNonInteractive(io, printFlags)

	assert.Equal(t, "jsonpath={.id}", *printFlags.OutputFormat)
}

// The spinner writes to stderr even on a TTY, so a non-interactive run has to
// silence it too.
func Test_ApplyNonInteractive_DisablesProgressIndicator(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetProgressIndicatorEnabled(true)

	ApplyNonInteractive(io, NewPrintFlags())

	assert.False(t, io.GetProgressIndicatorEnabled())
}

// Commands with no output flags (e.g. `auth signup`) pass a nil PrintFlags.
func Test_ApplyNonInteractive_NilPrintFlags(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	ApplyNonInteractive(io, nil)

	assert.False(t, io.CanPrompt())
}

// The narration has to survive somewhere: `auth login` prints the authorize URL
// through it, and nothing else can complete the flow.
func Test_RedirectHumanOutput(t *testing.T) {
	io, _, stdout, stderr := iostreams.Test()

	restore := RedirectHumanOutput(io)
	fmt.Fprintln(io.Out, "Waiting for authentication...")
	restore()
	fmt.Fprintln(io.Out, `{"success":true}`)

	assert.Equal(t, "{\"success\":true}\n", stdout.String())
	assert.Equal(t, "Waiting for authentication...\n", stderr.String())
}

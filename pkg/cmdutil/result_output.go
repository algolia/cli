package cmdutil

import (
	"fmt"
	"io"

	"github.com/algolia/cli/pkg/iostreams"
)

// PrintRunSummary prints a structured summary when an output format is
// requested, otherwise it falls back to a human-readable line.
func PrintRunSummary(
	ios *iostreams.IOStreams,
	printFlags *PrintFlags,
	summary interface{},
	human string,
) error {
	if printFlags != nil && printFlags.HasStructuredOutput() {
		return printFlags.Print(ios, summary)
	}
	_, err := fmt.Fprintln(ios.Out, human)
	return err
}

// DiscardHumanOutput drops the progress narration a command writes to stdout,
// so the only thing it emits is the structured document. Diagnostics written
// straight to stderr (warnings, errors) are untouched. The returned function
// restores the original writer and must run before the document is printed.
func DiscardHumanOutput(ios *iostreams.IOStreams) func() {
	original := ios.Out
	ios.Out = io.Discard

	return func() {
		ios.Out = original
	}
}

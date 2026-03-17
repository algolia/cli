package cmdutil

import (
	"fmt"

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

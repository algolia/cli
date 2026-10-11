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

// ApplyNonInteractive turns off every prompt and defaults the output to JSON,
// leaving an explicit --output untouched. printFlags may be nil, for commands
// that expose no output flags.
func ApplyNonInteractive(ios *iostreams.IOStreams, printFlags *PrintFlags) {
	ios.SetNeverPrompt(true)
	ios.SetProgressIndicatorEnabled(false)

	if printFlags != nil && printFlags.OutputFormat != nil && !printFlags.HasStructuredOutput() {
		*printFlags.OutputFormat = "json"
	}
}

// RedirectHumanOutput sends the progress narration a command writes to stdout
// to stderr instead, so stdout carries nothing but the structured document.
//
// The returned function restores the original writer and must run before the
// document is printed.
func RedirectHumanOutput(ios *iostreams.IOStreams) func() {
	original := ios.Out
	ios.Out = ios.ErrOut

	return func() {
		ios.Out = original
	}
}

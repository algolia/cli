// Package interactive builds request bodies by prompting the user for each
// field of a Go struct via reflection. Input is gathered through the Prompter
// interface so the traversal can be unit-tested without a terminal.
package interactive

import (
	"os"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

// Prompter is the input surface used by the reflective Builder. Production code
// uses SurveyPrompter; tests use a scripted fake.
type Prompter interface {
	// Input reads a free-text line. validate, when non-nil, is run on the entry;
	// the implementation re-prompts until it passes (SurveyPrompter) or surfaces
	// its error (ScriptedPrompter). Pass nil for no validation. An empty entry
	// means "skip" for optional fields and zero for required scalars.
	Input(label string, validate func(string) error) (string, error)
	// Confirm asks a yes/no question.
	Confirm(label string) (bool, error)
	// Select asks the user to pick exactly one option and returns the chosen
	// label.
	Select(label string, options []string) (string, error)
	// MultiSelect asks the user to pick zero or more options and returns the
	// chosen 0-based indexes.
	MultiSelect(label string, options []string) ([]int, error)
}

// SurveyPrompter implements Prompter using survey/v2 via the repo's pkg/prompt
// wrappers (which are swappable in tests).
type SurveyPrompter struct {
	io *iostreams.IOStreams
}

// NewSurveyPrompter returns a Prompter that reads and writes the given streams.
func NewSurveyPrompter(io *iostreams.IOStreams) *SurveyPrompter {
	return &SurveyPrompter{io: io}
}

func (s *SurveyPrompter) surveyOpts() []survey.AskOpt {
	return []survey.AskOpt{
		survey.WithStdio(
			fileReader{s.io.In},
			fileWriter{s.io.Out, os.Stdout.Fd()},
			fileWriter{s.io.ErrOut, os.Stderr.Fd()},
		),
	}
}

func (s *SurveyPrompter) Input(label string, validate func(string) error) (string, error) {
	var out string
	opts := s.surveyOpts()
	if validate != nil {
		// survey re-prompts in place (showing the error, with the line editable)
		// until the validator passes.
		opts = append(opts, survey.WithValidator(func(ans interface{}) error {
			str, ok := ans.(string)
			if !ok {
				return nil
			}
			return validate(str)
		}))
	}
	err := prompt.SurveyAskOne(&survey.Input{Message: label}, &out, opts...)
	return out, err
}

func (s *SurveyPrompter) Confirm(label string) (bool, error) {
	var out bool
	err := prompt.SurveyAskOne(&survey.Confirm{Message: label}, &out, s.surveyOpts()...)
	return out, err
}

func (s *SurveyPrompter) Select(label string, options []string) (string, error) {
	var out string
	err := prompt.SurveyAskOne(&survey.Select{Message: label, Options: options}, &out, s.surveyOpts()...)
	return out, err
}

func (s *SurveyPrompter) MultiSelect(label string, options []string) ([]int, error) {
	var out []int
	err := prompt.SurveyAskOne(&survey.MultiSelect{Message: label, Options: options}, &out, s.surveyOpts()...)
	return out, err
}

// fileReader/fileWriter adapt iostreams to survey's terminal file interfaces.
// Fd() returns the real stdio descriptor so survey can toggle raw mode on a
// genuine TTY; in tests the prompt vars are stubbed so Fd is never used.
type fileReader struct {
	r interface{ Read([]byte) (int, error) }
}

func (f fileReader) Read(p []byte) (int, error) { return f.r.Read(p) }
func (f fileReader) Fd() uintptr                { return os.Stdin.Fd() }

type fileWriter struct {
	w  interface{ Write([]byte) (int, error) }
	fd uintptr
}

func (f fileWriter) Write(p []byte) (int, error) { return f.w.Write(p) }
func (f fileWriter) Fd() uintptr                 { return f.fd }

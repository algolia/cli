// Package interactive builds request bodies by prompting the user for each
// field of a Go struct via reflection. Input is gathered through the Prompter
// interface so the traversal can be unit-tested without a terminal.
package interactive

import (
	"fmt"
	"os"
	"slices"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

// PromptConfig holds the optional help, default and completion of a prompt, and each prompt kind ignores the fields it cannot use.
type PromptConfig struct {
	Help    string
	Default string
	Suggest func(toComplete string) []string
}

// PromptOption sets one field of a PromptConfig.
type PromptOption func(*PromptConfig)

// WithHelp sets the description printed above the prompt.
func WithHelp(help string) PromptOption {
	return func(c *PromptConfig) { c.Help = help }
}

// WithDefault sets the value used when the user submits an empty line.
func WithDefault(def string) PromptOption {
	return func(c *PromptConfig) { c.Default = def }
}

// WithSuggest sets the tab-completion callback, and nil disables completion.
func WithSuggest(suggest func(toComplete string) []string) PromptOption {
	return func(c *PromptConfig) { c.Suggest = suggest }
}

// ApplyPromptOptions folds opts into a PromptConfig, in order, so a later option overrides an earlier one.
func ApplyPromptOptions(opts []PromptOption) PromptConfig {
	var c PromptConfig
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// Prompter is the input surface used by the reflective Builder. Production code
// uses SurveyPrompter; tests use a scripted fake.
type Prompter interface {
	// Input re-prompts (SurveyPrompter) or returns the error (ScriptedPrompter) when validate fails, and an empty entry means skip.
	Input(label string, validate func(string) error, opts ...PromptOption) (string, error)
	// Confirm asks a yes/no question.
	Confirm(label string, opts ...PromptOption) (bool, error)
	// Select asks the user to pick exactly one option and returns the chosen
	// label.
	Select(label string, options []string, opts ...PromptOption) (string, error)
	// MultiSelect asks the user to pick zero or more options and returns the
	// chosen 0-based indexes.
	MultiSelect(label string, options []string, opts ...PromptOption) ([]int, error)
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

func (s *SurveyPrompter) Input(label string, validate func(string) error, promptOpts ...PromptOption) (string, error) {
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
	cfg := s.printHelp(promptOpts)
	err := prompt.SurveyAskOne(&survey.Input{Message: label, Default: cfg.Default, Suggest: cfg.Suggest}, &out, opts...)
	return out, err
}

func (s *SurveyPrompter) Confirm(label string, promptOpts ...PromptOption) (bool, error) {
	var out bool
	s.printHelp(promptOpts)
	err := prompt.SurveyAskOne(&survey.Confirm{Message: label}, &out, s.surveyOpts()...)
	return out, err
}

func (s *SurveyPrompter) Select(label string, options []string, promptOpts ...PromptOption) (string, error) {
	var out string
	cfg := s.printHelp(promptOpts)
	err := prompt.SurveyAskOne(&survey.Select{Message: label, Options: options, Default: selectDefault(options, cfg.Default)}, &out, s.surveyOpts()...)
	return out, err
}

func (s *SurveyPrompter) MultiSelect(label string, options []string, promptOpts ...PromptOption) ([]int, error) {
	var out []int
	s.printHelp(promptOpts)
	err := prompt.SurveyAskOne(&survey.MultiSelect{Message: label, Options: options}, &out, s.surveyOpts()...)
	return out, err
}

// printHelp prints the description above the prompt because survey's own Help only adds a "[? for help]" hint.
func (s *SurveyPrompter) printHelp(promptOpts []PromptOption) PromptConfig {
	cfg := ApplyPromptOptions(promptOpts)
	if cfg.Help != "" {
		cs := s.io.ColorScheme()
		fmt.Fprintln(s.io.Out, cs.Gray(cfg.Help))
	}
	return cfg
}

// selectDefault returns the preselected option, or nil when def is not one of the options because survey rejects an unknown default.
func selectDefault(options []string, def string) any {
	if !slices.Contains(options, def) {
		return nil
	}
	return def
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

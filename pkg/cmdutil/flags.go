package cmdutil

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Pretty much the whole code in this file is copied from the GitHub CLI

// NilStringFlag defines a new flag with a string pointer receiver.
// This helps distinguishing `--flag ""` from not setting the flag at all.
func NilStringFlag(
	cmd *cobra.Command,
	p **string,
	name string,
	shorthand string,
	usage string,
) *pflag.Flag {
	return cmd.Flags().VarPF(newStringValue(p), name, shorthand, usage)
}

// StringEnumFlag defines a new string flag restricted to allowed options
func StringEnumFlag(
	cmd *cobra.Command,
	p *string,
	name, shorthand, defaultValue string,
	options []string,
	usage string,
) *pflag.Flag {
	*p = defaultValue
	val := &enumValue{string: p, options: options}
	f := cmd.Flags().
		VarPF(val, name, shorthand, fmt.Sprintf("%s: %s", usage, formatValuesForUsageDocs(options)))
	_ = cmd.RegisterFlagCompletionFunc(
		name,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return options, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return f
}

type enumValue struct {
	string  *string
	options []string
}

func (e *enumValue) Set(value string) error {
	if !isIncluded(value, e.options) {
		return fmt.Errorf("valid values are %s", formatValuesForUsageDocs(e.options))
	}
	*e.string = value
	return nil
}

func (e *enumValue) String() string {
	return *e.string
}

func (e *enumValue) Type() string {
	return "string"
}

func isIncluded(value string, opts []string) bool {
	for _, opt := range opts {
		if strings.EqualFold(opt, value) {
			return true
		}
	}
	return false
}

func formatValuesForUsageDocs(values []string) string {
	return fmt.Sprintf("{%s}", strings.Join(values, "|"))
}

type stringValue struct {
	string **string
}

func (s *stringValue) Set(value string) error {
	*s.string = &value
	return nil
}

func (s *stringValue) String() string {
	if s.string == nil || *s.string == nil {
		return ""
	}
	return **s.string
}

func (s *stringValue) Type() string {
	return "string"
}

func newStringValue(p **string) *stringValue {
	return &stringValue{p}
}

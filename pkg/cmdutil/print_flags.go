package cmdutil

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/printers"
)

// IsNoCompatiblePrinterError returns true if it is a not a compatible printer
// otherwise it will return false
func IsNoCompatiblePrinterError(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(NoCompatiblePrinterError)
	return ok
}

type PrintFlags struct {
	JSONPrintFlags     *JSONPrintFlags
	TemplatePrintFlags *GoTemplatePrintFlags

	OutputFormat        *string
	OutputFlagSpecified func() bool
}

// NoCompatiblePrinterError is a struct that contains error information.
// It will be constructed when a invalid printing format is provided
type NoCompatiblePrinterError struct {
	OutputFormat   *string
	AllowedFormats []string
	Options        interface{}
}

func (e NoCompatiblePrinterError) Error() string {
	output := ""
	if e.OutputFormat != nil {
		output = *e.OutputFormat
	}

	sort.Strings(e.AllowedFormats)
	return fmt.Sprintf("unable to match a printer suitable for the output format %q, allowed formats are: %s", output, strings.Join(e.AllowedFormats, ","))
}

func (f *PrintFlags) AllowedFormats() []string {
	ret := []string{}
	ret = append(ret, f.JSONPrintFlags.AllowedFormats()...)
	ret = append(ret, f.TemplatePrintFlags.AllowedFormats()...)
	return ret
}

func (f *PrintFlags) ToPrinter() (printers.Printer, error) {
	outputFormat := ""
	if f.OutputFormat != nil {
		outputFormat = *f.OutputFormat
	}

	if f.JSONPrintFlags != nil {
		if p, err := f.JSONPrintFlags.ToPrinter(outputFormat); !IsNoCompatiblePrinterError(err) {
			return p, err
		}
	}

	if f.TemplatePrintFlags != nil {
		if p, err := f.TemplatePrintFlags.ToPrinter(outputFormat); !IsNoCompatiblePrinterError(err) {
			return p, err
		}
	}

	return nil, NoCompatiblePrinterError{OutputFormat: f.OutputFormat, AllowedFormats: f.AllowedFormats()}
}

func (f *PrintFlags) AddFlags(cmd *cobra.Command) {
	f.JSONPrintFlags.AddFlags(cmd)
	f.TemplatePrintFlags.AddFlags(cmd)

	if f.OutputFormat != nil {
		cmd.Flags().StringVarP(f.OutputFormat, "output", "o", *f.OutputFormat, fmt.Sprintf(`Output format. One of: (%s).`, strings.Join(f.AllowedFormats(), ", ")))
		_ = cmd.Flags().SetAnnotation("output", "IsPrint", []string{"true"})
		if f.OutputFlagSpecified == nil {
			f.OutputFlagSpecified = func() bool {
				return cmd.Flag("output").Changed
			}
		}
	}
}

// WithDefaultOutput sets a default output format if one is not provided through a flag value
func (f *PrintFlags) WithDefaultOutput(output string) *PrintFlags {
	f.OutputFormat = &output
	return f
}

// NewPrintFlags returns a default *PrintFlags
func NewPrintFlags() *PrintFlags {
	outputFormat := ""

	return &PrintFlags{
		OutputFormat: &outputFormat,

		JSONPrintFlags:     NewJSONPrintFlags(),
		TemplatePrintFlags: NewGoTemplatePrintFlags(),
	}
}

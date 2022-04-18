package cmdutil

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/printers"
)

func (f *JSONPrintFlags) AllowedFormats() []string {
	if f == nil {
		return []string{}
	}
	return []string{"json"}
}

type JSONPrintFlags struct{}

func (f *JSONPrintFlags) ToPrinter(outputFormat string) (printers.Printer, error) {
	var printer printers.Printer

	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "json":
		printer = &printers.JSONPrinter{}
	default:
		return nil, NoCompatiblePrinterError{OutputFormat: &outputFormat, AllowedFormats: f.AllowedFormats()}
	}

	return printer, nil
}

func (f *JSONPrintFlags) AddFlags(c *cobra.Command) {}

func NewJSONPrintFlags() *JSONPrintFlags {
	return &JSONPrintFlags{}
}

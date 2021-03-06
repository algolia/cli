package cmdutil

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/algolia/cli/pkg/printers"
	"github.com/spf13/cobra"
)

var templateFormats = map[string]bool{
	"template":         true,
	"go-template":      true,
	"go-template-file": true,
	"templatefile":     true,
}

type GoTemplatePrintFlags struct {
	// indicates if it is OK to ignore missing keys for rendering
	// an output template.
	AllowMissingKeys *bool
	TemplateArgument *string
}

func (f *GoTemplatePrintFlags) AllowedFormats() []string {
	formats := make([]string, 0, len(templateFormats))
	for format := range templateFormats {
		formats = append(formats, format)
	}
	sort.Strings(formats)
	return formats
}

func (f *GoTemplatePrintFlags) ToPrinter(templateFormat string) (printers.Printer, error) {
	if (f.TemplateArgument == nil || len(*f.TemplateArgument) == 0) && len(templateFormat) == 0 {
		return nil, NoCompatiblePrinterError{Options: f, OutputFormat: &templateFormat}
	}

	templateValue := ""

	if f.TemplateArgument == nil || len(*f.TemplateArgument) == 0 {
		for format := range templateFormats {
			format = format + "="
			if strings.HasPrefix(templateFormat, format) {
				templateValue = templateFormat[len(format):]
				templateFormat = format[:len(format)-1]
				break
			}
		}
	} else {
		templateValue = *f.TemplateArgument
	}

	if _, supportedFormat := templateFormats[templateFormat]; !supportedFormat {
		return nil, NoCompatiblePrinterError{OutputFormat: &templateFormat, AllowedFormats: f.AllowedFormats()}
	}

	if len(templateValue) == 0 {
		return nil, fmt.Errorf("template format specified but no template given")
	}

	if templateFormat == "templatefile" || templateFormat == "go-template-file" {
		data, err := ioutil.ReadFile(templateValue)
		if err != nil {
			return nil, fmt.Errorf("error reading --template %s, %v", templateValue, err)
		}

		templateValue = string(data)
	}

	p, err := printers.NewGoTemplatePrinter([]byte(templateValue))
	if err != nil {
		return nil, fmt.Errorf("error parsing template %s, %v", templateValue, err)
	}

	allowMissingKeys := true
	if f.AllowMissingKeys != nil {
		allowMissingKeys = *f.AllowMissingKeys
	}

	p.AllowMissingKeys(allowMissingKeys)
	return p, nil
}

func (f *GoTemplatePrintFlags) AddFlags(c *cobra.Command) {
	if f.TemplateArgument != nil {
		c.Flags().StringVar(f.TemplateArgument, "template", *f.TemplateArgument, "Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].")
		_ = c.Flags().SetAnnotation("template", "IsPrint", []string{"true"})
		_ = c.MarkFlagFilename("template")
	}
	if f.AllowMissingKeys != nil {
		c.Flags().BoolVar(f.AllowMissingKeys, "allow-missing-template-keys", *f.AllowMissingKeys, "If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats.")
		_ = c.Flags().SetAnnotation("allow-missing-template-keys", "IsPrint", []string{"true"})
	}
}

func NewGoTemplatePrintFlags() *GoTemplatePrintFlags {
	allowMissingKeysPtr := true
	templateValuePtr := ""

	return &GoTemplatePrintFlags{
		TemplateArgument: &templateValuePtr,
		AllowMissingKeys: &allowMissingKeysPtr,
	}
}

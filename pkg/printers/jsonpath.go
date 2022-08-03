package printers

import (
	"bytes"
	"fmt"

	"github.com/algolia/cli/pkg/iostreams"
	"k8s.io/client-go/util/jsonpath"
)

// JSONPathPrinter is an implementation of ResourcePrinter which formats data with a JSONPath template.
var _ Printer = &JSONPathPrinter{}

type JSONPathPrinter struct {
	rawTemplate string
	JSONPath    *jsonpath.JSONPath
}

func NewJSONPathPrinter(tmpl string) (*JSONPathPrinter, error) {
	j := jsonpath.New("out")
	if err := j.Parse(tmpl); err != nil {
		return nil, err
	}
	return &JSONPathPrinter{
		rawTemplate: tmpl,
		JSONPath:    j,
	}, nil
}

// Print formats the interface with the JSONPath Template.
func (j *JSONPathPrinter) Print(ios *iostreams.IOStreams, data interface{}) error {
	err := j.JSONPath.Execute(ios.Out, data)
	if err == nil {
		_, _ = ios.Out.Write([]byte("\n"))
	} else {
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, "Error executing template: %v. Printing more information for debugging the template:\n", err)
		fmt.Fprintf(buf, "\ttemplate was:\n\t\t%v\n", j.rawTemplate)
		fmt.Fprintf(buf, "\tobject given to jsonpath engine was:\n\t\t%#v\n\n", data)
		return fmt.Errorf("error executing jsonpath %q: %v", j.rawTemplate, buf.String())
	}
	return nil
}

// AllowMissingKeys tells the template engine if missing keys are allowed.
func (j *JSONPathPrinter) AllowMissingKeys(allow bool) {
	j.JSONPath.AllowMissingKeys(allow)
}

// EnableJSONOutput enables JSON output.
func (j *JSONPathPrinter) EnableJSONOutput(enable bool) {
	j.JSONPath.EnableJSONOutput(enable)
}

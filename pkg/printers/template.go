package printers

import (
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/algolia/cli/pkg/iostreams"
)

// GoTemplatePrinter is an implementation of ResourcePrinter which formats data with a Go Template.
var _ Printer = &GoTemplatePrinter{}

type GoTemplatePrinter struct {
	rawTemplate string
	template    *template.Template
}

func NewGoTemplatePrinter(tmpl []byte) (*GoTemplatePrinter, error) {
	t, err := template.New("output").Parse(string(tmpl))
	if err != nil {
		return nil, err
	}
	return &GoTemplatePrinter{
		rawTemplate: string(tmpl),
		template:    t,
	}, nil
}

// AllowMissingKeys tells the template engine if missing keys are allowed.
func (p *GoTemplatePrinter) AllowMissingKeys(allow bool) {
	if allow {
		p.template.Option("missingkey=default")
	} else {
		p.template.Option("missingkey=error")
	}
}

// Print formats the interface with the Go Template.
func (p *GoTemplatePrinter) Print(ios *iostreams.IOStreams, data interface{}) error {
	dataM, err := json.Marshal(data)
	if err != nil {
		return err
	}

	out := map[string]interface{}{}
	if err := json.Unmarshal(dataM, &out); err != nil {
		return err
	}
	if err = p.safeExecute(ios.Out, out); err != nil {
		// It is way easier to debug this stuff when it shows up in
		// stdout instead of just stdin. So in addition to returning
		// a nice error, also print useful stuff with the writer.
		fmt.Fprintf(ios.ErrOut, "Error executing template: %v. Printing more information for debugging the template:\n", err)
		fmt.Fprintf(ios.ErrOut, "\ttemplate was:\n\t\t%v\n", p.rawTemplate)
		fmt.Fprintf(ios.ErrOut, "\traw data was:\n\t\t%v\n", string(dataM))
		fmt.Fprintf(ios.ErrOut, "\tobject given to template engine was:\n\t\t%+v\n\n", out)
		return fmt.Errorf("error executing template %q: %v", p.rawTemplate, err)
	}
	return nil
}

// safeExecute tries to execute the template, but catches panics and returns an error
// should the template engine panic.
func (p *GoTemplatePrinter) safeExecute(w io.Writer, obj interface{}) error {
	var panicErr error
	retErr := func() error {
		defer func() {
			if x := recover(); x != nil {
				panicErr = fmt.Errorf("caught panic: %+v", x)
			}
		}()
		return p.template.Execute(w, obj)
	}()
	if panicErr != nil {
		return panicErr
	}
	return retErr
}

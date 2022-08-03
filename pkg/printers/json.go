package printers

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/jsoncolor"
)

// JSONPrinter is an implementation of Printer which outputs an object as JSON.
var _ Printer = &JSONPrinter{}

type JSONPrinter struct{}

type JSONPrinterOptions struct {
	Template string
}

// Print is an implementation of Printer.Print which simply writes the object to the Writer.
func (p *JSONPrinter) Print(ios *iostreams.IOStreams, data interface{}) error {
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	w := ios.Out
	if ios.ColorEnabled() {
		return jsoncolor.Write(w, &buf, "  ")
	}

	_, err := io.Copy(w, &buf)
	return err
}

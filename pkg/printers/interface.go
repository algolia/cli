package printers

import (
	"github.com/algolia/cli/pkg/iostreams"
)

// PrinterFunc is a function that can print interfaces
type PrinterFunc func(interface{}, *iostreams.IOStreams) error

// Print implements Printer
func (fn PrinterFunc) Print(obj interface{}, io *iostreams.IOStreams) error {
	return fn(obj, io)
}

// Printer is an interface that knows how to print interfaces.
type Printer interface {
	// Print receives an interface, formats it and prints it.
	Print(*iostreams.IOStreams, interface{}) error
}

// PrintOptions struct defines a struct for various print options
type PrintOptions struct {
	NoHeaders    bool
	ColumnLabels []string
}

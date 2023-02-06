package iostreams

import (
	"errors"
	"golang.org/x/term"
)

func ttySize() (int, int, error) {
	// in case we are not in a terminal
	if !term.IsTerminal(0) {
		return -1, -1, errors.New("not a terminal")
	}

	return term.GetSize(0)
}

package iostreams

import (
	"errors"
	"os"
)

func (s *IOStreams) EnableVirtualTerminalProcessing() error {
	return nil
}

func hasAlternateScreenBuffer(hasTrueColor bool) bool {
	// on Windows we just assume that alternate screen buffer is supported if we
	// enabled virtual terminal processing, which in turn enables truecolor
	return hasTrueColor
}

func enableVirtualTerminalProcessing(f *os.File) error {
	return errors.New("not implemented")
}

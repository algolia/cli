package cmdutil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

const maxCapacity = 1024 * 5120 // 5MB

func ReadFile(filename string, stdin io.ReadCloser) ([]byte, error) {
	if filename == "-" {
		b, err := io.ReadAll(stdin)
		_ = stdin.Close()
		return b, err
	}

	return os.ReadFile(filename)
}

// JSONPayload returns inline --json or the contents of --file (use "-" for stdin).
// The two sources are mutually exclusive. Both empty returns "", nil.
func JSONPayload(jsonStr, filename string, stdin io.ReadCloser) (string, error) {
	if jsonStr != "" && filename != "" {
		return "", FlagErrorf("--json cannot be combined with --file")
	}
	if filename == "" {
		return jsonStr, nil
	}

	b, err := ReadFile(filename, stdin)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	return string(bytes.TrimSpace(b)), nil
}

func ScanFile(filename string, stdin io.ReadCloser) (*bufio.Scanner, error) {
	var scanner *bufio.Scanner

	if filename == "-" {
		scanner = bufio.NewScanner(stdin)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		scanner = bufio.NewScanner(f)
	}

	buffer := make([]byte, maxCapacity)
	scanner.Buffer(buffer, maxCapacity)
	return scanner, nil
}

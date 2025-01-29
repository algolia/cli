package cmdutil

import (
	"bufio"
	"io"
	"os"
)

const maxCapacity = 1024 * 1024 // 1MB

func ReadFile(filename string, stdin io.ReadCloser) ([]byte, error) {
	if filename == "-" {
		b, err := io.ReadAll(stdin)
		_ = stdin.Close()
		return b, err
	}

	return os.ReadFile(filename)
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

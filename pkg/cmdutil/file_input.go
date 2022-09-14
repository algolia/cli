package cmdutil

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

const maxCapacity = 100 * 1024

func ReadFile(filename string, stdin io.ReadCloser) ([]byte, error) {
	if filename == "-" {
		b, err := ioutil.ReadAll(stdin)
		_ = stdin.Close()
		return b, err
	}

	return ioutil.ReadFile(filename)
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

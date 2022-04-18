package cmdutil

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

func ReadFile(filename string, stdin io.ReadCloser) ([]byte, error) {
	if filename == "-" {
		b, err := ioutil.ReadAll(stdin)
		_ = stdin.Close()
		return b, err
	}

	return ioutil.ReadFile(filename)
}

func ScanFile(filename string, stdin io.ReadCloser) (*bufio.Scanner, error) {
	if filename == "-" {
		return bufio.NewScanner(stdin), nil
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return bufio.NewScanner(f), nil
}

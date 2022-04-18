package test

import (
	"bytes"
)

type CmdOut struct {
	OutBuf     *bytes.Buffer
	ErrBuf     *bytes.Buffer
	BrowsedURL string
}

func (c CmdOut) String() string {
	return c.OutBuf.String()
}

func (c CmdOut) Stderr() string {
	return c.ErrBuf.String()
}

type OutputStub struct {
	Out   []byte
	Error error
}

func (s OutputStub) Output() ([]byte, error) {
	if s.Error != nil {
		return s.Out, s.Error
	}
	return s.Out, nil
}

func (s OutputStub) Run() error {
	if s.Error != nil {
		return s.Error
	}
	return nil
}

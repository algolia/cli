package main

import (
	"os"

	"github.com/algolia/cli/pkg/cmd/root"
)

func main() {
	code := root.Execute()
	os.Exit(int(code))
}

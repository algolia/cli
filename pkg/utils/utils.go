package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Pluralize returns the plural form of a given string
func Pluralize(num int, thing string) string {
	if num == 1 {
		return fmt.Sprintf("%d %s", num, thing)
	}
	return fmt.Sprintf("%d %ss", num, thing)
}

// MakePath creates a path if it doesn't exist
func MakePath(path string) error {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// Contains check if a slice contains a given string
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ToKebabCase converts a string to kebab case
func ToKebabCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")

	return strings.ToLower(snake)
}

// Convert comma separated string values to slice
func StringToSlice(str string) []string {
	return strings.Split(strings.ReplaceAll(str, " ", ""), ",")
}

// Convert slice of string to comma separated string
func SliceToString(str []string) string {
	return strings.Join(str, ", ")
}

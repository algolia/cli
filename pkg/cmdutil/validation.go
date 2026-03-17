package cmdutil

import (
	"unicode"
)

// ValidateNoControlChars rejects identifier-like inputs containing control
// characters, which are easy to miss in automation and should never be valid.
func ValidateNoControlChars(field, value string) error {
	for _, r := range value {
		if unicode.IsControl(r) {
			return FlagErrorf("%s must not contain control characters", field)
		}
	}
	return nil
}

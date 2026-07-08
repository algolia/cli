package interactive

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// The validator builders below produce func(string) error values that are
// passed to Prompter.Input. With SurveyPrompter these become survey validators,
// so an invalid entry re-prompts in place (with the message shown and the line
// editable) instead of aborting the whole build. The functions are pure, so the
// validation rules can be unit-tested directly without a terminal.

// requiredString rejects an empty entry.
func requiredString(label string) func(string) error {
	return func(s string) error {
		if s == "" {
			return fmt.Errorf("%s is required", label)
		}
		return nil
	}
}

// intValidator accepts an integer that fits in bits. An empty entry is allowed
// only when optional (it means skip / zero); otherwise it is required.
func intValidator(label string, bits int, optional bool) func(string) error {
	return func(s string) error {
		if s == "" {
			if optional {
				return nil
			}
			return fmt.Errorf("%s is required", label)
		}
		if _, err := strconv.ParseInt(s, 10, bits); err != nil {
			return fmt.Errorf("must be a whole number")
		}
		return nil
	}
}

// floatValidator accepts a number that fits in bits (32 or 64), so a float32
// overflow is rejected at input rather than later becoming +Inf. An empty entry
// is allowed only when optional.
func floatValidator(label string, bits int, optional bool) func(string) error {
	return func(s string) error {
		if s == "" {
			if optional {
				return nil
			}
			return fmt.Errorf("%s is required", label)
		}
		if _, err := strconv.ParseFloat(s, bits); err != nil {
			return fmt.Errorf("must be a number")
		}
		return nil
	}
}

// countValidator accepts a non-negative whole number for "how many ..." prompts.
// An empty entry is allowed (it means none).
func countValidator() func(string) error {
	return func(s string) error {
		if s == "" {
			return nil
		}
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("must be a whole number")
		}
		if n < 0 {
			return fmt.Errorf("must be zero or greater")
		}
		return nil
	}
}

// boolValidator accepts a parseable boolean. It is used for optional *bool, so
// an empty entry (skip) is allowed.
func boolValidator() func(string) error {
	return func(s string) error {
		if s == "" {
			return nil
		}
		if _, err := strconv.ParseBool(s); err != nil {
			return fmt.Errorf("must be true or false")
		}
		return nil
	}
}

// jsonValidator accepts valid JSON. An empty entry (skip) is allowed.
func jsonValidator() func(string) error {
	return func(s string) error {
		if s == "" {
			return nil
		}
		if !json.Valid([]byte(s)) {
			return fmt.Errorf("must be valid JSON")
		}
		return nil
	}
}

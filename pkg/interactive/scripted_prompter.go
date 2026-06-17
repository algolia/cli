package interactive

import "strings"

// ScriptedPrompter is a deterministic Prompter for tests. Answers are keyed by a
// substring of the prompt label rather than by call order, so a test does not
// break when the target struct (for example a generated SDK type) adds or
// reorders fields. A prompt whose label matches no key returns a safe default:
// empty input, a false confirm, the first option, or no multi-selection. Unknown
// or newly added fields are therefore simply skipped.
//
// Matching rule: a key matches a label when the label equals the key or contains
// it as a substring. An exact match wins over any substring match; among
// substring matches the longest (most specific) key wins, so the result is
// deterministic regardless of map iteration order.
type ScriptedPrompter struct {
	Inputs       map[string]string
	Confirms     map[string]bool
	Selects      map[string]string
	MultiSelects map[string][]int
}

var _ Prompter = (*ScriptedPrompter)(nil)

func (p *ScriptedPrompter) Input(label string, validate func(string) error) (string, error) {
	v, _ := lookup(p.Inputs, label)
	// Run the validator on the scripted answer so invalid scripted data surfaces
	// deterministically (no infinite re-prompt loop). In production survey would
	// re-prompt the user instead.
	if validate != nil {
		if err := validate(v); err != nil {
			return "", err
		}
	}
	return v, nil
}

func (p *ScriptedPrompter) Confirm(label string) (bool, error) {
	if v, ok := lookup(p.Confirms, label); ok {
		return v, nil
	}
	return false, nil
}

func (p *ScriptedPrompter) Select(label string, options []string) (string, error) {
	if v, ok := lookup(p.Selects, label); ok {
		return v, nil
	}
	if len(options) > 0 {
		return options[0], nil
	}
	return "", nil
}

func (p *ScriptedPrompter) MultiSelect(label string, options []string) ([]int, error) {
	if v, ok := lookup(p.MultiSelects, label); ok {
		return v, nil
	}
	return nil, nil
}

// lookup returns the value whose key exactly equals label, or failing that the
// value whose key is the longest substring of label. The bool reports a match.
func lookup[T any](m map[string]T, label string) (T, bool) {
	var zero T
	if m == nil {
		return zero, false
	}
	if v, ok := m[label]; ok {
		return v, true
	}
	var best T
	bestLen := -1
	for k, v := range m {
		if k == "" || !strings.Contains(label, k) {
			continue
		}
		if len(k) > bestLen {
			best, bestLen = v, len(k)
		}
	}
	return best, bestLen >= 0
}

package shared

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FlagsToRule(t *testing.T) {
	tests := []struct {
		name        string
		ruleFlags   RuleFlags
		wantsErr    bool
		wantsErrMsg string
	}{
		// Correct Rules
		{
			name:      "",
			ruleFlags: RuleFlags{},
			wantsErr:  false,
		},
		// Wrong Rules
		{
			name:        "",
			ruleFlags:   RuleFlags{},
			wantsErr:    true,
			wantsErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := FlagsToRule(tt.ruleFlags)

			if tt.wantsErr {
				assert.EqualError(t, err, tt.wantsErrMsg)
				return
			}

			assert.Equal(t, err, nil)
			assert.Equal(t, reflect.TypeOf(rule).String(), "search.Rule")
		})
	}
}

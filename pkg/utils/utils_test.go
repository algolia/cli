package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluralize(t *testing.T) {
	tests := []struct {
		name     string
		num      int
		thing    string
		expected string
	}{
		{
			name:     "zero, no plural",
			num:      0,
			thing:    "sushi",
			expected: "0 sushi",
		},
		{
			name:     "one, no plural",
			num:      1,
			thing:    "sushi",
			expected: "1 sushi",
		},
		{
			name:     "negative, no plural",
			num:      -10,
			thing:    "sushi",
			expected: "-10 sushi",
		},
		{
			name:     "positive, plural",
			num:      10,
			thing:    "sushi",
			expected: "10 sushis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Pluralize(tt.num, tt.thing))
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name         string
		inputSlice   []string
		stringToTest string
		expected     bool
	}{
		{
			name:         "contains in slice",
			inputSlice:   []string{"maguro", "otoro", "unagi"},
			stringToTest: "otoro",
			expected:     true,
		},
		{
			name:         "contains in 1 element slice",
			inputSlice:   []string{"unagi"},
			stringToTest: "unagi",
			expected:     true,
		},
		{
			name:         "empty slice",
			inputSlice:   []string{},
			stringToTest: "otoro",
			expected:     false,
		},
		{
			name:         "missing in slice",
			inputSlice:   []string{"maguro", "otoro", "unagi"},
			stringToTest: "tamago",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Contains(tt.inputSlice, tt.stringToTest))
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "full lowercase",
			input:    "iaminlowercase",
			expected: "iaminlowercase",
		},
		{
			name:     "full uppercase",
			input:    "IAMINUPPERCASE",
			expected: "iaminuppercase",
		},
		{
			name:     "camelCase",
			input:    "camelCase",
			expected: "camel-case",
		},
		{
			name:     "PascalCase",
			input:    "PascalCase",
			expected: "pascal-case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToKebabCase(tt.input))
		})
	}
}

func TestStringToSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "no element",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "one element",
			input:    "maguro",
			expected: []string{"maguro"},
		},
		{
			name:     "comma separated",
			input:    "maguro,otoro",
			expected: []string{"maguro", "otoro"},
		},
		{
			name:     "with spaces",
			input:    "maguro    ",
			expected: []string{"maguro"},
		},
		{
			name:     "comma separated with spaces",
			input:    "maguro    , otoro,   tamago",
			expected: []string{"maguro", "otoro", "tamago"},
		},
		{
			name:     "space separated",
			input:    "maguro    otoro tamago",
			expected: []string{"magurootorotamago"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, StringToSlice(tt.input))
		})
	}
}

func TestSliceToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "no element",
			input:    []string{""},
			expected: "",
		},
		{
			name:     "one element",
			input:    []string{"maguro"},
			expected: "maguro",
		},
		{
			name:     "two elements",
			input:    []string{"maguro", "otoro"},
			expected: "maguro, otoro",
		},
		{
			name:     "three elements",
			input:    []string{"maguro", "otoro", "tamago"},
			expected: "maguro, otoro, tamago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SliceToString(tt.input))
		})
	}
}

func TestSliceToReadableString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "no element",
			input:    []string{""},
			expected: "",
		},
		{
			name:     "one element",
			input:    []string{"maguro"},
			expected: "maguro",
		},
		{
			name:     "two element",
			input:    []string{"maguro", "otoro"},
			expected: "maguro and otoro",
		},
		{
			name:     "three element",
			input:    []string{"maguro", "otoro", "tamago"},
			expected: "maguro, otoro and tamago",
		},
		{
			name:     "five element",
			input:    []string{"one", "two", "three", "four", "five"},
			expected: "one, two, three, four and five",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SliceToReadableString(tt.input))
		})
	}
}

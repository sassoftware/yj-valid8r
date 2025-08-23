package tests

import (
	"testing"

	yjvalid8r_lib "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-lib"
)

func TestCheckTabsAndWhitespacesFinder(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
		hasWarn  bool
	}{
		{
			name:     "No Issues",
			input:    "key: value\nanother: line\n",
			hasError: false,
			hasWarn:  false,
		},
		{
			name:     "Tab Present",
			input:    "\tkey: value\nanother: line\n",
			hasError: true,
			hasWarn:  false,
		},
		{
			name:     "Trailing Whitespace",
			input:    "key: value \nanother: line\t\n",
			hasError: false,
			hasWarn:  true,
		},
		{
			name:     "Tab and Trailing Whitespace",
			input:    "\tkey: value \nanother: line\t\n",
			hasError: true,
			hasWarn:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := yjvalid8r_lib.CheckTabsAndWhitespacesFinder([]byte(tt.input))

			if tt.hasError && len(result.Errors) == 0 {
				t.Errorf("Expected errors, but found none")
			}
			if !tt.hasError && len(result.Errors) != 0 {
				t.Errorf("Did not expect errors, but found some: %v", result.Errors)
			}
			if tt.hasWarn && len(result.Warnings) == 0 {
				t.Errorf("Expected warnings, but found none")
			}
			if !tt.hasWarn && len(result.Warnings) != 0 {
				t.Errorf("Did not expect warnings, but found some: %v", result.Warnings)
			}
		})
	}
}

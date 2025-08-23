package tests

import (
	"fmt"
	"testing"

	yjvalid8r_lib "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-lib"
)

func TestDetectDataType(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Valid JSON",
			input:    []byte(`{"name": "John", "age": 30}`),
			expected: yjvalid8r_lib.DataTypeJSON,
		},
		{
			name:     "Valid YAML",
			input:    []byte("name: John\nage: 30"),
			expected: yjvalid8r_lib.DataTypeYAML,
		},
		{
			name:     "Unknown Format - Plain Text",
			input:    []byte("Just a plain string"),
			expected: yjvalid8r_lib.DataTypeUNKNOWN,
		},
		{
			name:     "Unknown Format - Malformed JSON",
			input:    []byte(`{"name": "John", "age":}`),
			expected: yjvalid8r_lib.DataTypeUNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := yjvalid8r_lib.DetectDataType(tt.input)
			if result != tt.expected {
				fmt.Println(result)
				t.Errorf("DetectDataType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsUnknownDataType(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "Valid JSON",
			input:    []byte(`{"key": "value"}`),
			expected: false,
		},
		{
			name:     "Valid YAML",
			input:    []byte("key: value"),
			expected: false,
		},
		{
			name:     "Unknown Format",
			input:    []byte("invalid content"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := yjvalid8r_lib.IsUnknownDataType(tt.input)
			if result != tt.expected {
				t.Errorf("IsUnknownDataType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

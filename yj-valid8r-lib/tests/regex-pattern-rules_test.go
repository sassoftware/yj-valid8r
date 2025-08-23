package tests

import (
	"os"
	"testing"

	yjvalid8r_lib "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-lib"
)

func TestRegexPatternRulesFinder(t *testing.T) {
	// Set up an environment variable for testing
	os.Setenv("ENV_VAR", "dummy-value")
	defer os.Unsetenv("ENV_VAR") // Clean up

	data := `
key1: ${ENV_VAR}
key2: ${UNSET_VAR}
key3: value3
`

	regexRules := []yjvalid8r_lib.RegexPatternRules{
		{
			Name:  "Find ${VAR} Patterns",
			Regex: `\$\{(\w+)\}`,
			CheckEnv: &yjvalid8r_lib.RegexPatternRulesCheckEnvConfig{
				Enabled: true,
				Strict:  true,
			},
		},
	}

	results, hasStrictError := yjvalid8r_lib.RegexPatternRulesFinder(regexRules, []byte(data))

	// Validate result length
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	result := results[0]

	if result.Name != "Find ${VAR} Patterns" {
		t.Errorf("Unexpected rule name: %s", result.Name)
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(result.Data))
	}

	// ENV_VAR should be found
	found := false
	for _, msg := range result.EnvValues {
		if msg == "ENV_VAR=dummy-value" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ENV_VAR to be found in environment values")
	}

	// UNSET_VAR should produce an error
	errorFound := false
	for _, err := range result.Errors {
		if err == "Environment variable not found: UNSET_VAR" {
			errorFound = true
			break
		}
	}
	if !errorFound {
		t.Error("Expected missing UNSET_VAR to produce an error")
	}

	if !hasStrictError {
		t.Error("Expected strict error to be true due to missing environment variable")
	}
}

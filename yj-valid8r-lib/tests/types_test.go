package tests

import (
	"encoding/json"
	"reflect"
	"testing"

	yjvalid8r_lib "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-lib"
)

func TestConstants(t *testing.T) {
	if yjvalid8r_lib.DataTypeYAML != "yaml" {
		t.Errorf("DataTypeYAML = %v, want 'yaml'", yjvalid8r_lib.DataTypeYAML)
	}
	if yjvalid8r_lib.MessageTypeError != "error" {
		t.Errorf("MessageTypeError = %v, want 'error'", yjvalid8r_lib.MessageTypeError)
	}
	if yjvalid8r_lib.MessageTypeWarning != "warning" {
		t.Errorf("MessageTypeWarning = %v, want 'warning'", yjvalid8r_lib.MessageTypeWarning)
	}
}

func TestSchemaValidationMessageJSON(t *testing.T) {
	msg := yjvalid8r_lib.SchemaValidationMessage{
		Type:    yjvalid8r_lib.MessageTypeError,
		Message: "Invalid field",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal SchemaValidationMessage: %v", err)
	}

	var unmarshaled yjvalid8r_lib.SchemaValidationMessage
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SchemaValidationMessage: %v", err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Errorf("Unmarshaled message does not match original.\nGot: %+v\nWant: %+v", unmarshaled, msg)
	}
}

func TestRegexPatternRulesDefaultValues(t *testing.T) {
	rule := yjvalid8r_lib.RegexPatternRules{
		Name:  "Test Rule",
		Regex: `^abc$`,
		CheckEnv: &yjvalid8r_lib.RegexPatternRulesCheckEnvConfig{
			Enabled: true,
			Strict:  false,
		},
	}

	if !rule.CheckEnv.Enabled {
		t.Error("Expected CheckEnv.Enabled to be true")
	}
	if rule.CheckEnv.Strict {
		t.Error("Expected CheckEnv.Strict to be false")
	}
}

func TestSearchPathsDef(t *testing.T) {
	path := yjvalid8r_lib.SearchPathsDef{
		PathName: "Get Port",
		PathKey:  "spec.ports[].targetPort",
	}

	if path.PathName != "Get Port" {
		t.Errorf("Expected PathName = Get Port, got %s", path.PathName)
	}
	if path.PathKey != "spec.ports[].targetPort" {
		t.Errorf("Expected PathKey = spec.ports[].targetPort, got %s", path.PathKey)
	}
}

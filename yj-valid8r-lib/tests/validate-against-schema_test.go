package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	yjvalid8r_lib "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-lib"
)

func writeTempSchemaFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile := filepath.Join(t.TempDir(), "schema.json")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp schema: %v", err)
	}
	return "file://" + tmpFile
}

func TestValidateAgainstSchemaFinder_Valid(t *testing.T) {
	schema := `{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"age": { "type": "number" }
		},
		"required": ["name", "age"]
	}`

	data := `
name: John Doe
age: 30
`

	schemaURL := writeTempSchemaFile(t, schema)

	results, err := yjvalid8r_lib.ValidateAgainstSchemaFinder(schemaURL, []byte(data))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for _, msg := range results {
		t.Errorf("Expected no validation messages, got: %+v", msg)
	}
}

func TestValidateAgainstSchemaFinder_InvalidField(t *testing.T) {
	schema := `{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"age": { "type": "number" }
		},
		"required": ["name", "age"]
	}`

	data := `
name: John Doe
age: thirty
`

	schemaURL := writeTempSchemaFile(t, schema)

	results, err := yjvalid8r_lib.ValidateAgainstSchemaFinder(schemaURL, []byte(data))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	foundError := false
	for _, msg := range results {
		if msg.Type == "error" && strings.Contains(msg.Message, "age") {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Errorf("Expected validation error related to 'age', got: %+v", results)
	}
}

func TestValidateAgainstSchemaFinder_InvalidYAML(t *testing.T) {
	schema := `{
		"type": "object",
		"properties": {
			"name": { "type": "string" }
		}
	}`

	// invalid YAML (unclosed string)
	data := `
name: "John Doe
`

	schemaURL := writeTempSchemaFile(t, schema)

	_, err := yjvalid8r_lib.ValidateAgainstSchemaFinder(schemaURL, []byte(data))
	if err == nil || !strings.Contains(err.Error(), "parse yaml") {
		t.Errorf("Expected YAML parse error, got: %v", err)
	}
}

func TestValidateAgainstSchemaFinder_MissingSchema(t *testing.T) {
	// path to non-existing file
	schemaURL := "file:///tmp/does_not_exist_schema.json"

	_, err := yjvalid8r_lib.ValidateAgainstSchemaFinder(schemaURL, []byte(`name: Test`))
	if err == nil || !strings.Contains(err.Error(), "schema does not exist") {
		t.Errorf("Expected missing schema error, got: %v", err)
	}
}

package yjvalid8r_lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// ValidateAgainstSchemaFinder validates the input data against a JSON schema from a URL.
// It returns a slice of SchemaValidationMessage and any error encountered.
func ValidateAgainstSchemaFinder(schemaURL string, dataBytes []byte) ([]SchemaValidationMessage, error) {
	exists, normalizedURL := checkURLExists(schemaURL)
	if !exists {
		return nil, fmt.Errorf("schema does not exist or is unreachable: %s", normalizedURL)
	}

	// Parse YAML into yaml.Node to preserve line info
	var rootNode yaml.Node
	if err := yaml.Unmarshal(dataBytes, &rootNode); err != nil {
		return nil, fmt.Errorf("parse yaml/json into node: %w", err)
	}

	// Also unmarshal into map[string]interface{} for JSON schema validation
	var dataMap map[string]interface{}
	if err := yaml.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("parse yaml/json into map: %w", err)
	}

	props, err := extractTopLevelSchemaProperties(normalizedURL)
	if err != nil {
		return nil, fmt.Errorf("schema load failed: %w", err)
	}

	var messages []SchemaValidationMessage

	if topLevelFieldMismatch(props, dataMap) {
		messages = append(messages, SchemaValidationMessage{
			Type:    MessageTypeWarning,
			Message: "schema appears irrelevant: no overlapping top-level fields between schema and data.",
		})
	}

	jsonDataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("to json: %w", err)
	}

	schemaLoader := gojsonschema.NewReferenceLoader(normalizedURL)
	documentLoader := gojsonschema.NewBytesLoader(jsonDataBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		hint := `Note:
- For absolute $ref (file:///...), run from any folder.
- For relative $ref (file://...), run from the schema's folder.`
		return nil, fmt.Errorf("validation error: %w (%s)", err, hint)
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			dataType := DetectDataType(dataBytes)
			if dataType == DataTypeJSON || dataType == DataTypeYAML {
				// JSON path like: workloads.1.flows.0.processors.4.switch.cases.1.processors.5.log.level
				path := strings.Split(desc.Field(), ".")
				node := findNodeByPath(&rootNode, path)

				if node != nil {
					messages = append(messages, SchemaValidationMessage{
						Type:    "error",
						Message: fmt.Sprintf("Line %d: %s: %s", node.Line, desc.Field(), desc.Description()),
					})
				} else {
					messages = append(messages, SchemaValidationMessage{
						Type:    "error",
						Message: fmt.Sprintf("Line unknown: %s: %s", desc.Field(), desc.Description()),
					})
				}
			}

		}
	}

	return messages, nil
}

// Walk YAML node by JSON path segments, e.g. ["workloads", "1", "flows", "0", "processors", "4"]
func findNodeByPath(root *yaml.Node, path []string) *yaml.Node {
	if len(path) == 0 {
		return root
	}

	key := path[0]
	rest := path[1:]

	switch root.Kind {
	case yaml.DocumentNode:
		return findNodeByPath(root.Content[0], path)
	case yaml.MappingNode:
		// Mapping nodes: Content is [keyNode, valueNode, keyNode, valueNode, ...]
		for i := 0; i < len(root.Content); i += 2 {
			k := root.Content[i]
			v := root.Content[i+1]
			if k.Value == key {
				return findNodeByPath(v, rest)
			}
		}
	case yaml.SequenceNode:
		// Sequence nodes: Content is a list of nodes, key must be an index
		idx, err := strconv.Atoi(key)
		if err == nil && idx < len(root.Content) {
			return findNodeByPath(root.Content[idx], rest)
		}
	}
	return nil
}

func extractTopLevelSchemaProperties(schemaURL string) (map[string]interface{}, error) {
	var schemaBytes []byte
	var err error

	if strings.HasPrefix(schemaURL, "file://") {
		localPath := strings.TrimPrefix(schemaURL, "file://")
		schemaBytes, err = os.ReadFile(localPath)
		if err != nil {
			return nil, fmt.Errorf("read schema file: %w", err)
		}
	} else {
		resp, err := http.Get(schemaURL)
		if err != nil {
			return nil, fmt.Errorf("fetch schema: %w", err)
		}
		defer resp.Body.Close()

		schemaBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read schema response: %w", err)
		}
	}

	var schemaMap map[string]interface{}
	if err := json.Unmarshal(schemaBytes, &schemaMap); err != nil {
		return nil, fmt.Errorf("parse schema JSON: %w", err)
	}

	props, _ := schemaMap["properties"].(map[string]interface{})
	return props, nil
}

func topLevelFieldMismatch(schemaProps map[string]interface{}, data map[string]interface{}) bool {
	for key := range data {
		if _, ok := schemaProps[key]; ok {
			return false // ✅ At least one field matches
		}
	}
	return true // ❌ No matches found — schema likely irrelevant
}

func checkURLExists(pathOrURL string) (bool, string) {
	parsedURL, err := url.Parse(pathOrURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Scheme == "file" {
		// Handle local path
		normalizedPath := strings.TrimPrefix(pathOrURL, "file://")
		_, err := os.Stat(normalizedPath)
		if err == nil {
			// File exists, return with file:// prefix
			return true, "file://" + normalizedPath
		}
		return false, "file://" + normalizedPath
	}

	// Handle remote URL
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(pathOrURL)
	if err != nil {
		return false, pathOrURL
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, pathOrURL
}

package yjvalid8r_lib

import (
	"bytes"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// DetectDataType returns "json", "yaml", or "unknown".
func DetectDataType(data []byte) string {
	trimmed := bytes.TrimSpace(data)

	// 1) If it's valid JSON, we're done.
	if json.Valid(trimmed) {
		return DataTypeJSON
	}

	// 2) Otherwise, parse it into a YAML node.
	var doc yaml.Node
	if err := yaml.Unmarshal(trimmed, &doc); err != nil {
		return DataTypeUNKNOWN
	}

	// The actual root node is the first child of the document node.
	if len(doc.Content) == 0 {
		return DataTypeUNKNOWN
	}
	root := doc.Content[0]

	// 3) Accept only block-style mappings (objects) or sequences (arrays).
	switch root.Kind {
	case yaml.MappingNode:
		// reject if it’s flow-style (i.e. JSON-like `{ ... }`)
		if root.Style == yaml.FlowStyle {
			return DataTypeUNKNOWN
		}
		return DataTypeYAML

	case yaml.SequenceNode:
		// reject if it’s flow-style (i.e. JSON-like `[ ... ]`)
		if root.Style == yaml.FlowStyle {
			return DataTypeUNKNOWN
		}
		return DataTypeYAML

	default:
		// scalars, etc. → unknown
		return DataTypeUNKNOWN
	}
}

// IsUnknownDataType checks if the given data has an unrecognized format (neither JSON nor YAML).
func IsUnknownDataType(data []byte) bool {
	return DetectDataType(data) == DataTypeUNKNOWN
}

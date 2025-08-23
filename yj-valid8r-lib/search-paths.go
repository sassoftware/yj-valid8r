package yjvalid8r_lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// SearchPathsFinder searches for values in the input data at the specified paths.
// It returns structured results for each search path and any error encountered.
func SearchPathsFinder(input []byte, paths []SearchPathsDef) ([]SearchPathsOutput, error) {

	var data interface{}
	if err := yaml.Unmarshal(input, &data); err != nil {
		return nil, fmt.Errorf("parse yaml/json into map: %w", err)
	}

	var outputs []SearchPathsOutput

	for _, path := range paths {
		results := resolvePath(data, path.PathKey)
		outputs = append(outputs, SearchPathsOutput{
			PathName: path.PathName,
			PathKey:  path.PathKey,
			Results:  results,
		})
	}

	return outputs, nil
}

func resolvePath(data interface{}, path string) []SearchPathsOutputResultItem {
	segments := strings.Split(path, ".")
	return resolve(data, segments, "")
}

func resolve(data interface{}, segments []string, currentPath string) []SearchPathsOutputResultItem {
	if len(segments) == 0 {
		rawStr := marshalToString(data)
		return []SearchPathsOutputResultItem{{FullPath: currentPath, Raw: rawStr}}
	}

	current := segments[0]
	rest := segments[1:]

	if strings.HasSuffix(current, "[]") {
		key := strings.TrimSuffix(current, "[]")
		if m, ok := data.(map[string]interface{}); ok {
			if arr, ok := m[key].([]interface{}); ok {
				var results []SearchPathsOutputResultItem
				for i, item := range arr {
					path := fmt.Sprintf("%s%s[%d]", currentPathPrefix(currentPath), key, i)
					results = append(results, resolve(item, rest, path)...)
				}
				return results
			}
		}
	} else if strings.Contains(current, "[") && strings.HasSuffix(current, "]") {
		// Handles mixed key[index], e.g., "emails[0]"
		parts := strings.SplitN(current, "[", 2)
		key := parts[0]
		indexStr := strings.TrimSuffix(parts[1], "]")

		if m, ok := data.(map[string]interface{}); ok {
			if arr, ok := m[key].([]interface{}); ok {
				if index, ok := tryParseArrayIndex(indexStr); ok && index < len(arr) {
					path := fmt.Sprintf("%s%s[%d]", currentPathPrefix(currentPath), key, index)
					return resolve(arr[index], rest, path)
				}
			}
		}
	} else {
		if m, ok := data.(map[string]interface{}); ok {
			if val, exists := m[current]; exists {
				path := fmt.Sprintf("%s%s", currentPathPrefix(currentPath), current)
				return resolve(val, rest, path)
			}
		}
	}

	// fallback for recursive search
	var found []SearchPathsOutputResultItem
	switch v := data.(type) {
	case []interface{}:
		for i, item := range v {
			path := fmt.Sprintf("%s[%d]", currentPath, i)
			found = append(found, resolve(item, segments, path)...)
		}
	case map[string]interface{}:
		for key, val := range v {
			var path string
			if key == segments[0] {
				path = fmt.Sprintf("%s%s", currentPathPrefix(currentPath), key)
				found = append(found, resolve(val, segments[1:], path)...)
			} else {
				path = fmt.Sprintf("%s%s", currentPathPrefix(currentPath), key)
				found = append(found, resolve(val, segments, path)...)
			}
		}
	}
	return found
}

func currentPathPrefix(currentPath string) string {
	if currentPath == "" {
		return ""
	}
	return currentPath + "."
}

func tryParseArrayIndex(segment string) (int, bool) {
	var index int
	n, err := fmt.Sscanf(segment, "%d", &index)
	return index, n == 1 && err == nil
}

func marshalToString(val interface{}) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ") // Pretty print, optional
	if err := enc.Encode(val); err != nil {
		return fmt.Sprintf("%v", val) // fallback
	}
	return strings.TrimSpace(buf.String())
}

package yjvalid8r_lib

import (
	"fmt"
	"strings"
)

// CheckTabsAndWhitespacesFinder checks the input data for unwanted tabs or whitespace characters.
// It returns a WhitespaceCheckResult with validation status and messages.
func CheckTabsAndWhitespacesFinder(dataByte []byte) WhitespaceCheckResult {
	var result WhitespaceCheckResult
	lines := strings.Split(string(dataByte), "\n")
	for i, line := range lines {
		for _, ch := range line {
			if ch == '\t' {
				result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Tab character found.", i+1))
				break
			} else if ch != ' ' {
				break
			}
		}
		if len(line) > 0 && (strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t")) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Line %d: Trailing whitespace found.", i+1))
		}
	}

	if len(result.Errors) != 0 {
		result.Messages = append(result.Messages, "Tab issues found in document. Note: Tab characters are not allowed for indentation in either JSON or YAML. They are only valid within string values, not for structuring or formatting the document.")
	} else if len(result.Warnings) != 0 {
		result.Messages = append(result.Messages, "Whitespace issues found in document.")
	}
	// else {
	// 	result.Messages = append(result.Messages, "No whitespace or tab issues found in data.")
	// }

	return result
}

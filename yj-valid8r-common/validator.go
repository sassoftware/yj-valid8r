package internal

import (
	"fmt"
	"strings"

	validator "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-lib"
)

func InitValidation(
	schemas []string,
	dataBytes []byte,
	whitespace bool,
	regexPatterns []validator.RegexPatternRules,
	pathSearch []validator.SearchPathsDef,
	pluginPaths string,
) ValidationResponse {
	summary := ValidationSummary{}
	var regexFindings []validator.RegexPatternRulesOutput
	var pathSearchFindings []validator.SearchPathsOutput
	results := make([]SchemaResult, 0, len(schemas))
	hasError := false

	if whitespace {
		wsResult := validator.CheckTabsAndWhitespacesFinder(dataBytes)
		if len(wsResult.Errors) > 0 {
			summary.Valid = false
		}
		summary.Errors = append(summary.Errors, wsResult.Errors...)
		summary.Warnings = append(summary.Warnings, wsResult.Warnings...)
		summary.Messages = append(summary.Messages, wsResult.Messages...)
	}

	if len(regexPatterns) > 0 {
		regexFindings, hasError = validator.RegexPatternRulesFinder(regexPatterns, dataBytes)
		if hasError {
			summary.Errors = append(summary.Errors, "Environment variable(s) not set. Strict mode is true.")
		}
	}

	if len(pathSearch) > 0 {
		var err error
		pathSearchFindings, err = validator.SearchPathsFinder(dataBytes, pathSearch)
		if err != nil {
			hasError = true
			summary.Messages = append(summary.Messages, fmt.Sprintf("parse yaml/json into node: %v", err))
		}
	}

	if len(schemas) > 0 {
		for _, schemaPath := range schemas {
			messages, err := validator.ValidateAgainstSchemaFinder(schemaPath, dataBytes)
			if err != nil {
				results = append(results, SchemaResult{
					Schema: schemaPath,
					Valid:  false,
					Errors: []string{err.Error()},
				})
				hasError = true
				continue
			}

			var errors []string
			var warnings []string

			for _, msg := range messages {
				switch msg.Type {
				case validator.MessageTypeError:
					errors = append(errors, msg.Message)
				case validator.MessageTypeWarning:
					warnings = append(warnings, msg.Message)
				}
			}

			if len(errors) > 0 {
				hasError = true
				results = append(results, SchemaResult{
					Schema:   schemaPath,
					Valid:    false,
					Errors:   errors,
					Warnings: warnings,
				})
			} else {
				results = append(results, SchemaResult{
					Schema:   schemaPath,
					Valid:    true,
					Warnings: warnings,
				})
			}
		}
	} else {
		summary.Messages = append(summary.Messages, "No schema(s) provided.")
	}

	summary.Valid = !hasError
	summary.ValidationDataType = strings.ToUpper(validator.DetectDataType(dataBytes))

	pluginResults := UsePlugin(pluginPaths, dataBytes)

	resp := ValidationResponse{
		SchemaResults:     results,
		ValidationSummary: summary,
		RegexPatterns:     regexFindings,
		PathSearchOutput:  pathSearchFindings,
		PluginResults:     pluginResults,
	}

	return resp
}

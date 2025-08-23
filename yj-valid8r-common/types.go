package internal

import (
	"time"

	validator "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-lib"
)

// ValidationRequest: input from cli and web
type ValidationRequest struct {
	CLIOutputFormat         string                        `json:"-" yaml:"cliOutputFormat"` // omit from JSON
	Schemas                 []string                      `json:"schemas" yaml:"schemas"`
	Data                    string                        `json:"data" yaml:"data"`
	CheckTrailingWhitespace *bool                         `json:"checkTrailingWhitespace" yaml:"checkTrailingWhitespace"`
	StrictValidation        *bool                         `json:"-" yaml:"strictValidation"` // omit from JSON
	RegexPatternRules       []validator.RegexPatternRules `json:"regexPatternRules" yaml:"regexPatternRules"`
	SearchPaths             []validator.SearchPathsDef    `json:"searchPaths" yaml:"searchPaths"`
	Plugins                 string                        `json:"plugins" yaml:"plugins"`
}

// ValidationResponse: output from cli and web
type SchemaResult struct {
	Schema   string   `json:"schema"`
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type ValidationSummary struct {
	ValidationDataType string   `json:"validationDataType"`
	Valid              bool     `json:"valid"`
	Errors             []string `json:"errors,omitempty"`
	Warnings           []string `json:"warnings,omitempty"`
	Messages           []string `json:"messages,omitempty"`
}

type PluginResult struct {
	Name          string        `json:"name"`
	Messages      []string      `json:"messages,omitempty"`
	Warnings      []string      `json:"warnings,omitempty"`
	Errors        []string      `json:"errors,omitempty"`
	LoadError     string        `json:"load_error,omitempty"` // Load/init error
	ExecutionTime time.Duration `json:"execution_time"`
}

type ValidationResponse struct {
	ValidationSummary ValidationSummary                   `json:"validationSummary"`
	SchemaResults     []SchemaResult                      `json:"schemaResults,omitempty"`
	RegexPatterns     []validator.RegexPatternRulesOutput `json:"regexPatterns,omitempty"`
	PathSearchOutput  []validator.SearchPathsOutput       `json:"pathSearchOutput,omitempty"`
	PluginResults     []PluginResult                      `json:"pluginResults,omitempty"`
}

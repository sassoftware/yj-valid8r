package yjvalid8r_lib

const (
	// DataTypeYAML indicates that the input data is in YAML format.
	DataTypeYAML string = "yaml"
	// DataTypeJSON indicates that the input data is in JSON format.
	DataTypeJSON string = "json"
	// DataTypeUNKNOWN indicates that the input data format is unrecognized (neither JSON nor YAML).
	DataTypeUNKNOWN string = "unknown"
)

// ValidationMessageType represents the type of message produced during validation: error or warning.
type ValidationMessageType string

const (
	// MessageTypeError indicates a critical validation error.
	MessageTypeError ValidationMessageType = "error"

	// MessageTypeWarning indicates a non-critical issue or suggestion.
	MessageTypeWarning ValidationMessageType = "warning"
)

// SchemaValidationMessage represents a single validation result message.
type SchemaValidationMessage struct {
	Type    ValidationMessageType // The severity of the message: error or warning.
	Message string                // A human-readable description of the validation issue.
}

// RegexPatternRulesCheckEnvConfig defines configuration options for validating environment variables.
type RegexPatternRulesCheckEnvConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled"` // If true, enables environment variable validation.
	Strict  bool `json:"strict" yaml:"strict"`   // If true, enables strict matching of environment values.
}

// RegexPatternRules defines a rule for extracting or validating data using regular expressions.
type RegexPatternRules struct {
	Name     string                           `json:"name" yaml:"name"`         // Name of the regex rule.
	Regex    string                           `json:"regex" yaml:"regex"`       // Regular expression pattern.
	CheckEnv *RegexPatternRulesCheckEnvConfig `json:"checkEnv" yaml:"checkEnv"` // Optional environment variable validation config.
}

// RegexPatternRulesOutput contains the results of applying a single regex pattern rule.
type RegexPatternRulesOutput struct {
	Name               string   `json:"name"`                // Name of the regex rule.
	CheckEnv           bool     `json:"checkEnv"`            // Indicates if environment validation was enabled.
	CheckEnvStrictMode bool     `json:"checkEnvStrictMode"`  // Indicates if strict environment validation was used.
	Data               []string `json:"data"`                // Matched data from the input.
	Errors             []string `json:"errors,omitempty"`    // List of errors encountered during validation.
	EnvValues          []string `json:"envValues,omitempty"` // Environment variable values that matched the pattern.
	Messages           []string `json:"messages,omitempty"`  // Additional context or informational messages.
}

// SearchPathsDef defines a configuration for searching specific paths in structured data.
type SearchPathsDef struct {
	PathName string `json:"pathName" yaml:"pathName"` // User-friendly label for the search path.
	PathKey  string `json:"pathKey" yaml:"pathKey"`   // Dot notation key used to traverse the data structure.
}

// SearchPathsOutputResultItem represents a single match found during a search operation.
type SearchPathsOutputResultItem struct {
	FullPath string `json:"fullPath"` // Full dot-notated path to the matched value.
	Raw      string `json:"raw"`      // Raw value found at the specified path.
}

// SearchPathsOutput contains the complete result of a search operation for a specific path definition.
type SearchPathsOutput struct {
	PathName string                        `json:"pathName"` // Label of the search path used.
	PathKey  string                        `json:"pathKey"`  // Key used for data traversal.
	Results  []SearchPathsOutputResultItem `json:"results"`  // List of matched items found.
}

// WhitespaceCheckResult contains the result of checking for trailing whitespace or tab characters.
type WhitespaceCheckResult struct {
	Errors   []string `json:"errors,omitempty"`   // List of error messages related to whitespace.
	Warnings []string `json:"warnings,omitempty"` // List of warnings related to formatting.
	Messages []string `json:"messages,omitempty"` // General messages or suggestions.
}

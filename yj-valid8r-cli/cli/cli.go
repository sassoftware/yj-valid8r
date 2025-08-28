package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	internal "github.com/sassoftware/yj-valid8r/yj-valid8r-common"

	validator "github.com/sassoftware/yj-valid8r/yj-valid8r-lib"

	"gopkg.in/yaml.v3"
)

func StartCLI(configPath, flagData, flagCLIOutputFormat string, flagPlugins string,
	schemaList []string,
	flagStrictValidationMode *bool,
	flagWhitespace *bool,
	flagRegexPatterns []validator.RegexPatternRules,
	flagSearchPaths []validator.SearchPathsDef,
) {
	log.Println("Validation started")

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Apply overrides or defaults
	applyOverrides(cfg, schemaList, flagData, flagCLIOutputFormat, flagPlugins, flagRegexPatterns, flagSearchPaths, flagStrictValidationMode, flagWhitespace)

	if strings.TrimSpace(cfg.Data) == "" {
		log.Fatalf("Data file paths must be specified either via config file or flags")
	}

	dataBytes, err := os.ReadFile(cfg.Data)
	if err != nil {
		log.Fatal(err)
	}

	if validator.IsUnknownDataType(dataBytes) {
		log.Fatalf("Provided data is neither valid JSON nor valid YAML. Please check if your YAML/JSON is correct.")
	}

	results := internal.InitValidation(cfg.Schemas, dataBytes, *cfg.CheckTrailingWhitespace, cfg.RegexPatternRules, cfg.SearchPaths, cfg.Plugins)

	switch cfg.CLIOutputFormat {
	case string(CLIOutputFormatTypeJSON):
		jsonLogger(results)
	case string(CLIOutputFormatTypeYAML), string(CLIOutputFormatTypeYML):
		yamlLogger(results)
	case string(CLIOutputFormatTypeLegacy):
		legacyCLILogger(results)
	case string(CLIOutputFormatTypePretty):
		prettyCLILogger(results)
	default:
		fmt.Println("cliOutputFormat should be json | yaml | legacy | pretty (default)")
		prettyCLILogger(results)
	}

	log.Println("Validation finished.")

	if *cfg.StrictValidation && !results.ValidationSummary.Valid {
		os.Exit(1)
	}
}

// applyOverrides applies command line overrides or defaults to config
func applyOverrides(
	cfg *internal.ValidationRequest,
	schemaList []string,
	flagData, flagCLIOutputFormat, flagPlugins string,
	flagVarPatterns []validator.RegexPatternRules,
	flagSearchPaths []validator.SearchPathsDef,
	flagStrictValidationMode, flagWhitespace *bool,
) {
	if len(schemaList) > 0 {
		cfg.Schemas = schemaList
	}
	if flagData != "" {
		cfg.Data = flagData
	}
	if flagCLIOutputFormat != "" {
		cfg.CLIOutputFormat = string(CLIOutputFormatType(flagCLIOutputFormat))
	}
	if cfg.CLIOutputFormat == "" {
		cfg.CLIOutputFormat = string(CLIOutputFormatTypePretty)
	}
	if flagPlugins != "" {
		cfg.Plugins = flagPlugins
	}
	if len(flagVarPatterns) > 0 {
		cfg.RegexPatternRules = flagVarPatterns
	}
	if len(flagSearchPaths) > 0 {
		cfg.SearchPaths = flagSearchPaths
	}

	// Default StrictValidation = true
	if cfg.StrictValidation == nil {
		def := true
		cfg.StrictValidation = &def
	}
	if flagStrictValidationMode != nil {
		cfg.StrictValidation = flagStrictValidationMode
	}

	// Default CheckTrailingWhitespace = true
	if cfg.CheckTrailingWhitespace == nil {
		def := true
		cfg.CheckTrailingWhitespace = &def
	}
	if flagWhitespace != nil {
		cfg.CheckTrailingWhitespace = flagWhitespace
	}
}

func loadConfig(path string) (*internal.ValidationRequest, error) {
	if path == "" {
		return &internal.ValidationRequest{}, nil
	}

	// Check if the file exists
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("config file does not exist: %v", path)
	}

	// Read file content
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into struct
	var cfg internal.ValidationRequest
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

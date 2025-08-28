package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/sassoftware/yj-valid8r/yj-valid8r-cli/cli"
	validator "github.com/sassoftware/yj-valid8r/yj-valid8r-lib"
)

func main() {
	configPathFlag := flag.String("config", "", "Path to YAML config file")
	schemaPathsFlag := flag.String("schemas", "", "Comma-separated JSON schema files or urls")
	dataPathFlag := flag.String("data", "", "Path to YAML or JSON data file")
	cliOutputFormatFlag := flag.String("cliOutputFormat", "", "CLI output type: \"json\", \"yaml\", \"legacy\", \"pretty\"")
	strictValidationFlag := flag.Bool("strictValidation", true, "Fail if validation fails")
	checkTrailingWhitespaceFlag := flag.Bool("checkTrailingWhitespace", true, "Fail if whitespace errors")
	regexPatternRulesFlag := flag.String("regexPatternRules", "", "JSON array of regex pattern rule objects")
	searchPathsFlag := flag.String("searchPaths", "", "JSON array of path search objects")
	pluginsFlag := flag.String("plugins", "", "plugin file paths as a comma-separated or newline-separated")

	flag.Parse()

	schemaList := parseCommaList(*schemaPathsFlag)
	regexPatternRulesList := parseJSON[[]validator.RegexPatternRules](*regexPatternRulesFlag, "regexPatternRules")
	searchPathsRulesList := parseJSON[[]validator.SearchPathsDef](*searchPathsFlag, "searchPaths")

	cli.StartCLI(
		*configPathFlag,
		*dataPathFlag,
		*cliOutputFormatFlag,
		*pluginsFlag,
		schemaList,
		boolFlag(strictValidationFlag, "strictValidation"),
		boolFlag(checkTrailingWhitespaceFlag, "checkTrailingWhitespace"),
		regexPatternRulesList,
		searchPathsRulesList,
	)
}

func parseCommaList(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}

func parseJSON[T any](input string, flagName string) T {
	var result T
	if input == "" {
		// Return zero value of T (e.g., empty slice if T is a slice)
		return result
	}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		log.Fatalf("Error parsing --%s: %v\n", flagName, err)
	}
	return result
}

func boolFlag(ptr *bool, name string) *bool {
	if isFlagPassed(name) {
		return ptr
	}
	return nil
}

func isFlagPassed(name string) bool {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--"+name+"=") || arg == "--"+name {
			return true
		}
	}
	return false
}

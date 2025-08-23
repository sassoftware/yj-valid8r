package yjvalid8r_lib

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// RegexPatternRulesFinder finds and validates data using regex rules in the data.
// It returns rule-wise results and a boolean indicating if any strict errors were found.
func RegexPatternRulesFinder(data []RegexPatternRules, dataByte []byte) ([]RegexPatternRulesOutput, bool) {
	content := string(dataByte)
	var results []RegexPatternRulesOutput
	hasErrorStrictMode := false

	lines := strings.Split(content, "\n")

	for _, rule := range data {
		output := RegexPatternRulesOutput{Name: rule.Name}
		if rule.CheckEnv != nil {
			output.CheckEnv = rule.CheckEnv.Enabled
			output.CheckEnvStrictMode = rule.CheckEnv.Strict
		}

		re, err := regexp.Compile(rule.Regex)
		if err != nil {
			output.Errors = append(output.Errors, fmt.Sprintf("Invalid regex: %v", err))
		} else {
			for lineNum, line := range lines {
				matches := re.FindAllStringSubmatch(line, -1)

				for _, match := range matches {
					fullMatch := match[0]

					if len(match) > 1 {
						varName := match[1]
						fullVarName := match[0]
						output.Data = append(output.Data, fullVarName)
						output.Messages = append(output.Messages, fmt.Sprintf("%s found on line %d", fullVarName, lineNum+1))

						if rule.CheckEnv != nil && rule.CheckEnv.Enabled {
							if val, ok := os.LookupEnv(varName); ok {
								output.EnvValues = append(output.EnvValues, fmt.Sprintf("%s=%s", varName, val))
							} else {
								output.Errors = append(output.Errors, fmt.Sprintf("Environment variable not found: %s", varName))
								if rule.CheckEnv.Strict {
									hasErrorStrictMode = true
								}
							}
						}
					} else {
						output.Data = append(output.Data, fullMatch)
						output.Messages = append(output.Messages, fmt.Sprintf("%s found on line %d", fullMatch, lineNum+1))
					}
				}
			}
		}

		results = append(results, output)
	}

	return results, hasErrorStrictMode
}

package cli

import (
	"encoding/json"
	"fmt"
	"os"

	internal "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-common"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type CLIOutputFormatType string

const (
	CLIOutputFormatTypeJSON   CLIOutputFormatType = "json"
	CLIOutputFormatTypeYAML   CLIOutputFormatType = "yaml"
	CLIOutputFormatTypeYML    CLIOutputFormatType = "yml"
	CLIOutputFormatTypeLegacy CLIOutputFormatType = "legacy"
	CLIOutputFormatTypePretty CLIOutputFormatType = "pretty"
)

func jsonLogger(results internal.ValidationResponse) {
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}

func yamlLogger(results internal.ValidationResponse) {
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	var generic interface{}
	err = json.Unmarshal(jsonBytes, &generic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling JSON: %v\n", err)
		os.Exit(1)
	}

	yamlBytes, err := yaml.Marshal(generic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(yamlBytes))
}

func legacyCLILogger(results internal.ValidationResponse) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("🆗 %s", cyan(fmt.Sprintf("Validation going on for %s data type.\n", results.ValidationSummary.ValidationDataType)))

	for _, result := range results.SchemaResults {
		if !result.Valid {
			fmt.Printf("❌ %s ( ERRORS: %s | WARNINGS: %s)\n", red(fmt.Sprintf("Validated against schema '%s'", result.Schema)), red(len(result.Errors)), yellow(len(result.Warnings)))
		} else {
			fmt.Printf("✅ %s", green(fmt.Sprintf("Validated against schema '%s'\n", result.Schema)))
		}

		for _, err := range result.Errors {
			fmt.Printf("- [ERROR] %s\n", err)
		}
		for _, warning := range result.Warnings {
			fmt.Printf("- [WARNING] %s\n", warning)
		}
	}

	if len(results.RegexPatterns) > 0 {
		fmt.Println(green("✔ Regex Patterns:"))
		for _, r := range results.RegexPatterns {
			fmt.Printf("  ➡️  Name: %s | CheckEnv: %t | CheckEnvStrictMode: %t\n", r.Name, r.CheckEnv, r.CheckEnvStrictMode)

			if len(r.Data) > 0 {
				fmt.Println("     Values:")
				for _, m := range r.Data {
					fmt.Printf("        - %s\n", m)
				}
			} else {
				fmt.Printf("     - Values: %v | No patterns found.\n", r.Data)
			}

			for _, m := range r.Messages {
				fmt.Printf("     - [INFO] %v\n", m)
			}

			for _, e := range r.Errors {
				fmt.Printf("     - [ERROR] %v\n", e)
			}

			if len(r.EnvValues) > 0 {
				fmt.Println("     ENV Values:")
				for _, m := range r.EnvValues {
					fmt.Printf("     - %v\n", m)
				}
			}
		}
	}

	if len(results.PathSearchOutput) > 0 {
		fmt.Println(green("✔ Path Search Outputs:"))
		for _, r := range results.PathSearchOutput {
			fmt.Printf("  ➡️  PathName: %s | PathKey: %s\n", r.PathName, r.PathKey)

			if len(r.Results) > 0 {
				fmt.Println("     Results:")
				for _, m := range r.Results {
					fmt.Printf("        - FullPath: %s | RawData: %s\n", m.FullPath, m.Raw)
				}
			} else {
				fmt.Printf("     - Data: %v | No results found.\n", r.Results)
			}
		}
	}

	if len(results.PluginResults) > 0 {
		fmt.Println(green("✔ Plugin Results:"))
		for _, r := range results.PluginResults {
			fmt.Printf("  ➡️  Name: %s\n", r.Name)

			fmt.Printf("  %s %s\n", "   Execution Time:", r.ExecutionTime.String())

			if r.LoadError != "" {
				fmt.Printf("     Load Error: %s\n", r.LoadError)
				continue
			}

			if len(r.Messages) != 0 {
				fmt.Printf("     Messages: \n")
				for _, w := range r.Messages {
					fmt.Printf("       - %s\n", w)
				}
			}

			if len(r.Errors) != 0 {
				fmt.Printf("     Errors: \n")
				for _, w := range r.Errors {
					fmt.Printf("       - %s\n", w)
				}
			}

			if len(r.Warnings) != 0 {
				fmt.Printf("     Warnings: \n")
				for _, w := range r.Warnings {
					fmt.Printf("       - %s\n", w)
				}
			}

		}
	}

	for _, sumErr := range results.ValidationSummary.Errors {
		fmt.Printf("❌ %s\n", red(sumErr))
	}
	for _, sumWarn := range results.ValidationSummary.Warnings {
		fmt.Printf("⚠️  %s\n", yellow(sumWarn))
	}
	for _, sumMsg := range results.ValidationSummary.Messages {
		fmt.Printf("🆗 %s\n", cyan(sumMsg))
	}

	if !results.ValidationSummary.Valid {
		fmt.Printf("❌ %s\n", red("Validation failed."))
	} else {
		fmt.Printf("✅ %s\n", green("Validation successful."))
	}

}

func prettyCLILogger(results internal.ValidationResponse) {
	// Define color functions
	redBold := color.New(color.FgRed, color.Bold).SprintFunc()
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellowBold := color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	fmt.Println(white(fmt.Sprintf("Validation going on for %s data type.\n", results.ValidationSummary.ValidationDataType)))

	// Loop through each schema result
	for _, result := range results.SchemaResults {
		// Print schema header with colored status
		if !result.Valid {
			fmt.Printf("%s %s ( ERRORS: %s | WARNINGS: %s )\n",
				redBold("✖ Schema:"),
				cyan(result.Schema),
				redBold(len(result.Errors)),
				yellowBold(len(result.Warnings)),
			)
		} else {
			fmt.Printf("%s %s\n",
				greenBold("✔ Schema:"),
				cyan(result.Schema),
			)
		}

		// Print each error under this schema
		for _, errMsg := range result.Errors {
			fmt.Printf("  %s %s\n",
				redBold("ERROR"),
				white(errMsg),
			)
		}

		// Print each warning under this schema
		for _, warnMsg := range result.Warnings {
			fmt.Printf("  %s %s\n",
				yellowBold("WARNING"),
				white(warnMsg),
			)
		}

		// Add a blank line between schemas
		fmt.Println()
	}

	// If there are regex patterns, list them
	if len(results.RegexPatterns) > 0 {
		fmt.Println(greenBold("✔ Regex Patterns:"))
		for _, pattern := range results.RegexPatterns {
			fmt.Printf("  %s %s\n",
				cyan("➤ Name:"),
				white(pattern.Name),
			)
			fmt.Printf("    %s %t   %s %t\n",
				cyan("CheckEnv:"),
				pattern.CheckEnv,
				cyan("StrictMode:"),
				pattern.CheckEnvStrictMode,
			)

			// List found Data (or note if none)
			if len(pattern.Data) > 0 {
				fmt.Println("    Data:")
				for _, v := range pattern.Data {
					fmt.Printf("      - %s\n", white(v))
				}
			} else {
				fmt.Printf("    %s %v\n", cyan("Data:"), pattern.Data)
			}

			// Any informational messages
			for _, infoMsg := range pattern.Messages {
				fmt.Printf("    %s %s\n", cyan("ℹ INFO:"), white(infoMsg))
			}

			// Any errors in the pattern
			for _, errMsg := range pattern.Errors {
				fmt.Printf("    %s %s\n", redBold("ERROR:"), white(errMsg))
			}

			for _, env := range pattern.EnvValues {
				fmt.Printf("    %s %s\n", cyan("ℹ ENV VALUE:"), white(env))
			}

			fmt.Println()
		}
	}

	if len(results.PathSearchOutput) > 0 {
		fmt.Println(cyan("ℹ Path Search Outputs:"))
		for _, output := range results.PathSearchOutput {
			fmt.Printf("  %s %s\n", greenBold("PathName:"), white(output.PathName))
			fmt.Printf("  %s %s\n", greenBold("PathKey:"), white(output.PathKey))
			if len(output.Results) == 0 {
				fmt.Printf("    %s\n", yellowBold("No results found"))
			} else {
				fmt.Printf("    %s\n", greenBold("Results:"))
				for _, res := range output.Results {
					// FullPath and Raw in aligned fashion, indent nicely
					fmt.Printf("      %s %s\n", cyan("FullPath:"), white(res.FullPath))
					fmt.Printf("      %s %s\n", cyan("Raw:"), white(res.Raw))
				}
			}
			fmt.Println()
		}
	}

	if len(results.PluginResults) > 0 {
		fmt.Println(cyan("ℹ Plugin Results:"))
		for _, output := range results.PluginResults {
			fmt.Printf("  %s %s\n", greenBold("Name:"), white(output.Name))
			fmt.Printf("  %s %s\n", cyan("Execution Time:"), white(output.ExecutionTime.String()))

			if output.LoadError != "" {
				fmt.Printf("  %s %s\n", redBold("Load Error:"), white(output.LoadError))
				continue
			}

			if len(output.Messages) == 0 {
				fmt.Printf("  %s\n", yellowBold("No Messages found"))
			} else {
				fmt.Printf("  %s\n", greenBold("Messages:"))
				for _, w := range output.Messages {
					fmt.Printf("    - %s\n", white(w))
				}
			}

			if len(output.Errors) == 0 {
				fmt.Printf("  %s\n", yellowBold("No errors found"))
			} else {
				fmt.Printf("  %s\n", redBold("Errors:"))
				for _, w := range output.Errors {
					fmt.Printf("    - %s\n", white(w))
				}
			}

			if len(output.Warnings) == 0 {
				fmt.Printf("  %s\n", yellowBold("No warnings found"))
			} else {
				fmt.Printf("  %s\n", yellowBold("Warnings:"))
				for _, w := range output.Warnings {
					fmt.Printf("    - %s\n", white(w))
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Print summary-level errors
	if len(results.ValidationSummary.Errors) > 0 {
		fmt.Println(redBold("✖ Summary Errors:"))
		for _, sumErr := range results.ValidationSummary.Errors {
			fmt.Printf("  %s\n", white(sumErr))
		}
		fmt.Println()
	}

	// Print summary-level warnings
	if len(results.ValidationSummary.Warnings) > 0 {
		fmt.Println(yellowBold("⚠ Summary Warnings:"))
		for _, sumWarn := range results.ValidationSummary.Warnings {
			fmt.Printf("  %s\n", white(sumWarn))
		}
		fmt.Println()
	}

	// Print summary-level info/messages
	if len(results.ValidationSummary.Messages) > 0 {
		fmt.Println(cyan("ℹ Summary Messages:"))
		for _, sumMsg := range results.ValidationSummary.Messages {
			fmt.Printf("  %s\n", white(sumMsg))
		}
		fmt.Println()
	}

	// Final overall status
	if !results.ValidationSummary.Valid {
		fmt.Println(redBold("✖ Validation failed."))
	} else {
		fmt.Println(greenBold("✔ Validation successful."))
	}
}

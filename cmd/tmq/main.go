// Package main provides the tmq command-line interface.
//
// Tmq is a complete, standalone command-line tool for TOML.
// Like jq for JSON, yq for YAML — but for TOML.
//
// It reads TOML from a file or stdin, runs optional queries or modifications,
// and prints results to stdout. Exit codes: 0 success, 1 parse/runtime error, 2 usage error.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/azolfagharj/tmq/internal/converter"
	"github.com/azolfagharj/tmq/internal/modifier"
	"github.com/azolfagharj/tmq/internal/parser"
	"github.com/azolfagharj/tmq/internal/query"
)

// Global variables for operations
var (
	outputFormat converter.OutputFormat = converter.FormatTOML
	inplace      bool
	operation    string // "query", "set", or "delete"
	operationArg string
	dryRun       bool // Dry-run mode
)

var (
	AZ_VERSION string = "1.0.1"
	AZ_UPDATE  string = "2024-02-17"
)

// Exit codes for automation-friendly scripting
const (
	ExitSuccess       = 0 // Success
	ExitParseError    = 1 // TOML parsing or runtime error
	ExitUsageError    = 2 // Invalid arguments or usage
	ExitSecurityError = 3 // Security violation (path traversal, etc.)
	ExitFileError     = 4 // File operation error
)

func main() {
	// Parse command line flags and arguments
	var useStdin bool
	var filePaths []string

	// Validation & comparison flags
	var validateMode bool
	var compareMode bool
	var compareFile string
	var schemaFile string

	args := os.Args[1:]

	// First pass: collect positional arguments (files and operations)
	var positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			// This is a flag, will be handled in second pass
			continue
		}
		positional = append(positional, arg)
	}

	// Second pass: process flags
	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-h" || arg == "--help":
			printUsage()
			os.Exit(0)
		case arg == "--version":
			fmt.Printf("Version: %s\n", AZ_VERSION)
			fmt.Printf("Build Time: %s\n", AZ_UPDATE)
			os.Exit(0)
		case arg == "-i" || arg == "--inplace":
			inplace = true
		case arg == "--validate":
			validateMode = true
		case arg == "--compare":
			compareMode = true
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --compare flag requires a file path argument\n")
				os.Exit(2)
			}
			compareFile = args[i+1]
			i++ // Skip the file path value
		case arg == "--schema":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --schema flag requires a file path argument\n")
				os.Exit(2)
			}
			schemaFile = args[i+1]
			i++ // Skip the file path value
		case arg == "--dry-run":
			dryRun = true
		case arg == "-o" || arg == "--output":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: -o flag requires a format argument\n")
				os.Exit(2)
			}
			format, err := converter.ParseOutputFormat(args[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}
			outputFormat = format
			i++ // Skip the format value
		case strings.HasPrefix(arg, "-o="):
			formatStr := strings.TrimPrefix(arg, "-o=")
			format, err := converter.ParseOutputFormat(formatStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}
			outputFormat = format
		}
	}

	// Process positional arguments
	// Check if last argument looks like an operation
	if len(positional) > 0 {
		lastArg := positional[len(positional)-1]
		if strings.Contains(lastArg, "=") || strings.HasPrefix(lastArg, "del(") || strings.HasPrefix(lastArg, ".") || strings.Contains(lastArg, "[") || strings.Contains(lastArg, "]") {
			operationArg = lastArg
			operation = determineOperation(lastArg)
			// All preceding args are files
			filePaths = positional[:len(positional)-1]
		} else {
			// No operation specified, all args are files
			filePaths = positional
		}
	}

	// Set defaults
	if operationArg == "" {
		operation = "query"
	}
	if len(filePaths) == 0 {
		useStdin = true
	}

	// Handle bulk operations or single file/stdin
	if useStdin {
		handleSingleFile("", true, validateMode, compareMode, schemaFile, compareFile)
	} else if len(filePaths) == 1 {
		handleSingleFile(filePaths[0], false, validateMode, compareMode, schemaFile, compareFile)
	} else {
		handleBulkFiles(filePaths, validateMode)
	}
}

// determineOperation determines the type of operation from the argument
func determineOperation(arg string) string {
	if strings.Contains(arg, "=") {
		return "set"
	}
	if strings.HasPrefix(arg, "del(") && strings.HasSuffix(arg, ")") {
		return "delete"
	}
	return "query"
}

// outputData prints data in the specified format
func outputData(data interface{}, format converter.OutputFormat) {
	switch format {
	case converter.FormatTOML:
		fmt.Printf("%+v\n", data)
	case converter.FormatJSON:
		jsonStr, err := converter.ConvertToJSON(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to convert to JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(jsonStr)
	case converter.FormatYAML:
		yamlStr, err := converter.ConvertToYAML(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to convert to YAML: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(yamlStr)
	}
}

// writeTOMLFile writes TOML data back to a file
func writeTOMLFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(data)
}

// isFilePath checks if a string looks like a file path
func isFilePath(s string) bool {
	if strings.Contains(s, " ") {
		return false // Paths with spaces are unlikely
	}
	// Reject paths with directory traversal
	if strings.Contains(s, "..") {
		return false
	}
	if strings.HasPrefix(s, "/") {
		return true // Absolute paths
	}
	if strings.HasPrefix(s, "./") {
		return true // Relative paths
	}
	// Check for file extension - must have extension and look like filename
	if strings.Contains(s, ".") && !strings.HasPrefix(s, ".") && !strings.Contains(s, "(") && !strings.Contains(s, ")") {
		return true // Likely a file with extension, but not operations with parentheses
	}
	return false
}

// validateFilePath performs security validation on file paths
func validateFilePath(path string) error {
	// Prevent directory traversal
	if strings.Contains(path, "..") {
		return fmt.Errorf("directory traversal not allowed")
	}

	// Prevent absolute paths that might be dangerous
	if strings.HasPrefix(path, "/") {
		// Allow /tmp and /home but prevent /etc, /usr, /var, etc.
		if strings.HasPrefix(path, "/etc/") ||
			strings.HasPrefix(path, "/usr/") ||
			strings.HasPrefix(path, "/var/") ||
			strings.HasPrefix(path, "/proc/") ||
			strings.HasPrefix(path, "/sys/") ||
			strings.HasPrefix(path, "/dev/") {
			return fmt.Errorf("access to system directories not allowed")
		}
	}

	// Check file size (prevent extremely large files)
	if info, err := os.Stat(path); err == nil {
		const maxFileSize = 100 * 1024 * 1024 // 100MB limit
		if info.Size() > maxFileSize {
			return fmt.Errorf("file too large (max 100MB)")
		}
	}

	return nil
}

// handleValidation performs TOML validation
func handleValidation(data interface{}, schemaFile string) {
	if schemaFile != "" {
		// Schema validation (future feature)
		fmt.Fprintf(os.Stderr, "Error: Schema validation not yet implemented\n")
		os.Exit(1)
	}

	// Basic validation - if we reach here, TOML is valid
	fmt.Println("✓ TOML file is valid")
	os.Exit(0)
}

// handleComparison compares two TOML files
func handleComparison(data1 interface{}, file2 string) {
	// Validate second file path
	if err := validateFilePath(file2); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid second file path '%s'\n", file2)
		fmt.Fprintf(os.Stderr, "Details: %v\n", err)
		os.Exit(1)
	}

	// Parse second file
	p2 := parser.New()
	if err := p2.ParseFile(file2); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse second TOML file '%s'\n", file2)
		fmt.Fprintf(os.Stderr, "Details: %v\n", err)
		os.Exit(1)
	}

	data2 := p2.GetData()

	// Compare the data structures
	differences := compareTOML(data1, data2)

	if len(differences) == 0 {
		fmt.Println("✓ Files are identical")
		os.Exit(0)
	} else {
		fmt.Println("✗ Files differ:")
		for _, diff := range differences {
			fmt.Printf("  %s\n", diff)
		}
		os.Exit(1) // Non-zero exit for differences
	}
}

// compareTOML compares two TOML data structures and returns differences
func compareTOML(data1, data2 interface{}) []string {
	var differences []string

	// Convert to JSON strings for comparison (simple approach)
	json1, err1 := converter.ConvertToJSON(data1)
	json2, err2 := converter.ConvertToJSON(data2)

	if err1 != nil || err2 != nil {
		differences = append(differences, "Error converting to JSON for comparison")
		return differences
	}

	if json1 != json2 {
		differences = append(differences, "JSON representations differ")
	}

	return differences
}

// formatError formats error messages in a structured way for automation
func formatError(errorType, message, details, action string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
	if details != "" {
		fmt.Fprintf(os.Stderr, "DETAILS: %s\n", details)
	}
	if action != "" {
		fmt.Fprintf(os.Stderr, "ACTION: %s\n", action)
	}
}

// handleSingleFile processes a single file or stdin
func handleSingleFile(filePath string, useStdin bool, validateMode bool, compareMode bool, schemaFile string, compareFile string) {
	// Read TOML data
	var p *parser.Parser
	if useStdin {
		p = parser.New()
		if err := p.ParseReader(os.Stdin); err != nil {
			formatError("PARSE_ERROR", "Failed to parse TOML from stdin", err.Error(), "Check TOML syntax in piped input")
			os.Exit(ExitParseError)
		}
	} else {
		// Security: Validate file path before opening
		if err := validateFilePath(filePath); err != nil {
			formatError("SECURITY_ERROR", fmt.Sprintf("Invalid file path '%s'", filePath), err.Error(), "Use safe file paths without directory traversal")
			os.Exit(ExitSecurityError)
		}

		p = parser.New()
		if err := p.ParseFile(filePath); err != nil {
			formatError("PARSE_ERROR", fmt.Sprintf("Failed to parse TOML file '%s'", filePath), err.Error(), "Check TOML syntax, file permissions, or file existence")
			os.Exit(ExitParseError)
		}
	}

	data := p.GetData()
	if data == nil {
		data = make(map[string]interface{})
	}

	// Handle validation and comparison modes (special operations)
	if validateMode {
		handleValidation(data, schemaFile)
		return // Exit after validation
	}

	if compareMode {
		handleComparison(data, compareFile)
		return // Exit after comparison
	}

	// Convert to map for modification operations
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		formatError("DATA_ERROR", "TOML data must be a table/object for modification operations", "Data is not a valid TOML table", "Use table-based TOML structure")
		os.Exit(ExitParseError)
	}

	// Handle operations
	handleOperations(data, dataMap, filePath, useStdin)
}

// handleBulkFiles processes multiple files
func handleBulkFiles(filePaths []string, validateMode bool) {
	var hasErrors bool

	// For bulk operations, we need to handle each file
	for _, filePath := range filePaths {
		fmt.Fprintf(os.Stderr, "Processing: %s\n", filePath)

		// Validate file path
		if err := validateFilePath(filePath); err != nil {
			formatError("SECURITY_ERROR", fmt.Sprintf("Invalid file path '%s'", filePath), err.Error(), "Skipping file")
			hasErrors = true
			continue
		}

		// Parse file
		p := parser.New()
		if err := p.ParseFile(filePath); err != nil {
			formatError("PARSE_ERROR", fmt.Sprintf("Failed to parse TOML file '%s'", filePath), err.Error(), "Skipping file")
			hasErrors = true
			continue
		}

		data := p.GetData()
		if data == nil {
			data = make(map[string]interface{})
		}

		// Handle validation mode for bulk files
		if validateMode {
			fmt.Printf("✓ %s: TOML file is valid\n", filePath)
			continue
		}

		// For other operations, convert to map
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: %s: TOML data must be a table/object for modification operations\n", filePath)
			hasErrors = true
			continue
		}

		// Handle operations for this file
		if err := handleOperationsBulk(data, dataMap, filePath); err != nil {
			formatError("OPERATION_ERROR", fmt.Sprintf("Operation failed on '%s'", filePath), err.Error(), "Skipping file")
			hasErrors = true
		}
	}

	// Exit with appropriate code
	if hasErrors {
		os.Exit(1) // Some files had errors
	}
	os.Exit(0) // All files processed successfully
}

// handleOperations handles operations for single file
func handleOperations(data interface{}, dataMap map[string]interface{}, filePath string, useStdin bool) {
	// Handle operations
	switch operation {
	case "query":
		// If no operation arg, print all data
		if operationArg == "" {
			outputData(data, outputFormat)
			os.Exit(0)
		}

		// Execute query
		q, err := query.New(operationArg)
		if err != nil {
			formatError("QUERY_ERROR", fmt.Sprintf("Invalid query syntax '%s'", operationArg), err.Error(), "Use '.key' or '.nested.key' syntax")
			os.Exit(ExitUsageError)
		}

		result, err := q.Execute(data)
		if err != nil {
			formatError("RUNTIME_ERROR", fmt.Sprintf("Query execution failed for '%s'", operationArg), err.Error(), "Check query path exists in TOML data")
			os.Exit(ExitParseError)
		}

		outputData(result, outputFormat)

	case "set":
		if dryRun {
			// Dry-run mode: show what would be changed
			fmt.Printf("DRY RUN: Would set %s in %s\n", operationArg, filePath)
			// Simulate the operation to show the result
			m := modifier.New()
			err := m.SetValue(dataMap, operationArg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Set operation would fail: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Result:")
			outputData(dataMap, outputFormat)
		} else if inplace && !useStdin {
			// Modify file in-place
			m := modifier.New()
			err := m.SetValue(dataMap, operationArg)
			if err != nil {
				formatError("OPERATION_ERROR", "Set operation failed", err.Error(), "Check operation syntax and data types")
				os.Exit(ExitParseError)
			}

			// Write back to file
			err = writeTOMLFile(filePath, dataMap)
			if err != nil {
				formatError("FILE_ERROR", fmt.Sprintf("Failed to write file '%s'", filePath), err.Error(), "Check file permissions and disk space")
				os.Exit(ExitFileError)
			}
		} else {
			// Just modify and output
			m := modifier.New()
			err := m.SetValue(dataMap, operationArg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Set operation failed\n")
				fmt.Fprintf(os.Stderr, "Details: %v\n", err)
				os.Exit(1)
			}
			outputData(dataMap, outputFormat)
		}

	case "delete":
		if dryRun {
			// Dry-run mode: show what would be deleted
			fmt.Printf("DRY RUN: Would delete %s from %s\n", operationArg, filePath)
			// Simulate the operation to show the result
			m := modifier.New()
			err := m.DeleteValue(dataMap, operationArg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Delete operation would fail: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Result:")
			outputData(dataMap, outputFormat)
		} else if inplace && !useStdin {
			// Modify file in-place
			m := modifier.New()
			err := m.DeleteValue(dataMap, operationArg)
			if err != nil {
				formatError("OPERATION_ERROR", "Delete operation failed", err.Error(), "Check operation syntax and path exists")
				os.Exit(ExitParseError)
			}

			// Write back to file
			err = writeTOMLFile(filePath, dataMap)
			if err != nil {
				formatError("FILE_ERROR", fmt.Sprintf("Failed to write file '%s'", filePath), err.Error(), "Check file permissions and disk space")
				os.Exit(ExitFileError)
			}
		} else {
			// Just modify and output
			m := modifier.New()
			err := m.DeleteValue(dataMap, operationArg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Delete operation failed\n")
				fmt.Fprintf(os.Stderr, "Details: %v\n", err)
				os.Exit(1)
			}
			outputData(dataMap, outputFormat)
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown operation\n")
		os.Exit(2)
	}
}

// handleOperationsBulk handles operations for bulk files (limited operations)
func handleOperationsBulk(data interface{}, dataMap map[string]interface{}, filePath string) error {
	switch operation {
	case "query":
		// For bulk queries, just print the result with filename prefix
		if operationArg == "" {
			fmt.Printf("%s: ", filePath)
			outputData(data, outputFormat)
			return nil
		}

		// Execute query
		q, err := query.New(operationArg)
		if err != nil {
			return fmt.Errorf("invalid query syntax '%s': %v", operationArg, err)
		}

		result, err := q.Execute(data)
		if err != nil {
			return fmt.Errorf("query execution failed for '%s': %v", operationArg, err)
		}

		fmt.Printf("%s: ", filePath)
		outputData(result, outputFormat)
		return nil

	case "set":
		if dryRun {
			// Dry-run mode for bulk operations
			fmt.Printf("%s: DRY RUN: Would set %s\n", filePath, operationArg)
			// Simulate the operation to show the result
			m := modifier.New()
			err := m.SetValue(dataMap, operationArg)
			if err != nil {
				return fmt.Errorf("set operation would fail: %v", err)
			}
			fmt.Printf("%s: Result: ", filePath)
			outputData(dataMap, outputFormat)
			return nil
		} else if inplace {
			// Modify file in-place for bulk operations
			m := modifier.New()
			err := m.SetValue(dataMap, operationArg)
			if err != nil {
				return fmt.Errorf("set operation failed: %v", err)
			}

			// Write back to file
			err = writeTOMLFile(filePath, dataMap)
			if err != nil {
				return fmt.Errorf("failed to write file '%s': %v", filePath, err)
			}
			fmt.Printf("%s: updated\n", filePath)
			return nil
		} else {
			return fmt.Errorf("bulk set operations require -i (in-place) flag")
		}

	case "delete":
		if dryRun {
			// Dry-run mode for bulk operations
			fmt.Printf("%s: DRY RUN: Would delete %s\n", filePath, operationArg)
			// Simulate the operation to show the result
			m := modifier.New()
			err := m.DeleteValue(dataMap, operationArg)
			if err != nil {
				return fmt.Errorf("delete operation would fail: %v", err)
			}
			fmt.Printf("%s: Result: ", filePath)
			outputData(dataMap, outputFormat)
			return nil
		} else if inplace {
			// Modify file in-place for bulk operations
			m := modifier.New()
			err := m.DeleteValue(dataMap, operationArg)
			if err != nil {
				return fmt.Errorf("delete operation failed: %v", err)
			}

			// Write back to file
			err = writeTOMLFile(filePath, dataMap)
			if err != nil {
				return fmt.Errorf("failed to write file '%s': %v", filePath, err)
			}
			fmt.Printf("%s: updated\n", filePath)
			return nil
		} else {
			return fmt.Errorf("bulk delete operations require -i (in-place) flag")
		}

	default:
		return fmt.Errorf("operation '%s' not supported for bulk operations", operation)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "tmq - TOML Query Tool (like jq for TOML)\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [file] [operation]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       %s < file.toml | %s [options] [operation]\n", os.Args[0], os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -o, --output FORMAT    Output format: toml, json, yaml (default: toml)\n")
	fmt.Fprintf(os.Stderr, "  -i, --inplace          Modify file in-place (requires file argument)\n")
	fmt.Fprintf(os.Stderr, "      --dry-run          Preview changes without applying them\n")
	fmt.Fprintf(os.Stderr, "      --validate         Validate TOML syntax and structure\n")
	fmt.Fprintf(os.Stderr, "      --compare FILE     Compare with another TOML file\n")
	fmt.Fprintf(os.Stderr, "      --schema FILE      Validate against schema file (future)\n")
	fmt.Fprintf(os.Stderr, "  -h, --help             Show this help message\n")
	fmt.Fprintf(os.Stderr, "      --version          Show version information\n")
	fmt.Fprintf(os.Stderr, "\nArguments:\n")
	fmt.Fprintf(os.Stderr, "  file                   TOML file path (optional, reads from stdin if omitted)\n")
	fmt.Fprintf(os.Stderr, "  operation              Query: '.key' | Set: '.key = \"value\"' | Delete: 'del(.key)'\n")
	fmt.Fprintf(os.Stderr, "\nExit codes:\n")
	fmt.Fprintf(os.Stderr, "  0    Success\n")
	fmt.Fprintf(os.Stderr, "  1    Parse error or runtime error\n")
	fmt.Fprintf(os.Stderr, "  2    Usage error or invalid arguments\n")
	fmt.Fprintf(os.Stderr, "  3    Security error (unsafe file path)\n")
	fmt.Fprintf(os.Stderr, "  4    File operation error\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s config.toml '.project.version'\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  cat config.toml | %s '.database.host'\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s config.toml -o json\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s config.toml '.version = \"2.0\"' -i\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s config.toml '.version = \"2.0\"' --dry-run -i\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --validate config.toml\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --compare config1.toml config2.toml\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s config.toml 'del(.old_field)' -i\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --version\n", os.Args[0])
}

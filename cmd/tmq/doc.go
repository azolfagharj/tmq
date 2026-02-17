// Tmq is a complete, standalone command-line tool for TOML.
// Like jq for JSON, yq for YAML — but for TOML.
//
// It reads TOML from a file or stdin, runs optional dot-separated path queries,
// and prints results to stdout. Exit codes: 0 success, 1 parse/runtime error, 2 usage error.
//
// # Usage
//
// Basic usage:
//
//	tmq config.toml '.project.version'
//
// Read from stdin:
//
//	cat config.toml | tmq '.database.host'
//
// Print all data:
//
//	tmq config.toml
//
// # Exit Codes
//
//   - 0: Success
//   - 1: Parse error or runtime error
//   - 2: Usage error or invalid query syntax
//
// # Query Syntax
//
// Queries use dot-separated paths:
//
//   - '.key' - access top-level key
//   - '.nested.key' - access nested key
//   - '.array' - access array values
//
// # Examples
//
// Get project version:
//
//	tmq pyproject.toml '.project.version'
//
// Get database host:
//
//	cat config.toml | tmq '.database.host'
//
// List all configuration:
//
//	tmq config.toml
package main

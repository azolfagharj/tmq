// Package query provides TOML query functionality for tmq.
//
// This package handles parsing and executing dot-separated path queries
// to extract values from TOML data structures.
//
// # Basic Usage
//
// Create a query and execute it on TOML data:
//
//	q, err := query.New(".project.version")
//	if err != nil {
//		log.Fatal(err)
//	}
//	result, err := q.Execute(tomlData)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result)
//
// # Query Syntax
//
// Queries use dot-separated paths to navigate TOML structures:
//
//   - ".key" - access top-level key
//   - ".nested.key" - access nested key
//   - ".array" - access array values
//
// # Supported Data Types
//
// Queries work with all TOML data types:
//   - Strings: "hello world"
//   - Numbers: 42, 3.14
//   - Booleans: true, false
//   - Arrays: [1, 2, 3]
//   - Tables: {key = "value"}
//
// # Error Handling
//
// Query execution may fail if:
//   - The path doesn't exist in the data
//   - The data structure is incompatible with the query
//   - The path syntax is invalid
//
// Use query.New() to validate query syntax before execution.
package query
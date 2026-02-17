// Package parser provides TOML parsing functionality for tmq.
//
// This package handles parsing TOML files and provides a clean interface
// for accessing TOML data structures according to the TOML 1.0.0 specification.
//
// # Basic Usage
//
// Create a new parser and parse a TOML file:
//
//	p := parser.New()
//	err := p.ParseFile("config.toml")
//	if err != nil {
//		log.Fatal(err)
//	}
//	data := p.GetData()
//
// # TOML Specification Compliance
//
// This package uses BurntSushi/toml and supports TOML 1.0.0 features including:
//   - Basic data types (string, int, float, bool)
//   - Arrays and nested tables
//   - Inline tables
//   - Comments and whitespace handling
//
// # Parsing from Different Sources
//
// Parse from a file:
//
//	err := p.ParseFile("path/to/file.toml")
//
// Parse from an io.Reader (e.g., stdin, network):
//
//	err := p.ParseReader(reader)
//
// # Error Handling
//
// All methods return errors for invalid TOML syntax or I/O issues.
// The error messages provide detailed context about parsing failures.
package parser

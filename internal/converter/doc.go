// Package converter provides format conversion functionality for tmq.
//
// This package handles converting between TOML, JSON, and YAML formats
// for both input and output operations.
//
// # Supported Formats
//
//   - TOML: Tom's Obvious, Minimal Language
//   - JSON: JavaScript Object Notation
//   - YAML: YAML Ain't Markup Language
//
// # Usage
//
// Convert TOML data to JSON:
//
//	jsonStr, err := converter.ConvertToJSON(tomlData)
//
// Convert TOML data to YAML:
//
//	yamlStr, err := converter.ConvertToYAML(tomlData)
//
// Convert with specified format:
//
//	output, err := converter.ConvertData(data, converter.FormatJSON)
//
// # Output Formats
//
// Use the -o flag with tmq to specify output format:
//
//	tmq config.toml -o json
//	tmq config.toml -o yaml
package converter

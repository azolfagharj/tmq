// Package converter provides format conversion functionality for tmq.
//
// This package handles converting TOML data to JSON and YAML formats,
// and vice versa for input formats.
package converter

import (
	"encoding/json"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// OutputFormat represents supported output formats
type OutputFormat int

const (
	FormatTOML OutputFormat = iota
	FormatJSON
	FormatYAML
)

// String returns the string representation of the format
func (f OutputFormat) String() string {
	switch f {
	case FormatTOML:
		return "toml"
	case FormatJSON:
		return "json"
	case FormatYAML:
		return "yaml"
	default:
		return "unknown"
	}
}

// ParseOutputFormat parses a format string into OutputFormat
func ParseOutputFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(s) {
	case "toml":
		return FormatTOML, nil
	case "json":
		return FormatJSON, nil
	case "yaml", "yml":
		return FormatYAML, nil
	default:
		return FormatTOML, fmt.Errorf("unsupported output format: %s (supported: toml, json, yaml)", s)
	}
}

// ConvertToJSON converts TOML data to JSON string
func ConvertToJSON(data interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// ConvertToYAML converts TOML data to YAML string
func ConvertToYAML(data interface{}) (string, error) {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to convert to YAML: %w", err)
	}
	return string(yamlBytes), nil
}

// ConvertData converts TOML data to the specified output format
func ConvertData(data interface{}, format OutputFormat) (string, error) {
	switch format {
	case FormatTOML:
		// For TOML output, we could implement TOML encoding here
		// For now, return an error as it's not implemented yet
		return "", fmt.Errorf("TOML output format not yet implemented")
	case FormatJSON:
		return ConvertToJSON(data)
	case FormatYAML:
		return ConvertToYAML(data)
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}

// Package modifier provides TOML modification functionality for tmq.
//
// This package handles setting and deleting values in TOML data structures,
// with support for in-place file modifications.
package modifier

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/username/toml_query/internal/query"
)

// Modifier handles TOML modification operations
type Modifier struct{}

// New creates a new TOML modifier
func New() *Modifier {
	return &Modifier{}
}

// SetValue sets a value at the specified path in the TOML data
// Supports syntax like: .key = "value", .nested.key = 42
func (m *Modifier) SetValue(data map[string]interface{}, setExpr string) error {
	// Parse set expression: ".key = value"
	parts := strings.SplitN(setExpr, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid set expression: %s (expected: .key = value)", setExpr)
	}

	path := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])

	// Parse the path
	q, err := query.New(path)
	if err != nil {
		return fmt.Errorf("invalid path in set expression: %v", err)
	}

	// Parse the value
	value, err := parseValue(valueStr)
	if err != nil {
		return fmt.Errorf("invalid value in set expression: %v", err)
	}

	// Set the value
	return m.setValueAtPath(data, q.Parts(), value)
}

// DeleteValue deletes a value at the specified path
// Supports syntax like: del(.key), del(.nested.key)
func (m *Modifier) DeleteValue(data map[string]interface{}, deleteExpr string) error {
	// Parse delete expression: "del(.key)"
	if !strings.HasPrefix(deleteExpr, "del(") || !strings.HasSuffix(deleteExpr, ")") {
		return fmt.Errorf("invalid delete expression: %s (expected: del(.key))", deleteExpr)
	}

	pathStr := deleteExpr[4 : len(deleteExpr)-1] // Extract ".key" from "del(.key)"

	q, err := query.New(pathStr)
	if err != nil {
		return fmt.Errorf("invalid path in delete expression: %v", err)
	}

	return m.deleteValueAtPath(data, q.Parts())
}

// parseValue parses a string value into the appropriate Go type
func parseValue(s string) (interface{}, error) {
	// Remove quotes if present
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1], nil
	}

	// Try to parse as boolean
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}

	// Try to parse as number
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		// Check if it's an integer
		if strings.Contains(s, ".") {
			return num, nil
		}
		return int64(num), nil
	}

	// Default to string
	return s, nil
}

// setValueAtPath sets a value at the specified path in the data structure
func (m *Modifier) setValueAtPath(data map[string]interface{}, path []string, value interface{}) error {
	if len(path) == 0 {
		return fmt.Errorf("cannot set root value")
	}

	current := data
	for i, key := range path[:len(path)-1] {
		if next, exists := current[key]; exists {
			if nextMap, ok := next.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return fmt.Errorf("cannot navigate into %T at %s", next, strings.Join(path[:i+1], "."))
			}
		} else {
			// Create nested map
			newMap := make(map[string]interface{})
			current[key] = newMap
			current = newMap
		}
	}

	// Set the final value
	finalKey := path[len(path)-1]
	current[finalKey] = value
	return nil
}

// deleteValueAtPath deletes a value at the specified path
func (m *Modifier) deleteValueAtPath(data map[string]interface{}, path []string) error {
	if len(path) == 0 {
		return fmt.Errorf("cannot delete root value")
	}

	current := data
	for i, key := range path[:len(path)-1] {
		if next, exists := current[key]; exists {
			if nextMap, ok := next.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return fmt.Errorf("cannot navigate into %T at %s", next, strings.Join(path[:i+1], "."))
			}
		} else {
			return fmt.Errorf("path not found: %s", strings.Join(path[:i+1], "."))
		}
	}

	// Delete the final key
	finalKey := path[len(path)-1]
	if _, exists := current[finalKey]; !exists {
		return fmt.Errorf("key not found: %s", strings.Join(path, "."))
	}

	delete(current, finalKey)
	return nil
}
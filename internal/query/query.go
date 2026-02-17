package query

import (
	"fmt"
	"strings"
)

// Query represents a parsed query path
type Query struct {
	parts []string
}

// New creates a new query from a dot-separated path string
func New(path string) (*Query, error) {
	if path == "" {
		return nil, fmt.Errorf("query path cannot be empty")
	}

	// Remove leading dot if present
	path = strings.TrimPrefix(path, ".")

	if path == "" {
		// Root query (just ".")
		return &Query{parts: []string{}}, nil
	}

	parts := strings.Split(path, ".")
	return &Query{parts: parts}, nil
}

// Execute runs the query against the provided TOML data
func (q *Query) Execute(data interface{}) (interface{}, error) {
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}

	current := data

	// Navigate through the path parts
	for _, part := range q.parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, fmt.Errorf("key '%s' not found", part)
			}
		case map[interface{}]interface{}:
			// Handle cases where TOML parser returns map[interface{}]interface{}
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, fmt.Errorf("key '%s' not found", part)
			}
		case []interface{}:
			return nil, fmt.Errorf("cannot navigate into array at '%s'", part)
		default:
			return nil, fmt.Errorf("cannot navigate into %T at '%s'", current, part)
		}
	}

	return current, nil
}

// String returns the string representation of the query
func (q *Query) String() string {
	if len(q.parts) == 0 {
		return "."
	}
	return "." + strings.Join(q.parts, ".")
}

// Parts returns the individual parts of the query path
func (q *Query) Parts() []string {
	return q.parts
}

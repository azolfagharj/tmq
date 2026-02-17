// Package modifier provides TOML modification functionality for tmq.
//
// This package handles setting and deleting values in TOML data structures,
// with support for in-place file modifications and basic comment preservation.
//
// # Set Operations
//
// Set values using dot notation:
//
//	mod := modifier.New()
//	err := mod.SetValue(data, `.project.version = "1.0.0"`)
//	err := mod.SetValue(data, `.database.port = 5432`)
//
// # Delete Operations
//
// Delete values using del() syntax:
//
//	err := mod.DeleteValue(data, `del(.optional_field)`)
//	err := mod.DeleteValue(data, `del(.debug.enabled)`)
//
// # Data Types
//
// Supported value types in set operations:
//   - Strings: "hello", 'world'
//   - Numbers: 42, 3.14
//   - Booleans: true, false
//
// # Path Navigation
//
// Supports nested path navigation and automatic creation of intermediate maps:
//
//	// Creates nested structure automatically
//	mod.SetValue(data, `.config.database.host = "localhost"`)
//
// # Error Handling
//
// Operations return detailed errors for:
//   - Invalid syntax in expressions
//   - Path not found (for delete operations)
//   - Type conflicts in navigation
//   - Invalid value formats
package modifier
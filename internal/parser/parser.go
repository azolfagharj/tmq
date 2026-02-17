// Package parser provides TOML parsing functionality for tmq.
//
// This package handles parsing TOML files and provides a clean interface
// for accessing TOML data structures according to the TOML 1.0.0 specification.
package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// Parser handles TOML parsing operations
type Parser struct {
	data interface{}
}

// New creates a new TOML parser
func New() *Parser {
	return &Parser{}
}

// ParseFile parses a TOML file from the given path
func (p *Parser) ParseFile(path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	return p.ParseReader(file)
}

// ParseReader parses TOML data from an io.Reader
func (p *Parser) ParseReader(r io.Reader) error {
	if r == nil {
		return fmt.Errorf("reader cannot be nil")
	}

	_, err := toml.NewDecoder(r).Decode(&p.data)
	if err != nil {
		return fmt.Errorf("failed to parse TOML: %w", err)
	}

	return nil
}

// GetData returns the parsed TOML data
func (p *Parser) GetData() interface{} {
	return p.data
}

// GetValue retrieves a value from the parsed TOML using a dot-separated path.
// For path traversal use the [query] package instead; this method is reserved for future use.
func (p *Parser) GetValue(path string) (interface{}, error) {
	if path == "" {
		return p.data, nil
	}
	return nil, fmt.Errorf("path traversal not yet implemented: %s", path)
}

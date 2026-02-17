package parser

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test constants for TOML content
const (
	validTOMLContent   = `key = "value"`
	invalidTOMLContent = `invalid toml syntax [[[`
	complexTOMLContent = `[section]
key = "value"
number = 42`
	commentedTOMLContent = `# This is a comment
key = "value"
# Another comment`
	booleanTOMLContent = `enabled = true
disabled = false`
	emptyTOMLContent = ""
)

// Helper functions for test setup
func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.toml")

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	return tmpFile
}

func assertError(t *testing.T, err error, wantErr bool, errMsg string) {
	t.Helper()

	if wantErr {
		if err == nil {
			t.Errorf("expected error, got nil")
			return
		}
		if errMsg != "" && !strings.Contains(err.Error(), errMsg) {
			t.Errorf("expected error containing %q, got %q", errMsg, err.Error())
		}
	} else {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestNew(t *testing.T) {
	t.Run("creates new parser", func(t *testing.T) {
		p := New()
		if p == nil {
			t.Fatal("New() returned nil")
		}
		if p.data != nil {
			t.Errorf("expected data to be nil, got %v", p.data)
		}
	})
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (string, error) // Returns file path or error for special cases
		wantErr bool
		errMsg  string
	}{
		{
			name: "empty path",
			setup: func(t *testing.T) (string, error) {
				return "", nil // Special case handled in test
			},
			wantErr: true,
			errMsg:  "file path cannot be empty",
		},
		{
			name: "valid TOML",
			setup: func(t *testing.T) (string, error) {
				return createTempFile(t, validTOMLContent), nil
			},
			wantErr: false,
		},
		{
			name: "invalid TOML",
			setup: func(t *testing.T) (string, error) {
				return createTempFile(t, invalidTOMLContent), nil
			},
			wantErr: true,
			errMsg:  "failed to parse TOML",
		},
		{
			name: "complex TOML",
			setup: func(t *testing.T) (string, error) {
				return createTempFile(t, complexTOMLContent), nil
			},
			wantErr: false, // This is valid TOML now
		},
		{
			name: "TOML with comments",
			setup: func(t *testing.T) (string, error) {
				return createTempFile(t, commentedTOMLContent), nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()

			var filePath string
			var err error

			if tt.name == "empty path" {
				filePath = ""
			} else {
				filePath, err = tt.setup(t)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err = p.ParseFile(filePath)
			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}

	t.Run("nonexistent file", func(t *testing.T) {
		p := New()
		err := p.ParseFile("/nonexistent/file.toml")
		if err == nil {
			t.Error("expected error for nonexistent file, got nil")
		}
		if !strings.Contains(err.Error(), "failed to open file") {
			t.Errorf("expected file open error, got: %v", err)
		}
	})
}

func TestParseReader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil reader",
			input:   "",
			wantErr: true,
			errMsg:  "reader cannot be nil",
		},
		{
			name:    "valid TOML",
			input:   validTOMLContent,
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   emptyTOMLContent,
			wantErr: false,
		},
		{
			name:    "invalid TOML",
			input:   invalidTOMLContent,
			wantErr: true,
			errMsg:  "failed to parse TOML",
		},
		{
			name:    "complex TOML",
			input:   complexTOMLContent,
			wantErr: false,
		},
		{
			name:    "TOML with comments",
			input:   commentedTOMLContent,
			wantErr: false,
		},
		{
			name:    "TOML with boolean values",
			input:   booleanTOMLContent,
			wantErr: false,
		},
		{
			name:    "invalid backslash in TOML",
			input:   `\n\nkey = "value"\n\n\nother = 123\n\n`, // Backslash is invalid
			wantErr: true,
			errMsg:  "failed to parse TOML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()

			var reader io.Reader
			switch {
			case tt.name == "nil reader":
				reader = nil
			case tt.input != "":
				reader = bytes.NewReader([]byte(tt.input))
			default:
				reader = bytes.NewReader([]byte{})
			}

			err := p.ParseReader(reader)
			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestGetData(t *testing.T) {
	t.Run("no data initially", func(t *testing.T) {
		p := New()
		data := p.GetData()
		if data != nil {
			t.Errorf("expected nil data initially, got %v", data)
		}
	})

	t.Run("data after parsing", func(t *testing.T) {
		p := New()
		err := p.ParseReader(bytes.NewReader([]byte(validTOMLContent)))
		if err != nil {
			t.Fatalf("failed to parse TOML: %v", err)
		}

		data := p.GetData()
		if data == nil {
			t.Error("expected data after parsing, got nil")
		}
	})
}

func TestGetValue(t *testing.T) {
	t.Run("path traversal on empty data", func(t *testing.T) {
		p := New()
		_, err := p.GetValue("some.path")
		if err == nil {
			t.Error("expected error for path traversal on empty data")
		}
	})

	t.Run("empty path returns data", func(t *testing.T) {
		p := New()
		err := p.ParseReader(bytes.NewReader([]byte(validTOMLContent)))
		if err != nil {
			t.Fatalf("failed to parse TOML: %v", err)
		}

		val, err := p.GetValue("")
		if err != nil {
			t.Errorf("expected no error for empty path, got: %v", err)
		}
		if val == nil {
			t.Error("expected data for empty path, got nil")
		}
	})

	t.Run("unimplemented path traversal", func(t *testing.T) {
		p := New()
		err := p.ParseReader(bytes.NewReader([]byte(validTOMLContent)))
		if err != nil {
			t.Fatalf("failed to parse TOML: %v", err)
		}

		_, err = p.GetValue("key")
		if err == nil {
			t.Error("expected error for unimplemented path traversal")
		}
		if !strings.Contains(err.Error(), "not yet implemented") {
			t.Errorf("expected 'not yet implemented' error, got: %v", err)
		}
	})
}

// Benchmark tests for performance validation
func BenchmarkParseFile(b *testing.B) {
	// Create temp file with test content
	content := `title = "Test Document"
[database]
host = "localhost"
port = 5432
enabled = true

[[servers]]
name = "server1"
ip = "192.168.1.1"

[[servers]]
name = "server2"
ip = "192.168.1.2"`

	tmpDir := b.TempDir()
	tmpFile := filepath.Join(tmpDir, "benchmark.toml")

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("failed to create benchmark file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := New()
		err := p.ParseFile(tmpFile)
		if err != nil {
			b.Fatalf("parse error: %v", err)
		}
	}
}

func BenchmarkParseReader(b *testing.B) {
	content := `title = "Test Document"
[database]
host = "localhost"
port = 5432`

	reader := bytes.NewReader([]byte(content))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := reader.Seek(0, 0); err != nil {
			b.Fatalf("seek error: %v", err)
		}
		p := New()
		err := p.ParseReader(reader)
		if err != nil {
			b.Fatalf("parse error: %v", err)
		}
	}
}

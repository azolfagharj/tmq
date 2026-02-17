package query

import (
	"strings"
	"testing"
)

// Test constants for query paths
const (
	rootPath        = "."
	singleKeyPath   = ".key"
	nestedPath      = ".project.version"
	deepNestedPath  = ".a.b.c.d.e"
	noDotPath       = "key.subkey"
)

// Helper functions for query testing
func assertQueryCreation(t *testing.T, path string, wantErr bool, errMsg string) *Query {
	t.Helper()

	q, err := New(path)

	if wantErr {
		if err == nil {
			t.Errorf("expected error for path %q, got nil", path)
			return nil
		}
		if errMsg != "" && !strings.Contains(err.Error(), errMsg) {
			t.Errorf("expected error containing %q, got %q", errMsg, err.Error())
		}
		return nil
	}

	if err != nil {
		t.Errorf("unexpected error for path %q: %v", path, err)
		return nil
	}

	if q == nil {
		t.Errorf("expected non-nil query for path %q", path)
		return nil
	}

	return q
}

func assertQueryString(t *testing.T, q *Query, expected string) {
	t.Helper()

	if result := q.String(); result != expected {
		t.Errorf("expected string %q, got %q", expected, result)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
		wantStr string
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errMsg:  "query path cannot be empty",
		},
		{
			name:    "root path",
			path:    rootPath,
			wantErr: false,
			wantStr: rootPath,
		},
		{
			name:    "single key",
			path:    singleKeyPath,
			wantErr: false,
			wantStr: singleKeyPath,
		},
		{
			name:    "nested path",
			path:    nestedPath,
			wantErr: false,
			wantStr: nestedPath,
		},
		{
			name:    "deep nested path",
			path:    deepNestedPath,
			wantErr: false,
			wantStr: deepNestedPath,
		},
		{
			name:    "path without leading dot",
			path:    noDotPath,
			wantErr: false,
			wantStr: ".key.subkey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := assertQueryCreation(t, tt.path, tt.wantErr, tt.errMsg)
			if q != nil && tt.wantStr != "" {
				assertQueryString(t, q, tt.wantStr)
			}
		})
	}
}

func TestNew_EdgeCases(t *testing.T) {
	t.Run("dot only", func(t *testing.T) {
		q := assertQueryCreation(t, ".", false, "")
		if q != nil {
			assertQueryString(t, q, ".")
			if len(q.Parts()) != 0 {
				t.Errorf("expected empty parts for root query, got %v", q.Parts())
			}
		}
	})

	t.Run("multiple dots", func(t *testing.T) {
		q := assertQueryCreation(t, "..key", false, "")
		if q != nil {
			assertQueryString(t, q, "..key")
			expectedParts := []string{"", "key"}
			parts := q.Parts()
			if len(parts) != len(expectedParts) {
				t.Errorf("expected %d parts, got %d", len(expectedParts), len(parts))
			}
			for i, part := range parts {
				if i < len(expectedParts) && part != expectedParts[i] {
					t.Errorf("expected part[%d] = %q, got %q", i, expectedParts[i], part)
				}
			}
		}
	})

	t.Run("special characters", func(t *testing.T) {
		q := assertQueryCreation(t, ".key-with-dashes.and_underscores", false, "")
		if q != nil {
			assertQueryString(t, q, ".key-with-dashes.and_underscores")
		}
	})
}

func TestParts(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedParts []string
	}{
		{
			name:          "root",
			path:          rootPath,
			expectedParts: []string{},
		},
		{
			name:          "single key",
			path:          singleKeyPath,
			expectedParts: []string{"key"},
		},
		{
			name:          "nested path",
			path:          nestedPath,
			expectedParts: []string{"project", "version"},
		},
		{
			name:          "deep nested path",
			path:          deepNestedPath,
			expectedParts: []string{"a", "b", "c", "d", "e"},
		},
		{
			name:          "path without leading dot",
			path:          noDotPath,
			expectedParts: []string{"key", "subkey"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := assertQueryCreation(t, tt.path, false, "")
			if q != nil {
				parts := q.Parts()
				if len(parts) != len(tt.expectedParts) {
					t.Errorf("expected %d parts, got %d", len(tt.expectedParts), len(parts))
					return
				}

				for i, part := range parts {
					if part != tt.expectedParts[i] {
						t.Errorf("expected part[%d] = %q, got %q", i, tt.expectedParts[i], part)
					}
				}
			}
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkNew(b *testing.B) {
	testPaths := []string{
		rootPath,
		singleKeyPath,
		nestedPath,
		deepNestedPath,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			_, err := New(path)
			if err != nil {
				b.Fatalf("unexpected error for path %q: %v", path, err)
			}
		}
	}
}

func BenchmarkQueryString(b *testing.B) {
	q, err := New(deepNestedPath)
	if err != nil {
		b.Fatalf("failed to create query: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.String()
	}
}
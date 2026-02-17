package converter

import (
	"encoding/json"
	"strings"
	"testing"
)

// Test constants
const (
	testKey   = "key"
	testValue = "value"
	testTitle = "Test Document"
	testHost  = "localhost"
	testPort  = 5432
)

// Factory functions for creating test data
func createSimpleTestData() map[string]interface{} {
	return map[string]interface{}{
		testKey: testValue,
	}
}

func createComplexTestData() map[string]interface{} {
	return map[string]interface{}{
		"title": testTitle,
		"config": map[string]interface{}{
			"database": map[string]interface{}{
				"host": testHost,
				"port": testPort,
				"enabled": true,
			},
		},
		"servers": []interface{}{
			map[string]interface{}{
				"name": "server1",
				"ip":   "192.168.1.1",
			},
			map[string]interface{}{
				"name": "server2",
				"ip":   "192.168.1.2",
			},
		},
	}
}

// Helper functions
func assertFormatResult(t *testing.T, result string, shouldContain []string, shouldNotContain []string) {
	t.Helper()

	for _, expected := range shouldContain {
		if !strings.Contains(result, expected) {
			t.Errorf("expected result to contain %q, but it doesn't.\nResult: %s", expected, result)
		}
	}

	for _, unexpected := range shouldNotContain {
		if strings.Contains(result, unexpected) {
			t.Errorf("expected result to NOT contain %q, but it does.\nResult: %s", unexpected, result)
		}
	}
}

func assertValidJSON(t *testing.T, jsonStr string) {
	t.Helper()

	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Errorf("generated JSON is invalid: %v\nJSON: %s", err, jsonStr)
	}
}

func TestParseOutputFormat(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    OutputFormat
		wantErr     bool
		description string
	}{
		{
			name:        "toml format",
			input:       "toml",
			expected:    FormatTOML,
			wantErr:     false,
			description: "TOML format should be parsed correctly",
		},
		{
			name:        "json format",
			input:       "json",
			expected:    FormatJSON,
			wantErr:     false,
			description: "JSON format should be parsed correctly",
		},
		{
			name:        "yaml format",
			input:       "yaml",
			expected:    FormatYAML,
			wantErr:     false,
			description: "YAML format should be parsed correctly",
		},
		{
			name:        "yml format",
			input:       "yml",
			expected:    FormatYAML,
			wantErr:     false,
			description: "YML format should be parsed correctly",
		},
		{
			name:        "uppercase TOML",
			input:       "TOML",
			expected:    FormatTOML,
			wantErr:     false,
			description: "Format parsing should be case insensitive",
		},
		{
			name:        "invalid format",
			input:       "invalid",
			expected:    FormatTOML,
			wantErr:     true,
			description: "Invalid format should return error",
		},
		{
			name:        "empty string",
			input:       "",
			expected:    FormatTOML,
			wantErr:     true,
			description: "Empty string should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseOutputFormat(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("%s: expected error, got nil", tt.description)
				} else if !strings.Contains(err.Error(), "unsupported output format") {
					t.Errorf("%s: expected 'unsupported output format' error, got %q", tt.description, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
				return
			}

			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.description, tt.expected, result)
			}
		})
	}
}

func TestOutputFormatString(t *testing.T) {
	tests := []struct {
		format      OutputFormat
		expected    string
		description string
	}{
		{FormatTOML, "toml", "TOML format should stringify correctly"},
		{FormatJSON, "json", "JSON format should stringify correctly"},
		{FormatYAML, "yaml", "YAML format should stringify correctly"},
		{OutputFormat(999), "unknown", "Unknown format should return 'unknown'"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.format.String()
			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.description, tt.expected, result)
			}
		})
	}
}

func TestConvertToJSON(t *testing.T) {
	t.Run("simple data", func(t *testing.T) {
		data := createSimpleTestData()
		result, err := ConvertToJSON(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertFormatResult(t, result, []string{`"key": "value"`}, []string{})
		assertValidJSON(t, result)
	})

	t.Run("complex data", func(t *testing.T) {
		data := createComplexTestData()
		result, err := ConvertToJSON(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedParts := []string{
			`"title": "Test Document"`,
			`"config":`,
			`"database":`,
			`"host": "localhost"`,
			`"port": 5432`,
			`"enabled": true`,
			`"servers":`,
			`"name": "server1"`,
			`"ip": "192.168.1.1"`,
		}

		assertFormatResult(t, result, expectedParts, []string{})
		assertValidJSON(t, result)
	})

	t.Run("empty data", func(t *testing.T) {
		data := map[string]interface{}{}
		result, err := ConvertToJSON(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != "{}" {
			t.Errorf("expected empty JSON object, got %q", result)
		}
		assertValidJSON(t, result)
	})
}

func TestConvertToYAML(t *testing.T) {
	t.Run("simple data", func(t *testing.T) {
		data := createSimpleTestData()
		result, err := ConvertToYAML(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertFormatResult(t, result, []string{"key: value"}, []string{})
	})

	t.Run("complex data", func(t *testing.T) {
		data := createComplexTestData()
		result, err := ConvertToYAML(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedParts := []string{
			"title: Test Document",
			"config:",
			"database:",
			"host: localhost",
			"port: 5432",
			"enabled: true",
			"servers:",
			"name: server1",
			"ip: 192.168.1.1",
		}

		assertFormatResult(t, result, expectedParts, []string{})
	})

	t.Run("empty data", func(t *testing.T) {
		data := map[string]interface{}{}
		result, err := ConvertToYAML(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// YAML for empty map produces "{}" or similar representation
		if !strings.Contains(result, "{}") && result != "" && result != "\n" {
			t.Errorf("expected YAML representation of empty object, got %q", result)
		}
	})
}

func TestConvertData(t *testing.T) {
	data := createSimpleTestData()

	t.Run("TOML format", func(t *testing.T) {
		result, err := ConvertData(data, FormatTOML)
		if err == nil {
			t.Error("expected error for unimplemented TOML output, got nil")
		}
		if !strings.Contains(err.Error(), "TOML output format not yet implemented") {
			t.Errorf("expected implementation error, got %q", err.Error())
		}
		if result != "" {
			t.Errorf("expected empty result for error case, got %q", result)
		}
	})

	t.Run("JSON format", func(t *testing.T) {
		result, err := ConvertData(data, FormatJSON)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertFormatResult(t, result, []string{`"key": "value"`}, []string{})
		assertValidJSON(t, result)
	})

	t.Run("YAML format", func(t *testing.T) {
		result, err := ConvertData(data, FormatYAML)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertFormatResult(t, result, []string{"key: value"}, []string{})
	})

	t.Run("invalid format", func(t *testing.T) {
		result, err := ConvertData(data, OutputFormat(999))
		if err == nil {
			t.Error("expected error for invalid format, got nil")
		}
		if !strings.Contains(err.Error(), "unsupported output format") {
			t.Errorf("expected format error, got %q", err.Error())
		}
		if result != "" {
			t.Errorf("expected empty result for error case, got %q", result)
		}
	})
}

func TestConvertData_EdgeCases(t *testing.T) {
	t.Run("nil data", func(t *testing.T) {
		_, err := ConvertData(nil, FormatJSON)
		// Note: gopkg.in/yaml.v3 and encoding/json may handle nil differently
		// This test may need adjustment based on library behavior
		if err != nil {
			t.Logf("nil data produced error (expected): %v", err)
		}
		// For now, just ensure it doesn't panic
	})

	t.Run("complex nested data", func(t *testing.T) {
		data := map[string]interface{}{
			"deeply": map[string]interface{}{
				"nested": map[string]interface{}{
					"structure": map[string]interface{}{
						"with": map[string]interface{}{
							"many": map[string]interface{}{
								"levels": "of nesting",
							},
						},
					},
				},
			},
		}

		result, err := ConvertData(data, FormatJSON)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertFormatResult(t, result, []string{`"levels": "of nesting"`}, []string{})
		assertValidJSON(t, result)
	})
}

// Benchmark tests for performance validation
func BenchmarkConvertToJSON(b *testing.B) {
	data := createComplexTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ConvertToJSON(data)
		if err != nil {
			b.Fatalf("convert error: %v", err)
		}
	}
}

func BenchmarkConvertToYAML(b *testing.B) {
	data := createComplexTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ConvertToYAML(data)
		if err != nil {
			b.Fatalf("convert error: %v", err)
		}
	}
}

func BenchmarkParseOutputFormat(b *testing.B) {
	formats := []string{"toml", "json", "yaml", "yml"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, format := range formats {
			_, err := ParseOutputFormat(format)
			if err != nil {
				b.Fatalf("parse error for %s: %v", format, err)
			}
		}
	}
}
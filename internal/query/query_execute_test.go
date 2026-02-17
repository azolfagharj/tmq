package query

import (
	"reflect"
	"strings"
	"testing"
)

// Test constants
const (
	testString  = "hello"
	testNumber  = 42
	testBoolean = true
	testKey     = "key"
	testValue   = "value"
	testData    = "found"
)

// Factory functions for creating test data
func createSimpleTestData() map[string]interface{} {
	return map[string]interface{}{
		testKey: testString,
	}
}

func createComplexTestData() map[string]interface{} {
	return map[string]interface{}{
		"string":  testString,
		"number":  testNumber,
		"boolean": testBoolean,
		"nested": map[string]interface{}{
			testKey: testValue,
			"deep": map[string]interface{}{
				"data": testData,
			},
		},
		"array": []interface{}{"item1", "item2"},
	}
}

func createNestedTestData() map[string]interface{} {
	return map[string]interface{}{
		"config": map[string]interface{}{
			"database": map[string]interface{}{
				"host": "localhost",
				"port": 5432,
			},
		},
	}
}

// Helper functions
func assertQueryExecution(t *testing.T, q *Query, data interface{}, wantErr bool, errMsg string, expected interface{}) {
	t.Helper()

	result, err := q.Execute(data)

	if wantErr {
		if err == nil {
			t.Errorf("expected error, got nil")
			return
		}
		if errMsg != "" && !strings.Contains(err.Error(), errMsg) {
			t.Errorf("expected error containing %q, got %q", errMsg, err.Error())
		}
		return
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestExecute(t *testing.T) {
	data := createComplexTestData()

	tests := []struct {
		name     string
		path     string
		expected interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "root access",
			path:     ".",
			expected: data,
			wantErr:  false,
		},
		{
			name:     "string value",
			path:     ".string",
			expected: testString,
			wantErr:  false,
		},
		{
			name:     "number value",
			path:     ".number",
			expected: testNumber,
			wantErr:  false,
		},
		{
			name:     "boolean value",
			path:     ".boolean",
			expected: testBoolean,
			wantErr:  false,
		},
		{
			name:     "nested value",
			path:     ".nested.key",
			expected: testValue,
			wantErr:  false,
		},
		{
			name:     "deep nested value",
			path:     ".nested.deep.data",
			expected: testData,
			wantErr:  false,
		},
		{
			name:     "array access",
			path:     ".array",
			expected: []interface{}{"item1", "item2"},
			wantErr:  false,
		},
		{
			name:    "nonexistent key",
			path:    ".missing",
			wantErr: true,
			errMsg:  "key 'missing' not found",
		},
		{
			name:    "nonexistent nested key",
			path:    ".nested.missing",
			wantErr: true,
			errMsg:  "key 'missing' not found",
		},
		{
			name:    "cannot navigate into non-object",
			path:    ".string.key",
			wantErr: true,
			errMsg:  "cannot navigate into",
		},
		{
			name:    "navigate into array",
			path:    ".array.key",
			wantErr: true,
			errMsg:  "cannot navigate into array",
		},
		{
			name:    "navigate into number",
			path:    ".number.key",
			wantErr: true,
			errMsg:  "cannot navigate into",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.path)
			if err != nil {
				t.Fatalf("failed to create query: %v", err)
			}

			assertQueryExecution(t, q, data, tt.wantErr, tt.errMsg, tt.expected)
		})
	}
}

func TestExecute_EdgeCases(t *testing.T) {
	t.Run("nil data", func(t *testing.T) {
		q, err := New(".key")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		_, err = q.Execute(nil)
		if err == nil {
			t.Error("expected error for nil data, got nil")
		}
		if !strings.Contains(err.Error(), "data cannot be nil") {
			t.Errorf("expected 'data cannot be nil' error, got %q", err.Error())
		}
	})

	t.Run("empty map", func(t *testing.T) {
		q, err := New(".key")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		_, err = q.Execute(make(map[string]interface{}))
		if err == nil {
			t.Error("expected error for missing key, got nil")
		}
		if !strings.Contains(err.Error(), "key 'key' not found") {
			t.Errorf("expected key not found error, got %q", err.Error())
		}
	})

	t.Run("root query on simple value", func(t *testing.T) {
		q, err := New(".")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		simpleValue := "just a string"
		result, err := q.Execute(simpleValue)
		if err != nil {
			t.Errorf("unexpected error for root query: %v", err)
		}
		if result != simpleValue {
			t.Errorf("expected %q, got %q", simpleValue, result)
		}
	})
}

func TestExecute_WithMapInterface(t *testing.T) {
	// Test with map[interface{}]interface{} (sometimes returned by TOML parsers)
	data := map[interface{}]interface{}{
		testKey: testValue,
		"nested": map[interface{}]interface{}{
			"subkey": "subvalue",
		},
	}

	t.Run("direct key access", func(t *testing.T) {
		q, err := New(".key")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		result, err := q.Execute(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		expected := testValue
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("nested key access", func(t *testing.T) {
		q, err := New(".nested.subkey")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		result, err := q.Execute(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		expected := "subvalue"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("missing key error", func(t *testing.T) {
		q, err := New(".missing")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		_, err = q.Execute(data)
		if err == nil {
			t.Error("expected error for missing key, got nil")
		}
		if !strings.Contains(err.Error(), "key 'missing' not found") {
			t.Errorf("expected key not found error, got %q", err.Error())
		}
	})
}

func TestExecute_ArrayHandling(t *testing.T) {
	data := createComplexTestData()

	t.Run("access array", func(t *testing.T) {
		q, err := New(".array")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		result, err := q.Execute(data)
		if err != nil {
			t.Errorf("unexpected error accessing array: %v", err)
			return
		}

		expected := []interface{}{"item1", "item2"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("cannot navigate into array", func(t *testing.T) {
		q, err := New(".array.0")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		_, err = q.Execute(data)
		if err == nil {
			t.Error("expected error when navigating into array, got nil")
		}
		if !strings.Contains(err.Error(), "cannot navigate into") {
			t.Errorf("expected navigation error, got %q", err.Error())
		}
	})
}

func TestExecute_ComplexNesting(t *testing.T) {
	data := createNestedTestData()

	t.Run("deep config access", func(t *testing.T) {
		q, err := New(".config.database.host")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		result, err := q.Execute(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		expected := "localhost"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("deep config number access", func(t *testing.T) {
		q, err := New(".config.database.port")
		if err != nil {
			t.Fatalf("failed to create query: %v", err)
		}

		result, err := q.Execute(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		expected := 5432
		if result != expected {
			t.Errorf("expected %d, got %v", expected, result)
		}
	})
}

// Benchmark tests for performance validation
func BenchmarkExecute(b *testing.B) {
	data := createComplexTestData()
	q, err := New(".nested.deep.data")
	if err != nil {
		b.Fatalf("failed to create query: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := q.Execute(data)
		if err != nil {
			b.Fatalf("execute error: %v", err)
		}
	}
}

func BenchmarkExecuteRoot(b *testing.B) {
	data := createComplexTestData()
	q, err := New(".")
	if err != nil {
		b.Fatalf("failed to create query: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := q.Execute(data)
		if err != nil {
			b.Fatalf("execute error: %v", err)
		}
	}
}

func BenchmarkExecuteSimple(b *testing.B) {
	data := createSimpleTestData()
	q, err := New(".key")
	if err != nil {
		b.Fatalf("failed to create query: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := q.Execute(data)
		if err != nil {
			b.Fatalf("execute error: %v", err)
		}
	}
}

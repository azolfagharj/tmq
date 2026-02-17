package modifier

import (
	"reflect"
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
				"host":    testHost,
				"port":    testPort,
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

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
}

func TestSetValue(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]interface{}
		expr     string
		expected map[string]interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name:    "set simple string",
			initial: map[string]interface{}{},
			expr:    `.name = "John"`,
			expected: map[string]interface{}{
				"name": "John",
			},
			wantErr: false,
		},
		{
			name:     "overwrite existing key",
			initial:  createSimpleTestData(),
			expr:     `.key = "updated"`,
			expected: map[string]interface{}{testKey: "updated"},
			wantErr:  false,
		},
		{
			name:    "set number",
			initial: map[string]interface{}{},
			expr:    `.age = 30`,
			expected: map[string]interface{}{
				"age": int64(30),
			},
			wantErr: false,
		},
		{
			name:    "set boolean",
			initial: map[string]interface{}{},
			expr:    `.active = true`,
			expected: map[string]interface{}{
				"active": true,
			},
			wantErr: false,
		},
		{
			name: "set nested value",
			initial: map[string]interface{}{
				"config": map[string]interface{}{},
			},
			expr: `.config.host = "localhost"`,
			expected: map[string]interface{}{
				"config": map[string]interface{}{
					"host": "localhost",
				},
			},
			wantErr: false,
		},
		{
			name:    "create nested structure",
			initial: map[string]interface{}{},
			expr:    `.database.credentials.username = "admin"`,
			expected: map[string]interface{}{
				"database": map[string]interface{}{
					"credentials": map[string]interface{}{
						"username": "admin",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid expression - no equals",
			initial: map[string]interface{}{},
			expr:    `.name "John"`,
			wantErr: true,
			errMsg:  "invalid set expression",
		},
		{
			name:    "empty path",
			initial: map[string]interface{}{},
			expr:    ` = "value"`,
			wantErr: true,
			errMsg:  "query path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			data := make(map[string]interface{})

			// Copy initial data
			for k, v := range tt.initial {
				data[k] = v
			}

			err := m.SetValue(data, tt.expr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(data, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, data)
			}
		})
	}
}

func TestDeleteValue(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]interface{}
		expr     string
		expected map[string]interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name: "delete simple key",
			initial: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			expr: `del(.name)`,
			expected: map[string]interface{}{
				"age": 30,
			},
			wantErr: false,
		},
		{
			name: "delete nested key",
			initial: map[string]interface{}{
				"config": map[string]interface{}{
					"host": "localhost",
					"port": 5432,
				},
			},
			expr: `del(.config.host)`,
			expected: map[string]interface{}{
				"config": map[string]interface{}{
					"port": 5432,
				},
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent key",
			initial: map[string]interface{}{},
			expr:    `del(.missing)`,
			wantErr: true,
			errMsg:  "key not found",
		},
		{
			name:    "invalid delete syntax",
			initial: map[string]interface{}{},
			expr:    `delete(.key)`,
			wantErr: true,
			errMsg:  "invalid delete expression",
		},
		{
			name:    "delete with invalid path",
			initial: map[string]interface{}{},
			expr:    `del()`,
			wantErr: true,
			errMsg:  "query path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			data := make(map[string]interface{})

			// Copy initial data
			for k, v := range tt.initial {
				data[k] = v
			}

			err := m.DeleteValue(data, tt.expr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(data, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, data)
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	// Test parseValue indirectly through SetValue
	m := New()
	data := map[string]interface{}{}

	testCases := []struct {
		expr     string
		expected interface{}
	}{
		{`.str = "hello"`, "hello"},
		{`.str2 = 'world'`, "world"},
		{`.num = 42`, int64(42)},
		{`.float = 3.14`, 3.14},
		{`.bool = true`, true},
		{`.bool2 = false`, false},
		{`.bare = bareword`, "bareword"},
	}

	for _, tc := range testCases {
		err := m.SetValue(data, tc.expr)
		if err != nil {
			t.Errorf("failed to set value with %s: %v", tc.expr, err)
			continue
		}

		// Extract the key from expression
		key := strings.Trim(strings.Split(tc.expr, " = ")[0], ".")
		if !reflect.DeepEqual(data[key], tc.expected) {
			t.Errorf("for %s: expected %v (%T), got %v (%T)", tc.expr, tc.expected, tc.expected, data[key], data[key])
		}
	}
}

// BenchmarkSetValue benchmarks the SetValue operation with various data types
func BenchmarkSetValue(b *testing.B) {
	m := New()
	data := createComplexTestData()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Create a copy for each iteration to avoid modifying the same data
		testData := make(map[string]interface{})
		for k, v := range data {
			testData[k] = v
		}

		_ = m.SetValue(testData, ".benchmark.key = \"value\"")
	}
}

// BenchmarkSetValueSimple benchmarks simple set operations
func BenchmarkSetValueSimple(b *testing.B) {
	m := New()
	data := make(map[string]interface{})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = m.SetValue(data, ".key = \"value\"")
	}
}

// BenchmarkDeleteValue benchmarks the DeleteValue operation
func BenchmarkDeleteValue(b *testing.B) {
	m := New()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Create fresh data for each iteration
		data := createComplexTestData()
		_ = m.DeleteValue(data, ".database")
	}
}

// BenchmarkDeleteValueSimple benchmarks simple delete operations
func BenchmarkDeleteValueSimple(b *testing.B) {
	m := New()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		data := map[string]interface{}{"key": "value"}
		_ = m.DeleteValue(data, ".key")
	}
}
